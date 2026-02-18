package adapters

import (
	"context"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/domain/payments"
)

// StubPaymentProcessor is a no-op implementation of PaymentProcessor for local development and testing.
// CreateCheckoutSession returns a placeholder URL without calling any external provider.
type StubPaymentProcessor struct{}

var _ payments.PaymentProcessor = (*StubPaymentProcessor)(nil)

// NewStubPaymentProcessor returns a new stub payment processor.
func NewStubPaymentProcessor() *StubPaymentProcessor {
	return &StubPaymentProcessor{}
}

// CreateCustomer returns a placeholder external customer ID.
func (s *StubPaymentProcessor) CreateCustomer(_ context.Context, accountID, _, _ string) (string, error) {
	return fmt.Sprintf("cus_stub_%s", accountID), nil
}

// CreateCheckoutSession returns a placeholder URL for redirecting the user.
func (s *StubPaymentProcessor) CreateCheckoutSession(_ context.Context, params payments.CheckoutSessionParams) (sessionURL, sessionID string, err error) {
	return fmt.Sprintf("https://checkout.example.com/session/%s?product=%s", params.AccountID, params.ProductID),
		fmt.Sprintf("cs_stub_%s_%s", params.AccountID, params.ProductID),
		nil
}

// GetSubscriptionStatus returns "active" for any subscription.
func (s *StubPaymentProcessor) GetSubscriptionStatus(_ context.Context, _ string) (string, error) {
	return "active", nil
}

// CancelSubscription is a no-op.
func (s *StubPaymentProcessor) CancelSubscription(_ context.Context, _ string) error {
	return nil
}

// VerifyWebhookSignature always returns true for the stub (accepts all webhooks in dev).
func (s *StubPaymentProcessor) VerifyWebhookSignature(_ context.Context, _ []byte, _, _ string) bool {
	return true
}

// ParseWebhookEvent returns placeholder values. Real adapters would parse provider-specific payloads.
func (s *StubPaymentProcessor) ParseWebhookEvent(_ context.Context, _ []byte) (eventType, accountID, subscriptionID string, err error) {
	return "subscription.updated", "", "", nil
}
