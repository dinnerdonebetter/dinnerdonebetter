package workers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/logging/zerolog"
)

func newTestWritesWorker(t *testing.T) *WritesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	dbManager := &database.MockDatabase{}
	postWritesPublisher := &mockpublishers.Publisher{}

	worker, err := ProvideWritesWorker(
		ctx,
		logger,
		dbManager,
		postWritesPublisher,
		&email.MockEmailer{},
		&customerdata.MockCollector{},
		trace.NewNoopTracerProvider(),
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}

func newTestUpdatesWorker(t *testing.T) *UpdatesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	dbManager := &database.MockDatabase{}
	postUpdatesPublisher := &mockpublishers.Publisher{}

	worker, err := ProvideUpdatesWorker(
		ctx,
		logger,
		dbManager,
		postUpdatesPublisher,
		&email.MockEmailer{},
		&customerdata.MockCollector{},
		trace.NewNoopTracerProvider(),
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}

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

func newTestArchivesWorker(t *testing.T) *ArchivesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	dbManager := database.NewMockDatabase()
	postArchivesPublisher := &mockpublishers.Publisher{}

	worker, err := ProvideArchivesWorker(
		ctx,
		logger,
		dbManager,
		postArchivesPublisher,
		&customerdata.MockCollector{},
		trace.NewNoopTracerProvider(),
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}
