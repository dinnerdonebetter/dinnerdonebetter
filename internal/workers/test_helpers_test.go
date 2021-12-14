package workers

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func newTestWritesWorker(t *testing.T) *WritesWorker {
	t.Helper()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	dbManager := &database.MockDatabase{}
	postWritesPublisher := &mockpublishers.Publisher{}
	indexManagerProvider := &mocksearch.IndexManagerProvider{}
	indexManager := &mocksearch.IndexManager{}

	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_instruments"),
		[]string{"name", "variant", "description", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_ingredients"),
		[]string{"name", "variant", "description", "warning", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_preparations"),
		[]string{"name", "description", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_ingredient_preparations"),
		[]string{"notes", "validPreparationID", "validIngredientID"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("recipes"),
		[]string{"name", "source", "description", "inspiredByRecipeID"},
	).Return(indexManager, nil)

	worker, err := ProvideWritesWorker(
		ctx,
		logger,
		dbManager,
		postWritesPublisher,
		indexManagerProvider,
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
	indexManagerProvider := &mocksearch.IndexManagerProvider{}
	indexManager := &mocksearch.IndexManager{}

	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_instruments"),
		[]string{"name", "variant", "description", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_ingredients"),
		[]string{"name", "variant", "description", "warning", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_preparations"),
		[]string{"name", "description", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_ingredient_preparations"),
		[]string{"notes", "validPreparationID", "validIngredientID"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("recipes"),
		[]string{"name", "source", "description", "inspiredByRecipeID"},
	).Return(indexManager, nil)

	worker, err := ProvideUpdatesWorker(
		ctx,
		logger,
		dbManager,
		postUpdatesPublisher,
		indexManagerProvider,
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
		logging.NewZerologLogger(),
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
	indexManagerProvider := &mocksearch.IndexManagerProvider{}
	indexManager := &mocksearch.IndexManager{}

	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_instruments"),
		[]string{"name", "variant", "description", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_ingredients"),
		[]string{"name", "variant", "description", "warning", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_preparations"),
		[]string{"name", "description", "icon"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("valid_ingredient_preparations"),
		[]string{"notes", "validPreparationID", "validIngredientID"},
	).Return(indexManager, nil)
	indexManagerProvider.On(
		"ProvideIndexManager",
		testutils.ContextMatcher,
		logger,
		search.IndexName("recipes"),
		[]string{"name", "source", "description", "inspiredByRecipeID"},
	).Return(indexManager, nil)

	worker, err := ProvideArchivesWorker(
		ctx,
		logger,
		dbManager,
		postArchivesPublisher,
		indexManagerProvider,
		&customerdata.MockCollector{},
		trace.NewNoopTracerProvider(),
	)
	require.NotNil(t, worker)
	require.NoError(t, err)

	return worker
}
