package search

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type (
	Searchable interface {
		types.ValidPreparation |
			types.ValidIngredient |
			types.ValidInstrument |
			types.ValidIngredientState |
			types.ValidMeasurementUnit
	}

	// IndexSearcher is our wrapper interface for querying a text search index.
	IndexSearcher[T Searchable] interface {
		Search(ctx context.Context, query string) (ids []*T, err error)
	}

	// IndexManager is our wrapper interface for a text search index.
	IndexManager[T Searchable] interface {
		IndexSearcher[T]
		Index(ctx context.Context, id string, value any) error
		Delete(ctx context.Context, id string) (err error)
		Wipe(ctx context.Context) error
	}
)
