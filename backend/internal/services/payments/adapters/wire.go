package adapters

import (
	"github.com/dinnerdonebetter/backend/internal/domain/payments"

	"github.com/google/wire"
)

var (
	PaymentsAdapterProviders = wire.NewSet(
		NewStubPaymentProcessor,
		wire.Bind(new(payments.PaymentProcessor), new(*StubPaymentProcessor)),
	)
)
