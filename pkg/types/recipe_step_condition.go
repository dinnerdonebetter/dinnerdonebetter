package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepConditionDataType indicates an event is related to a recipe step condition.
	RecipeStepConditionDataType dataType = "recipe_step_condition"

	// RecipeStepConditionCreatedCustomerEventType indicates a recipe step condition was created.
	RecipeStepConditionCreatedCustomerEventType CustomerEventType = "recipe_step_condition_created"
	// RecipeStepConditionUpdatedCustomerEventType indicates a recipe step condition was updated.
	RecipeStepConditionUpdatedCustomerEventType CustomerEventType = "recipe_step_condition_updated"
	// RecipeStepConditionArchivedCustomerEventType indicates a recipe step condition was archived.
	RecipeStepConditionArchivedCustomerEventType CustomerEventType = "recipe_step_condition_archived"

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
	gob.Register(new(RecipeStepCondition))
	gob.Register(new(RecipeStepConditionCreationRequestInput))
	gob.Register(new(RecipeStepConditionUpdateRequestInput))
}

type (
	// RecipeStepCondition represents a recipe step condition. Effectively, this says "Ingredients must be in IngredientState".
	RecipeStepCondition struct {
		_                   struct{}
		CreatedAt           time.Time                        `json:"createdAt"`
		ArchivedAt          *time.Time                       `json:"archivedAt"`
		LastUpdatedAt       *time.Time                       `json:"lastUpdatedAt"`
		IngredientState     ValidIngredientState             `json:"ingredientState"`
		ID                  string                           `json:"id"`
		BelongsToRecipeStep string                           `json:"belongsToRecipeStep"`
		Notes               string                           `json:"notes"`
		Ingredients         []*RecipeStepConditionIngredient `json:"ingredients"`
		Optional            bool                             `json:"optional"`
	}

	RecipeStepConditionIngredient struct {
		_ struct{}

		CreatedAt                    time.Time  `json:"createdAt"`
		ArchivedAt                   *time.Time `json:"archivedAt"`
		LastUpdatedAt                *time.Time `json:"lastUpdatedAt"`
		ID                           string     `json:"id"`
		BelongsToRecipeStepCondition string     `json:"belongsToRecipeStepCondition"`
		RecipeStepIngredient         string     `json:"recipeStepIngredient"`
	}

	// RecipeStepConditionCreationRequestInput represents what a user could set as input for creating recipe step conditions.
	RecipeStepConditionCreationRequestInput struct {
		_                   struct{}
		IngredientStateID   string                                               `json:"ingredientState"`
		BelongsToRecipeStep string                                               `json:"belongsToRecipeStep"`
		Notes               string                                               `json:"notes"`
		Ingredients         []*RecipeStepConditionIngredientCreationRequestInput `json:"ingredients"`
		Optional            bool                                                 `json:"optional"`
	}

	RecipeStepConditionIngredientCreationRequestInput struct {
		_ struct{}

		BelongsToRecipeStepCondition string `json:"belongsToRecipeStepCondition"`
		RecipeStepIngredient         string `json:"recipeStepIngredient"`
	}

	// RecipeStepConditionDatabaseCreationInput represents what a user could set as input for creating recipe step conditions.
	RecipeStepConditionDatabaseCreationInput struct {
		_                   struct{}
		ID                  string
		IngredientStateID   string
		BelongsToRecipeStep string
		Notes               string
		Ingredients         []*RecipeStepConditionIngredientDatabaseCreationInput
		Optional            bool
	}

	RecipeStepConditionIngredientDatabaseCreationInput struct {
		_ struct{}

		ID                           string
		BelongsToRecipeStepCondition string
		RecipeStepIngredient         string
	}

	// RecipeStepConditionUpdateRequestInput represents what a user could set as input for updating recipe step conditions.
	RecipeStepConditionUpdateRequestInput struct {
		_ struct{}

		IngredientStateID   *string `json:"ingredientState"`
		BelongsToRecipeStep *string `json:"belongsToRecipeStep"`
		Notes               *string `json:"notes"`
		Optional            *bool   `json:"optional"`
	}

	// RecipeStepConditionDataManager describes a structure capable of storing recipe step conditions permanently.
	RecipeStepConditionDataManager interface {
		RecipeStepConditionExists(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (bool, error)
		GetRecipeStepCondition(ctx context.Context, recipeID, recipeStepID, recipeStepIngredientID string) (*RecipeStepCondition, error)
		GetRecipeStepConditions(ctx context.Context, recipeID, recipeStepID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStepCondition], error)
		CreateRecipeStepCondition(ctx context.Context, input *RecipeStepConditionDatabaseCreationInput) (*RecipeStepCondition, error)
		UpdateRecipeStepCondition(ctx context.Context, updated *RecipeStepCondition) error
		ArchiveRecipeStepCondition(ctx context.Context, recipeStepID, recipeStepIngredientID string) error
	}

	// RecipeStepConditionDataService describes a structure capable of serving traffic related to recipe step conditions.
	RecipeStepConditionDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepConditionUpdateRequestInput with a recipe step condition.
func (x *RecipeStepCondition) Update(input *RecipeStepConditionUpdateRequestInput) {
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

var _ validation.ValidatableWithContext = (*RecipeStepConditionCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepConditionCreationRequestInput.
func (x *RecipeStepConditionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.Ingredients, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepConditionIngredientCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepConditionIngredientCreationRequestInput.
func (x *RecipeStepConditionIngredientCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStepCondition, validation.Required),
		validation.Field(&x.RecipeStepIngredient, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepConditionDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepConditionDatabaseCreationInput.
func (x *RecipeStepConditionDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
		validation.Field(&x.Ingredients, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepConditionIngredientDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepConditionIngredientDatabaseCreationInput.
func (x *RecipeStepConditionIngredientDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.BelongsToRecipeStepCondition, validation.Required),
		validation.Field(&x.RecipeStepIngredient, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepConditionUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepConditionUpdateRequestInput.
func (x *RecipeStepConditionUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.IngredientStateID, validation.Required),
		validation.Field(&x.BelongsToRecipeStep, validation.Required),
	)
}
