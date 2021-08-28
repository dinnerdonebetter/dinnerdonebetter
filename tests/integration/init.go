package integration

import (
	"context"
	"fmt"
	"math"
	"net/url"
	"strings"
	"time"

	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"
)

const (
	debug                = true
	nonexistentID uint64 = math.MaxUint32
)

var (
	urlToUse       string
	parsedURLToUse *url.URL

	premadeAdminUser = &types.User{
		ID:              1,
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
	time.Sleep(3 * time.Second)
}
