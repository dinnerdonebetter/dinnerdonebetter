package adapters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/payments"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	paymentscfg "github.com/dinnerdonebetter/backend/internal/services/payments/config"

	"github.com/google/wire"
)

var (
	PaymentsAdapterProviders = wire.NewSet(
		ProvidePaymentProcessorRegistry,
		wire.Bind(new(payments.PaymentProcessorRegistry), new(*payments.MapProcessorRegistry)),
	)
)

// ProvidePaymentProcessorRegistry creates a registry with stripe and revenuecat processors.
// Uses stub when a provider is not configured.
func ProvidePaymentProcessorRegistry(
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
	cfg paymentscfg.Config,
) *payments.MapProcessorRegistry {
	processors := make(map[string]payments.PaymentProcessor)

	noopStubPaymentProcessor := NewStubPaymentProcessor(logger)

	// Stripe: use real adapter when configured, else stub
	if cfg.Stripe != nil && cfg.Stripe.APIKey != "" {
		processors["stripe"] = NewStripePaymentProcessor(logger, tracerProvider, &StripeConfig{
			APIKey:        cfg.Stripe.APIKey,
			WebhookSecret: cfg.Stripe.WebhookSecret,
		})
	} else {
		processors["stripe"] = noopStubPaymentProcessor
	}

	// RevenueCat: use real adapter when configured, else stub
	if cfg.RevenueCat != nil && cfg.RevenueCat.APIKey != "" {
		processors["revenuecat"] = NewRevenueCatPaymentProcessor(logger, tracerProvider, &RevenueCatConfig{
			APIKey:            cfg.RevenueCat.APIKey,
			WebhookAuthHeader: cfg.RevenueCat.WebhookAuthHeader,
		})
	} else {
		processors["revenuecat"] = noopStubPaymentProcessor
	}

	return payments.NewMapProcessorRegistry(processors)
}
