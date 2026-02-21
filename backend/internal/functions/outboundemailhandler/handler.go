package outboundemailhandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/analytics"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "outbound_email_handler"
)

var (
	errRequiredDataIsNil = errors.New("required data is nil")
)

type OutboundEmailHandler struct {
	logger                 logging.Logger
	tracer                 tracing.Tracer
	consumerProvider       messagequeue.ConsumerProvider
	emailer                email.Emailer
	analyticsEventReporter analytics.EventReporter
	executionTimeHistogram metrics.Float64Histogram
	queuesConfig           msgconfig.QueuesConfig
}

func NewOutboundEmailHandler(
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *config.OutboundEmailHandlerConfig,
	consumerProvider messagequeue.ConsumerProvider,
	analyticsEventReporter analytics.EventReporter,
	emailer email.Emailer,
	metricsProvider metrics.Provider,
) (*OutboundEmailHandler, error) {
	executionTimeHistogram, err := metricsProvider.NewFloat64Histogram("outbound_emails_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up outbound emails execution time histogram: %w", err)
	}

	return &OutboundEmailHandler{
		tracer:                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                 logging.EnsureLogger(logger).WithName(o11yName),
		consumerProvider:       consumerProvider,
		emailer:                emailer,
		analyticsEventReporter: analyticsEventReporter,
		executionTimeHistogram: executionTimeHistogram,
		queuesConfig:           cfg.Queues,
	}, nil
}

func (h *OutboundEmailHandler) HandleMessage(ctx context.Context, rawMsg []byte) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	start := time.Now()

	var emailMessage email.OutboundEmailMessage
	if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&emailMessage); err != nil {
		return fmt.Errorf("decoding JSON body: %w", err)
	}

	if err := h.handleEmailRequest(ctx, &emailMessage); err != nil {
		return fmt.Errorf("handling outbound email request: %w", err)
	}

	h.executionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

	return nil
}

func (h *OutboundEmailHandler) handleEmailRequest(
	ctx context.Context,
	mail *email.OutboundEmailMessage,
) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	if mail == nil {
		return errRequiredDataIsNil
	}

	if err := h.emailer.SendEmail(ctx, mail); err != nil {
		observability.AcknowledgeError(err, h.logger, span, "sending email")
	}

	if err := h.analyticsEventReporter.EventOccurred(ctx, email.SentEventType, mail.UserID, map[string]any{
		"toAddress":   mail.ToAddress,
		"toName":      mail.ToName,
		"fromAddress": mail.FromAddress,
		"fromName":    mail.FromName,
		"subject":     mail.Subject,
	}); err != nil {
		observability.AcknowledgeError(err, h.logger, span, "notifying customer data platform")
	}

	return nil
}

func (h *OutboundEmailHandler) ConsumeMessages(
	ctx context.Context,
	stopChan chan bool,
	errorsChan chan error,
) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	consumer, err := h.consumerProvider.ProvideConsumer(
		ctx,
		h.queuesConfig.OutboundEmailsTopicName,
		h.HandleMessage,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, h.logger, span, "configuring outbound emails consumer")
	}

	go consumer.Consume(stopChan, errorsChan)

	go func() {
		for e := range errorsChan {
			h.logger.Error("consuming message", e)
		}
	}()

	return nil
}
