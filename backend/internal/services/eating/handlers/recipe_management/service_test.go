package recipemanagement

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	mockencoding "github.com/dinnerdonebetter/backend/internal/lib/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/lib/routing/mock"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/lib/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/images"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                      logging.NewNoopLogger(),
		recipeManagementDataManager: RecipeManagementDataManagerMock{},
		recipeIDFetcher:             func(req *http.Request) string { return "" },
		encoderDecoder:              encoding.ProvideServerEncoderDecoder(nil, nil, encoding.ContentTypeJSON),
		tracer:                      tracing.NewTracerForTest("test"),
		cfg: &Config{
			UseSearchService: false,
		},
	}
}

func TestProvideRecipesService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()

		for _, k := range allURIKeys {
			rpm.On("BuildRouteParamStringIDFetcher", k).Return(func(*http.Request) string { return "" })
		}

		cfg := &Config{
			Uploads: uploads.Config{
				Storage: objectstorage.Config{
					BucketName: t.Name(),
					Provider:   objectstorage.MemoryProvider,
				},
				Debug: false,
			},
		}
		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		rpm.On(
			"BuildRouteParamStringIDFetcher",
			cfg.Uploads.Storage.UploadFilenameKey,
		).Return(func(*http.Request) string { return "" })

		pp := &mockpublishers.PublisherProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&textsearchcfg.Config{},
			// TODO: &recipeanalysis.MockRecipeAnalyzer{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			&images.MockImageUploadProcessor{},
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
			msgCfg,
			&RecipeManagementDataManagerMock{},
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Uploads: uploads.Config{
				Storage: objectstorage.Config{
					BucketName: t.Name(),
					Provider:   objectstorage.MemoryProvider,
				},
				Debug: false,
			},
		}
		msgCfg := &msgconfig.QueuesConfig{DataChangesTopicName: "data_changes"}

		pp := &mockpublishers.PublisherProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&textsearchcfg.Config{},
			// TODO: &recipeanalysis.MockRecipeAnalyzer{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&images.MockImageUploadProcessor{},
			tracing.NewNoopTracerProvider(),
			metrics.NewNoopMetricsProvider(),
			msgCfg,
			&RecipeManagementDataManagerMock{},
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
