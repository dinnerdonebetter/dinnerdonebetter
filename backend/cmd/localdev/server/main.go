package main

import (
	"context"
	"log"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/backend/internal/localdev"
)

const (
	apiConfigurationFilepath = "deploy/environments/testing/config_files/integration-tests-config.json"
)

func main() {
	ctx := context.Background()

	// create premade admin user
	premadeAdminUser := &identity.User{
		ID:              strings.Repeat("a", 20),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@example.email",
		Username:        "admin_user",
		HashedPassword:  "admin_pass",
	}

	apiConfig, err := config.LoadConfigFromPath[config.APIServiceConfig](ctx, apiConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	server, _, err := localdev.AllInOne(ctx, apiConfig, premadeAdminUser, &oauth.OAuth2ClientDatabaseCreationInput{
		ID:           strings.Repeat("b", 20),
		Name:         "localdev_admin_client",
		Description:  "localdev admin client",
		ClientID:     strings.Repeat("A", oauth.ClientIDSize),
		ClientSecret: strings.Repeat("A", oauth.ClientSecretSize),
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("starting server")
	server.Run()
}
