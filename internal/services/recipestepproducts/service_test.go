package recipestepproducts

import (
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
	recipestepsservice "github.com/prixfixeco/api_server/internal/services/recipesteps"
	mocktypes "github.com/prixfixeco/api_server/pkg/types/mock"
)

func buildTestService() *service {
	return &service{
		logger:                       logging.NewNoopLogger(),
		recipeStepProductDataManager: &mocktypes.RecipeStepProductDataManager{},
		recipeStepProductIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:               mockencoding.NewMockEncoderDecoder(),
		tracer:                       tracing.NewTracerForTest("test"),
	}
}

func TestProvideRecipeStepProductsService(T *testing.T) {
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
			recipestepsservice.RecipeStepIDURIParamKey,
		).Return(func(*http.Request) string { return "" })
		rpm.On(
			"BuildRouteParamStringIDFetcher",
			RecipeStepProductIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.RecipeStepProductDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			rpm,
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm, pp)
	})

	T.Run("with error providing data changes producer", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return((*mockpublishers.Publisher)(nil), errors.New("blah"))

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.RecipeStepProductDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			nil,
			pp,
			tracing.NewNoopTracerProvider(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, pp)
	})
}
