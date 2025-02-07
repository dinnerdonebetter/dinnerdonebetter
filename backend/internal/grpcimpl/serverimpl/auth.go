package serverimpl

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpcimpl/converters"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pquerna/otp/totp"
	passwordvalidator "github.com/wagslane/go-password-validator"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const (
	minimumPasswordEntropy = 60
	totpSecretSize         = 64
	passwordResetTokenSize = 32
)

type (
	contextKey string
)

const (
	SessionContextKey contextKey = "session_context"
)

func (s *Server) fetchSessionContext(ctx context.Context) *sessions.ContextData {
	sessionContext, ok := ctx.Value(SessionContextKey).(*sessions.ContextData)
	if !ok {
		return nil
	}

	return sessionContext
}

func (s *Server) AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		logger := s.logger.WithValue("grpc.method", info.FullMethod)

		switch info.FullMethod {
		// these methods don't require prior authentication
		case "/eating.EatingService/AdminLoginForToken",
			"/eating.EatingService/LoginForToken":
			logger.Info("skipping authentication for method")
			return handler(ctx, req)
		}

		_, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, Unauthenticated("missing metadata")
		}

		//authHeader := md.Get("authorization")
		//if len(authHeader) == 0 {
		//	return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		//}
		//
		//token := strings.TrimPrefix(authHeader[0], "Bearer ")

		//userID, err := s.tokenIssuer.ParseUserIDFromToken(ctx, token)
		//if err != nil {
		//	return nil, status.Error(codes.Unauthenticated, "invalid token")
		//}

		sessionContextData := &sessions.ContextData{
			HouseholdPermissions: nil,
			Requester: sessions.RequesterInfo{
				ServicePermissions:       authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
				AccountStatus:            "good",
				AccountStatusExplanation: "normal",
				UserID:                   "12345asdf",
				EmailAddress:             "fart@butts.com",
				Username:                 "example",
			},
			ActiveHouseholdID: "household123",
		}

		//sessionContextData, err := s.dataManager.BuildSessionContextDataForUser(ctx, userID)
		//if err != nil {
		//	return nil, status.Error(codes.Internal, "building session context data for user")
		//}

		ctx = context.WithValue(ctx, SessionContextKey, sessionContextData)

		return handler(ctx, req)
	}
}

func (s *Server) ExchangeToken(ctx context.Context, input *messages.ExchangeTokenRequest) (*messages.TokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	newToken, err := s.authManager.ExchangeTokenForUser(ctx, input.RefreshToken)
	if err != nil {
		return nil, Unauthenticated("invalid token")
	}

	output := &messages.TokenResponse{
		UserID:       newToken.UserID,
		HouseholdID:  newToken.HouseholdID,
		AccessToken:  newToken.AccessToken,
		RefreshToken: newToken.RefreshToken,
		ExpiresUTC:   timestamppb.New(newToken.ExpiresUTC),
	}

	return output, nil
}

func (s *Server) LoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return s.loginForToken(ctx, false, input)
}

func (s *Server) AdminLoginForToken(ctx context.Context, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	return s.loginForToken(ctx, true, input)
}

func (s *Server) loginForToken(ctx context.Context, admin bool, input *messages.UserLoginInput) (*messages.TokenResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	tokenResponse, err := s.authManager.ProcessLogin(ctx, admin, &types.UserLoginInput{
		Username:  input.Username,
		Password:  input.Password,
		TOTPToken: input.TOTPToken,
	})
	if err != nil {
		return nil, observability.PrepareError(err, span, "processing login")
	}

	output := &messages.TokenResponse{
		UserID:       tokenResponse.UserID,
		HouseholdID:  tokenResponse.HouseholdID,
		AccessToken:  tokenResponse.AccessToken,
		RefreshToken: tokenResponse.RefreshToken,
		ExpiresUTC:   timestamppb.New(tokenResponse.ExpiresUTC),
	}

	return output, nil
}

func (s *Server) CreateUser(ctx context.Context, input *messages.UserRegistrationInput) (*messages.UserCreationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	if err := validation.ValidateStructWithContext(
		ctx,
		input,
		validation.Field(&input.EmailAddress, validation.Required, is.EmailFormat),
		validation.Field(&input.Username, validation.Required, validation.Length(4, math.MaxInt8)),
		validation.Field(&input.Password, validation.Required, validation.Length(8, math.MaxInt8)),
	); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	input.Username = strings.TrimSpace(input.Username)
	tracing.AttachToSpan(span, keys.UsernameKey, input.Username)
	input.EmailAddress = strings.TrimSpace(strings.ToLower(input.EmailAddress))
	tracing.AttachToSpan(span, keys.UserEmailAddressKey, input.EmailAddress)
	input.Password = strings.TrimSpace(input.Password)

	logger := s.logger.WithValues(map[string]any{
		keys.UsernameKey:                 input.Username,
		keys.UserEmailAddressKey:         input.EmailAddress,
		keys.HouseholdInvitationIDKey:    input.InvitationID,
		keys.HouseholdInvitationTokenKey: input.InvitationToken,
	})

	// ensure the password is not garbage-tier
	if err := passwordvalidator.Validate(strings.TrimSpace(input.Password), minimumPasswordEntropy); err != nil {
		return nil, err
	}

	var invitation *types.HouseholdInvitation
	if input.InvitationID != "" && input.InvitationToken != "" {
		i, err := s.dataManager.GetHouseholdInvitationByTokenAndID(ctx, input.InvitationToken, input.InvitationID)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "getting household invitation")
		} else if err != nil {
			return nil, observability.PrepareError(err, span, "fetching household invitation")
		}

		invitation = i
		logger = logger.WithValue(keys.HouseholdInvitationIDKey, invitation.ID)
		logger.Debug("retrieved household invitation")
	}

	logger.Debug("completed invitation check")

	// hash the password
	hp, err := s.authenticator.HashPassword(ctx, strings.TrimSpace(input.Password))
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		return nil, err
	}

	// generate a two factor secret.
	tfs, err := s.secretGenerator.GenerateBase32EncodedString(ctx, totpSecretSize)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "generating two factor secret")
		return nil, err
	}

	creationInput := &types.UserDatabaseCreationInput{
		ID:              identifiers.New(),
		Username:        input.Username,
		FirstName:       input.FirstName,
		LastName:        input.LastName,
		EmailAddress:    input.EmailAddress,
		HashedPassword:  hp,
		TwoFactorSecret: tfs,
		InvitationToken: input.InvitationToken,
		Birthday:        converters.ConvertPBTimestampToTimePointer(input.Birthday),
		HouseholdName:   input.HouseholdName,
	}

	if invitation != nil {
		logger.Debug("supplementing user creation input with invitation data")
		creationInput.DestinationHouseholdID = invitation.DestinationHousehold.ID
		creationInput.InvitationToken = invitation.Token
	}

	// create the user.
	user, err := s.dataManager.CreateUser(ctx, creationInput)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "creating user")
		if errors.Is(err, database.ErrUserAlreadyExists) {
			return nil, err
		}
		return nil, err
	}

	logger.Debug("user created")

	defaultHouseholdID, err := s.dataManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching default household ID for user")
		return nil, err
	}

	emailVerificationToken, err := s.dataManager.GetEmailAddressVerificationTokenForUser(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// notify the relevant parties.
	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)

	dcm := &types.DataChangeMessage{
		HouseholdID:            defaultHouseholdID,
		EventType:              types.UserSignedUpServiceEventType,
		UserID:                 user.ID,
		EmailVerificationToken: emailVerificationToken,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	if err = s.analyticsReporter.AddUser(ctx, user.ID, map[string]any{
		"username":          user.Username,
		"default_household": defaultHouseholdID,
		"first_name":        user.FirstName,
		"last_name":         user.LastName,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user for analytics")
	}

	if err = s.featureFlagManager.Identify(ctx, user); err != nil {
		observability.AcknowledgeError(err, logger, span, "identifying user in feature flag manager")
	}

	userResponse := &messages.UserCreationResponse{
		CreatedAt:       converters.ConvertTimeToPBTimestamp(user.CreatedAt),
		Birthday:        converters.ConvertTimePointerToPBTimestamp(user.Birthday),
		TwoFactorQRCode: "TODO",
		Username:        user.Username,
		EmailAddress:    user.EmailAddress,
		CreatedUserID:   user.ID,
		AccountStatus:   user.AccountStatus,
		TwoFactorSecret: user.TwoFactorSecret,
		FirstName:       user.FirstName,
		LastName:        user.LastName,
	}

	return userResponse, nil
}

func (s *Server) VerifyTOTPSecret(ctx context.Context, input *messages.TOTPSecretVerificationInput) (*messages.TOTPSecretVerificationResponse, error) {
	ctx, span := s.tracer.StartSpan(ctx)
	defer span.End()

	logger := s.logger.Clone()

	// TODO: validate

	user, err := s.dataManager.GetUserWithUnverifiedTwoFactorSecret(ctx, input.UserID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "fetching user to verify two factor secret")
		return nil, err
	}

	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, keys.UsernameKey, user.Username)
	logger = logger.WithValue(keys.UsernameKey, user.Username)

	if user.TwoFactorSecretVerifiedAt != nil {
		// I suppose if this happens too many times, we might want to keep track of that
		logger.Debug("two factor secret already verified")
		return nil, err
	}

	totpValid := totp.Validate(input.TOTPToken, user.TwoFactorSecret)
	if !totpValid {
		logger = logger.WithValue(keys.ValidationErrorKey, err)
		observability.AcknowledgeError(err, logger, span, "TOTP code was invalid")
		return nil, err
	}

	if err = s.dataManager.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		observability.AcknowledgeError(err, logger, span, "verifying user two factor secret")
		return nil, err
	}

	dcm := &types.DataChangeMessage{
		EventType: types.TwoFactorSecretVerifiedServiceEventType,
		UserID:    user.ID,
	}

	if err = s.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		observability.AcknowledgeError(err, logger, span, "publishing data change message")
	}

	output := &messages.TOTPSecretVerificationResponse{
		Accepted: true,
	}

	return output, nil
}
