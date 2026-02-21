package main

import (
	"context"
	"net/http"
	"os"

	"github.com/dinnerdonebetter/backend/cmd/services/admin/components"
	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	phttp "github.com/dinnerdonebetter/backend/internal/platform/server/http"
)

const (
	o11yName                       = "admin_frontend"
	apiClientContextKey ContextKey = "api_client"
)

type ContextKey string

type AdminFrontendServer struct {
	tracer                                  tracing.Tracer
	logger                                  logging.Logger
	encoder                                 encoding.ServerEncoderDecoder
	cookieManager                           cookies.Manager
	server                                  phttp.Server
	validMeasurementUnitIDRouteParamFetcher func(req *http.Request) string
	validPreparationIDRouteParamFetcher     func(req *http.Request) string
	settingIDRouteParamFetcher              func(req *http.Request) string
	oauth2ClientIDRouteParamFetcher         func(req *http.Request) string
	validIngredientIDRouteParamFetcher      func(req *http.Request) string
	validInstrumentIDRouteParamFetcher      func(req *http.Request) string
	validVesselIDRouteParamFetcher          func(req *http.Request) string
	userIDRouteParamFetcher                 func(req *http.Request) string
	validIngredientStateIDRouteParamFetcher func(req *http.Request) string
	accountIDRouteParamFetcher              func(req *http.Request) string
	recipeIDRouteParamFetcher               func(req *http.Request) string
	waitlistIDRouteParamFetcher             func(req *http.Request) string
	issueReportIDRouteParamFetcher          func(req *http.Request) string
	validPrepTaskConfigIDRouteParamFetcher  func(req *http.Request) string
	productIDRouteParamFetcher              func(req *http.Request) string
	subscriptionIDRouteParamFetcher         func(req *http.Request) string
	config                                  *config.AdminWebappConfig
	componentRenderer                       *components.ComponentRenderer
	developingLocally                       bool
}

func NewAdminFrontendServer(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	encoder encoding.ServerEncoderDecoder,
	rpm routing.RouteParamManager,
	cfg *config.AdminWebappConfig,
) (*AdminFrontendServer, error) {
	cookieMan, err := cookies.NewCookieManager(&cfg.Cookies, tracerProvider)
	if err != nil {
		return nil, err
	}

	metricsProvider, err := cfg.Observability.Metrics.ProvideMetricsProvider(ctx, logger)
	if err != nil {
		return nil, err
	}

	router, err := cfg.Routing.ProvideRouter(logger, tracerProvider, metricsProvider)
	if err != nil {
		return nil, err
	}

	server, err := phttp.ProvideHTTPServer(cfg.HTTPServer, logger, router, tracerProvider)
	if err != nil {
		return nil, err
	}

	s := &AdminFrontendServer{
		tracer:                                  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                                  logging.EnsureLogger(logger).WithName(o11yName),
		developingLocally:                       os.Getenv("DEVELOPING_LOCALLY") == "true",
		componentRenderer:                       components.NewComponentRenderer(),
		userIDRouteParamFetcher:                 rpm.BuildRouteParamStringIDFetcher(userIDURLParamKey),
		accountIDRouteParamFetcher:              rpm.BuildRouteParamStringIDFetcher(accountIDURLParamKey),
		settingIDRouteParamFetcher:              rpm.BuildRouteParamStringIDFetcher(settingIDURLParamKey),
		oauth2ClientIDRouteParamFetcher:         rpm.BuildRouteParamStringIDFetcher(oauth2ClientIDURLParamKey),
		validIngredientIDRouteParamFetcher:      rpm.BuildRouteParamStringIDFetcher(validIngredientIDURLParamKey),
		validInstrumentIDRouteParamFetcher:      rpm.BuildRouteParamStringIDFetcher(validInstrumentIDURLParamKey),
		validVesselIDRouteParamFetcher:          rpm.BuildRouteParamStringIDFetcher(validVesselIDURLParamKey),
		validMeasurementUnitIDRouteParamFetcher: rpm.BuildRouteParamStringIDFetcher(validMeasurementUnitIDURLParamKey),
		validIngredientStateIDRouteParamFetcher: rpm.BuildRouteParamStringIDFetcher(validIngredientStateIDURLParamKey),
		validPreparationIDRouteParamFetcher:     rpm.BuildRouteParamStringIDFetcher(validPreparationIDURLParamKey),
		recipeIDRouteParamFetcher:               rpm.BuildRouteParamStringIDFetcher(recipeIDURLParamKey),
		waitlistIDRouteParamFetcher:             rpm.BuildRouteParamStringIDFetcher(waitlistIDURLParamKey),
		issueReportIDRouteParamFetcher:          rpm.BuildRouteParamStringIDFetcher(issueReportIDURLParamKey),
		validPrepTaskConfigIDRouteParamFetcher:  rpm.BuildRouteParamStringIDFetcher(validPrepTaskConfigIDURLParamKey),
		productIDRouteParamFetcher:              rpm.BuildRouteParamStringIDFetcher(productIDURLParamKey),
		subscriptionIDRouteParamFetcher:         rpm.BuildRouteParamStringIDFetcher(subscriptionIDURLParamKey),
		cookieManager:                           cookieMan,
		encoder:                                 encoder,
		config:                                  cfg,
		server:                                  server,
	}

	s.setupRoutes(router)

	return s, nil
}

func (s *AdminFrontendServer) Serve() {
	s.server.Serve()
}
