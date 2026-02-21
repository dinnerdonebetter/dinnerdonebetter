package webhookexecutionrequesthandler

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/messagequeue"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/metrics"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "webhook_execution_request_handler"
)

var (
	errRequiredDataIsNil = errors.New("required data is nil")
)

type WebhookExecutionRequestHandler struct {
	logger                 logging.Logger
	tracer                 tracing.Tracer
	consumerProvider       messagequeue.ConsumerProvider
	identityRepo           identity.Repository
	webhookRepo            webhooks.Repository
	decoder                encoding.ServerEncoderDecoder
	executionTimeHistogram metrics.Float64Histogram
	queuesConfig           msgconfig.QueuesConfig
}

func NewWebhookExecutionRequestHandler(
	_ context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg *config.WebhookExecutionRequestHandlerConfig,
	consumerProvider messagequeue.ConsumerProvider,
	identityRepo identity.Repository,
	webhookRepo webhooks.Repository,
	decoder encoding.ServerEncoderDecoder,
	metricsProvider metrics.Provider,
) (*WebhookExecutionRequestHandler, error) {
	executionTimeHistogram, err := metricsProvider.NewFloat64Histogram("webhook_requests_execution_time")
	if err != nil {
		return nil, fmt.Errorf("setting up webhook execution requests execution time histogram: %w", err)
	}

	return &WebhookExecutionRequestHandler{
		tracer:                 tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:                 logging.EnsureLogger(logger).WithName(o11yName),
		consumerProvider:       consumerProvider,
		identityRepo:           identityRepo,
		webhookRepo:            webhookRepo,
		decoder:                decoder,
		executionTimeHistogram: executionTimeHistogram,
		queuesConfig:           cfg.Queues,
	}, nil
}

func (h *WebhookExecutionRequestHandler) HandleMessage(ctx context.Context, rawMsg []byte) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	start := time.Now()

	var webhookExecutionRequest webhooks.WebhookExecutionRequest
	if err := h.decoder.DecodeBytes(ctx, rawMsg, &webhookExecutionRequest); err != nil {
		return fmt.Errorf("decoding JSON body: %w", err)
	}

	if err := h.handleWebhookExecutionRequest(ctx, &webhookExecutionRequest); err != nil {
		return fmt.Errorf("handling webhook execution request: %w", err)
	}

	h.executionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

	return nil
}

func (h *WebhookExecutionRequestHandler) handleWebhookExecutionRequest(
	ctx context.Context,
	webhookExecutionRequest *webhooks.WebhookExecutionRequest,
) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	if webhookExecutionRequest == nil {
		return errRequiredDataIsNil
	}

	logger := h.logger.WithValue(keys.RequestIDKey, webhookExecutionRequest.RequestID)

	account, err := h.identityRepo.GetAccount(ctx, webhookExecutionRequest.AccountID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting account")
	}

	webhook, err := h.webhookRepo.GetWebhook(ctx, webhookExecutionRequest.WebhookID, webhookExecutionRequest.AccountID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting webhook")
		return nil
	}

	var payloadBody []byte
	switch webhook.ContentType {
	case encoding.ContentTypeToString(encoding.ContentTypeJSON):
		payloadBody, err = json.Marshal(webhookExecutionRequest.Payload)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "marshaling webhook payload")
		}
	case encoding.ContentTypeToString(encoding.ContentTypeXML):
		payloadBody, err = xml.Marshal(webhookExecutionRequest.Payload)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "marshaling webhook payload")
		}
	}

	req, err := http.NewRequestWithContext(ctx, webhook.Method, webhook.URL, bytes.NewReader(payloadBody))
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating webhook request")
	}

	req.Header.Set("Content-Type", webhook.ContentType)

	decryptedKey, err := hex.DecodeString(account.WebhookEncryptionKey)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "decoding webhook encryption key")
	}

	digest := hmac.New(sha256.New, decryptedKey)
	digest.Write(payloadBody)
	req.Header.Set("X-Dinner-Done-Better-Signature", hex.EncodeToString(digest.Sum(nil)))

	res, err := tracing.BuildTracedHTTPClient().Do(req) //nolint:gosec // G704: webhook URL is admin-configured; webhooks intentionally deliver to external URLs
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "executing webhook request")
		return nil
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			logger.Error("closing response body", err)
		}
	}()

	logger = logger.WithResponse(res)
	tracing.AttachResponseToSpan(span, res)

	if res.StatusCode < 200 || res.StatusCode > 299 {
		observability.AcknowledgeError(err, logger, span, "invalid response type")
		return nil
	}

	return nil
}

func (h *WebhookExecutionRequestHandler) ConsumeMessages(
	ctx context.Context,
	stopChan chan bool,
	errorsChan chan error,
) error {
	ctx, span := h.tracer.StartSpan(ctx)
	defer span.End()

	consumer, err := h.consumerProvider.ProvideConsumer(
		ctx,
		h.queuesConfig.WebhookExecutionRequestsTopicName,
		h.HandleMessage,
	)
	if err != nil {
		return observability.PrepareAndLogError(err, h.logger, span, "configuring webhook execution requests consumer")
	}

	go consumer.Consume(stopChan, errorsChan)

	go func() {
		for e := range errorsChan {
			h.logger.Error("consuming message", e)
		}
	}()

	return nil
}
