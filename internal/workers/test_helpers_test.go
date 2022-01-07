package workers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

func newTestChoresWorker(t *testing.T) *ChoresWorker {
	t.Helper()

	worker := ProvideChoresWorker(
		zerolog.NewZerologLogger(),
		&database.MockDatabase{},
		&mockpublishers.Publisher{},
		&email.MockEmailer{},
		&customerdata.MockCollector{},
		trace.NewNoopTracerProvider(),
	)
	assert.NotNil(t, worker)

	return worker
}
