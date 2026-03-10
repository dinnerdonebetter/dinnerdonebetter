package secrets

import (
	"context"
	"os"
)

// SecretSource provides access to secrets.
type SecretSource interface {
	GetSecret(ctx context.Context, name string) (string, error)
	Close() error
}

type envSecretSource struct{}

// NewEnvSecretSource returns a SecretSource that reads from environment variables.
func NewEnvSecretSource() SecretSource {
	return &envSecretSource{}
}

func (e *envSecretSource) GetSecret(ctx context.Context, name string) (string, error) {
	return os.Getenv(name), nil
}

func (e *envSecretSource) Close() error {
	return nil
}
