package workers

import (
	"context"
	"fmt"

	"github.com/prixfixeco/api_server/internal/messagequeue"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/database"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// ArchivesWorker archives data from the pending archives topic to the database.
type ArchivesWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	dataChangesPublisher  messagequeue.Publisher
	dataManager           database.DataManager
	customerDataCollector customerdata.Collector
}

// ProvideArchivesWorker provides an ArchivesWorker.
func ProvideArchivesWorker(
	_ context.Context,
	logger logging.Logger,
	dataManager database.DataManager,
	postArchivesPublisher messagequeue.Publisher,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) (*ArchivesWorker, error) {
	const name = "pre_archives"

	w := &ArchivesWorker{
		logger:                logging.EnsureLogger(logger).WithName(name).WithValue("topic", name),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		dataChangesPublisher:  postArchivesPublisher,
		dataManager:           dataManager,
		customerDataCollector: customerDataCollector,
	}

	return w, nil
}

func (w *ArchivesWorker) determineArchiveMessageHandler(msg *types.PreArchiveMessage) func(context.Context, *types.PreArchiveMessage) error {
	funcMap := map[string]func(context.Context, *types.PreArchiveMessage) error{
		string(types.MealDataType):                w.archiveMeal,
		string(types.MealPlanDataType):            w.archiveMealPlan,
		string(types.MealPlanOptionDataType):      w.archiveMealPlanOption,
		string(types.MealPlanOptionVoteDataType):  w.archiveMealPlanOptionVote,
		string(types.UserMembershipDataType):      func(context.Context, *types.PreArchiveMessage) error { return nil },
		string(types.HouseholdInvitationDataType): func(context.Context, *types.PreArchiveMessage) error { return nil },
	}

	f, ok := funcMap[string(msg.DataType)]
	if ok {
		return f
	}

	return nil
}

// HandleMessage handles a pending archive.
func (w *ArchivesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.PreArchiveMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}
	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	logger := w.logger.WithValue("data_type", msg.DataType)

	logger.Debug("message read")

	f := w.determineArchiveMessageHandler(msg)

	if f == nil {
		return fmt.Errorf("no handler assigned to message type %q", msg.DataType)
	}

	return f(ctx, msg)
}
