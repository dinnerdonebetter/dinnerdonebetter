package main

import (
	"context"
	"log"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
)

const (
	apiConfigurationFilepath = "deploy/environments/testing/config_files/integration-tests-config.json"
)

func main() {
	ctx := context.Background()

	cfg, err := localdev.LoadServerConfig(ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	server, databaseClient, _, err := localdev.BuildInProcessServer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
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

	// set up a database container, migrate it, and build a connection client
	_, db, dbCfg, err := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
	if err != nil {
		log.Fatal(err)
	}
	cfg.Database = *dbCfg

	if err = migrations.NewMigrator(logger, tracing.NewNoopTracerProvider(), db, dbCfg).Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	// create premade admin user
	premadeAdminUser := &identity.User{
		ID:              identifiers.New(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  "admin_pass",
	}
	auditLogRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, databaseClient)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditLogRepo, databaseClient)
	if _, err = localdev.CreatePremadeAdminUser(ctx, logger, tracerProvider, identityRepo, databaseClient, premadeAdminUser); err != nil {
		log.Fatal(err)
	}

	createdClient, err := localdev.CreateOAuth2ClientForService(ctx, databaseClient, dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	logger.WithValue(keys.OAuth2ClientIDKey, createdClient.ClientID).WithValue("client_secret", createdClient.ClientSecret).Info("created oauth2 client")

	log.Println("starting server")
	server.Run()
}
