package householdinstrumentownerships

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/authentication"
	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "household_instrument_ownerships_service"
)

var _ types.HouseholdInstrumentOwnershipDataService = (*service)(nil)

type (
	// service handles household instrument ownerships.
	service struct {
		logger                                  logging.Logger
		householdInstrumentOwnershipDataManager types.HouseholdInstrumentOwnershipDataManager
		householdInstrumentOwnershipIDFetcher   func(*http.Request) string
		sessionContextDataFetcher               func(*http.Request) (*types.SessionContextData, error)
		dataChangesPublisher                    messagequeue.Publisher
		encoderDecoder                          encoding.ServerEncoderDecoder
		tracer                                  tracing.Tracer
	}
)

// ProvideService builds a new HouseholdInstrumentOwnershipsService.
func ProvideService(
	logger logging.Logger,
	cfg *Config,
	householdInstrumentOwnershipDataManager types.HouseholdInstrumentOwnershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
) (types.HouseholdInstrumentOwnershipDataService, error) {
	dataChangesPublisher, err := publisherProvider.ProvidePublisher(cfg.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	svc := &service{
		logger:                                  logging.EnsureLogger(logger).WithName(serviceName),
		householdInstrumentOwnershipIDFetcher:   routeParamManager.BuildRouteParamStringIDFetcher(HouseholdInstrumentOwnershipIDURIParamKey),
		sessionContextDataFetcher:               authentication.FetchContextFromRequest,
		householdInstrumentOwnershipDataManager: householdInstrumentOwnershipDataManager,
		dataChangesPublisher:                    dataChangesPublisher,
		encoderDecoder:                          encoder,
		tracer:                                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc, nil
}
