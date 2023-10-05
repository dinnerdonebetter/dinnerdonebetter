package capitalism

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/encoding"
	mockpublishers "github.com/dinnerdonebetter/backend/internal/messagequeue/mock"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/random"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/webhook"
)

const (
	fakeSigningSecret = "whsec_abcdefABCDEF0123456789abcdefABCD"
)

func TestValidInstrumentsService_StripeWebhookHandler(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		secret, err := random.GenerateHexEncodedString(ctx, 32)
		require.NoError(t, err)
		require.NotEmpty(t, secret)

		helper := buildTestHelper(t)
		helper.service.encoderDecoder = encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)
		// helper.service.cfg.SigningSecret = secret

		paymentIntent := &stripe.PaymentIntent{
			APIResource:      stripe.APIResource{},
			Amount:           0,
			AmountCapturable: 0,
			AmountDetails:    nil,
			AmountReceived:   0,
			Customer:         nil,
			ID:               "",
			Invoice:          nil,
			Metadata:         nil,
			PaymentMethod:    nil,
			ReceiptEmail:     "",
			Status:           "",
		}

		rawMessage, err := json.Marshal(paymentIntent)
		require.NoError(t, err)
		require.NotNil(t, rawMessage)

		exampleInput := &stripe.Event{
			APIResource: stripe.APIResource{},
			Account:     "",
			APIVersion:  "2023-08-16",
			Created:     0,
			Data: &stripe.EventData{
				Object:             nil,
				PreviousAttributes: nil,
				Raw:                json.RawMessage(rawMessage),
			},
			ID:              "",
			Livemode:        false,
			Object:          "",
			PendingWebhooks: 0,
			Request:         nil,
			Type:            stripe.EventTypePaymentIntentSucceeded,
		}
		jsonBytes := helper.service.encoderDecoder.MustEncode(helper.ctx, exampleInput)

		now := time.Now()
		signedPayload := webhook.GenerateTestSignedPayload(&webhook.UnsignedPayload{
			Payload:   jsonBytes,
			Secret:    secret,
			Timestamp: now,
		})

		event, err := webhook.ConstructEvent(signedPayload.Payload, signedPayload.Header, signedPayload.Secret)
		eventPayload := helper.service.encoderDecoder.MustEncode(helper.ctx, event)

		helper.req, err = http.NewRequestWithContext(helper.ctx, http.MethodPost, "https://whatever.whocares.gov", bytes.NewReader(eventPayload))
		require.NoError(t, err)
		require.NotNil(t, helper.req)
		helper.req.Header.Set(stripeSignatureHeaderKey, signedPayload.Header)

		dataChangesPublisher := &mockpublishers.Publisher{}
		helper.service.dataChangesPublisher = dataChangesPublisher

		helper.service.IncomingWebhookHandler(helper.res, helper.req)

		assert.Equal(t, http.StatusOK, helper.res.Code)

		mock.AssertExpectationsForObjects(t, dataChangesPublisher)
	})
}
