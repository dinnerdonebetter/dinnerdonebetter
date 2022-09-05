package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/messagequeue"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

const (
	advancedPrepStepCreationEnsurerWorkerName = "advanced_prep_step_creation_ensurer"
)

// AdvancedPrepStepCreationEnsurerWorker ensurers advanced prep steps are created.
type AdvancedPrepStepCreationEnsurerWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataManager           database.DataManager
	postUpdatesPublisher  messagequeue.Publisher
	customerDataCollector customerdata.Collector
}

// ProvideAdvancedPrepStepCreationEnsurerWorker provides a AdvancedPrepStepCreationEnsurerWorker.
func ProvideAdvancedPrepStepCreationEnsurerWorker(
	logger logging.Logger,
	dataManager database.DataManager,
	postUpdatesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *AdvancedPrepStepCreationEnsurerWorker {
	return &AdvancedPrepStepCreationEnsurerWorker{
		logger:                logging.EnsureLogger(logger).WithName(advancedPrepStepCreationEnsurerWorkerName),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(advancedPrepStepCreationEnsurerWorkerName)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataManager:           dataManager,
		postUpdatesPublisher:  postUpdatesPublisher,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a pending write.
func (w *AdvancedPrepStepCreationEnsurerWorker) HandleMessage(ctx context.Context, _ []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	return w.ensureAdvancedPrepStepsAreCreated(ctx)
}

func (w *AdvancedPrepStepCreationEnsurerWorker) ensureAdvancedPrepStepsAreCreated(ctx context.Context) error {
	_, span := w.tracer.StartSpan(ctx)
	defer span.End()

	logger := w.logger.Clone()

	logger.Info("ensureAdvancedPrepStepsAreCreated invoked")

	return nil
}
