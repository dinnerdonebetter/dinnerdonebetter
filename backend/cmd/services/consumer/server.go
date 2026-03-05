package main

import (
	"context"
	"os"

	"github.com/dinnerdonebetter/backend/cmd/services/consumer/components"
	"github.com/dinnerdonebetter/backend/internal/authentication/cookies"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/platform/routing"
	phttp "github.com/dinnerdonebetter/backend/internal/platform/server/http"
)

const o11yName = "consumer_frontend"

type ConsumerFrontendServer struct {
	tracer            tracing.Tracer
	logger            logging.Logger
	encoder           encoding.ServerEncoderDecoder
	cookieManager     cookies.Manager
	server            phttp.Server
	config            *config.ConsumerWebappConfig
	componentRenderer *components.ComponentRenderer
	developingLocally bool
}

func NewConsumerFrontendServer(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	encoder encoding.ServerEncoderDecoder,
	_ routing.RouteParamManager,
	cfg *config.ConsumerWebappConfig,
) (*ConsumerFrontendServer, error) {
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

	serviceName := "consumer_webapp"
	if cfg.Routing.Chi != nil && cfg.Routing.Chi.ServiceName != "" {
		serviceName = cfg.Routing.Chi.ServiceName
	}
	server, err := phttp.ProvideHTTPServer(cfg.HTTPServer, logger, router, tracerProvider, serviceName)
	if err != nil {
		return nil, err
	}

	s := &ConsumerFrontendServer{
		tracer:            tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:            logging.EnsureLogger(logger).WithName(o11yName),
		developingLocally: os.Getenv("DEVELOPING_LOCALLY") == "true",
		componentRenderer: components.NewComponentRenderer(),
		cookieManager:     cookieMan,
		encoder:           encoder,
		config:            cfg,
		server:            server,
	}

	s.setupRoutes(router)

	return s, nil
}

func (s *ConsumerFrontendServer) Serve() {
	s.server.Serve()
}
