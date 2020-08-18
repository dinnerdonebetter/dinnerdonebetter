package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStep represents a recipe step.
	RecipeStep struct {
		ID                        uint64  `json:"id"`
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparationID"`
		PrerequisiteStep          uint64  `json:"prerequisiteStep"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		TemperatureInCelsius      *uint16 `json:"temperatureInCelsius"`
		Notes                     string  `json:"notes"`
		RecipeID                  uint64  `json:"recipeID"`
		CreatedOn                 uint64  `json:"createdOn"`
		LastUpdatedOn             *uint64 `json:"lastUpdatedOn"`
		ArchivedOn                *uint64 `json:"archivedOn"`
		BelongsToRecipe           uint64  `json:"belongsToRecipe"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList struct {
		Pagination
		RecipeSteps []RecipeStep `json:"recipe_steps"`
	}

	// RecipeStepCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationInput struct {
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparationID"`
		PrerequisiteStep          uint64  `json:"prerequisiteStep"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		TemperatureInCelsius      *uint16 `json:"temperatureInCelsius"`
		Notes                     string  `json:"notes"`
		RecipeID                  uint64  `json:"recipeID"`
		BelongsToRecipe           uint64  `json:"-"`
	}

	// RecipeStepUpdateInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateInput struct {
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparationID"`
		PrerequisiteStep          uint64  `json:"prerequisiteStep"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		TemperatureInCelsius      *uint16 `json:"temperatureInCelsius"`
		Notes                     string  `json:"notes"`
		RecipeID                  uint64  `json:"recipeID"`
		BelongsToRecipe           uint64  `json:"belongsToRecipe"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*RecipeStep, error)
		GetAllRecipeStepsCount(ctx context.Context) (uint64, error)
		GetAllRecipeSteps(ctx context.Context, resultChannel chan []RecipeStep) error
		GetRecipeSteps(ctx context.Context, recipeID uint64, filter *QueryFilter) (*RecipeStepList, error)
		GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]RecipeStep, error)
		CreateRecipeStep(ctx context.Context, input *RecipeStepCreationInput) (*RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, updated *RecipeStep) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error
	}

	// RecipeStepDataServer describes a structure capable of serving traffic related to recipe steps.
	RecipeStepDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepInput with a recipe step.
func (x *RecipeStep) Update(input *RecipeStepUpdateInput) {
	if input.Index != x.Index {
		x.Index = input.Index
	}

	if input.PreparationID != x.PreparationID {
		x.PreparationID = input.PreparationID
	}

	if input.PrerequisiteStep != x.PrerequisiteStep {
		x.PrerequisiteStep = input.PrerequisiteStep
	}

	if input.MinEstimatedTimeInSeconds != x.MinEstimatedTimeInSeconds {
		x.MinEstimatedTimeInSeconds = input.MinEstimatedTimeInSeconds
	}

	if input.MaxEstimatedTimeInSeconds != x.MaxEstimatedTimeInSeconds {
		x.MaxEstimatedTimeInSeconds = input.MaxEstimatedTimeInSeconds
	}

	if input.TemperatureInCelsius != nil && input.TemperatureInCelsius != x.TemperatureInCelsius {
		x.TemperatureInCelsius = input.TemperatureInCelsius
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}

	if input.RecipeID != x.RecipeID {
		x.RecipeID = input.RecipeID
	}
}

// ToUpdateInput creates a RecipeStepUpdateInput struct for a recipe step.
func (x *RecipeStep) ToUpdateInput() *RecipeStepUpdateInput {
	return &RecipeStepUpdateInput{
		Index:                     x.Index,
		PreparationID:             x.PreparationID,
		PrerequisiteStep:          x.PrerequisiteStep,
		MinEstimatedTimeInSeconds: x.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: x.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      x.TemperatureInCelsius,
		Notes:                     x.Notes,
		RecipeID:                  x.RecipeID,
	}
}
