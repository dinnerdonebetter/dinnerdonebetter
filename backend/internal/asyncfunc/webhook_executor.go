package asyncfunc

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/database/postgres"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

func SendWebhook(
	ctx context.Context,
	cfg *config.InstanceConfig,
	logger logging.Logger,
	webhookExecutionRequest *types.WebhookExecutionRequest,
) error {
	tracerProvider, err := cfg.Observability.Tracing.ProvideTracerProvider(ctx, logger)
	if err != nil {
		logger.Error(err, "initializing tracer")
	}
	otel.SetTracerProvider(tracerProvider)

	tracer := tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("search_indexer_cloud_function"))
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	dbConnectionContext, cancel := context.WithTimeout(ctx, 15*time.Second)
	dataManager, err := postgres.ProvideDatabaseClient(dbConnectionContext, logger, tracerProvider, &cfg.Database)
	if err != nil {
		cancel()
		return observability.PrepareAndLogError(err, logger, span, "establishing database connection")
	}

	cancel()

	household, err := dataManager.GetHousehold(ctx, webhookExecutionRequest.HouseholdID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting household")
	}
	_ = household

	webhook, err := dataManager.GetWebhook(ctx, webhookExecutionRequest.WebhookID, webhookExecutionRequest.HouseholdID)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "getting webhook")
		return nil
	}

	var payloadBody []byte
	switch webhook.ContentType {
	case "application/json":
		payloadBody, err = json.Marshal(webhookExecutionRequest.Payload)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "marshaling webhook payload")
		}
	case "application/xml":
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

	decryptedKey, err := hex.DecodeString(household.WebhookEncryptionKey)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "decoding webhook encryption key")
	}

	digest := hmac.New(sha256.New, decryptedKey)
	digest.Write(payloadBody)
	req.Header.Set("X-Dinner-Done-Better-Signature", hex.EncodeToString(digest.Sum(nil)))

	res, err := otelhttp.DefaultClient.Do(req)
	if err != nil {
		observability.AcknowledgeError(err, logger, span, "executing webhook request")
		return nil
	}
	defer func() {
		if err = res.Body.Close(); err != nil {
			observability.AcknowledgeError(err, logger, span, "closing response body")
		}
	}()

	logger = logger.WithResponse(res)
	tracing.AttachResponseToSpan(span, res)

	if res.StatusCode < http.StatusOK || res.StatusCode < http.StatusMultipleChoices {
		observability.AcknowledgeError(err, logger, span, "invalid response type")
		return nil
	}

	return nil
}
