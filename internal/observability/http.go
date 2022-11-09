package observability

import (
	"net/http"
	"time"

	"github.com/prixfixeco/backend/internal/observability/tracing"
)

// HTTPClient returns a new *http.Client.
func HTTPClient() *http.Client {
	return &http.Client{Transport: tracing.BuildTracedHTTPTransport(time.Second)}
}
