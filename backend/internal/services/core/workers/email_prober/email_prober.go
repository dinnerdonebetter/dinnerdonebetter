package emailprober

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/email"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
)

const (
	serviceName = "email_prober"
)

type Job struct {
	logger  logging.Logger
	tracer  tracing.Tracer
	emailer email.Emailer
}

func NewEmailProber(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	emailer email.Emailer,
) (*Job, error) {
	return &Job{
		emailer: emailer,
		logger:  logging.EnsureLogger(logger).WithName(serviceName),
		tracer:  tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
	}, nil
}

func (j *Job) Do(ctx context.Context) error {
	ctx, span := j.tracer.StartSpan(ctx)
	defer span.End()

	if err := j.emailer.SendEmail(ctx, &email.OutboundEmailMessage{
		ToAddress:   "verygoodsoftwarenotvirus@protonmail.com",
		ToName:      "Jeffrey",
		FromAddress: "email@dinnerdonebetter.dev",
		FromName:    "Testing",
		Subject:     "Testing",
		HTMLContent: "Hi",
	}); err != nil {
		observability.AcknowledgeError(err, j.logger, span, "sending probe email")
	}

	return nil
}
