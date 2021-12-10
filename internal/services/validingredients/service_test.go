package validingredients

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/publishers/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	"github.com/prixfixeco/api_server/internal/search"
	mocksearch "github.com/prixfixeco/api_server/internal/search/mock"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func buildTestService() *service {
	return &service{
		logger:                     logging.NewNoopLogger(),
		validIngredientDataManager: &mocktypes.ValidIngredientDataManager{},
		validIngredientIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:             mockencoding.NewMockEncoderDecoder(),
		search:                     &mocksearch.IndexManager{},
		tracer:                     tracing.NewTracer("test"),
	}
}

func TestProvideValidIngredientsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			ValidIngredientIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return(&mockpublishers.Publisher{}, nil)

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_ingredients"),
			[]string{"name", "variant", "description", "warning"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			rpm,
			pp,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing pre-writes producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_ingredients"),
			[]string{"name", "variant", "description", "warning"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			nil,
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-updates producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_ingredients"),
			[]string{"name", "variant", "description", "warning"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			nil,
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing pre-archives producer", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.PreWritesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreUpdatesTopicName).Return(&mockpublishers.Publisher{}, nil)
		pp.On("ProviderPublisher", cfg.PreArchivesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		mockIndexManager := &mocksearch.IndexManager{}
		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_ingredients"),
			[]string{"name", "variant", "description", "warning"},
		).Return(mockIndexManager, nil)

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			nil,
			pp,
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})

	T.Run("with error providing index", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		logger := logging.NewNoopLogger()

		cfg := Config{
			SearchIndexPath:      "example/path",
			PreWritesTopicName:   "pre-writes",
			PreUpdatesTopicName:  "pre-updates",
			PreArchivesTopicName: "pre-archives",
		}

		mockIndexManagerProvider := &mocksearch.IndexManagerProvider{}
		mockIndexManagerProvider.On(
			"ProvideIndexManager",
			testutils.ContextMatcher,
			logger,
			search.IndexName("valid_ingredients"),
			[]string{"name", "variant", "description", "warning"},
		).Return(&mocksearch.IndexManager{}, errors.New("blah"))

		s, err := ProvideService(
			ctx,
			logger,
			&cfg,
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			mockIndexManagerProvider,
			mockrouting.NewRouteParamManager(),
			nil,
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
