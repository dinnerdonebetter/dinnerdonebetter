package managers

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/qrcodes"
	"github.com/dinnerdonebetter/backend/internal/platform/random"

	"github.com/pquerna/otp/totp"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const (
	passwordResetTokenSize = 32
)

type AuthManager struct {
	passwordResetTokenDataManager identity.PasswordResetTokenDataManager
	userDataManager               identity.UserDataManager
	tracer                        tracing.Tracer
	authenticator                 authentication.Authenticator
	logger                        logging.Logger
	dataChangesPublisher          messagequeue.Publisher
	secretGenerator               random.Generator
	qrCodeBuilder                 qrcodes.Builder
	sessionContextDataFetcher     func(context.Context) (*sessions.ContextData, error)
	minimumPasswordLength         uint8
}

func (l *AuthManager) Self(ctx context.Context) (*identity.User, error) {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "failed to get session context data")
	}
	tracing.AttachSessionContextDataToSpan(span, sessionContextData)
	logger := sessionContextData.AttachToLogger(l.logger)

	// figure out who this is all for.
	requester := sessionContextData.GetUserID()
	tracing.AttachToSpan(span, keys.RequesterIDKey, requester)

	// fetch user data.
	user, err := l.userDataManager.GetUser(ctx, requester)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no such user")
		return nil, observability.PrepareError(err, span, "no such user")
	} else if err != nil {
		return nil, observability.PrepareError(err, span, "fetching user")
	}

	return user, nil
}

func (l *AuthManager) UserPermissions(ctx context.Context, input *identity.UserPermissionsRequestInput) (*identity.UserPermissionsResponse, error) {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	body := &identity.UserPermissionsResponse{
		Permissions: make(map[string]bool),
	}

	for _, perm := range input.Permissions {
		p := authorization.Permission(perm)
		hasAccountPerm := sessionContextData.AccountPermissions[sessionContextData.ActiveAccountID].HasPermission(p)
		hasServicePerm := sessionContextData.Requester.ServicePermissions.HasPermission(p)
		body.Permissions[perm] = hasAccountPerm || hasServicePerm
	}

	return body, nil
}

func (l *AuthManager) TOTPSecretVerification(ctx context.Context, input *identity.TOTPSecretVerificationInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		logger = logger.WithValue(keys.ValidationErrorKey, err)
		return observability.PrepareError(err, span, "provided input was invalid")
	}

	logger = logger.WithValue(keys.UserIDKey, input.UserID)
	logger.Info("validated input, getting user")

	user, err := l.userDataManager.GetUserWithUnverifiedTwoFactorSecret(ctx, input.UserID)
	if err != nil {
		return observability.PrepareError(err, span, "fetching user to verify two factor secret")
	}

	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, keys.UsernameKey, user.Username)
	logger = logger.WithValue(keys.UsernameKey, user.Username)

	if user.TwoFactorSecretVerifiedAt != nil {
		// I suppose if this happens too many times, we might want to keep track of that?
		return errors.New("two factor secret already verified")
	}

	if totpValid := totp.Validate(input.TOTPToken, user.TwoFactorSecret); !totpValid {
		return observability.PrepareError(err, span, "TOTP code was invalid")
	}

	if err = l.userDataManager.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		return observability.PrepareError(err, span, "verifying user two factor secret")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.TwoFactorSecretVerifiedServiceEventType,
		UserID:    user.ID,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) NewTOTPSecret(ctx context.Context, input *identity.TOTPSecretRefreshInput) (*identity.TOTPSecretRefreshResponse, error) {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "provided input was invalid")
	}

	sessionCtxData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return nil, observability.PrepareError(err, span, "retrieving session context data")
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// fetch user
	user, err := l.userDataManager.GetUser(ctx, sessionCtxData.Requester.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "user does not exist")
		}
		return nil, observability.PrepareError(err, span, "retrieving user from database")
	}

	if user.TwoFactorSecretVerifiedAt != nil {
		// validate login.
		valid, validationErr := l.authenticator.CredentialsAreValid(ctx, user.HashedPassword, input.CurrentPassword, user.TwoFactorSecret, input.TOTPToken)
		if validationErr != nil {
			return nil, observability.PrepareError(validationErr, span, "validating credentials")
		} else if !valid {
			return nil, observability.PrepareError(validationErr, span, "invalid credentials")
		}
	} else {
		return nil, observability.PrepareError(err, span, "two factor secret not yet verified")
	}

	// document who this is for.
	tracing.AttachToSpan(span, keys.RequesterIDKey, sessionCtxData.Requester.UserID)
	tracing.AttachToSpan(span, keys.UsernameKey, user.Username)
	logger = logger.WithValue(keys.UserIDKey, user.ID)

	// set the two factor secret.
	tfs, err := l.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		return nil, observability.PrepareError(err, span, "generating 2FA secret")
	}

	// update the user in the database.
	if err = l.userDataManager.MarkUserTwoFactorSecretAsUnverified(ctx, user.ID, tfs); err != nil {
		return nil, observability.PrepareError(err, span, "updating 2FA secret")
	}

	user.TwoFactorSecret = tfs
	user.TwoFactorSecretVerifiedAt = nil

	dcm := &audit.DataChangeMessage{
		EventType: identity.TwoFactorSecretChangedServiceEventType,
		UserID:    user.ID,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	result := &identity.TOTPSecretRefreshResponse{
		TwoFactorSecret: user.TwoFactorSecret,
		TwoFactorQRCode: l.qrCodeBuilder.BuildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	return result, nil
}

func (l *AuthManager) UpdatePassword(ctx context.Context, input *identity.PasswordUpdateInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx, l.minimumPasswordLength); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}

	sessionCtxData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "retrieving session context data")
	}

	// determine relevant user ID.
	tracing.AttachToSpan(span, keys.RequesterIDKey, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything'l on the up-and-up
	user, err := l.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)
	if err != nil {
		return observability.PrepareError(err, span, "validating credentials")
	}
	tracing.AttachToSpan(span, keys.UsernameKey, user.Username)

	// ensure the password isn't garbage-tier
	if err = passwordvalidator.Validate(input.NewPassword, minimumPasswordEntropy); err != nil {
		return observability.PrepareError(err, span, "invalid password provided")
	}

	// hash the new password.
	newPasswordHash, err := l.authenticator.HashPassword(ctx, strings.TrimSpace(input.NewPassword))
	if err != nil {
		return observability.PrepareError(err, span, "hashing password")
	}

	// update the user.
	if err = l.userDataManager.UpdateUserPassword(ctx, user.ID, newPasswordHash); err != nil {
		return observability.PrepareError(err, span, "updating user")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.PasswordChangedEventType,
		UserID:    user.ID,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) UpdateUserEmailAddress(ctx context.Context, input *identity.UserEmailAddressUpdateInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}
	tracing.AttachToSpan(span, keys.UserEmailAddressKey, input.NewEmailAddress)

	sessionCtxData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "retrieving session context data")
	}

	// determine relevant user ID.
	tracing.AttachToSpan(span, keys.RequesterIDKey, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything'l on the up-and-up
	user, err := l.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)
	if err != nil {
		return observability.PrepareError(err, span, "validating credentials")
	}

	// update the user.
	if err = l.userDataManager.UpdateUserEmailAddress(ctx, user.ID, input.NewEmailAddress); err != nil {
		return observability.PrepareError(err, span, "updating user")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.EmailAddressChangedEventType,
		UserID:    user.ID,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) UpdateUserUsername(ctx context.Context, input *identity.UsernameUpdateInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}
	tracing.AttachToSpan(span, keys.UsernameKey, input.NewUsername)

	sessionCtxData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "retrieving session context data")
	}

	// determine relevant user ID.
	tracing.AttachToSpan(span, keys.RequesterIDKey, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything'l on the up-and-up
	user, err := l.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)
	if err != nil {
		return observability.PrepareError(err, span, "validating credentials")
	}

	// update the user.
	if err = l.userDataManager.UpdateUserUsername(ctx, user.ID, input.NewUsername); err != nil {
		return observability.PrepareError(err, span, "updating user")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UsernameChangedEventType,
		UserID:    user.ID,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) RequestUsernameReminder(ctx context.Context, input *identity.UsernameReminderRequestInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}

	u, err := l.userDataManager.GetUserByEmail(ctx, input.EmailAddress)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "user not found")
	} else if err != nil {
		return observability.PrepareError(err, span, "fetching user")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UsernameReminderRequestedEventType,
		UserID:    u.ID,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) CreatePasswordResetToken(ctx context.Context, input *identity.PasswordResetTokenCreationRequestInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}

	token, err := l.secretGenerator.GenerateBase32EncodedString(ctx, passwordResetTokenSize)
	if err != nil {
		return observability.PrepareError(err, span, "generating secret")
	}

	u, err := l.userDataManager.GetUserByEmail(ctx, input.EmailAddress)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "user not found")
	} else if err != nil {
		return observability.PrepareError(err, span, "fetching user")
	}

	dbInput := &identity.PasswordResetTokenDatabaseCreationInput{
		ID:            identifiers.New(),
		Token:         token,
		BelongsToUser: u.ID,
		ExpiresAt:     time.Now().Add(30 * time.Minute),
	}

	t, err := l.passwordResetTokenDataManager.CreatePasswordResetToken(ctx, dbInput)
	if err != nil {
		return observability.PrepareError(err, span, "creating password reset token")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.PasswordResetTokenCreatedEventType,
		UserID:    u.ID,
		Context: map[string]any{
			keys.PasswordResetTokenIDKey: t.ID,
		},
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) PasswordResetTokenRedemption(ctx context.Context, input *identity.PasswordResetTokenRedemptionRequestInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}

	t, err := l.passwordResetTokenDataManager.GetPasswordResetTokenByToken(ctx, input.Token)
	if errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "password reset token not found")
	} else if err != nil {
		return observability.PrepareError(err, span, "fetching password reset token")
	}

	u, err := l.userDataManager.GetUser(ctx, t.BelongsToUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return observability.PrepareError(err, span, "user not found")
		}
		return observability.PrepareError(err, span, "fetching user")
	}

	// ensure the password isn't garbage-tier
	if err = passwordvalidator.Validate(strings.TrimSpace(input.NewPassword), minimumPasswordEntropy); err != nil {
		return observability.PrepareError(err, span, "provided password was invalid")
	}

	// hash the new password.
	newPasswordHash, err := l.authenticator.HashPassword(ctx, strings.TrimSpace(input.NewPassword))
	if err != nil {

		return observability.PrepareError(err, span, "hashing password")
	}

	// update the user.
	if err = l.userDataManager.UpdateUserPassword(ctx, u.ID, newPasswordHash); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		if errors.Is(err, sql.ErrNoRows) {
			return observability.PrepareError(err, span, "user not found")
		}

		return observability.PrepareError(err, span, "updating user")
	}

	if redemptionErr := l.passwordResetTokenDataManager.RedeemPasswordResetToken(ctx, t.ID); redemptionErr != nil {
		observability.AcknowledgeError(err, logger, span, "redeeming password reset token")
		if errors.Is(err, sql.ErrNoRows) {
			return observability.PrepareError(err, span, "redeeming password reset token not found")
		}

		return observability.PrepareError(err, span, "redeeming password reset token")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.PasswordResetTokenRedeemedEventType,
		UserID:    t.BelongsToUser,
	}

	l.dataChangesPublisher.PublishAsync(ctx, dcm)

	return nil
}

func (l *AuthManager) RequestEmailVerificationEmail(ctx context.Context) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	verificationToken, err := l.userDataManager.GetEmailAddressVerificationTokenForUser(ctx, sessionContextData.Requester.UserID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "email verification token not found")
	} else if err != nil {
		return observability.PrepareError(err, span, "fetching email verification token")
	}

	l.dataChangesPublisher.PublishAsync(ctx, &audit.DataChangeMessage{
		EventType: identity.UserEmailAddressVerificationEmailRequestedEventType,
		UserID:    sessionContextData.Requester.UserID,
		Context: map[string]any{
			keys.UserEmailVerificationTokenKey: verificationToken,
		},
	})

	return nil
}

func (l *AuthManager) VerifyUserEmailAddress(ctx context.Context, input *identity.EmailAddressVerificationRequestInput) error {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithSpan(span)

	sessionContextData, err := l.sessionContextDataFetcher(ctx)
	if err != nil {
		return observability.PrepareError(err, span, "failed to get session context data")
	}
	logger = logger.WithValue(keys.UserIDKey, sessionContextData.GetUserID())

	if err = input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "provided input was invalid")
	}

	user, err := l.userDataManager.GetUserByEmailAddressVerificationToken(ctx, input.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return observability.PrepareError(err, span, "user not found")
		}
		return observability.PrepareError(err, span, "fetching user")
	}

	if err = l.userDataManager.MarkUserEmailAddressAsVerified(ctx, user.ID, input.Token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return observability.PrepareError(err, span, "user not found")
		}
		return observability.PrepareError(err, span, "marking user email as verified")
	}

	l.dataChangesPublisher.PublishAsync(ctx, &audit.DataChangeMessage{
		EventType: identity.UserEmailAddressVerifiedEventType,
		UserID:    user.ID,
	})

	return nil
}

// validateCredentialsForUpdateRequest takes a user's credentials and determines if they match what is on record.
func (l *AuthManager) validateCredentialsForUpdateRequest(ctx context.Context, userID, password, totpToken string) (*identity.User, error) {
	ctx, span := l.tracer.StartSpan(ctx)
	defer span.End()

	logger := l.logger.WithValue(keys.UserIDKey, userID)

	// fetch user data.
	user, err := l.userDataManager.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, err
		}

		logger.Error("error encountered fetching user", err)
		return nil, observability.PrepareError(err, span, "fetching user")
	}

	if user.TwoFactorSecretVerifiedAt != nil && totpToken == "" {
		return nil, observability.PrepareError(err, span, "two factor secret not provided")
	}

	tfs := user.TwoFactorSecret
	if user.TwoFactorSecretVerifiedAt == nil {
		tfs = ""
		totpToken = ""
	}

	// validate login.
	valid, err := l.authenticator.CredentialsAreValid(ctx, user.HashedPassword, password, tfs, totpToken)
	if err != nil {
		return nil, observability.PrepareError(err, span, "error validating credentials")
	} else if !valid {
		return nil, observability.PrepareError(err, span, "credentials are not valid")
	}

	return user, nil
}
