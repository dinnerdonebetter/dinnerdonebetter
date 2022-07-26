package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	maxIngredientsPerStep = 100
	maxProductsPerStep    = 100

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
		LastUpdatedOn             *uint64                 `json:"lastUpdatedOn"`
		MinTemperatureInCelsius   *float32                `json:"minTemperatureInCelsius"`
		MaxTemperatureInCelsius   *float32                `json:"maxTemperatureInCelsius"`
		ArchivedOn                *uint64                 `json:"archivedOn"`
		BelongsToRecipe           string                  `json:"belongsToRecipe"`
		Notes                     string                  `json:"notes"`
		ID                        string                  `json:"id"`
		Preparation               ValidPreparation        `json:"preparation"`
		Instruments               []*RecipeStepInstrument `json:"instruments"`
		Ingredients               []*RecipeStepIngredient `json:"ingredients"`
		Products                  []*RecipeStepProduct    `json:"products"`
		CreatedOn                 uint64                  `json:"createdOn"`
		Index                     uint32                  `json:"index"`
		MaxEstimatedTimeInSeconds uint32                  `json:"maxEstimatedTimeInSeconds"`
		MinEstimatedTimeInSeconds uint32                  `json:"minEstimatedTimeInSeconds"`
		Optional                  bool                    `json:"optional"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList struct {
		_           struct{}
		RecipeSteps []*RecipeStep `json:"data"`
		Pagination
	}

	// RecipeStepCreationRequestInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationRequestInput struct {
		_                         struct{}
		TemperatureInCelsius      *float32                                    `json:"temperatureInCelsius"`
		Instruments               []*RecipeStepInstrumentCreationRequestInput `json:"instruments"`
		Products                  []*RecipeStepProductCreationRequestInput    `json:"products"`
		Notes                     string                                      `json:"notes"`
		PreparationID             string                                      `json:"preparationID"`
		BelongsToRecipe           string                                      `json:"-"`
		ID                        string                                      `json:"-"`
		Ingredients               []*RecipeStepIngredientCreationRequestInput `json:"ingredients"`
		Index                     uint32                                      `json:"index"`
		MinEstimatedTimeInSeconds uint32                                      `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32                                      `json:"maxEstimatedTimeInSeconds"`
		Optional                  bool                                        `json:"optional"`
	}

	// RecipeStepDatabaseCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepDatabaseCreationInput struct {
		_                         struct{}
		TemperatureInCelsius      *float32                                     `json:"temperatureInCelsius"`
		Instruments               []*RecipeStepInstrumentDatabaseCreationInput `json:"instruments"`
		Products                  []*RecipeStepProductDatabaseCreationInput    `json:"products"`
		Notes                     string                                       `json:"notes"`
		PreparationID             string                                       `json:"preparationID"`
		BelongsToRecipe           string                                       `json:"belongsToRecipe"`
		ID                        string                                       `json:"id"`
		Ingredients               []*RecipeStepIngredientDatabaseCreationInput `json:"ingredients"`
		Index                     uint32                                       `json:"index"`
		MinEstimatedTimeInSeconds uint32                                       `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds uint32                                       `json:"maxEstimatedTimeInSeconds"`
		Optional                  bool                                         `json:"optional"`
	}

	// RecipeStepUpdateRequestInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateRequestInput struct {
		_                         struct{}
		TemperatureInCelsius      *float32           `json:"temperatureInCelsius"`
		Notes                     *string            `json:"notes"`
		Preparation               *ValidPreparation  `json:"preparation"`
		Index                     *uint32            `json:"index"`
		MinEstimatedTimeInSeconds *uint32            `json:"minEstimatedTimeInSeconds"`
		MaxEstimatedTimeInSeconds *uint32            `json:"maxEstimatedTimeInSeconds"`
		Optional                  *bool              `json:"optional"`
		BelongsToRecipe           string             `json:"belongsToRecipe"`
		Instruments               []*ValidInstrument `json:"instruments"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*RecipeStep, error)
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
	if input.Index != nil && *input.Index != x.Index {
		x.Index = *input.Index
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

	if input.MinEstimatedTimeInSeconds != nil && *input.MinEstimatedTimeInSeconds != x.MinEstimatedTimeInSeconds {
		x.MinEstimatedTimeInSeconds = *input.MinEstimatedTimeInSeconds
	}

	if input.MaxEstimatedTimeInSeconds != nil && *input.MaxEstimatedTimeInSeconds != x.MaxEstimatedTimeInSeconds {
		x.MaxEstimatedTimeInSeconds = *input.MaxEstimatedTimeInSeconds
	}

	if input.TemperatureInCelsius != nil && (x.MinTemperatureInCelsius == nil || (*input.TemperatureInCelsius != 0 && *input.TemperatureInCelsius != *x.MinTemperatureInCelsius)) {
		x.MinTemperatureInCelsius = input.TemperatureInCelsius
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	// TODO: do something about products here

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}

	if input.TemperatureInCelsius != nil && (x.MinTemperatureInCelsius == nil || (*input.TemperatureInCelsius != 0 && *input.TemperatureInCelsius != *x.MinTemperatureInCelsius)) {
		x.MinTemperatureInCelsius = input.TemperatureInCelsius
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCreationRequestInput.
func (x *RecipeStepCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.Products, validation.Required, validation.Length(1, maxProductsPerStep)),
		validation.Field(&x.Ingredients, validation.Required, validation.Length(1, maxIngredientsPerStep)),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepDatabaseCreationInput)(nil)

// ValidateWithContext validates a RecipeStepDatabaseCreationInput.
func (x *RecipeStepDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Products, validation.Required),
		validation.Field(&x.PreparationID, validation.Required),
	)
}

// RecipeStepUpdateRequestInputFromRecipeStep creates a DatabaseCreationInput from a CreationInput.
func RecipeStepUpdateRequestInputFromRecipeStep(input *RecipeStep) *RecipeStepUpdateRequestInput {
	x := &RecipeStepUpdateRequestInput{
		TemperatureInCelsius:      input.MinTemperatureInCelsius,
		Notes:                     &input.Notes,
		BelongsToRecipe:           input.BelongsToRecipe,
		Preparation:               &input.Preparation,
		Index:                     &input.Index,
		MinEstimatedTimeInSeconds: &input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: &input.MaxEstimatedTimeInSeconds,
		Optional:                  &input.Optional,
	}

	return x
}

// RecipeStepDatabaseCreationInputFromRecipeStepCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepDatabaseCreationInputFromRecipeStepCreationInput(input *RecipeStepCreationRequestInput) *RecipeStepDatabaseCreationInput {
	ingredients := []*RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		ingredients = append(ingredients, RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(ingredient))
	}

	products := []*RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		products = append(products, RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput(product))
	}

	x := &RecipeStepDatabaseCreationInput{
		Index:                     input.Index,
		PreparationID:             input.PreparationID,
		MinEstimatedTimeInSeconds: input.MinEstimatedTimeInSeconds,
		MaxEstimatedTimeInSeconds: input.MaxEstimatedTimeInSeconds,
		TemperatureInCelsius:      input.TemperatureInCelsius,
		Notes:                     input.Notes,
		Products:                  products,
		Optional:                  input.Optional,
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
		validation.Field(&x.MinEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.MaxEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.TemperatureInCelsius, validation.Required),
		// validation.Field(&x.Products, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
