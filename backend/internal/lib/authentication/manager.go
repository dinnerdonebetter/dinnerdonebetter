package authentication

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	authcfg "github.com/dinnerdonebetter/backend/internal/lib/authentication/config"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication/tokens"
	"github.com/dinnerdonebetter/backend/internal/lib/featureflags"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

type (
	Manager interface {
		ProcessLogin(ctx context.Context, adminOnly bool, loginData *types.UserLoginInput) (*tokens.TokenResponse, error)
	}

	UserAuthDataManager interface {
		GetUserByUsername(ctx context.Context, username string) (*types.User, error)
		GetAdminUserByUsername(ctx context.Context, username string) (*types.User, error)
		GetDefaultHouseholdIDForUser(ctx context.Context, userID string) (string, error)
	}

	manager struct {
		config               *authcfg.Config
		tokenIssuer          tokens.Issuer
		authenticator        Authenticator
		tracer               tracing.Tracer
		logger               logging.Logger
		featureFlagManager   featureflags.FeatureFlagManager
		dataChangesPublisher messagequeue.Publisher
		analyticsReporter    analytics.EventReporter
		userAuthDataManager  UserAuthDataManager
	}
)

func NewManager(
	cfg *authcfg.Config,
	queuesConfig msgconfig.QueuesConfig,
	tokenIssuer tokens.Issuer,
	authenticator Authenticator,
	tracingProvider tracing.TracerProvider,
	logger logging.Logger,
	featureFlagManager featureflags.FeatureFlagManager,
	publisherProvider messagequeue.PublisherProvider,
	analyticsReporter analytics.EventReporter,
	userAuthDataManager UserAuthDataManager,
) (Manager, error) {
	const name = "authentication_manager"

	if cfg == nil {
		return nil, internalerrors.NilConfigError(name)
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queuesConfig.DataChangesTopicName)
	if err != nil {
		return nil, observability.PrepareError(err, nil, "creating data changes publisher")
	}

	m := &manager{
		config:               cfg,
		tracer:               tracing.NewTracer(tracing.EnsureTracerProvider(tracingProvider).Tracer(name)),
		logger:               logging.EnsureLogger(logger).WithName(name),
		tokenIssuer:          tokenIssuer,
		authenticator:        authenticator,
		featureFlagManager:   featureFlagManager,
		dataChangesPublisher: dataChangesPublisher,
		analyticsReporter:    analyticsReporter,
		userAuthDataManager:  userAuthDataManager,
	}

	return m, nil
}

// validateLogin takes login information and returns whether the login is valid.
// In the event that there's an error, this function will return false and the error.
func (m *manager) validateLogin(ctx context.Context, user *types.User, loginInput *types.UserLoginInput) (bool, error) {
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

func (m *manager) ProcessLogin(ctx context.Context, adminOnly bool, loginData *types.UserLoginInput) (*tokens.TokenResponse, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := m.logger.Clone()

	if err := loginData.ValidateWithContext(ctx); err != nil {
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

	defaultHouseholdID, err := m.userAuthDataManager.GetDefaultHouseholdIDForUser(ctx, user.ID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "validating input")
	}

	dcm := &types.DataChangeMessage{
		EventType:   types.UserLoggedInServiceEventType,
		HouseholdID: defaultHouseholdID,
		UserID:      user.ID,
	}

	if err = m.dataChangesPublisher.Publish(ctx, dcm); err != nil {
		return nil, observability.PrepareError(err, span, "publishing data change")
	}

	response := &tokens.TokenResponse{
		UserID:      user.ID,
		HouseholdID: defaultHouseholdID,
		ExpiresUTC:  time.Now().Add(m.config.MaxAccessTokenLifetime).UTC(),
	}

	response.AccessToken, err = m.tokenIssuer.IssueToken(ctx, user, m.config.MaxAccessTokenLifetime)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating accessToken")
	}

	response.RefreshToken, err = m.tokenIssuer.IssueToken(ctx, user, m.config.MaxRefreshTokenLifetime)
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating accessToken")
	}

	return response, nil
}
