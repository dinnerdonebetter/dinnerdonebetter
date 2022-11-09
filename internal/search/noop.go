package search

import (
	"context"

	"github.com/prixfixeco/backend/internal/observability/logging"
)

// NoopIndexManagerProvider is a noop IndexManagerProvider.
type NoopIndexManagerProvider struct{}

// ProvideIndexManager is a no-op method.
func (*NoopIndexManagerProvider) ProvideIndexManager(context.Context, logging.Logger, IndexName, ...string) (IndexManager, error) {
	return &NoopIndexManager{}, nil
}

// NoopIndexManager is a noop IndexManager.
type NoopIndexManager struct{}

// Search is a no-op method.
func (*NoopIndexManager) Search(context.Context, string, string, string) ([]string, error) {
	return []string{}, nil
}

// Index is a no-op method.
func (*NoopIndexManager) Index(context.Context, string, any) error {
	return nil
}

// Delete is a no-op method.
func (*NoopIndexManager) Delete(context.Context, string) error {
	return nil
}
