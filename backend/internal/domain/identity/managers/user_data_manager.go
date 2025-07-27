package managers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"image/png"
	strings "strings"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

/*

TODO:
- all string and pointered params are checked for emptiness and return errors
- all methods write a data change message
- all values are observed in both logs and analytics messages
- defaults on query filters
- input parameters are validated and errors returned
- all parameters accounted for in o11y

*/

const (
	o11yName = "identity_data_manager"

	// UserIDURIParamKey is used to refer to user IDs in router params.
	UserIDURIParamKey = "userID"

	totpIssuer             = "DinnerDoneBetter"
	base64ImagePrefix      = "data:image/jpeg;base64,"
	minimumPasswordEntropy = 60
	totpSecretSize         = 64
	passwordResetTokenSize = 32
)

type (
	IdentityDataManager struct {
		tracer               tracing.Tracer
		logger               logging.Logger
		dataChangesPublisher messagequeue.Publisher
		identityRepo         identity.Repository
		secretGenerator      random.Generator
		authenticator        authentication.Hasher
		userSearchIndex      indexing.UserTextSearcher
	}

	Config struct {
		DataChangesTopicName string `json:"dataChangesTopicName" env:"DATA_CHANGES_TOPIC_NAME"`
	}
)

func NewIdentityDataManager(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	publisherProvider messagequeue.PublisherProvider,
	identityRepo identity.Repository,
	cfg *Config,
) (*IdentityDataManager, error) {
	publisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up data changes publisher: %w", err)
	}

	return &IdentityDataManager{
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:               logging.EnsureLogger(logger).WithName(o11yName),
		identityRepo:         identityRepo,
		dataChangesPublisher: publisher,
	}, nil
}

func (i *IdentityDataManager) AcceptAccountInvitation(ctx context.Context, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}, span, i.logger)

	if err := i.identityRepo.AcceptAccountInvitation(ctx, accountInvitationID, input.Token, input.Note); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "accepting account invitation")
	}

	i.dataChangesPublisher.PublishAsync(ctx, audit.BuildDataChangeMessageFromContext(ctx, logger, identity.AccountInvitationCanceledServiceEventType, map[string]any{
		keys.AccountInvitationIDKey: accountInvitationID,
	}))

	return nil
}

func (i *IdentityDataManager) ArchiveAccount(ctx context.Context, accountID, ownerID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, ownerID)

	if err := i.identityRepo.ArchiveAccount(ctx, accountID, ownerID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "ArchiveAccount")
	}

	return nil
}

func (i *IdentityDataManager) ArchiveUserMembership(ctx context.Context, userID, accountID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	if err := i.identityRepo.RemoveUserFromAccount(ctx, userID, accountID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "RemoveUserFromAccount")
	}

	return nil
}

func (i *IdentityDataManager) ArchiveUser(ctx context.Context, userID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	if err := i.identityRepo.ArchiveUser(ctx, userID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "ArchiveUser")
	}

	return nil
}

func (i *IdentityDataManager) CancelAccountInvitation(ctx context.Context, accountInvitationID, note string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	if err := i.identityRepo.CancelAccountInvitation(ctx, accountInvitationID, note); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating AccountInvitation")
	}

	return nil
}

func (i *IdentityDataManager) CreateAccount(ctx context.Context, input *identity.AccountCreationRequestInput) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	created, err := i.identityRepo.CreateAccount(ctx, converters.ConvertAccountCreationInputToAccountDatabaseCreationInput(input))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating Account")
	}

	return created, nil
}

func (i *IdentityDataManager) CreateAccountInvitation(ctx context.Context, input *identity.AccountInvitationCreationRequestInput) (*identity.AccountInvitation, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	created, err := i.identityRepo.CreateAccountInvitation(ctx, converters.ConvertAccountInvitationCreationInputToAccountInvitationDatabaseCreationInput(input))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating AccountInvitation")
	}

	return created, nil
}

func (i *IdentityDataManager) CreateUser(ctx context.Context, registrationInput *identity.UserRegistrationInput) (*identity.UserCreationResponse, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	registrationInput.Username = strings.TrimSpace(registrationInput.Username)
	tracing.AttachToSpan(span, keys.UsernameKey, registrationInput.Username)
	registrationInput.EmailAddress = strings.TrimSpace(strings.ToLower(registrationInput.EmailAddress))
	tracing.AttachToSpan(span, keys.UserEmailAddressKey, registrationInput.EmailAddress)
	registrationInput.Password = strings.TrimSpace(registrationInput.Password)

	logger = logger.WithValues(map[string]any{
		keys.UsernameKey:               registrationInput.Username,
		keys.UserEmailAddressKey:       registrationInput.EmailAddress,
		keys.AccountInvitationIDKey:    registrationInput.InvitationID,
		keys.AccountInvitationTokenKey: registrationInput.InvitationToken,
	})

	if err := registrationInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		return nil, observability.PrepareError(err, span, "invalid user creation input provided")
	}

	// ensure the password is not garbage-tier
	if err := passwordvalidator.Validate(strings.TrimSpace(registrationInput.Password), minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("weak password provided to user creation route")
		return nil, observability.PrepareAndLogError(err, logger, span, "weak password provided for user creation")
	}

	var invitation *identity.AccountInvitation
	if registrationInput.InvitationID != "" && registrationInput.InvitationToken != "" {
		invite, err := i.identityRepo.GetAccountInvitationByTokenAndID(ctx, registrationInput.InvitationToken, registrationInput.InvitationID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareAndLogError(err, logger, span, "no account invitation found")
		} else if err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "getting account invitation by token and ID")
		}

		invitation = invite
		logger.Debug("retrieved account invitation")
	}

	logger.Debug("completed invitation check")

	// hash the password
	hp, err := i.authenticator.HashPassword(ctx, strings.TrimSpace(registrationInput.Password))
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "hashing user creation password")
	}

	// generate a two factor secret.
	tfs, err := i.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "generating two factor secret")
	}

	input := &identity.UserDatabaseCreationInput{
		ID:              identifiers.New(),
		Username:        registrationInput.Username,
		FirstName:       registrationInput.FirstName,
		LastName:        registrationInput.LastName,
		EmailAddress:    registrationInput.EmailAddress,
		HashedPassword:  hp,
		TwoFactorSecret: tfs,
		InvitationToken: registrationInput.InvitationToken,
		Birthday:        registrationInput.Birthday,
		AccountName:     registrationInput.AccountName,
	}

	if invitation != nil {
		logger.Debug("supplementing user creation input with invitation data")
		input.DestinationAccountID = invitation.DestinationAccount.ID
		input.InvitationToken = invitation.Token
	}

	// create the user.
	user, err := i.identityRepo.CreateUser(ctx, input)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		if errors.Is(err, database.ErrUserAlreadyExists) {
			return nil, observability.PrepareAndLogError(err, logger, span, "user already exists")
		}
		return nil, observability.PrepareAndLogError(err, logger, span, "creating user in database")
	}

	logger.Debug("user created")

	defaultAccountID, err := i.identityRepo.GetDefaultAccountIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching default account ID for user")
	}

	emailVerificationToken, err := i.identityRepo.GetEmailAddressVerificationTokenForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching email verification token for user")
	}

	// notify the relevant parties.
	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)

	dcm := &audit.DataChangeMessage{
		AccountID: defaultAccountID,
		EventType: identity.UserSignedUpServiceEventType,
		UserID:    user.ID,
		Context: map[string]any{
			keys.UserEmailVerificationTokenKey: emailVerificationToken,
		},
	}

	if err = i.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	/* TODO: this should be done in a downstream handler

	if err = i.analyticsReporter.AddUser(ctx, user.ID, map[string]any{
		"username":        user.Username,
		"default_account": defaultAccountID,
		"first_name":      user.FirstName,
		"last_name":       user.LastName,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user for analytics")
	}

	if err = s.featureFlagManager.Identify(ctx, user); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user in feature flag manager")
	}
	*/

	// UserCreationResponse is a struct we can use to notify the user of their two factor secret, but ideally just this once and then never again.
	ucr := &identity.UserCreationResponse{
		CreatedUserID:   user.ID,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		EmailAddress:    user.EmailAddress,
		CreatedAt:       user.CreatedAt,
		TwoFactorSecret: user.TwoFactorSecret,
		Birthday:        user.Birthday,
		TwoFactorQRCode: i.buildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	return ucr, nil
}

func (i *IdentityDataManager) GetAccount(ctx context.Context, accountID string) (*identity.Account, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.AccountIDKey, accountID)

	account, err := i.identityRepo.GetAccount(ctx, accountID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching account")
	}

	return account, nil
}

func (i *IdentityDataManager) GetAccountInvitation(ctx context.Context, accountID, accountInvitationID string) (*identity.AccountInvitation, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.AccountIDKey, accountID).WithValue(keys.AccountInvitationIDKey, accountInvitationID)

	invitation, err := i.identityRepo.GetAccountInvitationByAccountAndID(ctx, accountID, accountInvitationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting invitation")
	}

	return invitation, nil
}

func (i *IdentityDataManager) GetAccounts(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.Account, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	accounts, err := i.identityRepo.GetAccounts(ctx, userID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching accounts")
	}

	return accounts.Data, "TODO", nil
}

func (i *IdentityDataManager) GetReceivedAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.AccountInvitation, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	invites, err := i.identityRepo.GetPendingAccountInvitationsForUser(ctx, userID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching invites")
	}

	return invites.Data, "TODO", nil
}

func (i *IdentityDataManager) GetSentAccountInvitations(ctx context.Context, userID string, filter *filtering.QueryFilter) ([]*identity.AccountInvitation, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	invites, err := i.identityRepo.GetPendingAccountInvitationsFromUser(ctx, userID, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching invites")
	}

	return invites.Data, "TODO", nil
}

func (i *IdentityDataManager) GetUser(ctx context.Context, userID string) (*identity.User, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	user, err := i.identityRepo.GetUser(ctx, userID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting user")
	}

	return user, nil
}

func (i *IdentityDataManager) RejectAccountInvitation(ctx context.Context, accountInvitationID string, input *identity.AccountInvitationUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "invalid input attached to request")
	}

	// note, this is where you would call input.ValidateWithContext, if that currently had any effect.

	invitation, err := i.identityRepo.GetAccountInvitationByTokenAndID(ctx, input.Token, accountInvitationID)
	if errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "account invitation not found")
	} else if err != nil {
		return observability.PrepareError(err, span, "retrieving invitation")
	}

	if err = i.identityRepo.RejectAccountInvitation(ctx, invitation.ID, input.Note); err != nil {
		return observability.PrepareError(err, span, "rejecting invitation")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.AccountInvitationRejectedServiceEventType,
		AccountID: invitation.DestinationAccount.ID,
		UserID:    "TODO",
		Context: map[string]any{
			keys.AccountInvitationIDKey: accountInvitationID,
		},
	}

	if err = i.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	return nil
}

func (i *IdentityDataManager) GetUsers(ctx context.Context, filter *filtering.QueryFilter) ([]*identity.User, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	users, err := i.identityRepo.GetUsers(ctx, filter)
	if err != nil {
		return nil, "", observability.PrepareAndLogError(err, logger, span, "fetching users")
	}

	return users.Data, "", nil
}

func (i *IdentityDataManager) SearchForUsers(ctx context.Context, query string, useDatabase bool, filter *filtering.QueryFilter) ([]*identity.User, string, error) {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	if useDatabase {
		users, err := i.identityRepo.SearchForUsersByUsername(ctx, query, filter)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, "", observability.PrepareError(err, span, "no users found")
			}
			return nil, "", observability.PrepareError(err, span, "searching for users")
		}

		return users, "TODO", nil
	} else {
		uss, err := i.userSearchIndex.Search(ctx, query)
		if err != nil {
			return nil, "", observability.PrepareError(err, span, "searching for users")
		}

		userIDs := []string{}
		for _, us := range uss {
			userIDs = append(userIDs, us.ID)
		}

		users, err := i.identityRepo.GetUsersWithIDs(ctx, userIDs)
		if err != nil {
			return nil, "", observability.PrepareError(err, span, "searching for users")
		}

		return users, "TODO", nil
	}
}

func (i *IdentityDataManager) SetDefaultAccount(ctx context.Context, userID, accountID string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID).WithValue(keys.AccountIDKey, accountID)

	// mark household as default in database.
	if err := i.identityRepo.MarkAccountAsUserDefault(ctx, userID, accountID); err != nil {
		return observability.PrepareError(err, span, "marking default account as user")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.AccountMemberRemovedServiceEventType,
		AccountID: accountID,
		UserID:    "TODO",
	}
	if err := i.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	return nil
}

func (i *IdentityDataManager) TransferAccountOwnership(ctx context.Context, accountID string, input *identity.AccountOwnershipTransferInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "")
	}

	// transfer ownership of household in database.
	if err := i.identityRepo.TransferAccountOwnership(ctx, accountID, input); err != nil {
		return observability.PrepareError(err, span, "")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.AccountOwnershipTransferredServiceEventType,
		AccountID: accountID,
		UserID:    "TODO",
	}
	if err := i.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created household")
	}

	return nil
}

func (i *IdentityDataManager) UpdateAccount(ctx context.Context, accountID string, input *identity.AccountUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating account update")
	}

	// determine account ID.
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	// fetch account from database.
	account, err := i.identityRepo.GetAccount(ctx, accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return observability.PrepareError(err, span, "no account found")
	} else if err != nil {
		return observability.PrepareError(err, span, "fetching account")
	}

	// update the data structure.
	account.Update(input)

	// update account in database.
	if err = i.identityRepo.UpdateAccount(ctx, account); err != nil {
		return observability.PrepareError(err, span, "updating account")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.AccountUpdatedServiceEventType,
		AccountID: account.ID,
		UserID:    "TODO",
	}

	if err = i.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created account")
	}

	return nil
}

func (i *IdentityDataManager) UpdateAccountMemberPermissions(ctx context.Context, userID, accountID string, input *identity.ModifyUserPermissionsInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithSpan(span).WithValue(keys.UserIDKey, userID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "invalid input attached to request")
	}

	// create account in database.
	if err := i.identityRepo.ModifyUserPermissions(ctx, accountID, userID, input); err != nil {
		return observability.PrepareError(err, span, "modifying user permissions")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.AccountMembershipPermissionsUpdatedServiceEventType,
		AccountID: accountID,
		UserID:    "TODO",
	}
	if err := i.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message for created account")
	}

	return nil
}

func (i *IdentityDataManager) UpdateUserDetails(ctx context.Context, input *identity.UserDetailsUpdateRequestInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateUserEmailAddress(ctx context.Context, newEmail string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UpdateUserUsername(ctx context.Context, newUsername string) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) UploadUserAvatar(ctx context.Context, newAvatar []byte) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

func (i *IdentityDataManager) AdminUpdateUserStatus(ctx context.Context, input identity.UserAccountStatusUpdateInput) error {
	ctx, span := i.tracer.StartSpan(ctx)
	defer span.End()

	return nil
}

// buildQRCode builds a QR code for a given username and secret.
func (i *IdentityDataManager) buildQRCode(ctx context.Context, username, twoFactorSecret string) string {
	_, span := i.tracer.StartSpan(ctx)
	defer span.End()

	logger := i.logger.WithValue(keys.UsernameKey, username)

	// "otpauth://totp/{{ .Issuer }}:{{ .EnsureUsername }}?secret={{ .Secret }}&issuer={{ .Issuer }}",
	otpString := fmt.Sprintf(
		"otpauth://totp/%s:%s?secret=%s&issuer=%s",
		totpIssuer,
		username,
		twoFactorSecret,
		totpIssuer,
	)

	// encode two factor secret as authenticator-friendly QR code
	qrCode, err := qr.Encode(otpString, qr.L, qr.Auto)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding OTP string")
		return ""
	}

	// scale the QR code so that it's not a PNG for ants.
	qrCode, err = barcode.Scale(qrCode, 256, 256)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "scaling QR code")
		return ""
	}

	// encode the QR code to PNG.
	var b bytes.Buffer
	if err = png.Encode(&b, qrCode); err != nil {
		observability.AcknowledgeError(err, logger, span, "encoding QR code to PNG")
		return ""
	}

	// base64 encode the image for easy HTML use.
	return fmt.Sprintf("%s%s", base64ImagePrefix, base64.StdEncoding.EncodeToString(b.Bytes()))
}
