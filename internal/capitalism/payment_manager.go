package capitalism

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/pkg/types"
)

type (
	// SubscriptionPlan describes a plan you pay on a recurring monthly basis for.
	SubscriptionPlan struct {
		ID    string
		Name  string
		Price uint32
	}

	// PaymentManager handles payments via 3rd-party providers.
	PaymentManager interface {
		CreateCustomerID(ctx context.Context, household *types.Household) (string, error)
		HandleSubscriptionEventWebhook(req *http.Request) error
		SubscribeToPlan(ctx context.Context, customerID, paymentMethodToken, planID string) (string, error)
		CreateCheckoutSession(ctx context.Context, subscriptionPlanID string) (string, error)
		UnsubscribeFromPlan(ctx context.Context, subscriptionID string) error
	}
)
