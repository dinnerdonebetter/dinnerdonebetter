package capitalism

import (
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

const (
	stripeSignatureHeaderKey = "Stripe-Signature"
)

// IncomingWebhookHandler is our webhook handling route.
func (s *service) IncomingWebhookHandler(res http.ResponseWriter, req *http.Request) {
	_, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	logger := s.logger.WithRequest(req)
	tracing.AttachRequestToSpan(span, req)

	statusToWrite := http.StatusUnprocessableEntity

	if err := s.paymentManager.HandleEventWebhook(req); err != nil {
		logger.Error(err, "handling subscription event webhook")
		statusToWrite = http.StatusInternalServerError
	}

	res.WriteHeader(statusToWrite)
}
