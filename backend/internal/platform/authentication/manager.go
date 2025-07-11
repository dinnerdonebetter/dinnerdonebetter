package authentication

import (
	"context"
	"database/sql"
	"errors"
	"math"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	name = "authentication_manager"
)

type (
	// UserLoginInput represents the payload used to log in a User.
	UserLoginInput struct {
		_ struct{} `json:"-"`

		Username  string `json:"username"`
		Password  string `json:"password"`
		TOTPToken string `json:"totpToken"`
	}

	Manager interface {
		ProcessLogin(ctx context.Context, adminOnly bool, loginData *UserLoginInput) (*identity.TokenResponse, error)
		ExchangeTokenForUser(ctx context.Context, refreshToken string) (*identity.TokenResponse, error)
	}

	UserAuthDataManager interface {
		GetUser(ctx context.Context, userID string) (*identity.User, error)
		GetUserByUsername(ctx context.Context, username string) (*identity.User, error)
		GetAdminUserByUsername(ctx context.Context, username string) (*identity.User, error)
		GetDefaultAccountIDForUser(ctx context.Context, userID string) (string, error)
	}

	manager struct {
		tokenIssuer             tokens.Issuer
		authenticator           Authenticator
		tracer                  tracing.Tracer
		logger                  logging.Logger
		dataChangesPublisher    messagequeue.Publisher
		userAuthDataManager     UserAuthDataManager
		maxAccessTokenLifetime  time.Duration
		maxRefreshTokenLifetime time.Duration
	}
)

func NewManager(
	queuesConfig *msgconfig.QueuesConfig,
	tokenIssuer tokens.Issuer,
	authenticator Authenticator,
	tracingProvider tracing.TracerProvider,
	logger logging.Logger,
	publisherProvider messagequeue.PublisherProvider,
	userAuthDataManager UserAuthDataManager,
	maxAccessTokenLifetime time.Duration,
	maxRefreshTokenLifetime time.Duration,
) (Manager, error) {

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queuesConfig.DataChangesTopicName)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "creating data changes publisher")
	}

	m := &manager{
		maxRefreshTokenLifetime: maxRefreshTokenLifetime,
		maxAccessTokenLifetime:  maxAccessTokenLifetime,
		tracer:                  tracing.NewTracer(tracing.EnsureTracerProvider(tracingProvider).Tracer(name)),
		logger:                  logging.EnsureLogger(logger).WithName(name),
		tokenIssuer:             tokenIssuer,
		authenticator:           authenticator,
		dataChangesPublisher:    dataChangesPublisher,
		userAuthDataManager:     userAuthDataManager,
	}

	return m, nil
}

// validateLogin takes login information and returns whether the login is valid.
// In the event that there's an error, this function will return false and the error.
func (m *manager) validateLogin(ctx context.Context, user *identity.User, loginInput *UserLoginInput) (bool, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	loginInput.TOTPToken = strings.TrimSpace(loginInput.TOTPToken)
	loginInput.Password = strings.TrimSpace(loginInput.Password)
	loginInput.Username = strings.TrimSpace(loginInput.Username)

	// alias the relevant data.
	logger := m.logger.WithValue(keys.UsernameKey, user.Username)

	// check for login validity.
	loginValid, err := m.authenticator.CredentialsAreValid(
		ctx,
		user.HashedPassword,
		loginInput.Password,
		user.TwoFactorSecret,
		loginInput.TOTPToken,
	)

	if errors.Is(err, ErrInvalidTOTPToken) || errors.Is(err, ErrPasswordDoesNotMatch) {
		return false, err
	}

	if err != nil {
		return false, observability.PrepareError(err, span, "validating login")
	}

	logger.Debug("login validated")

	return loginValid, nil
}

func (m *manager) ProcessLogin(ctx context.Context, adminOnly bool, loginData *UserLoginInput) (*identity.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.Clone()

	if err := validation.ValidateStructWithContext(ctx, loginData,
		validation.Field(&loginData.Username, validation.Required, validation.Length(4, math.MaxInt8)),
		validation.Field(&loginData.Password, validation.Required, validation.Length(8, math.MaxInt8)),
		validation.Field(&loginData.TOTPToken, is.Digit, validation.RuneLength(6, 6)),
	); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	logger = logger.WithValue(keys.UsernameKey, loginData.Username)

	userFunc := m.userAuthDataManager.GetUserByUsername
	if adminOnly {
		userFunc = m.userAuthDataManager.GetAdminUserByUsername
	}

	user, err := userFunc(ctx, loginData.Username)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "user does not exist")
		}

		return nil, observability.PrepareError(err, span, "fetching user")
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)

	if user.IsBanned() {
		return nil, observability.PrepareError(err, span, "user is banned")
	}

	loginValid, err := m.validateLogin(ctx, user, loginData)
	logger.WithValue("login_valid", loginValid)

	if err != nil {
		if errors.Is(err, ErrInvalidTOTPToken) {
			return nil, observability.PrepareError(err, span, "invalid TOTP AccessToken")
		}

		if errors.Is(err, ErrPasswordDoesNotMatch) {
			return nil, observability.PrepareError(err, span, "password did not match")
		}

		return nil, observability.PrepareError(err, span, "validating login")
	} else if !loginValid {
		return nil, observability.PrepareError(err, span, "login was invalid")
	}

	if user.TwoFactorSecretVerifiedAt != nil && loginData.TOTPToken == "" {
		return nil, observability.PrepareError(err, span, "TOTP code required but not provided")
	}

	defaultAccountID, err := m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: defaultAccountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	response := &identity.TokenResponse{
		UserID:     user.ID,
		AccountID:  defaultAccountID,
		ExpiresUTC: time.Now().Add(m.maxAccessTokenLifetime).UTC(),
	}

	response.AccessToken, err = m.tokenIssuer.IssueToken(ctx, user, m.maxAccessTokenLifetime)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating accessToken")
	}

	response.RefreshToken, err = m.tokenIssuer.IssueToken(ctx, user, m.maxRefreshTokenLifetime)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating accessToken")
	}

	return response, nil
}

func (m *manager) ExchangeTokenForUser(ctx context.Context, refreshToken string) (*identity.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.Clone()

	userID, err := m.tokenIssuer.ParseUserIDFromToken(ctx, refreshToken)
	if err != nil {
		return nil, observability.PrepareError(err, span, "parsing userID from token")
	}

	user, err := m.userAuthDataManager.GetUser(ctx, userID)
	if err != nil || user == nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, observability.PrepareError(err, span, "user does not exist")
		}

		return nil, observability.PrepareError(err, span, "fetching user")
	}

	logger = logger.WithValue(keys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, keys.UserIDKey, user.ID)

	if user.IsBanned() {
		return nil, observability.PrepareError(err, span, "user is banned")
	}

	defaultAccountID, err := m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: defaultAccountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	response := &identity.TokenResponse{
		UserID:     user.ID,
		AccountID:  defaultAccountID,
		ExpiresUTC: time.Now().Add(m.maxAccessTokenLifetime).UTC(),
	}

	response.AccessToken, err = m.tokenIssuer.IssueToken(ctx, user, m.maxAccessTokenLifetime)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating accessToken")
	}

	response.RefreshToken, err = m.tokenIssuer.IssueToken(ctx, user, m.maxRefreshTokenLifetime)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating accessToken")
	}

	return response, nil
}
