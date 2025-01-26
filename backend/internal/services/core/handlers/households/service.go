package households

import (
	"fmt"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/lib/authentication"
	"github.com/dinnerdonebetter/backend/internal/lib/encoding"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/lib/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/lib/random"
	"github.com/dinnerdonebetter/backend/internal/lib/routing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "households_service"
)

var _ types.HouseholdDataService = (*service)(nil)

type (
	// service handles to-do list households.
	service struct {
		logger                         logging.Logger
		householdDataManager           types.HouseholdDataManager
		householdMembershipDataManager types.HouseholdUserMembershipDataManager
		tracer                         tracing.Tracer
		encoderDecoder                 encoding.ServerEncoderDecoder
		dataChangesPublisher           messagequeue.Publisher
		secretGenerator                random.Generator
		sessionContextDataFetcher      func(*http.Request) (*types.SessionContextData, error)
		userIDFetcher                  func(*http.Request) string
		householdIDFetcher             func(*http.Request) string
	}
)

// ProvideService builds a new HouseholdsService.
func ProvideService(
	logger logging.Logger,
	householdDataManager types.HouseholdDataManager,
	householdMembershipDataManager types.HouseholdUserMembershipDataManager,
	encoder encoding.ServerEncoderDecoder,
	routeParamManager routing.RouteParamManager,
	publisherProvider messagequeue.PublisherProvider,
	tracerProvider tracing.TracerProvider,
	secretGenerator random.Generator,
	queueConfig *msgconfig.QueuesConfig,
) (types.HouseholdDataService, error) {
	if queueConfig == nil {
		return nil, fmt.Errorf("nil queue config provided")
	}

	dataChangesPublisher, err := publisherProvider.ProvidePublisher(queueConfig.DataChangesTopicName)
	if err != nil {
		return nil, fmt.Errorf("setting up %s data changes publisher: %w", serviceName, err)
	}

	s := &service{
		logger:                         logging.EnsureLogger(logger).WithName(serviceName),
		householdIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(HouseholdIDURIParamKey),
		userIDFetcher:                  routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher:      authentication.FetchContextFromRequest,
		householdDataManager:           householdDataManager,
		householdMembershipDataManager: householdMembershipDataManager,
		encoderDecoder:                 encoder,
		dataChangesPublisher:           dataChangesPublisher,
		secretGenerator:                secretGenerator,
		tracer:                         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return s, nil
}
