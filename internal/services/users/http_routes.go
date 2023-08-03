package users

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"image/png"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/pquerna/otp/totp"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

var _ types.UserDataService = (*service)(nil)

const (
	// UserIDURIParamKey is used to refer to user IDs in router params.
	UserIDURIParamKey = "userID"

	totpIssuer             = "DinnerDoneBetter"
	base64ImagePrefix      = "data:image/jpeg;base64,"
	minimumPasswordEntropy = 60
	totpSecretSize         = 64
	passwordResetTokenSize = 32
)

// validateCredentialsForUpdateRequest takes a user's credentials and determines if they match what is on record.
func (s *service) validateCredentialsForUpdateRequest(ctx context.Context, userID, password, totpToken string) (user *types.User, httpStatus int) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.UserIDKey, userID)

	// fetch user data.
	user, err := s.userDataManager.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, http.StatusNotFound
		}

		logger.Error(err, "error encountered fetching user")
		return nil, http.StatusInternalServerError
	}

	if user.TwoFactorSecretVerifiedAt != nil && totpToken == "" {
		return nil, http.StatusResetContent
	}

	tfs := user.TwoFactorSecret
	if user.TwoFactorSecretVerifiedAt == nil {
		tfs = ""
		totpToken = ""
	}

	// validate login.
	valid, err := s.authenticator.CredentialsAreValid(ctx, user.HashedPassword, password, tfs, totpToken)
	if err != nil {
		logger.WithValue("validation_error", err).Debug("error validating credentials")
		return nil, http.StatusBadRequest
	} else if !valid {
		logger.WithValue("valid", valid).Error(err, "invalid credentials")
		return nil, http.StatusUnauthorized
	}

	return user, http.StatusOK
}

// UsernameSearchHandler is a handler for responding to username queries.
func (s *service) UsernameSearchHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	query := req.URL.Query().Get(types.SearchQueryKey)
	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// fetch user data.
	users, err := s.userDataManager.SearchForUsersByUsername(ctx, query)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}

		observability.AcknowledgeError(err, logger, span, "searching for users")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response.
	s.encoderDecoder.RespondWithData(ctx, res, users)
}

// ListHandler is a handler for responding with a list of users.
func (s *service) ListHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// determine desired filter.
	qf := types.ExtractQueryFilterFromRequest(req)

	// fetch user data.
	users, err := s.userDataManager.GetUsers(ctx, qf)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching users")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response.
	s.encoderDecoder.RespondWithData(ctx, res, users)
}

// CreateHandler is our user creation route.
func (s *service) CreateHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// in the event that we don't want new users to be able to sign up (a config setting) just decline the request from the get-go
	if !s.authSettings.EnableUserSignup || os.Getenv("DISABLE_REGISTRATION") == "true" {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "user creation is disabled", http.StatusForbidden)
		return
	}

	// decode the request.
	registrationInput := new(types.UserRegistrationInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, registrationInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	registrationInput.Username = strings.TrimSpace(registrationInput.Username)
	tracing.AttachUsernameToSpan(span, registrationInput.Username)
	registrationInput.EmailAddress = strings.TrimSpace(strings.ToLower(registrationInput.EmailAddress))
	tracing.AttachEmailAddressToSpan(span, registrationInput.EmailAddress)
	registrationInput.Password = strings.TrimSpace(registrationInput.Password)

	logger = logger.WithValues(map[string]any{
		keys.UsernameKey:                 registrationInput.Username,
		keys.UserEmailAddressKey:         registrationInput.EmailAddress,
		keys.HouseholdInvitationIDKey:    registrationInput.InvitationID,
		keys.HouseholdInvitationTokenKey: registrationInput.InvitationToken,
	})

	if err := registrationInput.ValidateWithContext(ctx, s.authSettings.MinimumUsernameLength, s.authSettings.MinimumPasswordLength); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	// ensure the password is not garbage-tier
	if err := passwordvalidator.Validate(strings.TrimSpace(registrationInput.Password), minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("weak password provided to user creation route")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "password too weak", http.StatusBadRequest)
		return
	}

	var invitation *types.HouseholdInvitation
	if registrationInput.InvitationID != "" && registrationInput.InvitationToken != "" {
		i, err := s.householdInvitationDataManager.GetHouseholdInvitationByTokenAndID(ctx, registrationInput.InvitationToken, registrationInput.InvitationID)
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		} else if err != nil {
			observability.AcknowledgeError(err, logger, span, "retrieving invitation")
			s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
			return
		}

		invitation = i
		logger = logger.WithValue(keys.HouseholdInvitationIDKey, invitation.ID)
		logger.Debug("retrieved household invitation")
	}

	logger.Debug("completed invitation check")

	// hash the password
	hp, err := s.authenticator.HashPassword(ctx, strings.TrimSpace(registrationInput.Password))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// generate a two factor secret.
	tfs, err := s.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating two factor secret")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "internal error", http.StatusInternalServerError)
		return
	}

	input := &types.UserDatabaseCreationInput{
		ID:              identifiers.New(),
		Username:        registrationInput.Username,
		FirstName:       registrationInput.FirstName,
		LastName:        registrationInput.LastName,
		EmailAddress:    registrationInput.EmailAddress,
		HashedPassword:  hp,
		TwoFactorSecret: tfs,
		InvitationToken: registrationInput.InvitationToken,
		Birthday:        registrationInput.Birthday,
		HouseholdName:   registrationInput.HouseholdName,
	}

	if invitation != nil {
		logger.Debug("supplementing user creation input with invitation data")
		input.DestinationHouseholdID = invitation.DestinationHousehold.ID
		input.InvitationToken = invitation.Token
	}

	// create the user.
	user, userCreationErr := s.userDataManager.CreateUser(ctx, input)
	if userCreationErr != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		if errors.Is(userCreationErr, database.ErrUserAlreadyExists) {
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "username taken", http.StatusBadRequest)
			return
		}

		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger.Debug("user created")

	defaultHouseholdID, err := s.householdUserMembershipDataManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching default household ID for user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	emailVerificationToken, emailVerificationTokenErr := s.userDataManager.GetEmailAddressVerificationTokenForUser(ctx, user.ID)
	if emailVerificationTokenErr != nil {
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// notify the relevant parties.
	tracing.AttachUserIDToSpan(span, user.ID)

	dcm := &types.DataChangeMessage{
		HouseholdID:            defaultHouseholdID,
		EventType:              types.UserSignedUpCustomerEventType,
		UserID:                 user.ID,
		EmailVerificationToken: emailVerificationToken,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	if err = s.featureFlagManager.Identify(ctx, user); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user in feature flag manager")
	}

	// UserCreationResponse is a struct we can use to notify the user of their two factor secret, but ideally just this once and then never again.
	ucr := &types.UserCreationResponse{
		CreatedUserID:   user.ID,
		Username:        user.Username,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
		EmailAddress:    user.EmailAddress,
		CreatedAt:       user.CreatedAt,
		TwoFactorSecret: user.TwoFactorSecret,
		Birthday:        user.Birthday,
		TwoFactorQRCode: s.buildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	// encode and peace.
	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, ucr, http.StatusCreated)
}

// buildQRCode builds a QR code for a given username and secret.
func (s *service) buildQRCode(ctx context.Context, username, twoFactorSecret string) string {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.WithValue(keys.UsernameKey, username)

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

// SelfHandler returns information about the user making the request.
func (s *service) SelfHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// figure out who this is all for.
	requester := sessionCtxData.Requester.UserID
	logger = logger.WithValue(keys.RequesterIDKey, requester)
	tracing.AttachRequestingUserIDToSpan(span, requester)

	// fetch user data.
	user, err := s.userDataManager.GetUser(ctx, requester)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no such user")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, user)
}

// PermissionsHandler returns information about the user making the request.
func (s *service) PermissionsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, sessionCtxDataErr := s.sessionContextDataFetcher(req)
	if sessionCtxDataErr != nil {
		observability.AcknowledgeError(sessionCtxDataErr, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// decode the request.
	permissionsInput := new(types.UserPermissionsRequestInput)
	if decodeErr := s.encoderDecoder.DecodeRequest(ctx, req, permissionsInput); decodeErr != nil {
		observability.AcknowledgeError(decodeErr, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	body := &types.UserPermissionsResponse{
		Permissions: make(map[string]bool),
	}

	for _, perm := range permissionsInput.Permissions {
		p := authorization.Permission(perm)
		hasHouseholdPerm := sessionCtxData.HouseholdPermissions[sessionCtxData.ActiveHouseholdID].HasPermission(p)
		hasServicePerm := sessionCtxData.Requester.ServicePermissions.HasPermission(p)
		body.Permissions[perm] = hasHouseholdPerm || hasServicePerm
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, body)
}

// ReadHandler is our read route.
func (s *service) ReadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// figure out who this is all for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	// fetch user data.
	x, err := s.userDataManager.GetUser(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug("no such user")
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// encode response and peace.
	s.encoderDecoder.RespondWithData(ctx, res, x)
}

// TOTPSecretVerificationHandler accepts a TOTP token as input and returns 200 if the TOTP token
// is validated by the user's TOTP secret.
func (s *service) TOTPSecretVerificationHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.TOTPSecretVerificationInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	logger = logger.WithValue(keys.UserIDKey, input.UserID)

	user, fetchUserErr := s.userDataManager.GetUserWithUnverifiedTwoFactorSecret(ctx, input.UserID)
	if fetchUserErr != nil {
		observability.AcknowledgeError(fetchUserErr, logger, span, "fetching user to verify two factor secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	tracing.AttachUserIDToSpan(span, user.ID)
	tracing.AttachUsernameToSpan(span, user.Username)

	if user.TwoFactorSecretVerifiedAt != nil {
		// I suppose if this happens too many times, we might want to keep track of that
		logger.Debug("two factor secret already verified")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "TOTP secret already verified", http.StatusAlreadyReported)
		return
	}

	totpValid := totp.Validate(input.TOTPToken, user.TwoFactorSecret)
	if !totpValid {
		s.encoderDecoder.EncodeInvalidInputResponse(ctx, res)
		return
	}

	if err := s.userDataManager.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		observability.AcknowledgeError(err, logger, span, "verifying user two factor secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.TwoFactorSecretVerifiedCustomerEventType,
		UserID:    user.ID,
	}

	if err := s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// NewTOTPSecretHandler fetches a user, and issues them a new TOTP secret, after validating
// that information received from TOTPSecretRefreshInputContextMiddleware is valid.
func (s *service) NewTOTPSecretHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.TOTPSecretRefreshInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)

	// fetch user
	user, err := s.userDataManager.GetUser(ctx, sessionCtxData.Requester.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}
		observability.AcknowledgeError(err, logger, span, "fetching user from database")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if user.TwoFactorSecretVerifiedAt != nil {
		// validate login.
		valid, validationErr := s.authenticator.CredentialsAreValid(ctx, user.HashedPassword, input.CurrentPassword, user.TwoFactorSecret, input.TOTPToken)
		if validationErr != nil {
			observability.AcknowledgeError(validationErr, logger, span, "validating credentials")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
			return
		} else if !valid {
			observability.AcknowledgeError(validationErr, logger, span, "invalid credentials")
			s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusUnauthorized)
			return
		}
	} else {
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, []byte(""), http.StatusPreconditionFailed)
		return
	}

	// document who this is for.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	tracing.AttachUsernameToSpan(span, user.Username)
	logger = logger.WithValue(keys.UserIDKey, user.ID)

	// set the two factor secret.
	tfs, err := s.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating 2FA secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the user in the database.
	if err = s.userDataManager.MarkUserTwoFactorSecretAsUnverified(ctx, user.ID, tfs); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating 2FA secret")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	user.TwoFactorSecret = tfs
	user.TwoFactorSecretVerifiedAt = nil

	dcm := &types.DataChangeMessage{
		EventType: types.TwoFactorSecretChangedCustomerEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// let the requester know we're all good.
	result := &types.TOTPSecretRefreshResponse{
		TwoFactorSecret: user.TwoFactorSecret,
		TwoFactorQRCode: s.buildQRCode(ctx, user.Username, user.TwoFactorSecret),
	}

	s.encoderDecoder.EncodeResponseWithStatus(ctx, res, result, http.StatusAccepted)
}

// UpdatePasswordHandler updates a user's password.
func (s *service) UpdatePasswordHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.PasswordUpdateInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx, s.authSettings.MinimumPasswordLength); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	// determine relevant user ID.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything's on the up-and-up
	user, httpStatus := s.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)

	// if the above function returns something other than 200, it means some error occurred.
	if httpStatus != http.StatusOK {
		res.WriteHeader(httpStatus)
		return
	}

	tracing.AttachUsernameToSpan(span, user.Username)

	// ensure the password isn't garbage-tier
	if err = passwordvalidator.Validate(input.NewPassword, minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("invalid password provided")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "new password is too weak!", http.StatusBadRequest)
		return
	}

	// hash the new password.
	newPasswordHash, err := s.authenticator.HashPassword(ctx, strings.TrimSpace(input.NewPassword))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "hashing password")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the user.
	if err = s.userDataManager.UpdateUserPassword(ctx, user.ID, newPasswordHash); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.PasswordChangedEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	// we're all good, log the user out
	http.SetCookie(res, &http.Cookie{MaxAge: -1})
	res.WriteHeader(http.StatusAccepted)
}

// UpdateUserEmailAddressHandler updates a user's email address.
func (s *service) UpdateUserEmailAddressHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.UserEmailAddressUpdateInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}
	tracing.AttachEmailAddressToSpan(span, input.NewEmailAddress)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	// determine relevant user ID.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything's on the up-and-up
	user, httpStatus := s.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)

	// if the above function returns something other than 200, it means some error occurred.
	if httpStatus != http.StatusOK {
		res.WriteHeader(httpStatus)
		return
	}

	// update the user.
	if err = s.userDataManager.UpdateUserEmailAddress(ctx, user.ID, input.NewEmailAddress); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.EmailAddressChangedEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// UpdateUserUsernameHandler updates a user's username.
func (s *service) UpdateUserUsernameHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	input := new(types.UsernameUpdateInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided input was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}
	tracing.AttachUsernameToSpan(span, input.NewUsername)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	// determine relevant user ID.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything's on the up-and-up
	user, httpStatus := s.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		input.CurrentPassword,
		input.TOTPToken,
	)

	// if the above function returns something other than 200, it means some error occurred.
	if httpStatus != http.StatusOK {
		res.WriteHeader(httpStatus)
		return
	}

	// update the user.
	if err = s.userDataManager.UpdateUserUsername(ctx, user.ID, input.NewUsername); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.UsernameChangedEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// UpdateUserDetailsHandler updates a user's basic information.
func (s *service) UpdateUserDetailsHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// decode the request.
	providedInput := new(types.UserDetailsUpdateRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, providedInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := providedInput.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("provided providedInput was invalid")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	// determine relevant user ID.
	tracing.AttachRequestingUserIDToSpan(span, sessionCtxData.Requester.UserID)
	logger = sessionCtxData.AttachToLogger(logger)

	// make sure everything's on the up-and-up
	user, httpStatus := s.validateCredentialsForUpdateRequest(
		ctx,
		sessionCtxData.Requester.UserID,
		providedInput.CurrentPassword,
		providedInput.TOTPToken,
	)

	// if the above function returns something other than 200, it means some error occurred.
	if httpStatus != http.StatusOK {
		res.WriteHeader(httpStatus)
		return
	}

	dbInput := converters.ConvertUserDetailsUpdateRequestInputToUserDetailsUpdateInput(providedInput)

	// update the user.
	if err = s.userDataManager.UpdateUserDetails(ctx, user.ID, dbInput); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.UserDetailsChangedEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// AvatarUploadHandler updates a user's avatar.
func (s *service) AvatarUploadHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	// decode the request.
	input := new(types.AvatarUpdateInput)
	if err = s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	tracing.AttachSessionContextDataToSpan(span, sessionCtxData)
	logger = sessionCtxData.AttachToLogger(logger)
	logger.Debug("session context data extracted")

	user, err := s.userDataManager.GetUser(ctx, sessionCtxData.Requester.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching associated user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)
	logger.Debug("retrieved user from database")

	if err = s.userDataManager.UpdateUserAvatar(ctx, user.ID, input.Base64EncodedData); err != nil {
		observability.AcknowledgeError(err, logger, span, "updating user info")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}
}

// ArchiveHandler is a handler for archiving a user.
func (s *service) ArchiveHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	// figure out who this is for.
	userID := s.userIDFetcher(req)
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	// do the deed.
	err := s.userDataManager.ArchiveUser(ctx, userID)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "archiving user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// we're all good.
	res.WriteHeader(http.StatusNoContent)
}

// RequestUsernameReminderHandler checks for a user with a given email address and notifies them via email if there is a username associated with it.
func (s *service) RequestUsernameReminderHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.UsernameReminderRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	u, err := s.userDataManager.GetUserByEmail(ctx, input.EmailAddress)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.UsernameReminderRequestedEventType,
		UserID:    u.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// CreatePasswordResetTokenHandler rotates the cookie building secret with a new random secret.
func (s *service) CreatePasswordResetTokenHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.PasswordResetTokenCreationRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := s.secretGenerator.GenerateBase32EncodedString(ctx, passwordResetTokenSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating secret")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := s.userDataManager.GetUserByEmail(ctx, input.EmailAddress)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	dbInput := &types.PasswordResetTokenDatabaseCreationInput{
		ID:            identifiers.New(),
		Token:         token,
		BelongsToUser: u.ID,
		ExpiresAt:     time.Now().Add(30 * time.Minute),
	}

	t, err := s.passwordResetTokenDataManager.CreatePasswordResetToken(ctx, dbInput)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating password reset token")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:          types.PasswordResetTokenCreatedEventType,
		UserID:             u.ID,
		PasswordResetToken: t,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// PasswordResetTokenRedemptionHandler rotates the cookie building secret with a new random secret.
func (s *service) PasswordResetTokenRedemptionHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.PasswordResetTokenRedemptionRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	t, err := s.passwordResetTokenDataManager.GetPasswordResetTokenByToken(ctx, input.Token)
	if errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusNotFound)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching password reset token")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := s.userDataManager.GetUser(ctx, t.BelongsToUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}

		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	// ensure the password isn't garbage-tier
	if err = passwordvalidator.Validate(strings.TrimSpace(input.NewPassword), minimumPasswordEntropy); err != nil {
		logger.WithValue("password_validation_error", err).Debug("invalid password provided")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "new password is too weak!", http.StatusBadRequest)
		return
	}

	// hash the new password.
	newPasswordHash, err := s.authenticator.HashPassword(ctx, strings.TrimSpace(input.NewPassword))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "hashing password")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	// update the user.
	if err = s.userDataManager.UpdateUserPassword(ctx, u.ID, newPasswordHash); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}

		observability.AcknowledgeError(err, logger, span, "updating user")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	if redemptionErr := s.passwordResetTokenDataManager.RedeemPasswordResetToken(ctx, t.ID); redemptionErr != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}

		observability.AcknowledgeError(err, logger, span, "redeeming password reset token")
		s.encoderDecoder.EncodeUnspecifiedInternalServerErrorResponse(ctx, res)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.PasswordResetTokenRedeemedEventType,
		UserID:    t.BelongsToUser,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// VerifyUserEmailAddressHandler checks for a user with a given email address and notifies them via email if there is a username associated with it.
func (s *service) VerifyUserEmailAddressHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	input := new(types.EmailAddressVerificationRequestInput)
	if err := s.encoderDecoder.DecodeRequest(ctx, req, input); err != nil {
		observability.AcknowledgeError(err, logger, span, "decoding request body")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, "invalid request content", http.StatusBadRequest)
		return
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		logger.WithValue(keys.ValidationErrorKey, err).Debug("invalid input attached to request")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.userDataManager.GetUserByEmailAddressVerificationToken(ctx, input.Token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}

		observability.AcknowledgeError(err, logger, span, "fetching user")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = s.userDataManager.MarkUserEmailAddressAsVerified(ctx, user.ID, input.Token); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.encoderDecoder.EncodeNotFoundResponse(ctx, res)
			return
		}

		observability.AcknowledgeError(err, logger, span, "marking user email as verified")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType: types.UserEmailAddressVerifiedEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	res.WriteHeader(http.StatusAccepted)
}

// RequestEmailVerificationEmailHandler submits a request for an email verification email.
func (s *service) RequestEmailVerificationEmailHandler(res http.ResponseWriter, req *http.Request) {
	ctx, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	sessionCtxData, err := s.sessionContextDataFetcher(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "retrieving session context data")
		s.encoderDecoder.EncodeUnauthorizedResponse(ctx, res)
		return
	}

	verificationToken, err := s.userDataManager.GetEmailAddressVerificationTokenForUser(ctx, sessionCtxData.Requester.UserID)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		s.encoderDecoder.EncodeResponseWithStatus(ctx, res, nil, http.StatusAccepted)
		return
	} else if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching email address verification token")
		s.encoderDecoder.EncodeErrorResponse(ctx, res, err.Error(), http.StatusInternalServerError)
		return
	}

	dcm := &types.DataChangeMessage{
		EventType:              types.UserEmailAddressVerificationEmailRequestedEventType,
		UserID:                 sessionCtxData.Requester.UserID,
		EmailVerificationToken: verificationToken,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	logger.WithValue("data_change_message", dcm).Debug("published data change message")

	res.WriteHeader(http.StatusAccepted)
}
