package main

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

func handleEmailRequest(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	emailer email.Emailer,
	analyticsEventReporter analytics.EventReporter,
	mail *email.OutboundEmailMessage,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	if err := emailer.SendEmail(ctx, mail); err != nil {
		observability.AcknowledgeError(err, logger, span, "sending email")
	}

	if err := analyticsEventReporter.EventOccurred(ctx, email.SentEventType, mail.UserID, map[string]any{
		"toAddress":   mail.ToAddress,
		"toName":      mail.ToName,
		"fromAddress": mail.FromAddress,
		"fromName":    mail.FromName,
		"subject":     mail.Subject,
	}); err != nil {
		observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
	}

	return nil
}
