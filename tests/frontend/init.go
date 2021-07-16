package frontend

import (
	"context"
	"fmt"
	"strings"

	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	testutils "gitlab.com/prixfixe/prixfixe/tests/utils"
)

var urlToUse string

func init() {
	u := testutils.DetermineServiceURL()
	urlToUse = u.String()

	logger := logging.ProvideLogger(logging.Config{Provider: logging.ProviderZerolog})
	logger.WithValue(keys.URLKey, urlToUse).Info("checking server")
	testutils.EnsureServerIsUp(context.Background(), urlToUse)

	fiftySpaces := strings.Repeat("\n", 50)
	fmt.Printf("%s\tRunning tests%s", fiftySpaces, fiftySpaces)
}
