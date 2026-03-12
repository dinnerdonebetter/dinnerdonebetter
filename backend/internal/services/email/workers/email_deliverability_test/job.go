package emaildeliverabilitytest

import (
	"context"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/branding"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	serviceName = "email_deliverability_test"
)

// JobParams holds configuration for the email deliverability test job.
type JobParams struct {
	RecipientEmailAddress string
	ServiceEnvironment    string
}

// Job sends a test email to verify deliverability.
type Job struct {
	emailer            email.Emailer
	logger             logging.Logger
	tracer             tracing.Tracer
	recipientEmail     string
	serviceEnvironment string
}

// NewJob creates a new email deliverability test job.
func NewJob(
	emailer email.Emailer,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	params *JobParams,
) (*Job, error) {
	recipientEmail := params.RecipientEmailAddress
	serviceEnvironment := params.ServiceEnvironment
	if recipientEmail == "" {
		return nil, fmt.Errorf("recipient email is required")
	}
	if serviceEnvironment == "" {
		serviceEnvironment = "prod"
	}
	return &Job{
		emailer:            emailer,
		logger:             logging.EnsureLogger(logger).WithName(serviceName),
		tracer:             tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(serviceName)),
		recipientEmail:     recipientEmail,
		serviceEnvironment: serviceEnvironment,
	}, nil
}

// Do sends the deliverability test email.
func (j *Job) Do(ctx context.Context) error {
	ctx, span := j.tracer.StartSpan(ctx)
	defer span.End()

	envCfg := email.GetConfigForEnvironment(j.serviceEnvironment)
	if envCfg == nil {
		return fmt.Errorf("no email config for environment %q", j.serviceEnvironment)
	}

	fromAddress := envCfg.PasswordResetCreationEmailAddress()
	sentAt := time.Now().UTC().Format(time.RFC3339)

	msg := &email.OutboundEmailMessage{
		ToAddress:   j.recipientEmail,
		FromAddress: fromAddress,
		FromName:    branding.CompanyName,
		Subject:     fmt.Sprintf("%s – Email Deliverability Test", branding.CompanyName),
		HTMLContent: fmt.Sprintf("<p>Test sent at %s</p>", sentAt),
	}

	if err := j.emailer.SendEmail(ctx, msg); err != nil {
		j.logger.Error("sending deliverability test email", err)
		return fmt.Errorf("sending email: %w", err)
	}

	j.logger.Info("deliverability test email sent successfully")
	return nil
}
