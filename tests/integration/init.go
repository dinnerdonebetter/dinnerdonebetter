package integration

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

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
		ID:              "1",
		TwoFactorSecret: "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=",
		Username:        "exampleUser",
		HashedPassword:  "integration-tests-are-cool",
	}
)

func init() {
	ctx, span := tracing.StartSpan(context.Background())
	defer span.End()

	parsedURLToUse = testutils.DetermineServiceURL()
	urlToUse = parsedURLToUse.String()
	logger := logging.ProvideLogger(logging.Config{Provider: logging.ProviderZerolog})

	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(ctx, urlToUse)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
	time.Sleep(10 * time.Second)
}
