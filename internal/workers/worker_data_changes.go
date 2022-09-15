package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/customerdata"
	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// DataChangesWorker observes data changes in the database.
type DataChangesWorker struct {
	logger                logging.Logger
	tracer                tracing.Tracer
	encoder               encoding.ClientEncoder
	emailSender           email.Emailer
	customerDataCollector customerdata.Collector
}

// ProvideDataChangesWorker provides a DataChangesWorker.
func ProvideDataChangesWorker(
	logger logging.Logger,
	emailSender email.Emailer,
	customerDataCollector customerdata.Collector,
	tracerProvider tracing.TracerProvider,
) *DataChangesWorker {
	name := "data_changes"

	return &DataChangesWorker{
		logger:                logging.EnsureLogger(logger).WithName(name),
		tracer:                tracing.NewTracer(tracerProvider.Tracer(name)),
		encoder:               encoding.ProvideClientEncoder(logger, tracerProvider, encoding.ContentTypeJSON),
		emailSender:           emailSender,
		customerDataCollector: customerDataCollector,
	}
}

// HandleMessage handles a pending write.
func (w *DataChangesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.DataChangeMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, span, "unmarshalling message")
	}

	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	w.logger.WithValue("message", msg).Info("message received")

	return nil
}
