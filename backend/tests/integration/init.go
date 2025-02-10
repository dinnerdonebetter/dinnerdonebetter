package integration

import (
	"context"
	"crypto/tls"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net"
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

	"google.golang.org/grpc/credentials"
)

const (
	debug         = true
	nonexistentID = "_NOT_REAL_LOL_"

	grpcServiceURLEnvVarKey = "TARGET_GRPC_ADDRESS"
	httpServiceURLEnvVarKey = "TARGET_HTTP_ADDRESS"
)

var (
	dbManager database.DataManager

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

	grpcServerAddress := os.Getenv(grpcServiceURLEnvVarKey)
	httpServerAddress := os.Getenv(grpcServiceURLEnvVarKey)

	logger.WithValue(keys.URLKey, grpcServerAddress).Info("checking server")
	ensureServerIsUp(ctx, grpcServerAddress)

	dbAddr := os.Getenv("TARGET_DATABASE")
	if dbAddr == "" {
		panic("empty database grpcAddress provided")
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

	simpleClient := buildUnauthedGRPCClient(grpcServerAddress)

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

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}

func ensureServerIsUp(ctx context.Context, address string) {
	c := buildUnauthedGRPCClient(address)

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

// fakeTLSCreds is a dummy TransportCredentials implementation that
// simulates a TLS connection (even though it doesn’t perform any encryption).
type fakeTLSCreds struct{}

// ClientHandshake simply returns the underlying connection along with a fake TLSInfo.
func (c fakeTLSCreds) ClientHandshake(ctx context.Context, authority string, conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	// We’re not performing a real TLS handshake; we just return a fake TLSInfo.
	fakeState := tls.ConnectionState{}
	return conn, credentials.TLSInfo{State: fakeState}, nil
}

// ServerHandshake does the same on the server side.
func (c fakeTLSCreds) ServerHandshake(conn net.Conn) (net.Conn, credentials.AuthInfo, error) {
	fakeState := tls.ConnectionState{}
	return conn, credentials.TLSInfo{State: fakeState}, nil
}

// Info returns a ProtocolInfo that advertises TLS.
func (c fakeTLSCreds) Info() credentials.ProtocolInfo {
	return credentials.ProtocolInfo{
		// Advertise as TLS.
		SecurityProtocol: "tls",
		SecurityVersion:  "1.2", // you can set an appropriate version string
	}
}

// Clone returns a new instance of fakeTLSCreds.
func (c fakeTLSCreds) Clone() credentials.TransportCredentials {
	return fakeTLSCreds{}
}

// OverrideServerName is a no-op for our fake TLS.
func (c fakeTLSCreds) OverrideServerName(serverName string) error {
	return nil
}

// RequireTransportSecurity indicates that this transport “requires” security.
func (c fakeTLSCreds) RequireTransportSecurity() bool {
	return true
}

// NewFakeTransportCredentials creates a new instance of our fake credentials
// that wrap insecure credentials.
func NewFakeTransportCredentials() credentials.TransportCredentials {
	return fakeTLSCreds{}
}
