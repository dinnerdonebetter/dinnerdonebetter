package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	encryptioncfg "github.com/dinnerdonebetter/backend/internal/platform/cryptography/encryption/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"

	"github.com/spf13/cobra"
)

const (
	// Placeholder TOTP secret for bootstrap admin (2FA is marked verified without real TOTP).
	twoFactorSecretPlaceholder = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
	bootstrapEncryptionKey     = "bootstrap-placeholder-encryption-key-32chars"
)

func main() {
	var (
		dbHost        string
		dbPort        uint16
		dbUser        string
		dbPassword    string
		dbName        string
		dbSSLDisable  bool
		adminUsername string
		adminPassword string
		adminEmail    string
	)

	root := &cobra.Command{
		Use:   "bootstrap",
		Short: "Bootstrap an empty database with an admin user and OAuth2 clients",
		Long:  "Creates an admin user (with 2FA pre-verified) and three OAuth2 clients for Admin Webapp, Consumer Webapp, and iOS App.",
		RunE: func(_ *cobra.Command, _ []string) error {
			return runBootstrap(dbHost, dbPort, dbUser, dbPassword, dbName, dbSSLDisable, adminUsername, adminPassword, adminEmail)
		},
	}

	root.Flags().StringVar(&dbHost, "db-host", "", "Postgres host (or DB_HOST)")
	root.Flags().Uint16Var(&dbPort, "db-port", 5432, "Postgres port (or DB_PORT)")
	root.Flags().StringVar(&dbUser, "db-user", "", "Postgres username (or DB_USER)")
	root.Flags().StringVar(&dbPassword, "db-password", "", "Postgres password (or DB_PASSWORD)")
	root.Flags().StringVar(&dbName, "db-name", "", "Postgres database name (or DB_NAME)")
	root.Flags().BoolVar(&dbSSLDisable, "db-ssl-disable", true, "Disable SSL for DB connection (default: true for local/proxy)")
	root.Flags().StringVar(&adminUsername, "username", "", "Admin username to create")
	root.Flags().StringVar(&adminPassword, "password", "", "Admin password (will be hashed with Argon2)")
	root.Flags().StringVar(&adminEmail, "email", "", "Admin email")

	requiredFlags := []string{"username", "email", "password", "db-host", "db-user", "db-password", "db-name"}
	for _, flag := range requiredFlags {
		if err := root.MarkFlagRequired(flag); err != nil {
			log.Fatalln(err)
		}
	}

	if err := root.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func runBootstrap(
	dbHost string, dbPort uint16, dbUser, dbPassword, dbName string, dbSSLDisable bool,
	adminUsername, adminPassword, adminEmail string,
) error {
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		return errors.New("database connection requires --db-host, --db-user, --db-password, --db-name (or DB_HOST, DB_USER, DB_PASSWORD, DB_NAME env vars)")
	}
	if adminEmail == "" {
		adminEmail = adminUsername + "@bootstrap.local"
	}

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	tracerProvider := tracing.NewNoopTracerProvider()

	connDetails := databasecfg.ConnectionDetails{
		Host:       dbHost,
		Port:       dbPort,
		Username:   dbUser,
		Password:   dbPassword,
		Database:   dbName,
		DisableSSL: dbSSLDisable,
	}

	dbConfig := &databasecfg.Config{
		Provider:                 databasecfg.ProviderPostgres,
		MaxPingAttempts:          10,
		PingWaitPeriod:           time.Second,
		ReadConnection:           connDetails,
		WriteConnection:          connDetails,
		Encryption:               encryptioncfg.Config{Provider: encryptioncfg.ProviderSalsa20},
		OAuth2TokenEncryptionKey: bootstrapEncryptionKey,
	}

	// Use client config that includes sslmode=disable when DisableSSL is true
	clientConfig := &bootstrapClientConfig{
		connDetails: connDetails,
	}
	client, err := postgres.ProvideDatabaseClient(ctx, logger, tracerProvider, clientConfig)
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer func() {
		if closeErr := client.Close(); closeErr != nil {
			logger.Error("closing database client", closeErr)
		}
	}()

	auditRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, client)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditRepo, client)
	oauthRepo := oauthrepo.ProvideOAuthRepository(logger, tracerProvider, auditRepo, dbConfig, client)

	hasher := authentication.ProvideArgon2Authenticator(logger, tracerProvider)
	hashedPassword, err := hasher.HashPassword(ctx, adminPassword)
	if err != nil {
		return fmt.Errorf("hashing password: %w", err)
	}

	userInput := &identity.UserDatabaseCreationInput{
		ID:              identifiers.New(),
		Username:        strings.TrimSpace(adminUsername),
		EmailAddress:    strings.TrimSpace(strings.ToLower(adminEmail)),
		FirstName:       "Admin",
		LastName:        "",
		HashedPassword:  hashedPassword,
		TwoFactorSecret: twoFactorSecretPlaceholder,
		AccountName:     "Bootstrap account",
	}

	user, err := identityRepo.CreateUser(ctx, userInput)
	if err != nil {
		if errors.Is(err, database.ErrUserAlreadyExists) {
			return fmt.Errorf("user %q already exists: %w", adminUsername, err)
		}
		return fmt.Errorf("creating user: %w", err)
	}

	_, err = client.WriteDB().ExecContext(ctx, "UPDATE users SET service_role = $1 WHERE id = $2", "service_admin", user.ID)
	if err != nil {
		return fmt.Errorf("promoting user to admin: %w", err)
	}

	if err = identityRepo.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		return fmt.Errorf("marking 2FA as verified: %w", err)
	}

	oauthClients := []*struct {
		name string
		desc string
	}{
		{"Admin Webapp", "Admin web application OAuth2 client"},
		{"Consumer Webapp", "Consumer web application OAuth2 client"},
		{"iOS App", "iOS mobile application OAuth2 client"},
	}

	var createdClients []*oauth.OAuth2Client
	for _, oc := range oauthClients {
		clientID, clientIDErr := random.GenerateHexEncodedString(ctx, oauth.ClientIDSize)
		if clientIDErr != nil {
			return fmt.Errorf("generating client ID for %s: %w", oc.name, clientIDErr)
		}

		clientSecret, clientSecErr := random.GenerateHexEncodedString(ctx, oauth.ClientSecretSize)
		if clientSecErr != nil {
			return fmt.Errorf("generating client secret for %s: %w", oc.name, clientSecErr)
		}

		created, creationErr := oauthRepo.CreateOAuth2Client(ctx, &oauth.OAuth2ClientDatabaseCreationInput{
			ID:           identifiers.New(),
			Name:         oc.name,
			Description:  oc.desc,
			ClientID:     clientID,
			ClientSecret: clientSecret,
		})
		if creationErr != nil {
			return fmt.Errorf("creating OAuth2 client %s: %w", oc.name, creationErr)
		}
		createdClients = append(createdClients, created)
	}

	fmt.Println("Bootstrap complete.")
	fmt.Println()
	fmt.Printf("Admin user: %s created\n", adminUsername)
	fmt.Println("  - Log in at the admin interface with the provided username and password")
	fmt.Println()
	fmt.Println("OAuth2 clients created:")
	for i, c := range createdClients {
		fmt.Printf("  %s:    client_id=%s client_secret=%s\n", oauthClients[i].name, c.ClientID, c.ClientSecret)
	}

	return nil
}

// bootstrapClientConfig implements database.ClientConfig for bootstrap, using URI format
// (with sslmode=disable) when DisableSSL is true.
type bootstrapClientConfig struct {
	connDetails databasecfg.ConnectionDetails
}

var _ database.ClientConfig = (*bootstrapClientConfig)(nil)

func (b *bootstrapClientConfig) GetReadConnectionString() string {
	if b.connDetails.DisableSSL {
		return b.connDetails.URI()
	}
	return b.connDetails.String()
}

func (b *bootstrapClientConfig) GetWriteConnectionString() string {
	return b.GetReadConnectionString()
}

func (b *bootstrapClientConfig) GetMaxPingAttempts() uint64 {
	return 10
}

func (b *bootstrapClientConfig) GetPingWaitPeriod() time.Duration {
	return time.Second
}
