package search

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

type NoopIndexManagerProvider struct{}

func (*NoopIndexManagerProvider) ProvideIndexManager(context.Context, logging.Logger, IndexName, ...string) (IndexManager, error) {
	return &NoopIndexManager{}, nil
}

type NoopIndexManager struct{}

func (*NoopIndexManager) Search(context.Context, string, string, string) ([]string, error) {
	return []string{}, nil
}

func (*NoopIndexManager) Index(context.Context, string, interface{}) error {
	return nil
}

func (*NoopIndexManager) Delete(context.Context, string) error {
	return nil
}
