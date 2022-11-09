package recipes

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	"github.com/prixfixeco/backend/internal/features/recipeanalysis"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/backend/internal/routing/mock"
	"github.com/prixfixeco/backend/internal/storage"
	"github.com/prixfixeco/backend/internal/uploads"
	"github.com/prixfixeco/backend/internal/uploads/images"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:            logging.NewNoopLogger(),
		recipeDataManager: &mocktypes.RecipeDataManager{},
		recipeIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:    mockencoding.NewMockEncoderDecoder(),
		tracer:            tracing.NewTracerForTest("test"),
	}
}

func TestProvideRecipesService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			RecipeIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
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
			&mocktypes.RecipeDataManager{},
			&mocktypes.MockRecipeMediaDataManager{},
			&recipeanalysis.MockRecipeAnalyzer{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			&images.MockImageUploadProcessor{},
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
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
			&mocktypes.RecipeDataManager{},
			&mocktypes.MockRecipeMediaDataManager{},
			&recipeanalysis.MockRecipeAnalyzer{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			&images.MockImageUploadProcessor{},
			tracing.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
