package models

import (
	"context"
	"net/http"
)

type (
	// RecipeIterationStep represents a recipe iteration step.
	RecipeIterationStep struct {
		ID              uint64  `json:"id"`
		StartedOn       *uint64 `json:"started_on"`
		EndedOn         *uint64 `json:"ended_on"`
		State           string  `json:"state"`
		CreatedOn       uint64  `json:"created_on"`
		UpdatedOn       *uint64 `json:"updated_on"`
		ArchivedOn      *uint64 `json:"archived_on"`
		BelongsToRecipe uint64  `json:"belongs_to_recipe"`
	}

	// RecipeIterationStepList represents a list of recipe iteration steps.
	RecipeIterationStepList struct {
		Pagination
		RecipeIterationSteps []RecipeIterationStep `json:"recipe_iteration_steps"`
	}

	// RecipeIterationStepCreationInput represents what a user could set as input for creating recipe iteration steps.
	RecipeIterationStepCreationInput struct {
		StartedOn       *uint64 `json:"started_on"`
		EndedOn         *uint64 `json:"ended_on"`
		State           string  `json:"state"`
		BelongsToRecipe uint64  `json:"-"`
	}

	// RecipeIterationStepUpdateInput represents what a user could set as input for updating recipe iteration steps.
	RecipeIterationStepUpdateInput struct {
		StartedOn       *uint64 `json:"started_on"`
		EndedOn         *uint64 `json:"ended_on"`
		State           string  `json:"state"`
		BelongsToRecipe uint64  `json:"belongs_to_recipe"`
	}

	// RecipeIterationStepDataManager describes a structure capable of storing recipe iteration steps permanently.
	RecipeIterationStepDataManager interface {
		RecipeIterationStepExists(ctx context.Context, recipeID, recipeIterationStepID uint64) (bool, error)
		GetRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) (*RecipeIterationStep, error)
		GetAllRecipeIterationStepsCount(ctx context.Context) (uint64, error)
		GetRecipeIterationSteps(ctx context.Context, recipeID uint64, filter *QueryFilter) (*RecipeIterationStepList, error)
		CreateRecipeIterationStep(ctx context.Context, input *RecipeIterationStepCreationInput) (*RecipeIterationStep, error)
		UpdateRecipeIterationStep(ctx context.Context, updated *RecipeIterationStep) error
		ArchiveRecipeIterationStep(ctx context.Context, recipeID, recipeIterationStepID uint64) error
	}

	// RecipeIterationStepDataServer describes a structure capable of serving traffic related to recipe iteration steps.
	RecipeIterationStepDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ExistenceHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeIterationStepInput with a recipe iteration step.
func (x *RecipeIterationStep) Update(input *RecipeIterationStepUpdateInput) {
	if input.StartedOn != nil && input.StartedOn != x.StartedOn {
		x.StartedOn = input.StartedOn
	}

	if input.EndedOn != nil && input.EndedOn != x.EndedOn {
		x.EndedOn = input.EndedOn
	}

	if input.State != "" && input.State != x.State {
		x.State = input.State
	}
}

// ToUpdateInput creates a RecipeIterationStepUpdateInput struct for a recipe iteration step.
func (x *RecipeIterationStep) ToUpdateInput() *RecipeIterationStepUpdateInput {
	return &RecipeIterationStepUpdateInput{
		StartedOn: x.StartedOn,
		EndedOn:   x.EndedOn,
		State:     x.State,
	}
}
