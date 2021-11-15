package workers

import (
	"context"

	"github.com/prixfixeco/api_server/internal/encoding"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// DataChangesWorker observes data changes in the database.
type DataChangesWorker struct {
	logger  logging.Logger
	tracer  tracing.Tracer
	encoder encoding.ClientEncoder
}

// ProvideDataChangesWorker provides a DataChangesWorker.
func ProvideDataChangesWorker(logger logging.Logger) *DataChangesWorker {
	name := "post_writes"

	return &DataChangesWorker{
		logger:  logging.EnsureLogger(logger).WithName(name),
		tracer:  tracing.NewTracer(name),
		encoder: encoding.ProvideClientEncoder(logger, encoding.ContentTypeJSON),
	}
}

// HandleMessage handles a pending write.
func (w *DataChangesWorker) HandleMessage(ctx context.Context, message []byte) error {
	ctx, span := w.tracer.StartSpan(ctx)
	defer span.End()

	var msg *types.DataChangeMessage

	if err := w.encoder.Unmarshal(ctx, message, &msg); err != nil {
		return observability.PrepareError(err, w.logger, span, "unmarshalling message")
	}

	tracing.AttachUserIDToSpan(span, msg.AttributableToUserID)
	w.logger.WithValue("message", message).Info("message received")

	return nil
}
