package search

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

type (
	RecipeSearchSubset struct {
		ID          string                   `json:"id,omitempty"`
		Name        string                   `json:"name,omitempty"`
		Description string                   `json:"description,omitempty"`
		Steps       []RecipeStepSearchSubset `json:"steps,omitempty"`
	}

	RecipeStepSearchSubset struct {
		Preparation string   `json:"preparation,omitempty"`
		Ingredients []string `json:"ingredients,omitempty"`
		Instruments []string `json:"instruments,omitempty"`
		Vessels     []string `json:"vessels,omitempty"`
	}

	MealSearchSubset struct {
		ID          string   `json:"id,omitempty"`
		Name        string   `json:"name,omitempty"`
		Description string   `json:"description,omitempty"`
		Recipes     []string `json:"recipes,omitempty"`
	}

	Searchable interface {
		RecipeSearchSubset |
			MealSearchSubset |
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
