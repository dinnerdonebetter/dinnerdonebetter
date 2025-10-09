package integration

import (
	"context"
	"fmt"
	"log"
	"time"

	apiserver "github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	notificationsrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/notifications"
)

const (
	apiConfigurationFilepath = "../../deploy/environments/testing/config_files/integration-tests-config.json"
)

var (
	dbConnStr                            string
	createdClientID, createdClientSecret string
	databaseClient                       database.Client
	apiServiceConfig                     *config.APIServiceConfig
	notifsRepo                           notifications.Repository
)

func init() {
	ctx := context.Background()

	cfg, err := config.LoadConfigFromPath[config.APIServiceConfig](ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}
	apiServiceConfig = cfg

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var (
		server *apiserver.Server
		dbCfg  *databasecfg.Config
	)

	server, databaseClient, dbCfg, err = localdev.BuildInProcessServer(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	dbConnStr = dbCfg.ConnectionDetails.String()

	// create premade admin user
	auditLogRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, databaseClient)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditLogRepo, databaseClient)
	notifsRepo = notificationsrepo.ProvideNotificationsRepository(nil, nil, auditLogRepo, databaseClient)
	adminUser, err := localdev.CreatePremadeAdminUser(ctx, logger, tracerProvider, identityRepo, databaseClient, premadeAdminUser)
	if err != nil {
		log.Fatal(err)
	}

	createdClient, err := localdev.CreateOAuth2ClientForService(ctx, databaseClient, dbCfg)
	if err != nil {
		log.Fatal(err)
	}
	createdClientID, createdClientSecret = createdClient.ClientID, createdClient.ClientSecret

	go server.Run()

	fmt.Printf("DB conn str: %s", dbCfg.ConnectionDetails.String())
	dbConnStr = dbCfg.ConnectionDetails.String()
	fmt.Println("db conn str: " + dbConnStr)

	// accursed, but nevertheless we ball.
	time.Sleep(1 * time.Second)

	adminClient, err = createClientForUser(ctx, adminUser)
	if err != nil {
		log.Fatal(err)
	}
}
