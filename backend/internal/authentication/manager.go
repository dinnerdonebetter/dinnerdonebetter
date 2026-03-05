package authentication

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication/tokens"
	tokenscfg "github.com/dinnerdonebetter/backend/internal/authentication/tokens/config"
	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitykeys "github.com/dinnerdonebetter/backend/internal/domain/identity/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	name = "authentication_manager"
)

type (
	Manager interface {
		ProcessLogin(ctx context.Context, adminOnly bool, loginData *auth.UserLoginInput) (*auth.TokenResponse, error)
		ExchangeTokenForUser(ctx context.Context, refreshToken, desiredAccountID string) (*auth.TokenResponse, error)
	}

	manager struct {
		tokenIssuer             tokens.Issuer
		authenticator           Authenticator
		tracer                  tracing.Tracer
		logger                  logging.Logger
		dataChangesPublisher    messagequeue.Publisher
		userAuthDataManager     identity.Repository
		maxAccessTokenLifetime  time.Duration
		maxRefreshTokenLifetime time.Duration
	}
)

func NewManager(
	ctx context.Context,
	queuesConfig *msgconfig.QueuesConfig,
	tokenIssuer tokens.Issuer,
	authenticator Authenticator,
	tracingProvider tracing.TracerProvider,
	logger logging.Logger,
	publisherProvider messagequeue.PublisherProvider,
	userAuthDataManager identity.Repository,
	cfg *tokenscfg.Config,
) (Manager, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(ctx, queuesConfig.DataChangesTopicName)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "creating data changes publisher")
	}

	m := &manager{
		maxRefreshTokenLifetime: cfg.MaxRefreshTokenLifetime,
		maxAccessTokenLifetime:  cfg.MaxAccessTokenLifetime,
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
func (m *manager) validateLogin(ctx context.Context, user *identity.User, loginInput *auth.UserLoginInput) (bool, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	loginInput.TOTPToken = strings.TrimSpace(loginInput.TOTPToken)
	loginInput.Password = strings.TrimSpace(loginInput.Password)
	loginInput.Username = strings.TrimSpace(loginInput.Username)

	// alias the relevant data.
	logger := m.logger.WithValue(identitykeys.UsernameKey, user.Username)

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

func (m *manager) ProcessLogin(ctx context.Context, adminOnly bool, loginData *auth.UserLoginInput) (*auth.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.Clone()

	if err := loginData.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	logger = logger.WithValue(identitykeys.UsernameKey, loginData.Username)

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

	logger = logger.WithValue(identitykeys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, user.ID)

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
		return nil, observability.PrepareError(errors.New("TOTP code required but not provided"), span, "processing login")
	}

	var accountID string
	if loginData.DesiredAccountID != "" {
		var isMember bool
		isMember, err = m.userAuthDataManager.UserIsMemberOfAccount(ctx, user.ID, loginData.DesiredAccountID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating account membership")
		}
		if !isMember {
			return nil, observability.PrepareError(errors.New("user does not have access to account"), span, "user does not have access to the desired account")
		}
		accountID = loginData.DesiredAccountID
	} else {
		var defaultAccountID string
		defaultAccountID, err = m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating input")
		}
		accountID = defaultAccountID
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: accountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	response := &auth.TokenResponse{
		UserID:     user.ID,
		AccountID:  accountID,
		ExpiresUTC: time.Now().Add(m.maxAccessTokenLifetime).UTC(),
	}

	response.AccessToken, err = m.tokenIssuer.IssueTokenWithAccount(ctx, user, m.maxAccessTokenLifetime, accountID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating access token")
	}

	response.RefreshToken, err = m.tokenIssuer.IssueTokenWithAccount(ctx, user, m.maxRefreshTokenLifetime, accountID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating refresh token")
	}

	return response, nil
}

func (m *manager) ExchangeTokenForUser(ctx context.Context, refreshToken, desiredAccountID string) (*auth.TokenResponse, error) {
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

	logger = logger.WithValue(identitykeys.UserIDKey, user.ID)
	tracing.AttachToSpan(span, identitykeys.UserIDKey, user.ID)

	if user.IsBanned() {
		return nil, observability.PrepareError(err, span, "user is banned")
	}

	var accountID string
	if desiredAccountID != "" {
		var isMember bool
		isMember, err = m.userAuthDataManager.UserIsMemberOfAccount(ctx, user.ID, desiredAccountID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating account membership")
		}
		if !isMember {
			return nil, observability.PrepareError(errors.New("user does not have access to account"), span, "user does not have access to the desired account")
		}
		accountID = desiredAccountID
	} else {
		var defaultAccountID string
		defaultAccountID, err = m.userAuthDataManager.GetDefaultAccountIDForUser(ctx, user.ID)
		if err != nil {
			return nil, observability.PrepareError(err, span, "validating input")
		}
		accountID = defaultAccountID
	}

	dcm := &audit.DataChangeMessage{
		EventType: identity.UserLoggedInServiceEventType,
		AccountID: accountID,
		UserID:    user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	response := &auth.TokenResponse{
		UserID:     user.ID,
		AccountID:  accountID,
		ExpiresUTC: time.Now().Add(m.maxAccessTokenLifetime).UTC(),
	}

	response.AccessToken, err = m.tokenIssuer.IssueTokenWithAccount(ctx, user, m.maxAccessTokenLifetime, accountID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating accessToken")
	}

	response.RefreshToken, err = m.tokenIssuer.IssueTokenWithAccount(ctx, user, m.maxRefreshTokenLifetime, accountID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating accessToken")
	}

	return response, nil
}
