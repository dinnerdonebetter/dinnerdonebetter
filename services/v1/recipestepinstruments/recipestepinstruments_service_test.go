package recipestepinstruments

import (
	"errors"
	"net/http"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/v1/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/v1/metrics"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/v1/metrics/mock"
	mockmodels "gitlab.com/prixfixe/prixfixe/models/v1/mock"

	"github.com/stretchr/testify/assert"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
)

func buildTestService() *Service {
	return &Service{
		logger:                          noop.ProvideNoopLogger(),
		recipeStepInstrumentCounter:     &mockmetrics.UnitCounter{},
		recipeDataManager:               &mockmodels.RecipeDataManager{},
		recipeStepDataManager:           &mockmodels.RecipeStepDataManager{},
		recipeStepInstrumentDataManager: &mockmodels.RecipeStepInstrumentDataManager{},
		recipeIDFetcher:                 func(req *http.Request) uint64 { return 0 },
		recipeStepIDFetcher:             func(req *http.Request) uint64 { return 0 },
		recipeStepInstrumentIDFetcher:   func(req *http.Request) uint64 { return 0 },
		userIDFetcher:                   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:                  &mockencoding.EncoderDecoder{},
		reporter:                        nil,
	}
}

func TestProvideRecipeStepInstrumentsService(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return &mockmetrics.UnitCounter{}, nil
		}

		s, err := ProvideRecipeStepInstrumentsService(
			noop.ProvideNoopLogger(),
			&mockmodels.RecipeDataManager{},
			&mockmodels.RecipeStepDataManager{},
			&mockmodels.RecipeStepInstrumentDataManager{},
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			&mockencoding.EncoderDecoder{},
			ucp,
			nil,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)
	})

	T.Run("with error providing unit counter", func(t *testing.T) {
		var ucp metrics.UnitCounterProvider = func(counterName metrics.CounterName, description string) (metrics.UnitCounter, error) {
			return nil, errors.New("blah")
		}

		s, err := ProvideRecipeStepInstrumentsService(
			noop.ProvideNoopLogger(),
			&mockmodels.RecipeDataManager{},
			&mockmodels.RecipeStepDataManager{},
			&mockmodels.RecipeStepInstrumentDataManager{},
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			func(req *http.Request) uint64 { return 0 },
			&mockencoding.EncoderDecoder{},
			ucp,
			nil,
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
