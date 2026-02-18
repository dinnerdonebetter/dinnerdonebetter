package manager

import (
	"context"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitymanager "github.com/dinnerdonebetter/backend/internal/domain/identity/manager"
	"github.com/dinnerdonebetter/backend/internal/domain/payments"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

const (
	o11yName = "payments_data_manager"
)

var _ PaymentsDataManager = (*paymentsManager)(nil)

type paymentsManager struct {
	tracer      tracing.Tracer
	logger      logging.Logger
	repo        payments.Repository
	processor   payments.PaymentProcessor
	identityMgr identitymanager.IdentityDataManager
}

// NewPaymentsDataManager returns a new PaymentsDataManager.
func NewPaymentsDataManager(
	tracerProvider tracing.TracerProvider,
	logger logging.Logger,
	repo payments.Repository,
	processor payments.PaymentProcessor,
	identityMgr identitymanager.IdentityDataManager,
) PaymentsDataManager {
	return &paymentsManager{
		tracer:      tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer(o11yName)),
		logger:      logging.EnsureLogger(logger).WithName(o11yName),
		repo:        repo,
		processor:   processor,
		identityMgr: identityMgr,
	}
}

func (m *paymentsManager) CreateProduct(ctx context.Context, input *payments.ProductCreationRequestInput) (*payments.Product, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, observability.PrepareError(nil, span, "nil product creation input")
	}
	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating product creation input")
	}

	dbInput := &payments.ProductDatabaseCreationInput{
		ID:                    identifiers.New(),
		Name:                  input.Name,
		Description:           input.Description,
		Kind:                  input.Kind,
		AmountCents:           input.AmountCents,
		Currency:              input.Currency,
		BillingIntervalMonths: input.BillingIntervalMonths,
		ExternalProductID:     input.ExternalProductID,
	}
	return m.repo.CreateProduct(ctx, dbInput)
}

func (m *paymentsManager) GetProduct(ctx context.Context, id string) (*payments.Product, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetProduct(ctx, id)
}

func (m *paymentsManager) GetProducts(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Product], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetProducts(ctx, filter)
}

func (m *paymentsManager) UpdateProduct(ctx context.Context, id string, input *payments.ProductUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	product, err := m.repo.GetProduct(ctx, id)
	if err != nil {
		return observability.PrepareError(err, span, "fetching product")
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Description != nil {
		product.Description = *input.Description
	}
	if input.Kind != nil {
		product.Kind = *input.Kind
	}
	if input.AmountCents != nil {
		product.AmountCents = *input.AmountCents
	}
	if input.Currency != nil {
		product.Currency = *input.Currency
	}
	if input.BillingIntervalMonths != nil {
		product.BillingIntervalMonths = input.BillingIntervalMonths
	}
	if input.ExternalProductID != nil {
		product.ExternalProductID = *input.ExternalProductID
	}

	return m.repo.UpdateProduct(ctx, product)
}

func (m *paymentsManager) ArchiveProduct(ctx context.Context, id string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.ArchiveProduct(ctx, id)
}

func (m *paymentsManager) CreateSubscription(ctx context.Context, input *payments.SubscriptionCreationRequestInput) (*payments.Subscription, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, observability.PrepareError(nil, span, "nil subscription creation input")
	}
	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "validating subscription creation input")
	}

	dbInput := &payments.SubscriptionDatabaseCreationInput{
		ID:                     identifiers.New(),
		BelongsToAccount:       input.BelongsToAccount,
		ProductID:              input.ProductID,
		ExternalSubscriptionID: input.ExternalSubscriptionID,
		Status:                 input.Status,
		CurrentPeriodStart:     input.CurrentPeriodStart,
		CurrentPeriodEnd:       input.CurrentPeriodEnd,
	}
	return m.repo.CreateSubscription(ctx, dbInput)
}

func (m *paymentsManager) GetSubscription(ctx context.Context, id string) (*payments.Subscription, error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetSubscription(ctx, id)
}

func (m *paymentsManager) GetSubscriptionsForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Subscription], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetSubscriptionsForAccount(ctx, accountID, filter)
}

func (m *paymentsManager) UpdateSubscription(ctx context.Context, id string, input *payments.SubscriptionUpdateRequestInput) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	sub, err := m.repo.GetSubscription(ctx, id)
	if err != nil {
		return observability.PrepareError(err, span, "fetching subscription")
	}

	if input.Status != nil {
		sub.Status = *input.Status
	}
	if input.CurrentPeriodStart != nil {
		sub.CurrentPeriodStart = *input.CurrentPeriodStart
	}
	if input.CurrentPeriodEnd != nil {
		sub.CurrentPeriodEnd = *input.CurrentPeriodEnd
	}

	return m.repo.UpdateSubscription(ctx, sub)
}

func (m *paymentsManager) ArchiveSubscription(ctx context.Context, id string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.ArchiveSubscription(ctx, id)
}

func (m *paymentsManager) CancelSubscription(ctx context.Context, id string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	sub, err := m.repo.GetSubscription(ctx, id)
	if err != nil {
		return observability.PrepareError(err, span, "fetching subscription")
	}
	if sub.ExternalSubscriptionID == "" {
		return observability.PrepareError(nil, span, "subscription has no external ID (cannot cancel)")
	}

	if err = m.processor.CancelSubscription(ctx, sub.ExternalSubscriptionID); err != nil {
		return observability.PrepareError(err, span, "cancelling subscription with provider")
	}

	return m.repo.UpdateSubscriptionStatus(ctx, id, payments.SubscriptionStatusCancelled)
}

func (m *paymentsManager) CreateCheckoutSession(ctx context.Context, input *payments.CheckoutSessionRequestInput) (sessionURL, sessionID string, err error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	product, err := m.repo.GetProduct(ctx, input.ProductID)
	if err != nil {
		return "", "", observability.PrepareError(err, span, "fetching product")
	}

	account, err := m.identityMgr.GetAccount(ctx, input.AccountID)
	if err != nil {
		return "", "", observability.PrepareError(err, span, "fetching account")
	}

	var email, name string
	for _, member := range account.Members {
		if member.BelongsToUser != nil && member.BelongsToUser.ID == account.BelongsToUser {
			email = member.BelongsToUser.EmailAddress
			if member.BelongsToUser.FirstName != "" || member.BelongsToUser.LastName != "" {
				name = member.BelongsToUser.FirstName + " " + member.BelongsToUser.LastName
			} else {
				name = account.Name
			}
			break
		}
	}
	if name == "" {
		name = account.Name
	}

	if account.PaymentProcessorCustomerID == "" {
		var extCustomerID string
		extCustomerID, err = m.processor.CreateCustomer(ctx, input.AccountID, email, name)
		if err != nil {
			return "", "", observability.PrepareError(err, span, "creating payment provider customer")
		}
		syncNow := time.Now()

		if err = m.identityMgr.UpdateAccountBillingFields(ctx, input.AccountID, nil, nil, &extCustomerID, &syncNow); err != nil {
			return "", "", observability.PrepareError(err, span, "updating account with customer ID")
		}
	}

	params := payments.CheckoutSessionParams{
		ProductID:   input.ProductID,
		AccountID:   input.AccountID,
		SuccessURL:  input.SuccessURL,
		CancelURL:   input.CancelURL,
		IsRecurring: product.Kind == payments.ProductKindRecurring,
	}

	return m.processor.CreateCheckoutSession(ctx, params)
}

func (m *paymentsManager) GetPurchasesForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Purchase], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetPurchasesForAccount(ctx, accountID, filter)
}

func (m *paymentsManager) GetPaymentTransactionsForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.PaymentTransaction], error) {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	return m.repo.GetPaymentTransactionsForAccount(ctx, accountID, filter)
}

func (m *paymentsManager) ProcessWebhookEvent(ctx context.Context, payload []byte, signature, accountID string) error {
	ctx, span := m.tracer.StartSpan(ctx)
	defer span.End()

	logger := observability.ObserveValues(map[string]any{
		keys.AccountIDKey: accountID,
	}, span, m.logger)

	if !m.processor.VerifyWebhookSignature(ctx, payload, signature, accountID) {
		return observability.PrepareAndLogError(nil, logger, span, "invalid webhook signature")
	}

	eventType, _, subscriptionID, err := m.processor.ParseWebhookEvent(ctx, payload)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "parsing webhook event")
	}

	syncNow := time.Now()

	switch eventType {
	case "subscription.updated", "subscription.created", "customer.subscription.updated":
		if subscriptionID == "" {
			return nil
		}
		var sub *payments.Subscription
		sub, err = m.repo.GetSubscriptionByExternalID(ctx, subscriptionID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "fetching subscription by external ID")
		}

		var status string
		status, err = m.processor.GetSubscriptionStatus(ctx, subscriptionID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "fetching subscription status from provider")
		}

		if err = m.repo.UpdateSubscriptionStatus(ctx, sub.ID, status); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "updating subscription status")
		}

		billingStatus := subscriptionStatusToBillingStatus(status)
		if err = m.identityMgr.UpdateAccountBillingFields(ctx, sub.BelongsToAccount, &billingStatus, &sub.ProductID, nil, &syncNow); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "updating account billing fields")
		}
	case "subscription.deleted", "customer.subscription.deleted":
		if subscriptionID == "" {
			return nil
		}

		var sub *payments.Subscription
		sub, err = m.repo.GetSubscriptionByExternalID(ctx, subscriptionID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "fetching subscription by external ID")
		}

		if err = m.repo.UpdateSubscriptionStatus(ctx, sub.ID, payments.SubscriptionStatusCancelled); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "updating subscription status")
		}

		unpaid := identity.UnpaidAccountBillingStatus
		if err = m.identityMgr.UpdateAccountBillingFields(ctx, sub.BelongsToAccount, &unpaid, nil, nil, &syncNow); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "updating account billing fields")
		}
	default:
		// Unknown event type - no-op
	}

	return nil
}

func subscriptionStatusToBillingStatus(status string) string {
	switch status {
	case payments.SubscriptionStatusActive:
		return identity.PaidAccountBillingStatus
	case payments.SubscriptionStatusTrialing:
		return identity.TrialAccountBillingStatus
	case payments.SubscriptionStatusCancelled, payments.SubscriptionStatusPastDue, payments.SubscriptionStatusIncomplete:
		return identity.UnpaidAccountBillingStatus
	default:
		return identity.UnpaidAccountBillingStatus
	}
}
