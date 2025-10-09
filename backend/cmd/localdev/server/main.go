package main

import (
	"context"
	"log"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
)

const (
	apiConfigurationFilepath = "deploy/environments/testing/config_files/integration-tests-config.json"
)

func main() {
	ctx := context.Background()

	// create premade admin user
	premadeAdminUser := &identity.User{
		ID:              identifiers.New(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  "admin_pass",
	}

	cfg, err := config.LoadConfigFromPath[config.APIServiceConfig](ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	server, createdClient, err := localdev.AllInOne(ctx, cfg, premadeAdminUser)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf(`
created oauth2 client ID:     %s
created oauth2 client Secret: %s
`, createdClient.ClientID, createdClient.ClientSecret)

	log.Println("starting server")
	server.Run()
}
