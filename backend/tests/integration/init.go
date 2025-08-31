package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	apiserver "github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres/migrations"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
)

const (
	apiConfigurationFilepath = "../../deploy/environments/testing/config_files/integration-tests-config.json"
)

var (
	dbConnStr                            string
	createdClientID, createdClientSecret string
	databaseClient                       database.Client
	shutdownFunc                         func()
)

func init() {
	ctx := context.Background()

	cfg, err := deriveServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// setup Pubsub queue infra
	redisConfig, redisShutdownFunc, err := redis.BuildContainerBackedRedisConfig(ctx)
	if err != nil {
		log.Fatal(err)
	}
	cfg.Events.Publisher.Provider = msgconfig.ProviderRedis
	cfg.Events.Publisher.Redis = *redisConfig
	cfg.Events.Consumer.Redis = *redisConfig

	// set up a database container, migrate it, and build a connection client
	dbContainer, db, dbCfg, err := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
	if err != nil {
		log.Fatal(err)
	}
	cfg.Database = *dbCfg

	if err = migrations.NewMigrator(logger, tracing.NewNoopTracerProvider(), db, dbCfg).Migrate(ctx); err != nil {
		log.Fatal(err)
	}

	databaseClient, err = postgres.ProvideDatabaseClient(ctx, logger, tracing.NewNoopTracerProvider(), dbCfg)
	if err != nil {
		log.Fatal(err)
	}

	// create premade admin user
	auditLogRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, databaseClient)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditLogRepo, databaseClient)
	adminUser, err := createPremadeAdminUser(ctx, logger, tracerProvider, identityRepo, databaseClient)
	if err != nil {
		log.Fatal(err)
	}

	if err = createOAuth2ClientForTests(ctx, databaseClient, dbCfg); err != nil {
		log.Fatal(err)
	}

	server, err := apiserver.NewServer(ctx, logger, cfg)
	if err != nil {
		log.Fatal(err)
	}

	shutdownFunc = func() {
		databaseClient.Close()
		redisShutdownFunc(context.Background())
		dbContainer.Stop(context.Background(), pointer.To(10*time.Second))
	}

	go server.Run()

	fmt.Printf("DB conn str: %s", dbCfg.ConnectionDetails.String())
	dbConnStr = dbCfg.ConnectionDetails.String()

	time.Sleep(1 * time.Second)

	adminClient, err = createClientForUser(ctx, []string{"service_admin"}, adminUser)
	if err != nil {
		log.Fatal(err)
	}
}
