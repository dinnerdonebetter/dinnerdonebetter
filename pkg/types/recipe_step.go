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
		_                             struct{}
		LastUpdatedAt                 *uint64                 `json:"lastUpdatedAt"`
		MaximumTemperatureInCelsius   *uint16                 `json:"maximumTemperatureInCelsius"`
		ArchivedAT                    *uint64                 `json:"archivedAt"`
		MinimumTemperatureInCelsius   *uint16                 `json:"minimumTemperatureInCelsius"`
		Notes                         string                  `json:"notes"`
		BelongsToRecipe               string                  `json:"belongsToRecipe"`
		ID                            string                  `json:"id"`
		ExplicitInstructions          string                  `json:"explicitInstructions"`
		Products                      []*RecipeStepProduct    `json:"products"`
		Ingredients                   []*RecipeStepIngredient `json:"ingredients"`
		Instruments                   []*RecipeStepInstrument `json:"instruments"`
		Preparation                   ValidPreparation        `json:"preparation"`
		CreatedAt                     uint64                  `json:"createdAt"`
		Index                         uint32                  `json:"index"`
		MaximumEstimatedTimeInSeconds uint32                  `json:"maximumEstimatedTimeInSeconds"`
		MinimumEstimatedTimeInSeconds uint32                  `json:"minimumEstimatedTimeInSeconds"`
		Optional                      bool                    `json:"optional"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList struct {
		_           struct{}
		RecipeSteps []*RecipeStep `json:"data"`
		Pagination
	}

	// RecipeStepCreationRequestInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationRequestInput struct {
		_                             struct{}
		MinimumTemperatureInCelsius   *uint16                                     `json:"minimumTemperatureInCelsius"`
		MaximumTemperatureInCelsius   *uint16                                     `json:"maximumTemperatureInCelsius"`
		BelongsToRecipe               string                                      `json:"-"`
		Notes                         string                                      `json:"notes"`
		PreparationID                 string                                      `json:"preparationID"`
		ID                            string                                      `json:"-"`
		ExplicitInstructions          string                                      `json:"explicitInstructions"`
		Products                      []*RecipeStepProductCreationRequestInput    `json:"products"`
		Instruments                   []*RecipeStepInstrumentCreationRequestInput `json:"instruments"`
		Ingredients                   []*RecipeStepIngredientCreationRequestInput `json:"ingredients"`
		Index                         uint32                                      `json:"index"`
		MinimumEstimatedTimeInSeconds uint32                                      `json:"minimumEstimatedTimeInSeconds"`
		MaximumEstimatedTimeInSeconds uint32                                      `json:"maximumEstimatedTimeInSeconds"`
		Optional                      bool                                        `json:"optional"`
	}

	// RecipeStepDatabaseCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepDatabaseCreationInput struct {
		_                             struct{}
		MinimumTemperatureInCelsius   *uint16                                      `json:"minimumTemperatureInCelsius"`
		MaximumTemperatureInCelsius   *uint16                                      `json:"maximumTemperatureInCelsius"`
		PreparationID                 string                                       `json:"preparationID"`
		ID                            string                                       `json:"id"`
		Notes                         string                                       `json:"notes"`
		BelongsToRecipe               string                                       `json:"belongsToRecipe"`
		ExplicitInstructions          string                                       `json:"explicitInstructions"`
		Instruments                   []*RecipeStepInstrumentDatabaseCreationInput `json:"instruments"`
		Ingredients                   []*RecipeStepIngredientDatabaseCreationInput `json:"ingredients"`
		Products                      []*RecipeStepProductDatabaseCreationInput    `json:"products"`
		Index                         uint32                                       `json:"index"`
		MinimumEstimatedTimeInSeconds uint32                                       `json:"minimumEstimatedTimeInSeconds"`
		MaximumEstimatedTimeInSeconds uint32                                       `json:"maximumEstimatedTimeInSeconds"`
		Optional                      bool                                         `json:"optional"`
	}

	// RecipeStepUpdateRequestInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateRequestInput struct {
		_                             struct{}
		MinimumTemperatureInCelsius   *uint16           `json:"minimumTemperatureInCelsius"`
		MaximumTemperatureInCelsius   *uint16           `json:"maximumTemperatureInCelsius"`
		Notes                         *string           `json:"notes"`
		Preparation                   *ValidPreparation `json:"preparation"`
		Index                         *uint32           `json:"index"`
		MinimumEstimatedTimeInSeconds *uint32           `json:"minimumEstimatedTimeInSeconds"`
		MaximumEstimatedTimeInSeconds *uint32           `json:"maximumEstimatedTimeInSeconds"`
		Optional                      *bool             `json:"optional"`
		ExplicitInstructions          *string           `json:"explicitInstructions"`
		BelongsToRecipe               string            `json:"belongsToRecipe"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*RecipeStep, error)
		GetRecipeSteps(ctx context.Context, recipeID string, filter *QueryFilter) (*RecipeStepList, error)
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

	if input.MinimumEstimatedTimeInSeconds != nil && *input.MinimumEstimatedTimeInSeconds != x.MinimumEstimatedTimeInSeconds {
		x.MinimumEstimatedTimeInSeconds = *input.MinimumEstimatedTimeInSeconds
	}

	if input.MaximumEstimatedTimeInSeconds != nil && *input.MaximumEstimatedTimeInSeconds != x.MaximumEstimatedTimeInSeconds {
		x.MaximumEstimatedTimeInSeconds = *input.MaximumEstimatedTimeInSeconds
	}

	if input.MinimumTemperatureInCelsius != nil && (x.MinimumTemperatureInCelsius == nil || (*input.MinimumTemperatureInCelsius != 0 && *input.MinimumTemperatureInCelsius != *x.MinimumTemperatureInCelsius)) {
		x.MinimumTemperatureInCelsius = input.MinimumTemperatureInCelsius
	}

	if input.MaximumTemperatureInCelsius != nil && (x.MaximumTemperatureInCelsius == nil || (*input.MaximumTemperatureInCelsius != 0 && *input.MaximumTemperatureInCelsius != *x.MaximumTemperatureInCelsius)) {
		x.MaximumTemperatureInCelsius = input.MaximumTemperatureInCelsius
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	// TODO: do something about products here

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}

	if input.MinimumTemperatureInCelsius != nil && (x.MinimumTemperatureInCelsius == nil || (*input.MinimumTemperatureInCelsius != 0 && *input.MinimumTemperatureInCelsius != *x.MinimumTemperatureInCelsius)) {
		x.MinimumTemperatureInCelsius = input.MinimumTemperatureInCelsius
	}

	if input.ExplicitInstructions != nil && *input.ExplicitInstructions != x.ExplicitInstructions {
		x.ExplicitInstructions = *input.ExplicitInstructions
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
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         &input.Notes,
		BelongsToRecipe:               input.BelongsToRecipe,
		Preparation:                   &input.Preparation,
		Index:                         &input.Index,
		MinimumEstimatedTimeInSeconds: &input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: &input.MaximumEstimatedTimeInSeconds,
		Optional:                      &input.Optional,
		ExplicitInstructions:          &input.ExplicitInstructions,
	}

	return x
}

// RecipeStepDatabaseCreationInputFromRecipeStepCreationInput creates a DatabaseCreationInput from a CreationInput.
func RecipeStepDatabaseCreationInputFromRecipeStepCreationInput(input *RecipeStepCreationRequestInput) *RecipeStepDatabaseCreationInput {
	ingredients := []*RecipeStepIngredientDatabaseCreationInput{}
	for _, ingredient := range input.Ingredients {
		ingredients = append(ingredients, RecipeStepIngredientDatabaseCreationInputFromRecipeStepIngredientCreationInput(ingredient))
	}

	instruments := []*RecipeStepInstrumentDatabaseCreationInput{}
	for _, instrument := range input.Instruments {
		instruments = append(instruments, RecipeStepInstrumentDatabaseCreationInputFromRecipeStepInstrumentCreationInput(instrument))
	}

	products := []*RecipeStepProductDatabaseCreationInput{}
	for _, product := range input.Products {
		products = append(products, RecipeStepProductDatabaseCreationInputFromRecipeStepProductCreationInput(product))
	}

	x := &RecipeStepDatabaseCreationInput{
		Index:                         input.Index,
		PreparationID:                 input.PreparationID,
		MinimumEstimatedTimeInSeconds: input.MinimumEstimatedTimeInSeconds,
		MaximumEstimatedTimeInSeconds: input.MaximumEstimatedTimeInSeconds,
		MinimumTemperatureInCelsius:   input.MinimumTemperatureInCelsius,
		MaximumTemperatureInCelsius:   input.MaximumTemperatureInCelsius,
		Notes:                         input.Notes,
		Products:                      products,
		Optional:                      input.Optional,
		Ingredients:                   ingredients,
		Instruments:                   instruments,
		ExplicitInstructions:          input.ExplicitInstructions,
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
		validation.Field(&x.MinimumEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.MaximumEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.MinimumTemperatureInCelsius, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
