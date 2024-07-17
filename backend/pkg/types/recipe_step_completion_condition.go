package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepCompletionConditionCreatedCustomerEventType indicates a recipe step completion condition was created.
	RecipeStepCompletionConditionCreatedCustomerEventType ServiceEventType = "recipe_step_completion_condition_created"
	// RecipeStepCompletionConditionUpdatedCustomerEventType indicates a recipe step completion condition was updated.
	RecipeStepCompletionConditionUpdatedCustomerEventType ServiceEventType = "recipe_step_completion_condition_updated"
	// RecipeStepCompletionConditionArchivedCustomerEventType indicates a recipe step completion condition was archived.
	RecipeStepCompletionConditionArchivedCustomerEventType ServiceEventType = "recipe_step_completion_condition_archived"
)

func init() {
	gob.Register(new(RecipeStepCompletionCondition))
	gob.Register(new(RecipeStepCompletionConditionCreationRequestInput))
	gob.Register(new(RecipeStepCompletionConditionUpdateRequestInput))
}

type (
	// RecipeStepCompletionCondition represents a recipe step completion condition. Effectively, this says "Ingredients must be in IngredientState".
	RecipeStepCompletionCondition struct {
		_ struct{} `json:"-"`

		CreatedAt           time.Time                                  `json:"createdAt"`
		ArchivedAt          *time.Time                                 `json:"archivedAt"`
		LastUpdatedAt       *time.Time                                 `json:"lastUpdatedAt"`
		IngredientState     ValidIngredientState                       `json:"ingredientState"`
		ID                  string                                     `json:"id"`
		BelongsToRecipeStep string                                     `json:"belongsToRecipeStep"`
		Notes               string                                     `json:"notes"`
		Ingredients         []*RecipeStepCompletionConditionIngredient `json:"ingredients"`
		Optional            bool                                       `json:"optional"`
	}

	RecipeStepCompletionConditionIngredient struct {
		_ struct{} `json:"-"`

		CreatedAt                              time.Time  `json:"createdAt"`
		ArchivedAt                             *time.Time `json:"archivedAt"`
		LastUpdatedAt                          *time.Time `json:"lastUpdatedAt"`
		ID                                     string     `json:"id"`
		BelongsToRecipeStepCompletionCondition string     `json:"belongsToRecipeStepCompletionCondition"`
		RecipeStepIngredient                   string     `json:"recipeStepIngredient"`
	}

	// RecipeStepCompletionConditionCreationRequestInput represents what a user could set as input for creating recipe step completion conditions.
	RecipeStepCompletionConditionCreationRequestInput struct {
		_ struct{} `json:"-"`

		IngredientStateID   string   `json:"ingredientState"`
		BelongsToRecipeStep string   `json:"belongsToRecipeStep"`
		Notes               string   `json:"notes"`
		Ingredients         []uint64 `json:"ingredients"`
		Optional            bool     `json:"optional"`
	}

	// RecipeStepCompletionConditionIngredientCreationRequestInput represents what a user could set as input for creating recipe step completion condition ingredient.
	RecipeStepCompletionConditionIngredientCreationRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepIngredient string `json:"recipeStepIngredient"`
	}

	// RecipeStepCompletionConditionForExistingRecipeCreationRequestInput represents what a user could set as input for creating recipe step completion conditions for existing recipes.
	RecipeStepCompletionConditionForExistingRecipeCreationRequestInput struct {
		_ struct{} `json:"-"`

		IngredientStateID   string                                                                          `json:"ingredientState"`
		BelongsToRecipeStep string                                                                          `json:"belongsToRecipeStep"`
		Notes               string                                                                          `json:"notes"`
		Ingredients         []*RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput `json:"ingredients"`
		Optional            bool                                                                            `json:"optional"`
	}

	// RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput represents what a user could set as input for creating recipe step completion condition for existing recipes.
	RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput struct {
		_ struct{} `json:"-"`

		RecipeStepIngredient string `json:"recipeStepIngredient"`
	}

	// RecipeStepCompletionConditionDatabaseCreationInput represents what a user could set as input for creating recipe step completion conditions.
	RecipeStepCompletionConditionDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                  string
		IngredientStateID   string
		BelongsToRecipeStep string
		Notes               string
		Ingredients         []*RecipeStepCompletionConditionIngredientDatabaseCreationInput
		Optional            bool
	}

	RecipeStepCompletionConditionIngredientDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                                     string
		BelongsToRecipeStepCompletionCondition string
		RecipeStepIngredient                   string
	}

	// RecipeStepCompletionConditionUpdateRequestInput represents what a user could set as input for updating recipe step completion conditions.
	RecipeStepCompletionConditionUpdateRequestInput struct {
		_ struct{} `json:"-"`

		IngredientStateID   *string `json:"ingredientState,omitempty"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep,omitempty"`
		Notes               *string `json:"notes,omitempty"`
		Optional            *bool   `json:"optional,omitempty"`
	}

	// RecipeStepCompletionConditionDataManager describes a structure capable of storing recipe step completion conditions permanently.
	RecipeStepCompletionConditionDataManager interface {
		RecipeStepCompletionConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error)
		GetRecipeStepCompletionCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*RecipeStepCompletionCondition, error)
		GetRecipeStepCompletionConditions(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepCompletionCondition], error)
		CreateRecipeStepCompletionCondition(ctx context.Context, input *RecipeStepCompletionConditionDatabaseCreationInput) (*RecipeStepCompletionCondition, error)
		UpdateRecipeStepCompletionCondition(ctx context.Context, updated *RecipeStepCompletionCondition) error
		ArchiveRecipeStepCompletionCondition(ctx context.Context, recipeStepID, recipeStepIngredientID string) error
	}

	// RecipeStepCompletionConditionDataService describes a structure capable of serving traffic related to recipe step completion conditions.
	RecipeStepCompletionConditionDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an RecipeStepCompletionConditionUpdateRequestInput with a recipe step completion condition.
func (x *RecipeStepCompletionCondition) Update(input *RecipeStepCompletionConditionUpdateRequestInput) {
	if input.IngredientStateID != nil && *input.IngredientStateID != x.IngredientState.ID {
		x.IngredientState = ValidIngredientState{ID: *input.IngredientStateID}
	}

	if input.BelongsToRecipeStep != nil && *input.BelongsToRecipeStep != x.BelongsToRecipeStep {
		x.BelongsToRecipeStep = *input.BelongsToRecipeStep
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionCreationRequestInput.
func (x *RecipeStepCompletionConditionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.Ingredients, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionIngredientCreationRequestInput.
func (x *RecipeStepCompletionConditionIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.RecipeStepIngredient, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionForExistingRecipeCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionForExistingRecipeCreationRequestInput.
func (x *RecipeStepCompletionConditionForExistingRecipeCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.Ingredients, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput.
func (x *RecipeStepCompletionConditionIngredientForExistingRecipeCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.RecipeStepIngredient, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionDatabaseCreationInput.
func (x *RecipeStepCompletionConditionDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
		validation.Field(&x.Ingredients, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionIngredientDatabaseCreationInput.
func (x *RecipeStepCompletionConditionIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStepCompletionCondition, validation.Required),
		validation.Field(&x.RecipeStepIngredient, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepCompletionConditionUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCompletionConditionUpdateRequestInput.
func (x *RecipeStepCompletionConditionUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}
