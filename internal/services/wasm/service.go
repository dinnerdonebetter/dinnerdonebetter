package wasm

import (
	"context"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	authservice "github.com/dinnerdonebetter/backend/internal/services/authentication"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "wasm_service"
)

type (
	// Service handles valid ingredients.
	Service struct {
		cfg                       *Config
		logger                    logging.Logger
		sessionContextDataFetcher func(*http.Request) (*types.SessionContextData, error)
		encoderDecoder            encoding.ServerEncoderDecoder
		tracer                    tracing.Tracer
	}
)

// ProvideService builds a new ValidIngredientsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	cfg *Config,
	encoder encoding.ServerEncoderDecoder,
	tracerProvider tracing.TracerProvider,
) (*Service, error) {
	svc := &Service{
		cfg:                       cfg,
		logger:                    logging.EnsureLogger(logger).WithName(serviceName),
		sessionContextDataFetcher: authservice.FetchContextFromRequest,
		encoderDecoder:            encoder,
		tracer:                    tracing.NewTracer(tracerProvider.Tracer(serviceName)),
	}

	return svc, nil
}
