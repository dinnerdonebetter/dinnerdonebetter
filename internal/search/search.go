package search

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability/logging"
)

type (
	// IndexPath is a type alias for dependency injection's sake.
	IndexPath string

	// IndexName is a type alias for dependency injection's sake.
	IndexName string

	// IndexManager is our wrapper interface for a text search index.
	IndexManager interface {
		Index(ctx context.Context, id string, value interface{}) error
		Search(ctx context.Context, byField, query, householdID string) (ids []string, err error)
		Delete(ctx context.Context, id string) (err error)
	}

	// IndexManagerProvider is a struct that provides a IndexManager for a given index.
	IndexManagerProvider interface {
		ProvideIndexManager(ctx context.Context, logger logging.Logger, name IndexName, fields ...string) (IndexManager, error)
	}
)
