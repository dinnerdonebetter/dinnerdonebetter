package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	name = "meal_plan_finalizer"
)

// ChoresWorker performs meal_plan_finalizer.
type ChoresWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	emailSender           email.Emailer
	customerDataCollector customerdata.Collector
}

// ProvideChoresWorker provides a ChoresWorker.
func ProvideChoresWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *ChoresWorker {
	return &ChoresWorker{
		logger:                logging.EnsureLogger(logger).WithName(name),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		postUpdatesPublisher:  postUpdatesPublisher,
		emailSender:           emailSender,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a pending write.
func (w *ChoresWorker) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	return w.finalizeExpiredMealPlans(ctx)
}
