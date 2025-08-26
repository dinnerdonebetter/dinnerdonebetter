package integration

import (
	"context"
	"fmt"
	grpcapi "github.com/dinnerdonebetter/backend/internal/build/services/api/grpc"
	types "github.com/dinnerdonebetter/backend/internal/domain/oauth"
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
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/internal/config"
)

const (
	apiConfigurationFilepath = "../../deploy/environments/testing/config_files/integration-tests-config.json"
)

var (
	createdClientID, createdClientSecret string
)

func deriveServerConfig() (*config.APIServiceConfig, error) {
	wd, _ := os.Getwd()
	fmt.Println(wd)

	content, err := os.ReadFile(apiConfigurationFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read api configuration file: %w", err)
	}

	decoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	var x *config.APIServiceConfig
	if err = decoder.DecodeBytes(context.Background(), content, &x); err != nil {
		return nil, fmt.Errorf("failed to decode api configuration file: %w", err)
	}

	return x, nil
}

func createOAuth2ClientForTests(ctx context.Context, pgc database.Client, dbCfg *databasecfg.Config) error {
	auditRepo := auditlogentries.ProvideAuditLogRepository(nil, nil, pgc)
	oauth2ClientManager := oauth.ProvideOAuthRepository(nil, nil, auditRepo, *dbCfg, pgc)

	clientID, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		return fmt.Errorf("failed to generate client ID: %w", err)
	}

	clientSecret, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		return fmt.Errorf("failed to generate client secret: %w", err)
	}

	createdClient, err := oauth2ClientManager.CreateOAuth2Client(ctx, &types.OAuth2ClientDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "integration_client",
		Description:  "integration test client",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return fmt.Errorf("failed to create oauth2 client: %w", err)
	}

	createdClientID, createdClientSecret = createdClient.ClientID, createdClient.ClientSecret

	return nil
}

func init() {
	ctx := context.Background()

	cfg, err := deriveServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, _, dbCfg, err := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
	if err != nil {
		log.Fatal(err)
	}
	cfg.Database = *dbCfg

	pgc, err := postgres.ProvideDatabaseClient(ctx, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	redisConfig, _, err := redis.BuildContainerBackedRedisConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Events.Publisher.Provider = msgconfig.ProviderRedis
	cfg.Events.Publisher.Redis = *redisConfig
	cfg.Events.Consumer.Redis = *redisConfig

	if err = createOAuth2ClientForTests(ctx, pgc, dbCfg); err != nil {
		log.Fatal(err)
	}

	grpcServer, err := grpcapi.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go grpcServer.Serve()
}
