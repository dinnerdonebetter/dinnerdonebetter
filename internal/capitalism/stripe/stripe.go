package stripe

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/capitalism"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
	"github.com/stripe/stripe-go/v72/webhook"
)

const (
	implementationName = "stripe_payment_manager"

	webhookHeaderName                    = "Stripe-Signature"
	webhookEventTypeCheckoutCompleted    = "checkout.session.completed"
	webhookEventTypeInvoicePaid          = "invoice.paid"
	webhookEventTypeInvoicePaymentFailed = "invoice.payment_failed"
)

type (
	// WebhookSecret is a string alias for dependency injection.
	WebhookSecret string
	// APIKey is a string alias for dependency injection.
	APIKey string

	stripePaymentManager struct {
		logger        logging.Logger
		tracer        tracing.Tracer
		client        *client.API
		successURL    string
		cancelURL     string
		webhookSecret string
	}
)

// ProvideStripePaymentManager builds a Stripe-backed stripePaymentManager.
func ProvideStripePaymentManager(logger logging.Logger, tracerProvider tracing.TracerProvider, cfg *Config) capitalism.PaymentManager {
	if cfg == nil {
		return &capitalism.NoopPaymentManager{}
	}

	spm := &stripePaymentManager{
		client:        client.New(cfg.APIKey, nil),
		webhookSecret: cfg.WebhookSecret,
		successURL:    cfg.SuccessURL,
		cancelURL:     cfg.CancelURL,
		logger:        logging.EnsureLogger(logger),
		tracer:        tracing.NewTracer(tracerProvider.Tracer(implementationName)),
	}

	return spm
}

func buildCustomerName(household *types.Household) string {
	return fmt.Sprintf("%s (%s)", household.Name, household.ID)
}

func buildGetCustomerParams(a *types.Household) *stripe.CustomerParams {
	p := &stripe.CustomerParams{
		Name:    stripe.String(buildCustomerName(a)),
		Email:   stripe.String("UNKNOWN"),
		Phone:   stripe.String(a.ContactPhone),
		Address: &stripe.AddressParams{},
	}
	p.AddMetadata(keys.HouseholdIDKey, a.ID)

	return p
}

func (s *stripePaymentManager) buildCheckoutSessionParams(subscriptionPlanID string) *stripe.CheckoutSessionParams {
	return &stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String(s.successURL),
		CancelURL:          stripe.String(s.cancelURL),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				{
					Plan:     stripe.String(subscriptionPlanID),
					Quantity: stripe.Int64(1), // For metered billing, do not pass quantity
				},
			},
		},
	}
}

func (s *stripePaymentManager) CreateCheckoutSession(ctx context.Context, subscriptionPlanID string) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	params := &stripe.CheckoutSessionParams{
		SuccessURL:         stripe.String(s.successURL),
		CancelURL:          stripe.String(s.cancelURL),
		Mode:               stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			Items: []*stripe.CheckoutSessionSubscriptionDataItemsParams{
				{
					Plan:     stripe.String(subscriptionPlanID),
					Quantity: stripe.Int64(1), // For metered billing, do not pass quantity
				},
			},
		},
	}

	sess, err := s.client.CheckoutSessions.New(params)
	if err != nil {
		return "", observability.PrepareError(err, span, "creating checkout session")
	}

	return sess.ID, nil
}

func (s *stripePaymentManager) HandleSubscriptionEventWebhook(req *http.Request) error {
	_, span := s.tracer.StartSpan(req.Context())
	defer span.End()

	b, err := io.ReadAll(req.Body)
	if err != nil {
		return observability.PrepareError(err, span, "parsing received webhook content")
	}

	event, err := webhook.ConstructEvent(b, req.Header.Get(webhookHeaderName), s.webhookSecret)
	if err != nil {
		return observability.PrepareError(err, span, "constructing webhook event")
	}

	switch event.Type {
	case webhookEventTypeCheckoutCompleted:
		// Payment is successful, and the subscription is created.
		// You should provision the subscription and save the customer ID to your database.
	case webhookEventTypeInvoicePaid:
		// Continue to provision the subscription as payments continue to be made.
		// Store the status in your database and check when a user accesses your service.
		// This approach helps you avoid hitting rate limits.
	case webhookEventTypeInvoicePaymentFailed:
		// The payment failed, or the customer does not have a valid payment method.
		// The subscription becomes past_due. Notify your customer and send them to the
		// customer portal to update their payment information.
	default:
		// unhandled event type
	}

	return nil
}

func (s *stripePaymentManager) CreateCustomerID(ctx context.Context, household *types.Household) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	params := buildGetCustomerParams(household)

	c, err := s.client.Customers.New(params)
	if err != nil {
		return "", observability.PrepareError(err, span, "creating customer")
	}

	return c.ID, nil
}

var errSubscriptionNotFound = errors.New("subscription not found")

func (s *stripePaymentManager) findSubscriptionID(ctx context.Context, customerID, planID string) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	cus, err := s.client.Customers.Get(customerID, nil)
	if err != nil {
		return "", observability.PrepareError(err, span, "fetching customer")
	}

	for _, sub := range cus.Subscriptions.Data {
		if sub.Plan.ID == planID {
			return sub.ID, nil
		}
	}

	return "", errSubscriptionNotFound
}

func (s *stripePaymentManager) SubscribeToPlan(ctx context.Context, customerID, paymentMethodToken, planID string) (string, error) {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	// check first if the plan is already implemented.
	subscriptionID, err := s.findSubscriptionID(ctx, customerID, planID)
	if err != nil && !errors.Is(err, errSubscriptionNotFound) {
		return "", observability.PrepareError(err, span, "checking subscription status")
	} else if subscriptionID != "" {
		return subscriptionID, nil
	}

	params := &stripe.SubscriptionParams{
		Customer:      stripe.String(customerID),
		Plan:          stripe.String(planID),
		DefaultSource: stripe.String(paymentMethodToken),
	}

	subscription, err := s.client.Subscriptions.New(params)
	if err != nil {
		return "", observability.PrepareError(err, span, "subscribing to plan")
	}

	return subscription.ID, nil
}

func buildCancellationParams() *stripe.SubscriptionCancelParams {
	return &stripe.SubscriptionCancelParams{
		InvoiceNow: stripe.Bool(true),
		Prorate:    stripe.Bool(true),
	}
}

func (s *stripePaymentManager) UnsubscribeFromPlan(ctx context.Context, subscriptionID string) error {
	_, span := s.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachStringToSpan(span, "subscription.id", subscriptionID)

	params := buildCancellationParams()

	if _, err := s.client.Subscriptions.Cancel(subscriptionID, params); err != nil {
		return observability.PrepareError(err, span, "unsubscribing from plan")
	}

	return nil
}
