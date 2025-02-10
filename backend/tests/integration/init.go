package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"time"

	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/database/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/grpc/messages"
	"github.com/dinnerdonebetter/backend/internal/grpc/service"
	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/identifiers"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	debug         = true
	nonexistentID = "_NOT_REAL_LOL_"

	targetDatabaseEnvVarKey = "TARGET_DATABASE"
	grpcServiceURLEnvVarKey = "TARGET_GRPC_ADDRESS"
	httpServiceURLEnvVarKey = "TARGET_HTTP_ADDRESS"
)

var (
	dbManager database.DataManager

	createdClientID, createdClientSecret string

	premadeAdminUser = &types.User{
		ID:              identifiers.New(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@whatever.gov",
		Username:        "exampleUser",
		HashedPassword:  "integration-tests-are-cool",
	}

	premadeAdminClient service.EatingServiceClient
)

func init() {
	ctx := context.Background()
	logger, err := (&loggingcfg.Config{Provider: loggingcfg.ProviderSlog}).ProvideLogger(ctx)
	if err != nil {
		panic("could not create logger: " + err.Error())
	}

	grpcServerAddress := os.Getenv(grpcServiceURLEnvVarKey)
	httpServerAddress := os.Getenv(httpServiceURLEnvVarKey)

	logger.WithValue(keys.URLKey, grpcServerAddress).Info("checking server")
	ensureServerIsUp(ctx, grpcServerAddress)

	dbAddr := os.Getenv(targetDatabaseEnvVarKey)
	if dbAddr == "" {
		panic("empty database grpcAddress provided")
	}

	cfg := &databasecfg.Config{
		OAuth2TokenEncryptionKey: "                                ", // enough characters to validate
		Debug:                    false,
		RunMigrations:            false,
		MaxPingAttempts:          50,
	}
	if err = cfg.LoadConnectionDetailsFromURL(dbAddr); err != nil {
		panic(err)
	}

	dbManager, err = postgres.ProvideDatabaseClient(ctx, logger, tracing.NewNoopTracerProvider(), cfg)
	if err != nil {
		panic(err)
	}

	hasher := authentication.ProvideArgon2Authenticator(logger, tracing.NewNoopTracerProvider())
	actuallyHashedPass, err := hasher.HashPassword(ctx, premadeAdminUser.HashedPassword)
	if err != nil {
		panic(err)
	}

	if _, err = dbManager.GetUserByUsername(ctx, premadeAdminUser.Username); err != nil {
		_, creationErr := dbManager.CreateUser(ctx, &types.UserDatabaseCreationInput{
			ID:              premadeAdminUser.ID,
			Username:        premadeAdminUser.Username,
			EmailAddress:    premadeAdminUser.EmailAddress,
			HashedPassword:  actuallyHashedPass,
			TwoFactorSecret: premadeAdminUser.TwoFactorSecret,
		})
		if creationErr != nil {
			panic(creationErr)
		}
	}

	if err = dbManager.MarkUserTwoFactorSecretAsVerified(ctx, premadeAdminUser.ID); err != nil && !errors.Is(err, sql.ErrNoRows) {
		panic(err)
	}

	clientID, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		panic(err)
	}

	clientSecret, err := random.GenerateHexEncodedString(ctx, 16)
	if err != nil {
		panic(err)
	}

	createdClient, err := dbManager.CreateOAuth2Client(ctx, &types.OAuth2ClientDatabaseCreationInput{
		ID:           identifiers.New(),
		Name:         "integration_client",
		Description:  "integration test client",
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		panic(err)
	}

	createdClientID, createdClientSecret = createdClient.ClientID, createdClient.ClientSecret

	db, err := sql.Open("pgx", dbAddr)
	if err != nil {
		panic(err)
	}

	if _, err = db.Exec(`UPDATE users SET service_role = $1 WHERE id = $2`, authorization.ServiceAdminRole.String(), premadeAdminUser.ID); err != nil {
		panic(err)
	}

	simpleClient := buildUnauthenticatedGRPCClient(grpcServerAddress)

	code, err := generateTOTPTokenForUserWithoutTest(premadeAdminUser)
	if err != nil {
		panic(err)
	}

	jwtRes, err := simpleClient.AdminLoginForToken(ctx, &messages.AdminLoginForTokenRequest{
		Input: &messages.UserLoginInput{
			Username:  premadeAdminUser.Username,
			Password:  premadeAdminUser.HashedPassword,
			TOTPToken: code,
		},
	})
	if err != nil {
		panic(err)
	}

	premadeAdminClient = buildAuthedGRPCClient(ctx, []string{"household_admin"}, httpServerAddress, grpcServerAddress, jwtRes.Result.AccessToken)
	if _, err = premadeAdminClient.Ping(ctx, nil); err != nil {
		panic(err)
	}

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}

func ensureServerIsUp(ctx context.Context, address string) {
	c := buildUnauthenticatedGRPCClient(address)

	var (
		isDown           = true
		interval         = time.Second
		maxAttempts      = 50
		numberOfAttempts = 0
	)

	for isDown {
		if _, err := c.Ping(ctx, nil); err != nil {
			log.Printf("waiting %s before pinging again", interval)
			time.Sleep(interval)

			numberOfAttempts++
			if numberOfAttempts >= maxAttempts {
				log.Fatal("Maximum number of attempts made, something's gone awry")
			}
		} else {
			isDown = false
		}
	}
}

// Custom insecure OAuth2 credentials that skip security checks
type insecureOAuth struct {
	TokenSource oauth2.TokenSource
}

func (i *insecureOAuth) GetRequestMetadata(_ context.Context, _ ...string) (map[string]string, error) {
	token, err := i.TokenSource.Token()
	if err != nil {
		return nil, err
	}

	return map[string]string{"authorization": token.Type() + " " + token.AccessToken}, nil
}

func (i *insecureOAuth) RequireTransportSecurity() bool {
	return false // Explicitly allow insecure transport
}
