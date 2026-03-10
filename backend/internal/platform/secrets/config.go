package secrets

import (
	"fmt"
	"strings"
)

const (
	ProviderEnv  = "env"
	ProviderNoop = "noop"
)

// Config configures secret source selection.
type Config struct {
	Provider string `env:"PROVIDER" json:"provider"`
}

// ProvideSecretSource returns a SecretSource from config.
func (cfg *Config) ProvideSecretSource() (SecretSource, error) {
	if cfg == nil {
		return NewEnvSecretSource(), nil
	}

	switch strings.TrimSpace(strings.ToLower(cfg.Provider)) {
	case "", ProviderEnv:
		return NewEnvSecretSource(), nil
	case ProviderNoop:
		return NewNoopSecretSource(), nil
	default:
		return nil, fmt.Errorf("unknown secret source provider: %q", cfg.Provider)
	}
}
