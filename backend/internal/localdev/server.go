package localdev

import (
	"context"
	"fmt"
	"os"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	apiserver "github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
)

func CreatePremadeAdminUser(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepo identity.Repository,
	dbClient database.Client,
	premadeAdminUser *identity.User,
) (*identity.User, error) {
	hasher := authentication.ProvideArgon2Authenticator(logger, tracerProvider)

	actuallyHashedPass, err := hasher.HashPassword(ctx, premadeAdminUser.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	premadeAdminUser.HashedPassword = actuallyHashedPass

	var user *identity.User
	if user, err = identityRepo.GetUserByUsername(ctx, premadeAdminUser.Username); err == nil {
		return user, nil
	}

	user, err = identityRepo.CreateUser(ctx, identityconverters.ConvertUserToUserDatabaseCreationInput(premadeAdminUser))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// one-off query because I really don't want to make this functionality concrete
	if _, err = dbClient.DB().Exec(fmt.Sprintf("UPDATE users SET service_role='service_admin' WHERE id='%s'", user.ID)); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if err = identityRepo.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to mark user as verified: %w", err)
	}

	adminUser, err := identityRepo.GetAdminUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return adminUser, nil
}

func LoadServerConfig(ctx context.Context, apiConfigurationFilepath string) (*config.APIServiceConfig, error) {
	content, err := os.ReadFile(apiConfigurationFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read api configuration file: %w", err)
	}

	decoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	var x *config.APIServiceConfig
	if err = decoder.DecodeBytes(ctx, content, &x); err != nil {
		return nil, fmt.Errorf("failed to decode api configuration file: %w", err)
	}

	return x, nil
}

func CreateOAuth2ClientForService(ctx context.Context, pgc database.Client, dbCfg *databasecfg.Config) (*oauth.OAuth2Client, error) {
	auditRepo := auditlogentries.ProvideAuditLogRepository(nil, nil, pgc)
	oauth2ClientManager := oauthrepo.ProvideOAuthRepository(nil, nil, auditRepo, dbCfg, pgc)

	clientID, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client ID: %w", err)
	}

	clientSecret, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		return nil, fmt.Errorf("failed to generate client secret: %w", err)
	}

	createdClient, err := oauth2ClientManager.CreateOAuth2Client(ctx, &oauth.OAuth2ClientDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "integration_client",
		Description:  "integration test client",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 client: %w", err)
	}

	return createdClient, nil
}

func BuildInProcessServer(ctx context.Context, cfg *config.APIServiceConfig) (server *apiserver.Server, databaseClient database.Client, dbCfg *databasecfg.Config, err error) {
	logger, _, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("setting up logger: %w", err)
	}

	redisConfig, _, err := redis.BuildContainerBackedRedisConfig(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connecting to redis: %w", err)
	}
	cfg.Events.Publisher.Provider = msgconfig.ProviderRedis
	cfg.Events.Publisher.Redis = *redisConfig
	cfg.Events.Consumer.Redis = *redisConfig

	// set up a database container, migrate it, and build a connection client
	_, db, dbCfg, err := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connecting to postgres: %w", err)
	}
	cfg.Database = *dbCfg

	if err = migrations.NewMigrator(logger, tracing.NewNoopTracerProvider(), db, dbCfg).Migrate(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf(": %w", err)
	}

	databaseClient, err = postgres.ProvideDatabaseClient(ctx, logger, tracing.NewNoopTracerProvider(), dbCfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("initializing database client: %w", err)
	}

	// create premade admin user
	server, err = apiserver.NewServer(ctx, logger, cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("building API server: %w", err)
	}

	return server, databaseClient, dbCfg, nil
}
