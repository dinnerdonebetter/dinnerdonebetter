package frontend

import (
	"context"
	"fmt"
	"strings"

	"github.com/prixfixeco/api_server/internal/observability/keys"
	logcfg "github.com/prixfixeco/api_server/internal/observability/logging/config"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

var urlToUse string

func init() {
	u := testutils.DetermineServiceURL()
	urlToUse = u.String()

	logger := (&logcfg.Config{Provider: logcfg.ProviderZerolog}).ProvideLogger()
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(context.Background(), urlToUse)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}
