package models

import (
	"context"
	"net/http"
)

type (
	// RecipeStep represents a recipe step
	RecipeStep struct {
		ID                        uint64  `json:"id"`
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparation_id"`
		PrerequisiteStep          uint64  `json:"prerequisite_step"`
		MinEstimatedTimeInSeconds uint32  `json:"min_estimated_time_in_seconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"max_estimated_time_in_seconds"`
		TemperatureInCelsius      *uint16 `json:"temperature_in_celsius"`
		Notes                     string  `json:"notes"`
		RecipeID                  uint64  `json:"recipe_id"`
		CreatedOn                 uint64  `json:"created_on"`
		UpdatedOn                 *uint64 `json:"updated_on"`
		ArchivedOn                *uint64 `json:"archived_on"`
		BelongsTo                 uint64  `json:"belongs_to"`
	}

	// RecipeStepList represents a list of recipe steps
	RecipeStepList struct {
		Pagination
		RecipeSteps []RecipeStep `json:"recipe_steps"`
	}

	// RecipeStepCreationInput represents what a user could set as input for creating recipe steps
	RecipeStepCreationInput struct {
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparation_id"`
		PrerequisiteStep          uint64  `json:"prerequisite_step"`
		MinEstimatedTimeInSeconds uint32  `json:"min_estimated_time_in_seconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"max_estimated_time_in_seconds"`
		TemperatureInCelsius      *uint16 `json:"temperature_in_celsius"`
		Notes                     string  `json:"notes"`
		RecipeID                  uint64  `json:"recipe_id"`
		BelongsTo                 uint64  `json:"-"`
	}

	// RecipeStepUpdateInput represents what a user could set as input for updating recipe steps
	RecipeStepUpdateInput struct {
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparation_id"`
		PrerequisiteStep          uint64  `json:"prerequisite_step"`
		MinEstimatedTimeInSeconds uint32  `json:"min_estimated_time_in_seconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"max_estimated_time_in_seconds"`
		TemperatureInCelsius      *uint16 `json:"temperature_in_celsius"`
		Notes                     string  `json:"notes"`
		RecipeID                  uint64  `json:"recipe_id"`
		BelongsTo                 uint64  `json:"-"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently
	RecipeStepDataManager interface {
		GetRecipeStep(ctx context.Context, recipeStepID, userID uint64) (*RecipeStep, error)
		GetRecipeStepCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllRecipeStepsCount(ctx context.Context) (uint64, error)
		GetRecipeSteps(ctx context.Context, filter *QueryFilter, userID uint64) (*RecipeStepList, error)
		GetAllRecipeStepsForUser(ctx context.Context, userID uint64) ([]RecipeStep, error)
		CreateRecipeStep(ctx context.Context, input *RecipeStepCreationInput) (*RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, updated *RecipeStep) error
		ArchiveRecipeStep(ctx context.Context, id, userID uint64) error
	}

	// RecipeStepDataServer describes a structure capable of serving traffic related to recipe steps
	RecipeStepDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an RecipeStepInput with a recipe step
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

// ToInput creates a RecipeStepUpdateInput struct for a recipe step
func (x *RecipeStep) ToInput() *RecipeStepUpdateInput {
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
