package grpc

import (
	paymentsmanager "github.com/dinnerdonebetter/backend/internal/domain/payments/manager"
	paymentssvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/payments"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "payments_service"
)

var _ paymentssvc.PaymentsServiceServer = (*serviceImpl)(nil)

type serviceImpl struct {
	paymentssvc.UnimplementedPaymentsServiceServer
	tracer          tracing.Tracer
	logger          logging.Logger
	paymentsManager paymentsmanager.PaymentsDataManager
}

func NewService(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	paymentsManager paymentsmanager.PaymentsDataManager,
) paymentssvc.PaymentsServiceServer {
	return &serviceImpl{
		logger:          logging.EnsureLogger(logger).WithName(o11yName),
		tracer:          tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		paymentsManager: paymentsManager,
	}
}
