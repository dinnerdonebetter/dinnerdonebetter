package validenumerations

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/internalerrors"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/internal/lib/search/text"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "valid_ingredients_service"
)

var _ types.ValidEnumerationDataService = (*service)(nil)

type (
	// service handles valid enumerations.
	service struct {
		logger                                  logging.Logger
		validIngredientStatesSearchIndex        textsearch.IndexSearcher[types.ValidIngredientStateSearchSubset]
		validInstrumentSearchIndex              textsearch.IndexSearcher[types.ValidInstrumentSearchSubset]
		validMeasurementUnitSearchIndex         textsearch.IndexSearcher[types.ValidMeasurementUnitSearchSubset]
		validIngredientSearchIndex              textsearch.IndexSearcher[types.ValidIngredientSearchSubset]
		validPreparationsSearchIndex            textsearch.IndexSearcher[types.ValidPreparationSearchSubset]
		validVesselsSearchIndex                 textsearch.IndexSearcher[types.ValidVesselSearchSubset]
		validEnumerationDataManager             types.ValidEnumerationDataManager
		validPreparationVesselIDFetcher         func(*http.Request) string
		validVesselIDFetcher                    func(*http.Request) string
		validPreparationInstrumentIDFetcher     func(*http.Request) string
		validMeasurementUnitConversionIDFetcher func(*http.Request) string
		validInstrumentIDFetcher                func(*http.Request) string
		validIngredientStateIDFetcher           func(*http.Request) string
		validIngredientStateIngredientIDFetcher func(*http.Request) string
		validIngredientPreparationIDFetcher     func(*http.Request) string
		validPreparationIDFetcher               func(*http.Request) string
		validIngredientMeasurementUnitIDFetcher func(*http.Request) string
		validIngredientIDFetcher                func(*http.Request) string
		validMeasurementUnitIDFetcher           func(*http.Request) string
		validIngredientGroupIDFetcher           func(*http.Request) string
		sessionContextDataFetcher               func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                    messagequeue.Publisher
		encoderDecoder                          encoding.ServerEncoderDecoder
		tracer                                  tracing.Tracer
		useSearchService                        bool
	}
)

// ProvideService builds a new ValidEnumerationDataService.
func ProvideService(
	cfg *Config,
	logger logging.Logger,
	dataManager types.ValidEnumerationDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	queueConfig *msgconfig.QueuesConfig,
) (types.ValidEnumerationDataService, error) {
	if cfg == nil {
		return nil, internalerrors.NilConfigError("valid enumerations config")
	}

	if queueConfig == nil {
		return nil, internalerrors.NilConfigError("valid enumerations queue config")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                                  logging.EnsureLogger(logger).WithName(serviceName),
		validIngredientGroupIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientGroupIDURIParamKey),
		validPreparationVesselIDFetcher:         routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationVesselIDURIParamKey),
		validVesselIDFetcher:                    routeParamManager.BuildRouteParamStringIDFetcher(ValidVesselIDURIParamKey),
		validPreparationInstrumentIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationInstrumentIDURIParamKey),
		validMeasurementUnitConversionIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitConversionIDURIParamKey),
		validInstrumentIDFetcher:                routeParamManager.BuildRouteParamStringIDFetcher(ValidInstrumentIDURIParamKey),
		validIngredientStateIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIDURIParamKey),
		validIngredientStateIngredientIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientStateIngredientIDURIParamKey),
		validIngredientPreparationIDFetcher:     routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientPreparationIDURIParamKey),
		validPreparationIDFetcher:               routeParamManager.BuildRouteParamStringIDFetcher(ValidPreparationIDURIParamKey),
		validIngredientMeasurementUnitIDFetcher: routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientMeasurementUnitIDURIParamKey),
		validIngredientIDFetcher:                routeParamManager.BuildRouteParamStringIDFetcher(ValidIngredientIDURIParamKey),
		validMeasurementUnitIDFetcher:           routeParamManager.BuildRouteParamStringIDFetcher(ValidMeasurementUnitIDURIParamKey),
		sessionContextDataFetcher:               authentication.FetchContextFromRequest,
		validEnumerationDataManager:             dataManager,
		dataChangesPublisher:                    dataChangesPublisher,
		encoderDecoder:                          encoder,
		tracer:                                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		useSearchService:                        cfg.UseSearchService,
	}

	return svc, nil
}

func (s *service) searchFromDatabase(req *http.Request) bool {
	return strings.TrimSpace(strings.ToLower(req.URL.Query().Get(types.QueryKeySearchWithDatabase))) == "true"
}
