package invitations

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	routing "gitlab.com/prixfixe/prixfixe/internal/routing"
	"gitlab.com/prixfixe/prixfixe/internal/search"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	counterName        metrics.CounterName = "invitations"
	counterDescription string              = "the number of invitations managed by the invitations service"
	serviceName        string              = "invitations_service"
)

var _ types.InvitationDataService = (*service)(nil)

type (
	// SearchIndex is a type alias for dependency injection's sake.
	SearchIndex search.IndexManager

	// service handles invitations.
	service struct {
		logger                    logging.Logger
		invitationDataManager     types.InvitationDataManager
		invitationIDFetcher       func(*http.Request) uint64
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		invitationCounter         metrics.UnitCounter
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new InvitationsService.
func ProvideService(
	logger logging.Logger,
	cfg Config,
	invitationDataManager types.InvitationDataManager,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	routeParamManager routing.RouteParamManager,
) (types.InvitationDataService, error) {
	svc := &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		invitationIDFetcher:       routeParamManager.BuildRouteParamIDFetcher(logger, InvitationIDURIParamKey, "invitation"),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		invitationDataManager:     invitationDataManager,
		encoderDecoder:            encoder,
		invitationCounter:         metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		tracer:                    tracing.NewTracer(serviceName),
	}

	return svc, nil
}
