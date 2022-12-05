package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepCompletionConditionDataType indicates an event is related to a recipe step completion condition.
	RecipeStepCompletionConditionDataType dataType = "recipe_step_completion_condition"

	// RecipeStepCompletionConditionCreatedCustomerEventType indicates a recipe step completion condition was created.
	RecipeStepCompletionConditionCreatedCustomerEventType CustomerEventType = "recipe_step_completion_condition_created"
	// RecipeStepCompletionConditionUpdatedCustomerEventType indicates a recipe step completion condition was updated.
	RecipeStepCompletionConditionUpdatedCustomerEventType CustomerEventType = "recipe_step_completion_condition_updated"
	// RecipeStepCompletionConditionArchivedCustomerEventType indicates a recipe step completion condition was archived.
	RecipeStepCompletionConditionArchivedCustomerEventType CustomerEventType = "recipe_step_completion_condition_archived"

	// IngredientStateAttributeTypeTexture represents the ingredient attribute type for texture.
	IngredientStateAttributeTypeTexture = "texture"
	// IngredientStateAttributeTypeConsistency represents the ingredient attribute type for consistency.
	IngredientStateAttributeTypeConsistency = "consistency"
	// IngredientStateAttributeTypeColor represents the ingredient attribute type for color.
	IngredientStateAttributeTypeColor = "color"
	// IngredientStateAttributeTypeAppearance represents the ingredient attribute type for appearance.
	IngredientStateAttributeTypeAppearance = "appearance"
	// IngredientStateAttributeTypeOdor represents the ingredient attribute type for odor.
	IngredientStateAttributeTypeOdor = "odor"
	// IngredientStateAttributeTypeTaste represents the ingredient attribute type for taste.
	IngredientStateAttributeTypeTaste = "taste"
	// IngredientStateAttributeTypeSound represents the ingredient attribute type for sound.
	IngredientStateAttributeTypeSound = "sound"
	// IngredientStateAttributeTypeOther represents the ingredient attribute type for other.
	IngredientStateAttributeTypeOther = "other"
)

func init() {
	gob.Register(new(RecipeStepCompletionCondition))
	gob.Register(new(RecipeStepCompletionConditionCreationRequestInput))
	gob.Register(new(RecipeStepCompletionConditionUpdateRequestInput))
}

type (
	// RecipeStepCompletionCondition represents a recipe step completion condition. Effectively, this says "Ingredients must be in IngredientState".
	RecipeStepCompletionCondition struct {
		_                   struct{}
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
		_ struct{}

		CreatedAt                              time.Time  `json:"createdAt"`
		ArchivedAt                             *time.Time `json:"archivedAt"`
		LastUpdatedAt                          *time.Time `json:"lastUpdatedAt"`
		ID                                     string     `json:"id"`
		BelongsToRecipeStepCompletionCondition string     `json:"belongsToRecipeStepCompletionCondition"`
		RecipeStepIngredient                   string     `json:"recipeStepIngredient"`
	}

	// RecipeStepCompletionConditionCreationRequestInput represents what a user could set as input for creating recipe step completion conditions.
	RecipeStepCompletionConditionCreationRequestInput struct {
		_                   struct{}
		IngredientStateID   string                                                         `json:"ingredientState"`
		BelongsToRecipeStep string                                                         `json:"belongsToRecipeStep"`
		Notes               string                                                         `json:"notes"`
		Ingredients         []*RecipeStepCompletionConditionIngredientCreationRequestInput `json:"ingredients"`
		Optional            bool                                                           `json:"optional"`
	}

	RecipeStepCompletionConditionIngredientCreationRequestInput struct {
		_ struct{}

		RecipeStepIngredient string `json:"recipeStepIngredient"`
	}

	// RecipeStepCompletionConditionDatabaseCreationInput represents what a user could set as input for creating recipe step completion conditions.
	RecipeStepCompletionConditionDatabaseCreationInput struct {
		_                   struct{}
		ID                  string
		IngredientStateID   string
		BelongsToRecipeStep string
		Notes               string
		Ingredients         []*RecipeStepCompletionConditionIngredientDatabaseCreationInput
		Optional            bool
	}

	RecipeStepCompletionConditionIngredientDatabaseCreationInput struct {
		_ struct{}

		ID                                     string
		BelongsToRecipeStepCompletionCondition string
		RecipeStepIngredient                   string
	}

	// RecipeStepCompletionConditionUpdateRequestInput represents what a user could set as input for updating recipe step completion conditions.
	RecipeStepCompletionConditionUpdateRequestInput struct {
		_ struct{}

		IngredientStateID   *string `json:"ingredientState"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
		Notes               *string `json:"notes"`
		Optional            *bool   `json:"optional"`
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
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
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
