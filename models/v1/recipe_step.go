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
		ValidPreparationID        uint64  `json:"validPreparationID"`
		PrerequisiteStepID        *uint64 `json:"prerequisiteStepID"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		YieldsProductName         string  `json:"yieldsProductName"`
		YieldsQuantity            uint    `json:"yieldsQuantity"`
		Notes                     string  `json:"notes"`
		CreatedOn                 uint64  `json:"createdOn"`
		UpdatedOn                 *uint64 `json:"updatedOn"`
		ArchivedOn                *uint64 `json:"archivedOn"`
		BelongsToRecipe           uint64  `json:"belongsToRecipe"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList struct {
		Pagination
		RecipeSteps []RecipeStep `json:"recipeSteps"`
	}

	// RecipeStepCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationInput struct {
		Index                     uint    `json:"index"`
		ValidPreparationID        uint64  `json:"validPreparationID"`
		PrerequisiteStepID        *uint64 `json:"prerequisiteStepID"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		YieldsProductName         string  `json:"yieldsProductName"`
		YieldsQuantity            uint    `json:"yieldsQuantity"`
		Notes                     string  `json:"notes"`
		BelongsToRecipe           uint64  `json:"-"`
	}

	// RecipeStepUpdateInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateInput struct {
		Index                     uint    `json:"index"`
		ValidPreparationID        uint64  `json:"validPreparationID"`
		PrerequisiteStepID        *uint64 `json:"prerequisiteStepID"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		YieldsProductName         string  `json:"yieldsProductName"`
		YieldsQuantity            uint    `json:"yieldsQuantity"`
		Notes                     string  `json:"notes"`
		BelongsToRecipe           uint64  `json:"belongsToRecipe"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*RecipeStep, error)
		GetAllRecipeStepsCount(ctx context.Context) (uint64, error)
		GetRecipeSteps(ctx context.Context, recipeID uint64, filter *QueryFilter) (*RecipeStepList, error)
		CreateRecipeStep(ctx context.Context, input *RecipeStepCreationInput) (*RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, updated *RecipeStep) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) error
	}

	// RecipeStepDataServer describes a structure capable of serving traffic related to recipe steps.
	RecipeStepDataServer interface {
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

// Update merges an RecipeStepInput with a recipe step.
func (x *RecipeStep) Update(input *RecipeStepUpdateInput) {
	if input.Index != x.Index {
		x.Index = input.Index
	}

	if input.ValidPreparationID != x.ValidPreparationID {
		x.ValidPreparationID = input.ValidPreparationID
	}

	if input.PrerequisiteStepID != nil && input.PrerequisiteStepID != x.PrerequisiteStepID {
		x.PrerequisiteStepID = input.PrerequisiteStepID
	}

	if input.MinEstimatedTimeInSeconds != x.MinEstimatedTimeInSeconds {
		x.MinEstimatedTimeInSeconds = input.MinEstimatedTimeInSeconds
	}

	if input.MaxEstimatedTimeInSeconds != x.MaxEstimatedTimeInSeconds {
		x.MaxEstimatedTimeInSeconds = input.MaxEstimatedTimeInSeconds
	}

	if input.YieldsProductName != "" && input.YieldsProductName != x.YieldsProductName {
		x.YieldsProductName = input.YieldsProductName
	}

	if input.YieldsQuantity != x.YieldsQuantity {
		x.YieldsQuantity = input.YieldsQuantity
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

// ToUpdateInput creates a RecipeStepUpdateInput struct for a recipe step.
func (x *RecipeStep) ToUpdateInput() *RecipeStepUpdateInput {
	return &RecipeStepUpdateInput{
		Index:                     x.Index,
		ValidPreparationID:        x.ValidPreparationID,
		PrerequisiteStepID:        x.PrerequisiteStepID,
		MinEstimatedTimeInSeconds: x.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: x.MaxEstimatedTimeInSeconds,
		YieldsProductName:         x.YieldsProductName,
		YieldsQuantity:            x.YieldsQuantity,
		Notes:                     x.Notes,
	}
}
