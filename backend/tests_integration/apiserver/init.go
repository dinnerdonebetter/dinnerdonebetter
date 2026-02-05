package integration

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	apiserver "github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/notifications"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
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
	httpTestServerAddress                string
)

// getFreePort asks the OS for a free open port that is ready to use.
func getFreePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}

	if err = l.Close(); err != nil {
		return 0, err
	}

	return l.Addr().(*net.TCPAddr).Port, nil
}

func init() {
	ctx := context.Background()

	cfg, err := config.LoadConfigFromPath[config.APIServiceConfig](ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	// Use random ports to avoid conflicts with other running instances
	httpPort, err := getFreePort()
	if err != nil {
		log.Fatal(err)
	}
	grpcPort, err := getFreePort()
	if err != nil {
		log.Fatal(err)
	}

	cfg.HTTPServer.Port = uint16(httpPort)
	cfg.GRPCServer.Port = uint16(grpcPort)
	httpTestServerAddress = fmt.Sprintf("http://localhost:%d", httpPort)

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
	dbConnStr = dbCfg.ReadConnection.String()

	// create premade admin user
	auditLogRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, databaseClient)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditLogRepo, databaseClient)
	notifsRepo = notificationsrepo.ProvideNotificationsRepository(nil, nil, auditLogRepo, databaseClient)
	adminUser, err := localdev.CreatePremadeAdminUser(ctx, logger, tracerProvider, identityRepo, databaseClient, premadeAdminUser)
	if err != nil {
		log.Fatal(err)
	}

	createdClient, err := localdev.CreateOAuth2ClientForService(ctx, databaseClient, dbCfg, &oauth.OAuth2ClientDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "integration_client",
		Description:  "integration test client",
		ClientID:     random.MustGenerateHexEncodedString(ctx, oauth.ClientIDSize),
		ClientSecret: random.MustGenerateHexEncodedString(ctx, oauth.ClientSecretSize),
	})
	if err != nil {
		log.Fatal(err)
	}
	createdClientID, createdClientSecret = createdClient.ClientID, createdClient.ClientSecret

	go server.Run()

	fmt.Printf("DB conn str: %s", dbCfg.ReadConnection.String())
	dbConnStr = dbCfg.ReadConnection.String()
	fmt.Println("db conn str: " + dbConnStr)

	// accursed, but nevertheless we ball.
	time.Sleep(1 * time.Second)

	adminClient, err = createClientForUser(ctx, adminUser)
	if err != nil {
		log.Fatal(err)
	}
}
