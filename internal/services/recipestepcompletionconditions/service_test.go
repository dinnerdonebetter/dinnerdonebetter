package recipestepingredients

import (
	"errors"
	"net/http"
	"testing"

	mockencoding "github.com/prixfixeco/backend/internal/encoding/mock"
	mockpublishers "github.com/prixfixeco/backend/internal/messagequeue/mock"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	mockrouting "github.com/prixfixeco/backend/internal/routing/mock"
	recipesservice "github.com/prixfixeco/backend/internal/services/recipes"
	recipestepsservice "github.com/prixfixeco/backend/internal/services/recipesteps"
	mocktypes "github.com/prixfixeco/backend/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                                   logging.NewNoopLogger(),
		recipeStepCompletionConditionDataManager: &mocktypes.RecipeStepCompletionConditionDataManager{},
		recipeStepCompletionConditionIDFetcher:   func(req *http.Request) string { return "" },
		encoderDecoder:                           mockencoding.NewMockEncoderDecoder(),
		tracer:                                   tracing.NewTracerForTest("test"),
	}
}

func TestProvideRecipeStepCompletionConditionsService(T *testing.T) {
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
			RecipeStepCompletionConditionIDURIParamKey,
		).Return(func(*http.Request) string { return "" })

		cfg := &Config{
			DataChangesTopicName: "data_changes",
		}

		pp := &mockpublishers.ProducerProvider{}
		pp.On("ProviderPublisher", cfg.DataChangesTopicName).Return(&mockpublishers.Publisher{}, nil)

		s, err := ProvideService(
			logging.NewNoopLogger(),
			cfg,
			&mocktypes.RecipeStepCompletionConditionDataManager{},
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
			&mocktypes.RecipeStepCompletionConditionDataManager{},
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
