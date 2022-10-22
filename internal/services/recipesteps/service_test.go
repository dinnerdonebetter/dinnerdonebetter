package recipesteps

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockencoding "github.com/prixfixeco/api_server/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/api_server/internal/messagequeue/mock"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/api_server/internal/routing/mock"
	recipesservice "github.com/prixfixeco/api_server/internal/services/recipes"
	"github.com/prixfixeco/api_server/internal/storage"
	"github.com/prixfixeco/api_server/internal/uploads"
	"github.com/prixfixeco/api_server/internal/uploads/images"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                logging.NewNoopLogger(),
		recipeStepDataManager: &mocktypes.RecipeStepDataManager{},
		recipeStepIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:        mockencoding.NewMockEncoderDecoder(),
		tracer:                tracing.NewTracerForTest("test"),
	}
}

func TestProvideRecipeStepsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			recipesservice.RecipeIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			RecipeStepIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: "data_changes",
			Uploads: uploads.Config{
				Storage: storage.Config{
					FilesystemConfig: &storage.FilesystemConfig{RootDirectory: t.Name()},
					BucketName:       t.Name(),
					Provider:         storage.FilesystemProvider,
				},
				Debug: false,
			},
		}

		rpm.On(
			"BuildRouteParamStringIDFetcher",
			cfg.Uploads.Storage.UploadFilenameKey,
		).Return(func(*http.Request) string { return "" })

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.RecipeStepDataManager{},
			&mocktypes.RecipeMediaDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
			&images.MockImageUploadProcessor{},
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
			Uploads: uploads.Config{
				Storage: storage.Config{
					FilesystemConfig: &storage.FilesystemConfig{RootDirectory: t.Name()},
					BucketName:       t.Name(),
					Provider:         storage.FilesystemProvider,
				},
				Debug: false,
			},
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			context.Background(),
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.RecipeStepDataManager{},
			&mocktypes.RecipeMediaDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			tracing.NewNoopTracerProvider(),
			&images.MockImageUploadProcessor{},
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
