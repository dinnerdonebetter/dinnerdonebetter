package users

import (
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/encoding"
	"gitlab.com/prixfixe/prixfixe/internal/observability/logging"
	"gitlab.com/prixfixe/prixfixe/internal/observability/metrics"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/internal/random"
	"gitlab.com/prixfixe/prixfixe/internal/routing"
	authservice "gitlab.com/prixfixe/prixfixe/internal/services/authentication"
	"gitlab.com/prixfixe/prixfixe/internal/uploads"
	"gitlab.com/prixfixe/prixfixe/internal/uploads/images"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	serviceName        = "users_service"
	counterDescription = "number of users managed by the users service"
	counterName        = metrics.CounterName("users")
)

var _ types.UserDataService = (*service)(nil)

type (
	// RequestValidator validates request.
	RequestValidator interface {
		Validate(req *http.Request) (bool, error)
	}

	// service handles our users.
	service struct {
		userDataManager           types.UserDataManager
		householdDataManager      types.HouseholdDataManager
		authSettings              *authservice.Config
		authenticator             authentication.Authenticator
		logger                    logging.Logger
		encoderDecoder            encoding.ServerEncoderDecoder
		userIDFetcher             func(*http.Request) string
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		userCounter               metrics.UnitCounter
		secretGenerator           random.Generator
		imageUploadProcessor      images.ImageUploadProcessor
		uploadManager             uploads.UploadManager
		tracer                    tracing.Tracer
	}
)

// ProvideUsersService builds a new UsersService.
func ProvideUsersService(
	authSettings *authservice.Config,
	logger logging.Logger,
	userDataManager types.UserDataManager,
	householdDataManager types.HouseholdDataManager,
	authenticator authentication.Authenticator,
	encoder encoding.ServerEncoderDecoder,
	counterProvider metrics.UnitCounterProvider,
	imageUploadProcessor images.ImageUploadProcessor,
	uploadManager uploads.UploadManager,
	routeParamManager routing.RouteParamManager,
) types.UserDataService {
	return &service{
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		userDataManager:           userDataManager,
		householdDataManager:      householdDataManager,
		authenticator:             authenticator,
		userIDFetcher:             routeParamManager.BuildRouteParamStringIDFetcher(UserIDURIParamKey),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		encoderDecoder:            encoder,
		authSettings:              authSettings,
		userCounter:               metrics.EnsureUnitCounter(counterProvider, logger, counterName, counterDescription),
		secretGenerator:           random.NewGenerator(logger),
		tracer:                    tracing.NewTracer(serviceName),
		imageUploadProcessor:      imageUploadProcessor,
		uploadManager:             uploadManager,
	}
}
