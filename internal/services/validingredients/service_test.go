package validingredients

import (
	"errors"
	"net/http"
	"testing"

	mockencoding "gitlab.com/prixfixe/prixfixe/internal/encoding/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	mockmetrics "gitlab.com/prixfixe/prixfixe/internal/observability/metrics/mock"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	mockrouting "gitlab.com/prixfixe/prixfixe/internal/routing/mock"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	mocksearch "gitlab.com/prixfixe/prixfixe/internal/search/mock"
	mocktypes "gitlab.com/prixfixe/prixfixe/pkg/types/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func buildTestService() *service {
	return &service{
		logger:                     logging.NewNoopLogger(),
		validIngredientCounter:     &mockmetrics.UnitCounter{},
		validIngredientDataManager: &mocktypes.ValidIngredientDataManager{},
		validIngredientIDFetcher:   func(req *http.Request) uint64 { return 0 },
		encoderDecoder:             mockencoding.NewMockEncoderDecoder(),
		search:                     &mocksearch.IndexManager{},
		tracer:                     tracing.NewTracer("test"),
	}
}

func TestProvideValidIngredientsService(T *testing.T) {
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
			ValidIngredientIDURIParamKey,
			"valid_ingredient",
		).Return(func(*http.Request) uint64 { return 0 })

		s, err := ProvideService(
			logging.NewNoopLogger(),
			Config{SearchIndexPath: "example/path"},
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			func(path search.IndexPath, name search.IndexName, logger logging.Logger) (search.IndexManager, error) {
				return &mocksearch.IndexManager{}, nil
			},
			rpm,
		)

		assert.NotNil(t, s)
		assert.NoError(t, err)

		mock.AssertExpectationsForObjects(t, rpm)
	})

	T.Run("with error providing index", func(t *testing.T) {
		t.Parallel()

		var ucp metrics.UnitCounterProvider = func(counterName, description string) metrics.UnitCounter {
			return &mockmetrics.UnitCounter{}
		}

		s, err := ProvideService(
			logging.NewNoopLogger(),
			Config{SearchIndexPath: "example/path"},
			&mocktypes.ValidIngredientDataManager{},
			mockencoding.NewMockEncoderDecoder(),
			ucp,
			func(path search.IndexPath, name search.IndexName, logger logging.Logger) (search.IndexManager, error) {
				return nil, errors.New("blah")
			},
			mockrouting.NewRouteParamManager(),
		)

		assert.Nil(t, s)
		assert.Error(t, err)
	})
}
