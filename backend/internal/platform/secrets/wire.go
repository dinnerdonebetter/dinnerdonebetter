package secrets

import (
	"fmt"

	"github.com/google/wire"
)

var (
	// Providers provides secret source construction for dependency injection.
	Providers = wire.NewSet(
		ProvideSecretSourceFromConfig,
	)
)

// ProvideSecretSourceFromConfig provides a SecretSource from config.
func ProvideSecretSourceFromConfig(cfg *Config) (SecretSource, error) {
	if cfg == nil {
		return NewEnvSecretSource(), nil
	}
	source, err := cfg.ProvideSecretSource()
	if err != nil {
		return nil, fmt.Errorf("provide secret source: %w", err)
	}
	return source, nil
}
