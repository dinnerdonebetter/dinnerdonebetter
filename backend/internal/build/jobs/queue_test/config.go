package queuetest

import (
	"github.com/dinnerdonebetter/backend/internal/config"
	msgconfig "github.com/dinnerdonebetter/backend/internal/platform/messagequeue/config"
	queuetest "github.com/dinnerdonebetter/backend/internal/services/internalops/workers/queue_test"

	"github.com/google/wire"
)

var (
	// ConfigProviders represents this package's offering to the dependency injector.
	ConfigProviders = wire.NewSet(
		wire.FieldsOf(
			new(*config.QueueTestJobConfig),
			"Queues",
			"Observability",
			"Database",
		),
		ProvideJobParams,
		ProvideEventsConfig,
	)
)

// ProvideJobParams builds JobParams from the config.
func ProvideJobParams(cfg *config.QueueTestJobConfig) *queuetest.JobParams {
	return &queuetest.JobParams{
		Queues: cfg.Queues,
	}
}

// ProvideEventsConfig returns the message queue config for the publisher.
func ProvideEventsConfig(cfg *config.QueueTestJobConfig) *msgconfig.Config {
	return &cfg.Events
}
