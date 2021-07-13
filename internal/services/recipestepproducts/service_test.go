package recipestepproducts

import (
	"net/http"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	recipesservice "gitlab.com/prixfixe/prixfixe/internal/services/recipes"
	recipestepsservice "gitlab.com/prixfixe/prixfixe/internal/services/recipesteps"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                       logging.NewNoopLogger(),
		recipeStepProductCounter:     &mockmetrics.UnitCounter{},
		recipeStepProductDataManager: &mocktypes.RecipeStepProductDataManager{},
		recipeStepProductIDFetcher:   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:               mockencoding.NewMockEncoderDecoder(),
		tracer:                       tracing.NewTracer("test"),
	}
}

func TestProvideRecipeStepProductsService(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		var ucp metrics.UnitCounterProvider = func(counterName, description string) metrics.UnitCounter {
			return &mockmetrics.UnitCounter{}
		}

		rpm := mockrouting.NewRouteParamManager()
		rpm.On(
			"BuildRouteParamIDFetcher",
			mock.IsType(logging.NewNoopLogger()),
			recipesservice.RecipeIDURIParamKey,
			"recipe",
		).Return(func(*http.Request) uint64 { return 0 })
		rpm.On(
			"BuildRouteParamIDFetcher",
			mock.IsType(logging.NewNoopLogger()),
			recipestepsservice.RecipeStepIDURIParamKey,
			"recipe_step",
		).Return(func(*http.Request) uint64 { return 0 })
		rpm.On(
			"BuildRouteParamIDFetcher",
			mock.IsType(logging.NewNoopLogger()),
			RecipeStepProductIDURIParamKey,
			"recipe_step_product",
		).Return(func(*http.Request) uint64 { return 0 })

		s, err := ProvideService(
			logging.NewNoopLogger(),
			Config{},
			&mocktypes.RecipeStepProductDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			rpm,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})
}
