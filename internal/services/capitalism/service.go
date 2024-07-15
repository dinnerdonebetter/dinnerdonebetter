package capitalism

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/capitalism"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	serviceName string = "capitalism_service"
)

var _ types.CapitalismService = (*service)(nil)

type (
	// service handles valid instruments.
	service struct {
		logger         logging.Logger
		tracer         tracing.Tracer
		paymentManager capitalism.PaymentManager
	}
)

// ProvideService builds a new ValidInstrumentsService.
func ProvideService(
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	paymentManager capitalism.PaymentManager,
) types.CapitalismService {
	svc := &service{
		logger:         logging.EnsureLogger(logger).WithName(serviceName),
		paymentManager: paymentManager,
		tracer:         tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}

	return svc
}
