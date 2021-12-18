package workers

import (
	"context"
	"errors"
	"testing"

	"go.opentelemetry.io/otel/trace"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestProvideWritesWorker(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
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

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			indexManagerProvider,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.NotNil(t, actual)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing first search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
		postArchivesPublisher := &mockpublishers.Publisher{}
		indexManagerProvider := &mocksearch.IndexManagerProvider{}

		indexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_instruments"),
			[]string{"name", "variant", "description", "icon"},
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			indexManagerProvider,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing second search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
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
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			indexManagerProvider,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing third search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
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
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			indexManagerProvider,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing fourth search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
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
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			indexManagerProvider,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})

	T.Run("with error providing fifth search index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()
		dbManager := &database.MockDatabase{}
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
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		actual, err := ProvideWritesWorker(
			ctx,
			logger,
			dbManager,
			postArchivesPublisher,
			indexManagerProvider,
			&email.MockEmailer{},
			&customerdata.MockCollector{},
			trace.NewNoopTracerProvider(),
		)
		assert.Nil(t, actual)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, dbManager, postArchivesPublisher)
	})
}
