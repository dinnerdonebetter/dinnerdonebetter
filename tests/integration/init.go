package integration

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/lib/pq"
	"github.com/segmentio/ksuid"

	"github.com/prixfixeco/api_server/internal/authentication"
	"github.com/prixfixeco/api_server/internal/authorization"
	"github.com/prixfixeco/api_server/internal/database"
	dbconfig "github.com/prixfixeco/api_server/internal/database/config"
	"github.com/prixfixeco/api_server/internal/database/queriers/postgres"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

const (
	debug         = true
	nonexistentID = "_NOT_REAL_LOL_"
)

var (
	urlToUse       string
	parsedURLToUse *url.URL

	premadeAdminUser = &types.User{
		ID:              ksuid.New().String(),
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		EmailAddress:    "integration_tests@prixfixe.email",
		Username:        "exampleUser",
		HashedPassword:  "integration-tests-are-cool",
	}
)

func init() {
	ctx, span := tracing.StartSpan(context.Background())
	defer span.End()

	logger := logging.ProvideLogger(logging.Config{Provider: logging.ProviderZerolog})

	parsedURLToUse = testutils.DetermineServiceURL()
	urlToUse = parsedURLToUse.String()

	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(ctx, urlToUse)

	dbAddr := os.Getenv("DATABASE_ADDRESS")
	if dbAddr == "" {
		panic("empty database address provided")
	}

	cfg := &dbconfig.Config{
		Provider:          dbconfig.PostgresProvider,
		ConnectionDetails: database.ConnectionDetails(dbAddr),
		Debug:             false,
		RunMigrations:     false,
		MaxPingAttempts:   50,
	}
	dbmanager, err := postgres.ProvideDatabaseClient(ctx, logger, cfg, false)
	if err != nil {
		panic(err)
	}

	hasher := authentication.ProvideArgon2Authenticator(logger)
	actuallyHashedPass, err := hasher.HashPassword(ctx, premadeAdminUser.HashedPassword)
	if err != nil {
		panic(err)
	}

	dbmanager.CreateUser(ctx, &types.UserDataStoreCreationInput{
		ID:              premadeAdminUser.ID,
		Username:        premadeAdminUser.Username,
		EmailAddress:    premadeAdminUser.EmailAddress,
		HashedPassword:  actuallyHashedPass,
		TwoFactorSecret: premadeAdminUser.TwoFactorSecret,
	})

	if err = dbmanager.MarkUserTwoFactorSecretAsVerified(ctx, premadeAdminUser.ID); err != nil {
		panic(err)
	}

	db, err := sql.Open("postgres", dbAddr)
	if err != nil {
		panic(err)
	}

	if _, err = db.Exec(`UPDATE users SET service_roles = $1 WHERE id = $2`, authorization.ServiceAdminRole.String(), premadeAdminUser.ID); err != nil {
		panic(err)
	}

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
	time.Sleep(10 * time.Second)
}
