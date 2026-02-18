package payments

import "context"

// CheckoutSessionRequestInput represents user input for creating a checkout session.
type CheckoutSessionRequestInput struct {
	ProductID   string `json:"productId"`
	AccountID   string `json:"accountId"`
	SuccessURL  string `json:"successUrl"`
	CancelURL   string `json:"cancelUrl"`
	IsRecurring bool   `json:"isRecurring"`
}

// CheckoutSessionParams contains parameters for creating a checkout session.
type CheckoutSessionParams struct {
	ProductID   string
	AccountID   string
	SuccessURL  string
	CancelURL   string
	IsRecurring bool
}

// PaymentProcessor defines the interface for payment provider integrations.
// Implementations (e.g., Stripe, Paddle) live in adapters and are wired via dependency injection.
type PaymentProcessor interface {
	CreateCustomer(ctx context.Context, accountID, email, name string) (externalCustomerID string, err error)
	CreateCheckoutSession(ctx context.Context, params CheckoutSessionParams) (sessionURL, sessionID string, err error)
	GetSubscriptionStatus(ctx context.Context, externalSubscriptionID string) (status string, err error)
	CancelSubscription(ctx context.Context, externalSubscriptionID string) error
	VerifyWebhookSignature(ctx context.Context, payload []byte, signature string, accountID string) bool
	ParseWebhookEvent(ctx context.Context, payload []byte) (eventType string, accountID string, subscriptionID string, err error)
}
