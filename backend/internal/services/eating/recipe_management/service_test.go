package recipemanagement

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/business/recipeanalysis"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockencoding "github.com/dinnerdonebetter/backend/internal/encoding/mock"
	msgconfig "github.com/dinnerdonebetter/backend/internal/messagequeue/config"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	mockrouting "github.com/dinnerdonebetter/backend/internal/routing/mock"
	textsearchcfg "github.com/dinnerdonebetter/backend/internal/search/text/config"
	"github.com/dinnerdonebetter/backend/internal/uploads"
	"github.com/dinnerdonebetter/backend/internal/uploads/images"
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"

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

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&textsearchcfg.Config{},
			&recipeanalysis.MockRecipeAnalyzer{},
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

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProvidePublisher", msgCfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&textsearchcfg.Config{},
			&recipeanalysis.MockRecipeAnalyzer{},
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
