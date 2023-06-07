package search

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type (
	Searchable interface {
		any | // TODO: just make these the subset types
			types.ValidIngredient |
			types.ValidInstrument |
			types.ValidMeasurementUnit |
			types.ValidPreparation |
			types.ValidIngredientState |
			types.ValidIngredientMeasurementUnit |
			types.ValidMeasurementUnitConversion |
			types.ValidPreparationInstrument |
			types.ValidIngredientPreparation
	}

	// IndexSearcher is our wrapper interface for querying a text search index.
	IndexSearcher[T Searchable] interface {
		Search(ctx context.Context, query string) (ids []*T, err error)
	}

	// IndexManager is our wrapper interface for a text search index.
	IndexManager interface {
		Index(ctx context.Context, id string, value any) error
		Delete(ctx context.Context, id string) (err error)
		Wipe(ctx context.Context) error
	}

	// Index is our wrapper interface for a text search index.
	Index[T Searchable] interface {
		IndexSearcher[T]
		IndexManager
	}
)
