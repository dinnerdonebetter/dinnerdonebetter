package secrets

import "context"

// NoopSecretSource returns empty string for all secrets.
type NoopSecretSource struct{}

// GetSecret returns empty string.
func (n *NoopSecretSource) GetSecret(ctx context.Context, name string) (string, error) {
	return "", nil
}

// Close is a no-op.
func (n *NoopSecretSource) Close() error {
	return nil
}

// NewNoopSecretSource returns a SecretSource that returns empty strings.
func NewNoopSecretSource() SecretSource {
	return &NoopSecretSource{}
}
