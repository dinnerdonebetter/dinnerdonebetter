package search

import (
	"context"

	"github.com/dinnerdonebetter/backend/pkg/types"
)

const (
	IndexTypeRecipes = "recipes"
)

type (
	Searchable interface {
		any |
			types.Recipe |
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

	IndexRequest struct {
		Recipe    *types.Recipe `json:"recipe,omitempty"`
		IndexType string        `json:"type"`
	}

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
)

func SubsetFromRecipe(r *types.Recipe) *RecipeSearchSubset {
	x := &RecipeSearchSubset{
		ID:          r.ID,
		Name:        r.Name,
		Description: r.Description,
	}

	for _, step := range r.Steps {
		stepSubset := RecipeStepSearchSubset{
			Preparation: step.Preparation.Name,
		}

		for _, ingredient := range step.Ingredients {
			stepSubset.Ingredients = append(stepSubset.Ingredients, ingredient.Name)
		}

		for _, instrument := range step.Instruments {
			stepSubset.Instruments = append(stepSubset.Instruments, instrument.Name)
		}

		for _, vessel := range step.Vessels {
			stepSubset.Vessels = append(stepSubset.Vessels, vessel.Name)
		}

		x.Steps = append(x.Steps, stepSubset)
	}

	return x
}
