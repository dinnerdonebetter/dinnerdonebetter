package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// RecipeStepDataType indicates an event is related to a recipe step.
	RecipeStepDataType dataType = "recipe_step"

	// RecipeStepCreatedCustomerEventType indicates a recipe step was created.
	RecipeStepCreatedCustomerEventType CustomerEventType = "recipe_step_created"
	// RecipeStepUpdatedCustomerEventType indicates a recipe step was updated.
	RecipeStepUpdatedCustomerEventType CustomerEventType = "recipe_step_updated"
	// RecipeStepArchivedCustomerEventType indicates a recipe step was archived.
	RecipeStepArchivedCustomerEventType CustomerEventType = "recipe_step_archived"
)

func init() {
	gob.Register(new(RecipeStep))
	gob.Register(new(RecipeStepList))
	gob.Register(new(RecipeStepCreationRequestInput))
	gob.Register(new(RecipeStepUpdateRequestInput))
}

type (
	// RecipeStep represents a recipe step.
	RecipeStep struct {
		_                         struct{}
		ArchivedOn                *uint64                 `json:"archivedOn"`
		LastUpdatedOn             *uint64                 `json:"lastUpdatedOn"`
		TemperatureInCelsius      *uint16                 `json:"temperatureInCelsius"`
		Notes                     string                  `json:"notes"`
		ID                        string                  `json:"id"`
		Yields                    string                  `json:"yields"`
		BelongsToRecipe           string                  `json:"belongsToRecipe"`
		Preparation               ValidPreparation        `json:"preparation"`
		Ingredients               []*RecipeStepIngredient `json:"ingredients"`
		PrerequisiteStep          uint64                  `json:"prerequisiteStep"`
		Index                     uint                    `json:"index"`
		CreatedOn                 uint64                  `json:"createdOn"`
		MaxEstimatedTimeInSeconds uint32                  `json:"maxEstimatedTimeInSeconds"`
		MinEstimatedTimeInSeconds uint32                  `json:"minEstimatedTimeInSeconds"`
		Optional                  bool                    `json:"optional"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList struct {
		_           struct{}
		RecipeSteps []*RecipeStep `json:"recipeSteps"`
		Pagination
	}

	// RecipeStepCreationRequestInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationRequestInput struct {
		_                         struct{}
		TemperatureInCelsius      *uint16                                     `json:"temperatureInCelsius"`
		Yields                    string                                      `json:"yields"`
		Notes                     string                                      `json:"notes"`
		PreparationID             string                                      `json:"preparationID"`
		BelongsToRecipe           string                                      `json:"-"`
		ID                        string                                      `json:"-"`
		Ingredients               []*RecipeStepIngredientCreationRequestInput `json:"ingredients"`
		Index                     uint                                        `json:"index"`
		PrerequisiteStep          uint64                                      `json:"prerequisiteStep"`
		MinEstimatedTimeInSeconds uint32                                      `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32                                      `json:"maxEstimatedTimeInSeconds"`
		Optional                  bool                                        `json:"optional"`
	}

	// RecipeStepDatabaseCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepDatabaseCreationInput struct {
		_                         struct{}
		TemperatureInCelsius      *uint16                                      `json:"temperatureInCelsius"`
		Yields                    string                                       `json:"yields"`
		Notes                     string                                       `json:"notes"`
		PreparationID             string                                       `json:"preparationID"`
		BelongsToRecipe           string                                       `json:"belongsToRecipe"`
		ID                        string                                       `json:"id"`
		Ingredients               []*RecipeStepIngredientDatabaseCreationInput `json:"ingredients"`
		Index                     uint                                         `json:"index"`
		PrerequisiteStep          uint64                                       `json:"prerequisiteStep"`
		MinEstimatedTimeInSeconds uint32                                       `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32                                       `json:"maxEstimatedTimeInSeconds"`
		Optional                  bool                                         `json:"optional"`
	}

	// RecipeStepUpdateRequestInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateRequestInput struct {
		_                         struct{}
		TemperatureInCelsius      *uint16          `json:"temperatureInCelsius"`
		Notes                     string           `json:"notes"`
		Why                       string           `json:"why"`
		BelongsToRecipe           string           `json:"belongsToRecipe"`
		Yields                    string           `json:"yields"`
		Preparation               ValidPreparation `json:"preparation"`
		Index                     uint             `json:"index"`
		PrerequisiteStep          uint64           `json:"prerequisiteStep"`
		MinEstimatedTimeInSeconds uint32           `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32           `json:"maxEstimatedTimeInSeconds"`
		Optional                  bool             `json:"optional"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*RecipeStep, error)
		GetTotalRecipeStepCount(ctx context.Context) (uint64, error)
		GetRecipeSteps(ctx context.Context, recipeID string, filter *QueryFilter) (*RecipeStepList, error)
		GetRecipeStepsWithIDs(ctx context.Context, recipeID string, limit uint8, ids []string) ([]*RecipeStep, error)
		CreateRecipeStep(ctx context.Context, input *RecipeStepDatabaseCreationInput) (*RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, updated *RecipeStep) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error
	}

	// RecipeStepDataService describes a structure capable of serving traffic related to recipe steps.
	RecipeStepDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepUpdateRequestInput with a recipe step.
func (x *RecipeStep) Update(input *RecipeStepUpdateRequestInput) {
	if input.Index != 0 && input.Index != x.Index {
		x.Index = input.Index
	}

	if input.Preparation.Name != "" && input.Preparation.Name != x.Preparation.Name {
		x.Preparation.Name = input.Preparation.Name
	}

	if input.Preparation.Description != "" && input.Preparation.Description != x.Preparation.Description {
		x.Preparation.Description = input.Preparation.Description
	}

	if input.Preparation.IconPath != "" && input.Preparation.IconPath != x.Preparation.IconPath {
		x.Preparation.IconPath = input.Preparation.IconPath
	}

	if input.PrerequisiteStep != 0 && input.PrerequisiteStep != x.PrerequisiteStep {
		x.PrerequisiteStep = input.PrerequisiteStep
	}

	if input.MinEstimatedTimeInSeconds != 0 && input.MinEstimatedTimeInSeconds != x.MinEstimatedTimeInSeconds {
		x.MinEstimatedTimeInSeconds = input.MinEstimatedTimeInSeconds
	}

	if input.MaxEstimatedTimeInSeconds != 0 && input.MaxEstimatedTimeInSeconds != x.MaxEstimatedTimeInSeconds {
		x.MaxEstimatedTimeInSeconds = input.MaxEstimatedTimeInSeconds
	}

	if input.TemperatureInCelsius != nil && (x.TemperatureInCelsius == nil || (*input.TemperatureInCelsius != 0 && *input.TemperatureInCelsius != *x.TemperatureInCelsius)) {
		x.TemperatureInCelsius = input.TemperatureInCelsius
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCreationRequestInput.
func (x *RecipeStepCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.Ingredients, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepDatabaseCreationInput.
func (x *RecipeStepDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.PreparationID, validation.Required),
	)
}

// RecipeStepDatabaseCreationInputFromRecipeStepCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepDatabaseCreationInputFromRecipeStepCreationInput(input *RecipeStepCreationRequestInput) *RecipeStepDatabaseCreationInput {
	ingredients := []*RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		ingredients = append(ingredients, RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(ingredient))
	}

	x := &RecipeStepDatabaseCreationInput{
		Index:                     input.Index,
		PreparationID:             input.PreparationID,
		PrerequisiteStep:          input.PrerequisiteStep,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      input.TemperatureInCelsius,
		Notes:                     input.Notes,
		Ingredients:               ingredients,
	}

	return x
}

var _ validation.ValidatableWithContext = (*RecipeStepUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepUpdateRequestInput.
func (x *RecipeStepUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Index, validation.Required),
		validation.Field(&x.Preparation, validation.Required),
		validation.Field(&x.PrerequisiteStep, validation.Required),
		validation.Field(&x.MinEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.MaxEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.TemperatureInCelsius, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
