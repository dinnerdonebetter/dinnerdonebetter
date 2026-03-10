package emaildeliverabilitytest

import (
	"github.com/dinnerdonebetter/backend/internal/config"
	emaildeliverabilitytest "github.com/dinnerdonebetter/backend/internal/services/email/workers/email_deliverability_test"

	"github.com/google/wire"
)

var (
	// ConfigProviders represents this package's offering to the dependency injector.
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.EmailDeliverabilityTestConfig),
			"Observability",
			"Email",
		),
		ProvideJobParams,
	)
)

// ProvideJobParams builds JobParams from the config.
func ProvideJobParams(cfg *config.EmailDeliverabilityTestConfig) *emaildeliverabilitytest.JobParams {
	return &emaildeliverabilitytest.JobParams{
		RecipientEmailAddress: cfg.RecipientEmailAddress,
		ServiceEnvironment:    cfg.ServiceEnvironment,
	}
}
