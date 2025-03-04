package main

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func handleWebhookExecutionRequest(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	webhookExecutionRequest *types.WebhookExecutionRequest,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	household, err := dataManager.GetHousehold(ctx, webhookExecutionRequest.HouseholdID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting household")
	}

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

	res, err := tracing.BuildTracedHTTPClient().Do(req)
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
