package managers

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	mockauthn "github.com/dinnerdonebetter/backend/internal/authentication/mock"
	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	authfakes "github.com/dinnerdonebetter/backend/internal/domain/auth/fakes"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityfakes "github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	identitymock "github.com/dinnerdonebetter/backend/internal/domain/identity/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/qrcodes"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	randommock "github.com/dinnerdonebetter/backend/internal/platform/random/mock"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/pquerna/otp/totp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockPasswordResetTokenDataManager is a test double for auth.PasswordResetTokenDataManager.
type mockPasswordResetTokenDataManager struct {
	mock.Mock
}

func (m *mockPasswordResetTokenDataManager) GetPasswordResetTokenByToken(ctx context.Context, token string) (*auth.PasswordResetToken, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.PasswordResetToken), args.Error(1)
}

func (m *mockPasswordResetTokenDataManager) CreatePasswordResetToken(ctx context.Context, input *auth.PasswordResetTokenDatabaseCreationInput) (*auth.PasswordResetToken, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.PasswordResetToken), args.Error(1)
}

func (m *mockPasswordResetTokenDataManager) RedeemPasswordResetToken(ctx context.Context, passwordResetTokenID string) error {
	return m.Called(ctx, passwordResetTokenID).Error(0)
}

func TestProvideAuthManager(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		queueCfg := &msgconfig.QueuesConfig{DataChangesTopicName: t.Name()}

		mpp := &mockpublishers.PublisherProvider{}
		mpp.On(reflection.GetMethodName(mpp.ProvidePublisher), queueCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		m, err := ProvideAuthManager(
			ctx,
			logging.NewNoopLogger(),
			tracing.NewNoopTracerProvider(),
			&mockPasswordResetTokenDataManager{},
			&identitymock.RepositoryMock{},
			&mockauthn.Authenticator{},
			mpp,
			random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
			qrcodes.NewBuilder(tracing.NewNoopTracerProvider(), logging.NewNoopLogger()),
			queueCfg,
		)

		require.NoError(t, err)
		assert.NotNil(t, m)
		mock.AssertExpectationsForObjects(t, mpp)
	})
}

func TestAuthManager_Self(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := identityfakes.BuildFakeID()
		expectedUser := identityfakes.BuildFakeUser()
		expectedUser.ID = userID

		userDataManager := &identitymock.RepositoryMock{}
		userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, userID).Return(expectedUser, nil)

		sessionData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{UserID: userID},
		}
		sessionFetcher := func(context.Context) (*sessions.ContextData, error) {
			return sessionData, nil
		}

		manager := &AuthManager{
			userDataManager:           userDataManager,
			sessionContextDataFetcher: sessionFetcher,
			logger:                    logging.NewNoopLogger().WithName("auth_manager"),
			tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
		}

		result, err := manager.Self(ctx)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, userID, result.ID)
		assert.Equal(t, expectedUser.Username, result.Username)
		mock.AssertExpectationsForObjects(t, userDataManager)
	})
}

func TestAuthManager_CheckUserPermissions(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		userID := identityfakes.BuildFakeID()
		accountID := identityfakes.BuildFakeID()

		sessionData := &sessions.ContextData{
			Requester: sessions.RequesterInfo{
				UserID:             userID,
				ServicePermissions: authorization.NewServiceRolePermissionChecker(authorization.ServiceUserRole.String()),
			},
			ActiveAccountID: accountID,
			AccountPermissions: map[string]authorization.AccountRolePermissionsChecker{
				accountID: authorization.NewAccountRolePermissionChecker(authorization.AccountMemberRoleName),
			},
		}
		sessionFetcher := func(context.Context) (*sessions.ContextData, error) {
			return sessionData, nil
		}

		manager := &AuthManager{
			sessionContextDataFetcher: sessionFetcher,
			logger:                    logging.NewNoopLogger().WithName("auth_manager"),
			tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
		}

		input := &auth.UserPermissionsRequestInput{
			Permissions: []string{"meal_planning:read"},
		}

		result, err := manager.CheckUserPermissions(ctx, input)

		require.NoError(t, err)
		assert.NotNil(t, result)
		assert.NotNil(t, result.Permissions)
	})

	t.Run("session fetch error", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		sessionFetcher := func(context.Context) (*sessions.ContextData, error) {
			return nil, errors.New("session error")
		}

		manager := &AuthManager{
			sessionContextDataFetcher: sessionFetcher,
			logger:                    logging.NewNoopLogger().WithName("auth_manager"),
			tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
		}

		result, err := manager.CheckUserPermissions(ctx, &auth.UserPermissionsRequestInput{Permissions: []string{"read"}})

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProvideAuthManager_NilConfig(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	mpp := &mockpublishers.PublisherProvider{}

	m, err := ProvideAuthManager(
		ctx,
		logging.NewNoopLogger(),
		tracing.NewNoopTracerProvider(),
		&mockPasswordResetTokenDataManager{},
		&identitymock.RepositoryMock{},
		&mockauthn.Authenticator{},
		mpp,
		random.NewGenerator(logging.NewNoopLogger(), tracing.NewNoopTracerProvider()),
		qrcodes.NewBuilder(tracing.NewNoopTracerProvider(), logging.NewNoopLogger()),
		nil, // nil config
	)

	assert.Error(t, err)
	assert.Nil(t, m)
}

func TestAuthManager_Self_SessionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	sessionFetcher := func(context.Context) (*sessions.ContextData, error) {
		return nil, errors.New("session error")
	}

	manager := &AuthManager{
		sessionContextDataFetcher: sessionFetcher,
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	result, err := manager.Self(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestAuthManager_Self_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := identityfakes.BuildFakeID()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, userID).Return((*identity.User)(nil), sql.ErrNoRows)

	sessionFetcher := func(context.Context) (*sessions.ContextData, error) {
		return &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: userID}}, nil
	}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		sessionContextDataFetcher: sessionFetcher,
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	result, err := manager.Self(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
	mock.AssertExpectationsForObjects(t, userDataManager)
}

func TestAuthManager_TOTPSecretVerification_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	key, err := totp.Generate(totp.GenerateOpts{Issuer: "test", AccountName: "user"})
	require.NoError(t, err)

	user := identityfakes.BuildFakeUser()
	user.TwoFactorSecret = key.Secret()
	user.TwoFactorSecretVerifiedAt = nil

	token, err := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	require.NoError(t, err)

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserWithUnverifiedTwoFactorSecret), testutils.ContextMatcher, user.ID).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.MarkUserTwoFactorSecretAsVerified), testutils.ContextMatcher, user.ID).Return(nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	manager := &AuthManager{
		userDataManager:           userDataManager,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	input := &auth.TOTPSecretVerificationInput{UserID: user.ID, TOTPToken: token}
	err = manager.TOTPSecretVerification(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, publisher)
}

func TestAuthManager_TOTPSecretVerification_InvalidInput(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	manager := &AuthManager{
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.TOTPSecretVerification(ctx, &auth.TOTPSecretVerificationInput{UserID: "", TOTPToken: "123"})

	assert.Error(t, err)
}

func TestAuthManager_TOTPSecretVerification_AlreadyVerified(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	verifiedAt := time.Now()
	user := identityfakes.BuildFakeUser()
	user.TwoFactorSecretVerifiedAt = &verifiedAt

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserWithUnverifiedTwoFactorSecret), testutils.ContextMatcher, user.ID).Return(user, nil)

	manager := &AuthManager{
		userDataManager:           userDataManager,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	input := &auth.TOTPSecretVerificationInput{UserID: user.ID, TOTPToken: "123456"}
	err := manager.TOTPSecretVerification(ctx, input)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "already verified")
	mock.AssertExpectationsForObjects(t, userDataManager)
}

func TestAuthManager_RequestUsernameReminder_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	input := authfakes.BuildFakeUsernameReminderRequestInput()
	input.EmailAddress = user.EmailAddress

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmail), testutils.ContextMatcher, input.EmailAddress).Return(user, nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	manager := &AuthManager{
		userDataManager:           userDataManager,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.RequestUsernameReminder(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, publisher)
}

func TestAuthManager_RequestUsernameReminder_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := authfakes.BuildFakeUsernameReminderRequestInput()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmail), testutils.ContextMatcher, input.EmailAddress).Return((*identity.User)(nil), sql.ErrNoRows)

	manager := &AuthManager{
		userDataManager:           userDataManager,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.RequestUsernameReminder(ctx, input)

	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager)
}

func TestAuthManager_CreatePasswordResetToken_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	input := authfakes.BuildFakePasswordResetTokenCreationRequestInput()
	input.EmailAddress = user.EmailAddress

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmail), testutils.ContextMatcher, input.EmailAddress).Return(user, nil)

	secretGen := &randommock.Generator{}
	secretGen.On(reflection.GetMethodName(secretGen.GenerateBase32EncodedString), testutils.ContextMatcher, 32).Return("faketoken123", nil)

	prtManager := &mockPasswordResetTokenDataManager{}
	createdToken := authfakes.BuildFakePasswordResetToken()
	createdToken.BelongsToUser = user.ID
	prtManager.On(reflection.GetMethodName(prtManager.CreatePasswordResetToken), testutils.ContextMatcher, mock.Anything).Return(createdToken, nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	manager := &AuthManager{
		userDataManager:               userDataManager,
		passwordResetTokenDataManager: prtManager,
		secretGenerator:               secretGen,
		dataChangesPublisher:          publisher,
		sessionContextDataFetcher:     func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                        logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                        tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.CreatePasswordResetToken(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, secretGen, prtManager, publisher)
}

func TestAuthManager_CreatePasswordResetToken_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := authfakes.BuildFakePasswordResetTokenCreationRequestInput()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmail), testutils.ContextMatcher, input.EmailAddress).Return((*identity.User)(nil), sql.ErrNoRows)

	manager := &AuthManager{
		userDataManager:           userDataManager,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.CreatePasswordResetToken(ctx, input)

	// Returns success without sending email to avoid email enumeration.
	assert.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager)
}

func TestAuthManager_RequestEmailVerificationEmail_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := identityfakes.BuildFakeID()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetEmailAddressVerificationTokenForUser), testutils.ContextMatcher, userID).Return("verification-token-123", nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: userID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.RequestEmailVerificationEmail(ctx)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, publisher)
}

func TestAuthManager_VerifyUserEmailAddress_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	input := authfakes.BuildFakeEmailAddressVerificationRequestInput()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmailAddressVerificationToken), testutils.ContextMatcher, input.Token).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.MarkUserEmailAddressAsVerified), testutils.ContextMatcher, user.ID, input.Token).Return(nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	manager := &AuthManager{
		userDataManager:           userDataManager,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.VerifyUserEmailAddress(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, publisher)
}

func TestAuthManager_VerifyUserEmailAddressByToken_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	token := "verification-token"

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmailAddressVerificationToken), testutils.ContextMatcher, token).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.MarkUserEmailAddressAsVerified), testutils.ContextMatcher, user.ID, token).Return(nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	manager := &AuthManager{
		userDataManager:       userDataManager,
		dataChangesPublisher:  publisher,
		logger:                logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.VerifyUserEmailAddressByToken(ctx, token)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, publisher)
}

func TestAuthManager_VerifyUserEmailAddressByToken_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	token := "invalid-token"

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmailAddressVerificationToken), testutils.ContextMatcher, token).Return((*identity.User)(nil), sql.ErrNoRows)

	manager := &AuthManager{
		userDataManager:      userDataManager,
		logger:               logging.NewNoopLogger().WithName("auth_manager"),
		tracer:               tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.VerifyUserEmailAddressByToken(ctx, token)

	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager)
}

func TestAuthManager_UpdatePassword_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	user.TwoFactorSecretVerifiedAt = nil
	password := authfakes.BuildFakePasswordUpdateInput()
	password.CurrentPassword = "current"
	password.NewPassword = "Abcdefghij123!@#$%^&*()"
	password.TOTPToken = ""

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.UpdateUserPassword), testutils.ContextMatcher, user.ID, mock.AnythingOfType("string")).Return(nil)

	authenticator := &mockauthn.Authenticator{}
	authenticator.On(reflection.GetMethodName(authenticator.CredentialsAreValid), testutils.ContextMatcher, user.HashedPassword, "current", "", "").Return(true, nil)
	authenticator.On(reflection.GetMethodName(authenticator.HashPassword), testutils.ContextMatcher, "Abcdefghij123!@#$%^&*()").Return("hashed", nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: user.ID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.UpdatePassword(ctx, password)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, authenticator, publisher)
}

func TestAuthManager_UpdateUserEmailAddress_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	user.TwoFactorSecretVerifiedAt = nil
	input := authfakes.BuildFakeUserEmailAddressUpdateInput()
	input.CurrentPassword = "current"

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.UpdateUserEmailAddress), testutils.ContextMatcher, user.ID, input.NewEmailAddress).Return(nil)

	authenticator := &mockauthn.Authenticator{}
	authenticator.On(reflection.GetMethodName(authenticator.CredentialsAreValid), testutils.ContextMatcher, user.HashedPassword, "current", "", "").Return(true, nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: user.ID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.UpdateUserEmailAddress(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, authenticator, publisher)
}

func TestAuthManager_UpdateUserUsername_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	user.TwoFactorSecretVerifiedAt = nil
	input := authfakes.BuildFakeUsernameUpdateInput()
	input.CurrentPassword = "current"

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.UpdateUserUsername), testutils.ContextMatcher, user.ID, input.NewUsername).Return(nil)

	authenticator := &mockauthn.Authenticator{}
	authenticator.On(reflection.GetMethodName(authenticator.CredentialsAreValid), testutils.ContextMatcher, user.HashedPassword, "current", "", "").Return(true, nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: user.ID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.UpdateUserUsername(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, authenticator, publisher)
}

func TestAuthManager_PasswordResetTokenRedemption_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	token := authfakes.BuildFakePasswordResetToken()
	token.BelongsToUser = user.ID
	token.Token = "reset-token-123"
	input := authfakes.BuildFakePasswordResetTokenRedemptionRequestInput()
	input.Token = token.Token
	input.NewPassword = "Abcdefghij123!@#$%^&*()"

	prtManager := &mockPasswordResetTokenDataManager{}
	prtManager.On(reflection.GetMethodName(prtManager.GetPasswordResetTokenByToken), testutils.ContextMatcher, token.Token).Return(token, nil)
	prtManager.On(reflection.GetMethodName(prtManager.RedeemPasswordResetToken), testutils.ContextMatcher, token.ID).Return(nil)

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.UpdateUserPassword), testutils.ContextMatcher, user.ID, mock.AnythingOfType("string")).Return(nil)

	authenticator := &mockauthn.Authenticator{}
	authenticator.On(reflection.GetMethodName(authenticator.HashPassword), testutils.ContextMatcher, "Abcdefghij123!@#$%^&*()").Return("hashed", nil)

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	manager := &AuthManager{
		passwordResetTokenDataManager: prtManager,
		userDataManager:               userDataManager,
		authenticator:                 authenticator,
		dataChangesPublisher:          publisher,
		sessionContextDataFetcher:     func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                        logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                        tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.PasswordResetTokenRedemption(ctx, input)

	require.NoError(t, err)
	mock.AssertExpectationsForObjects(t, prtManager, userDataManager, authenticator, publisher)
}

func TestAuthManager_NewTOTPSecret_Success(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	verifiedAt := time.Now()
	user.TwoFactorSecretVerifiedAt = &verifiedAt
	input := authfakes.BuildFakeTOTPSecretRefreshInput()
	input.CurrentPassword = "current"
	token, _ := totp.GenerateCode(user.TwoFactorSecret, time.Now().UTC())
	input.TOTPToken = token

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)
	userDataManager.On(reflection.GetMethodName(userDataManager.MarkUserTwoFactorSecretAsUnverified), testutils.ContextMatcher, user.ID, mock.AnythingOfType("string")).Return(nil)

	authenticator := &mockauthn.Authenticator{}
	authenticator.On(reflection.GetMethodName(authenticator.CredentialsAreValid), testutils.ContextMatcher, user.HashedPassword, "current", user.TwoFactorSecret, token).Return(true, nil)

	secretGen := &randommock.Generator{}
	secretGen.On(reflection.GetMethodName(secretGen.GenerateBase32EncodedString), testutils.ContextMatcher, 64).Return("newsecretencoded", nil)

	qrBuilder := qrcodes.NewBuilder(tracing.NewNoopTracerProvider(), logging.NewNoopLogger())

	publisher := &mockpublishers.Publisher{}
	publisher.On(reflection.GetMethodName(publisher.PublishAsync), testutils.ContextMatcher, mock.Anything).Return()

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: user.ID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		secretGenerator:           secretGen,
		qrCodeBuilder:             qrBuilder,
		dataChangesPublisher:      publisher,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	result, err := manager.NewTOTPSecret(ctx, input)

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "newsecretencoded", result.TwoFactorSecret)
	assert.NotEmpty(t, result.TwoFactorQRCode)
	mock.AssertExpectationsForObjects(t, userDataManager, authenticator, secretGen, publisher)
}

func TestAuthManager_PasswordResetTokenRedemption_TokenNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := authfakes.BuildFakePasswordResetTokenRedemptionRequestInput()

	prtManager := &mockPasswordResetTokenDataManager{}
	prtManager.On(reflection.GetMethodName(prtManager.GetPasswordResetTokenByToken), testutils.ContextMatcher, input.Token).Return((*auth.PasswordResetToken)(nil), sql.ErrNoRows)

	manager := &AuthManager{
		passwordResetTokenDataManager: prtManager,
		sessionContextDataFetcher:     func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                        logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                        tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.PasswordResetTokenRedemption(ctx, input)

	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, prtManager)
}

func TestAuthManager_PasswordResetTokenRedemption_InvalidPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	token := authfakes.BuildFakePasswordResetToken()
	token.BelongsToUser = user.ID
	input := authfakes.BuildFakePasswordResetTokenRedemptionRequestInput()
	input.Token = token.Token
	input.NewPassword = "a" // too weak for entropy 60

	prtManager := &mockPasswordResetTokenDataManager{}
	prtManager.On(reflection.GetMethodName(prtManager.GetPasswordResetTokenByToken), testutils.ContextMatcher, token.Token).Return(token, nil)

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)

	manager := &AuthManager{
		passwordResetTokenDataManager: prtManager,
		userDataManager:               userDataManager,
		sessionContextDataFetcher:     func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                        logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                        tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.PasswordResetTokenRedemption(ctx, input)

	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, prtManager, userDataManager)
}

func TestAuthManager_VerifyUserEmailAddress_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	input := authfakes.BuildFakeEmailAddressVerificationRequestInput()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUserByEmailAddressVerificationToken), testutils.ContextMatcher, input.Token).Return((*identity.User)(nil), sql.ErrNoRows)

	manager := &AuthManager{
		userDataManager:           userDataManager,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return &sessions.ContextData{}, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.VerifyUserEmailAddress(ctx, input)

	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager)
}

func TestAuthManager_UpdatePassword_InvalidNewPassword(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	user := identityfakes.BuildFakeUser()
	user.TwoFactorSecretVerifiedAt = nil
	password := authfakes.BuildFakePasswordUpdateInput()
	password.CurrentPassword = "current"
	password.NewPassword = "a" // too weak for entropy 60
	password.TOTPToken = ""

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, user.ID).Return(user, nil)

	authenticator := &mockauthn.Authenticator{}
	authenticator.On(reflection.GetMethodName(authenticator.CredentialsAreValid), testutils.ContextMatcher, user.HashedPassword, "current", "", "").Return(true, nil)

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: user.ID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		authenticator:             authenticator,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	err := manager.UpdatePassword(ctx, password)

	assert.Error(t, err)
	mock.AssertExpectationsForObjects(t, userDataManager, authenticator)
}

func TestAuthManager_NewTOTPSecret_UserNotFound(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	userID := identityfakes.BuildFakeID()
	input := authfakes.BuildFakeTOTPSecretRefreshInput()

	userDataManager := &identitymock.RepositoryMock{}
	userDataManager.On(reflection.GetMethodName(userDataManager.GetUser), testutils.ContextMatcher, userID).Return((*identity.User)(nil), sql.ErrNoRows)

	sessionData := &sessions.ContextData{Requester: sessions.RequesterInfo{UserID: userID}}

	manager := &AuthManager{
		userDataManager:           userDataManager,
		sessionContextDataFetcher: func(context.Context) (*sessions.ContextData, error) { return sessionData, nil },
		logger:                    logging.NewNoopLogger().WithName("auth_manager"),
		tracer:                    tracing.NewTracer(tracing.NewNoopTracerProvider().Tracer("auth_manager")),
	}

	result, err := manager.NewTOTPSecret(ctx, input)

	assert.Error(t, err)
	assert.Nil(t, result)
	mock.AssertExpectationsForObjects(t, userDataManager)
}
