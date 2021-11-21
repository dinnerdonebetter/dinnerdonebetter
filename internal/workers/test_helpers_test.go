package workers

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
)

func newTestWritesWorker(t *testing.T) *WritesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	client := &http.Client{}
	dbManager := &database.MockDatabase{}
	postArchivesPublisher := &mockpublishers.Publisher{}
	searchIndexLocation := search.IndexPath(t.Name())
	searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
		return nil, nil
	}

	worker, err := ProvideWritesWorker(
		ctx,
		logger,
		client,
		dbManager,
		postArchivesPublisher,
		searchIndexLocation,
		searchIndexProvider,
		&email.MockEmailer{},
		&customerdata.MockCollector{},
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}

func newTestUpdatesWorker(t *testing.T) *UpdatesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	client := &http.Client{}
	dbManager := &database.MockDatabase{}
	postArchivesPublisher := &mockpublishers.Publisher{}
	searchIndexLocation := search.IndexPath(t.Name())
	searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
		return nil, nil
	}

	worker, err := ProvideUpdatesWorker(
		ctx,
		logger,
		client,
		dbManager,
		postArchivesPublisher,
		searchIndexLocation,
		searchIndexProvider,
		&email.MockEmailer{},
		&customerdata.MockCollector{},
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}

func newTestChoresWorker(t *testing.T) *ChoresWorker {
	t.Helper()

	worker := ProvideChoresWorker(
		logging.NewZerologLogger(),
		&database.MockDatabase{},
		&mockpublishers.Publisher{},
		&email.MockEmailer{},
		&customerdata.MockCollector{},
	)
	assert.NotNil(t, worker)

	return worker
}

func newTestArchivesWorker(t *testing.T) *ArchivesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	client := &http.Client{}
	dbManager := database.NewMockDatabase()
	postArchivesPublisher := &mockpublishers.Publisher{}
	searchIndexLocation := search.IndexPath(t.Name())
	searchIndexProvider := func(context.Context, logging.Logger, *http.Client, search.IndexPath, search.IndexName, ...string) (search.IndexManager, error) {
		return nil, nil
	}

	worker, err := ProvideArchivesWorker(
		ctx,
		logger,
		client,
		dbManager,
		postArchivesPublisher,
		searchIndexLocation,
		searchIndexProvider,
		&customerdata.MockCollector{},
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}
