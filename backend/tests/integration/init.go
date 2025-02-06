package integration

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
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
	loggingcfg "github.com/dinnerdonebetter/backend/internal/lib/observability/logging/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	debug         = true
	nonexistentID = "_NOT_REAL_LOL_"
)

var (
	urlToUse       string
	parsedURLToUse *url.URL
	dbManager      database.DataManager

	createdClientID, createdClientSecret string

	premadeAdminUser = &types.User{
		ID:              identifiers.New(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@dinnerdonebetter.email",
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

	//parsedURLToUse = serverutils.DetermineServiceURL()
	//urlToUse = parsedURLToUse.String()
	urlToUse = "api_server:8000"

	//logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	//serverutils.EnsureServerIsUp(ctx, urlToUse)

	time.Sleep(15 * time.Second)

	dbAddr := os.Getenv("TARGET_DATABASE")
	if dbAddr == "" {
		panic("empty database address provided")
	}

	cfg := &databasecfg.Config{
		OAuth2TokenEncryptionKey: "                                ", // enough characters to validate
		Debug:                    false,
		RunMigrations:            false,
		MaxPingAttempts:          500,
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

	simpleClient := buildUnauthedGRPCClient()

	code, err := generateTOTPTokenForUserWithoutTest(premadeAdminUser)
	if err != nil {
		panic(err)
	}

	jwtRes, err := simpleClient.AdminLoginForToken(ctx, &messages.UserLoginInput{
		Username:  premadeAdminUser.Username,
		Password:  premadeAdminUser.HashedPassword,
		TOTPToken: code,
	})
	if err != nil {
		panic(err)
	}

	// TODO: replace with grpc-based API client builder

	opts := []grpc.DialOption{
		grpc.WithPerRPCCredentials(&tokenCreds{token: jwtRes.AccessToken}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.NewClient(urlToUse, opts...)
	if err != nil {
		panic(err)
	}

	premadeAdminClient = service.NewEatingServiceClient(conn)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}

type tokenCreds struct {
	token string
}

func (t *tokenCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", t.token),
	}, nil
}

func (t *tokenCreds) RequireTransportSecurity() bool {
	return false
}
