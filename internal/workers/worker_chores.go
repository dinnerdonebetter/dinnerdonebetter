package workers

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// ChoresWorker performs chores.
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
	name := "chores"

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

func (w *ChoresWorker) determineChoreHandler(msg *types.ChoreMessage) func(context.Context, *types.ChoreMessage) error {
	funcMap := map[string]func(context.Context, *types.ChoreMessage) error{
		string(types.FinalizeMealPlansWithExpiredVotingPeriodsChoreType): w.finalizeExpiredMealPlans,
	}

	f, ok := funcMap[string(msg.ChoreType)]
	if ok {
		return f
	}

	return nil
}

// HandleMessage handles a pending write.
func (w *ChoresWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.ChoreMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	logger := w.logger.WithValue("chore_type", msg.ChoreType)

	logger.Debug("message read")

	f := w.determineChoreHandler(msg)

	if f == nil {
		return fmt.Errorf("no handler assigned to chore type %q", msg.ChoreType)
	}

	return f(ctx, msg)
}
