package localdev

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	apiserver "github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identityconverters "github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	databasecfg "github.com/dinnerdonebetter/backend/internal/platform/database/config"
	"github.com/dinnerdonebetter/backend/internal/platform/database/postgres"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue/redis"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/random"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/auditlogentries"
	identityrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/identity"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/migrations"
	oauthrepo "github.com/dinnerdonebetter/backend/internal/repositories/postgres/oauth"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreatePremadeAdminUser(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	identityRepo identity.Repository,
	dbClient database.Client,
	premadeAdminUser *identity.User,
) (*identity.User, error) {
	hasher := authentication.ProvideArgon2Authenticator(logger, tracerProvider)

	actuallyHashedPass, err := hasher.HashPassword(ctx, premadeAdminUser.HashedPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}
	premadeAdminUser.HashedPassword = actuallyHashedPass

	var user *identity.User
	if user, err = identityRepo.GetUserByUsername(ctx, premadeAdminUser.Username); err == nil {
		return user, nil
	}

	user, err = identityRepo.CreateUser(ctx, identityconverters.ConvertUserToUserDatabaseCreationInput(premadeAdminUser))
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// one-off query because I really don't want to make this functionality concrete
	if _, err = dbClient.DB().Exec(fmt.Sprintf("UPDATE users SET service_role='service_admin' WHERE id='%s'", user.ID)); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	if err = identityRepo.MarkUserTwoFactorSecretAsVerified(ctx, user.ID); err != nil {
		return nil, fmt.Errorf("failed to mark user as verified: %w", err)
	}

	adminUser, err := identityRepo.GetAdminUserByUsername(ctx, user.Username)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin user: %w", err)
	}

	return adminUser, nil
}

func LoadServerConfig(ctx context.Context, apiConfigurationFilepath string) (*config.APIServiceConfig, error) {
	content, err := os.ReadFile(apiConfigurationFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read api configuration file: %w", err)
	}

	decoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	var x *config.APIServiceConfig
	if err = decoder.DecodeBytes(ctx, content, &x); err != nil {
		return nil, fmt.Errorf("failed to decode api configuration file: %w", err)
	}

	return x, nil
}

func CreateOAuth2ClientForService(ctx context.Context, pgc database.Client, dbCfg *databasecfg.Config, oauth2Input *oauth.OAuth2ClientDatabaseCreationInput) (*oauth.OAuth2Client, error) {
	auditRepo := auditlogentries.ProvideAuditLogRepository(nil, nil, pgc)
	oauth2ClientManager := oauthrepo.ProvideOAuthRepository(nil, nil, auditRepo, dbCfg, pgc)

	createdClient, err := oauth2ClientManager.CreateOAuth2Client(ctx, oauth2Input)
	if err != nil {
		return nil, fmt.Errorf("failed to create oauth2 client: %w", err)
	}

	return createdClient, nil
}

func BuildInProcessServer(ctx context.Context, cfg *config.APIServiceConfig) (server *apiserver.Server, databaseClient database.Client, dbCfg *databasecfg.Config, err error) {
	logger, _, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("setting up logger: %w", err)
	}

	redisConfig, _, err := redis.BuildContainerBackedRedisConfig(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connecting to redis: %w", err)
	}
	cfg.Events.Publisher.Provider = msgconfig.ProviderRedis
	cfg.Events.Publisher.Redis = *redisConfig
	cfg.Events.Consumer.Redis = *redisConfig

	// set up a database container, migrate it, and build a connection client
	_, db, dbCfg, err := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
	if err != nil {
		return nil, nil, nil, fmt.Errorf("connecting to postgres: %w", err)
	}
	cfg.Database = *dbCfg

	if err = migrations.NewMigrator(logger, tracing.NewNoopTracerProvider(), db, dbCfg).Migrate(ctx); err != nil {
		return nil, nil, nil, fmt.Errorf(": %w", err)
	}

	databaseClient, err = postgres.ProvideDatabaseClient(ctx, logger, tracing.NewNoopTracerProvider(), dbCfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("initializing database client: %w", err)
	}

	// create premade admin user
	server, err = apiserver.NewServer(ctx, logger, cfg)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("building API server: %w", err)
	}

	return server, databaseClient, dbCfg, nil
}

func AllInOne(ctx context.Context, cfg *config.APIServiceConfig, adminUser *identity.User, oauth2Input *oauth.OAuth2ClientDatabaseCreationInput) (*apiserver.Server, *oauth.OAuth2Client, error) {
	server, databaseClient, _, err := BuildInProcessServer(ctx, cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("building in-process server: %w", err)
	}

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("getting o11y pillars: %w", err)
	}

	redisConfig, _, err := redis.BuildContainerBackedRedisConfig(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("starting redis container: %w", err)
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
		return nil, nil, fmt.Errorf("migrating database: %w", err)
	}

	if adminUser == nil {
		adminUser = &identity.User{
			ID:              identifiers.New(),
			TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
			EmailAddress:    "integration_tests@example.email",
			Username:        "admin_user",
			HashedPassword:  "admin_pass",
		}
	}

	// create premade admin user
	auditLogRepo := auditlogentries.ProvideAuditLogRepository(logger, tracerProvider, databaseClient)
	identityRepo := identityrepo.ProvideIdentityRepository(logger, tracerProvider, auditLogRepo, databaseClient)
	if _, err = CreatePremadeAdminUser(ctx, logger, tracerProvider, identityRepo, databaseClient, adminUser); err != nil {
		return nil, nil, fmt.Errorf("creating admin user: %w", err)
	}

	if oauth2Input == nil {
		oauth2Input = &oauth.OAuth2ClientDatabaseCreationInput{
			ID:          identifiers.New(),
			Name:        "integration_client",
			Description: "integration test client",
		}

		oauth2Input.ClientID, err = random.GenerateHexEncodedString(ctx, oauth.ClientIDSize)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate client ID: %w", err)
		}

		oauth2Input.ClientSecret, err = random.GenerateHexEncodedString(ctx, oauth.ClientSecretSize)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate client secret: %w", err)
		}
	}

	createdClient, err := CreateOAuth2ClientForService(ctx, databaseClient, dbCfg, oauth2Input)
	if err != nil {
		return nil, nil, fmt.Errorf("creating oauth2 client: %w", err)
	}

	logger.WithValues(map[string]any{
		keys.OAuth2ClientIDKey: createdClient.ClientID,
		"client_sec":           createdClient.ClientSecret,
	}).Info("created oauth2 client")

	return server, createdClient, nil
}

func BuildInsecureOAuthedGRPCClient(
	ctx context.Context,
	createdClientID,
	createdClientSecret,
	httpTestServerAddress,
	grpcServerAddress,
	token string,
) (client.Client, error) {
	state, err := random.GenerateBase64EncodedString(ctx, 32)
	if err != nil {
		return nil, fmt.Errorf("generating state: %w", err)
	}

	oauth2Config := oauth2.Config{
		ClientID:     createdClientID,
		ClientSecret: createdClientSecret,
		Scopes:       []string{"anything"}, // TODO: This should be nil-able
		RedirectURL:  httpTestServerAddress,
		Endpoint: oauth2.Endpoint{
			AuthStyle: oauth2.AuthStyleInParams,
			AuthURL:   httpTestServerAddress + "/oauth2/authorize",
			TokenURL:  httpTestServerAddress + "/oauth2/token",
		},
	}

	authCodeURL := oauth2Config.AuthCodeURL(
		state,
		oauth2.SetAuthURLParam("code_challenge_method", "plain"),
	)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		authCodeURL,
		http.NoBody,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth2 code: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Location", "localhost")

	httpClient := tracing.BuildTracedHTTPClient()
	httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth2 code: %w", err)
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			log.Println("failed to close oauth2 response body", err)
		}
	}()

	const (
		codeKey = "code"
	)

	rl, err := res.Location()
	if err != nil {
		return nil, fmt.Errorf("getting location value from response: %w", err)
	}

	code := rl.Query().Get(codeKey)
	if code == "" {
		return nil, fmt.Errorf("code not returned from oauth2 redirect")
	}

	oauth2Token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("exchanging OAuth2 code: %w", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithPerRPCCredentials(&insecureOAuth{
			TokenSource: oauth2Config.TokenSource(ctx, oauth2Token),
		}),
	}

	c, err := client.BuildClient(grpcServerAddress, opts...)
	if err != nil {
		return nil, fmt.Errorf("building client: %w", err)
	}

	return c, nil
}

// Custom insecure OAuth2 credentials that skip security checks.
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

func FetchLoginTokenForUser(ctx context.Context, grpcServerAddr string, loginInput *authsvc.UserLoginInput) (string, error) {
	unauthedClient, err := client.BuildUnauthenticatedGRPCClient(grpcServerAddr)
	if err != nil {
		return "", fmt.Errorf("initializing client: %w", err)
	}

	tokenRes, err := unauthedClient.LoginForToken(ctx, &authsvc.LoginForTokenRequest{
		Input: loginInput,
	})
	if err != nil {
		return "", fmt.Errorf("fetching login token: %w", err)
	}

	return tokenRes.Result.AccessToken, nil
}
