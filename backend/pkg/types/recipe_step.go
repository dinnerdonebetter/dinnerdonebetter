package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hashicorp/go-multierror"
)

const (
	maxIngredientsPerStep = 100

	// RecipeStepCreatedCustomerEventType indicates a recipe step was created.
	RecipeStepCreatedCustomerEventType ServiceEventType = "recipe_step_created"
	// RecipeStepUpdatedCustomerEventType indicates a recipe step was updated.
	RecipeStepUpdatedCustomerEventType ServiceEventType = "recipe_step_updated"
	// RecipeStepArchivedCustomerEventType indicates a recipe step was archived.
	RecipeStepArchivedCustomerEventType ServiceEventType = "recipe_step_archived"
)

func init() {
	gob.Register(new(RecipeStep))
	gob.Register(new(RecipeStepCreationRequestInput))
	gob.Register(new(RecipeStepUpdateRequestInput))
}

type (
	// RecipeStep represents a recipe step.
	RecipeStep struct {
		_ struct{} `json:"-"`

		CreatedAt               time.Time                        `json:"createdAt"`
		EstimatedTimeInSeconds  OptionalUint32Range              `json:"estimatedTimeInSeconds"`
		TemperatureInCelsius    OptionalFloat32Range             `json:"temperatureInCelsius"`
		ArchivedAt              *time.Time                       `json:"archivedAt"`
		LastUpdatedAt           *time.Time                       `json:"lastUpdatedAt"`
		BelongsToRecipe         string                           `json:"belongsToRecipe"`
		ConditionExpression     string                           `json:"conditionExpression"`
		ID                      string                           `json:"id"`
		Notes                   string                           `json:"notes"`
		ExplicitInstructions    string                           `json:"explicitInstructions"`
		Media                   []*RecipeMedia                   `json:"media"`
		Products                []*RecipeStepProduct             `json:"products"`
		Instruments             []*RecipeStepInstrument          `json:"instruments"`
		Vessels                 []*RecipeStepVessel              `json:"vessels"`
		CompletionConditions    []*RecipeStepCompletionCondition `json:"completionConditions"`
		Ingredients             []*RecipeStepIngredient          `json:"ingredients"`
		Preparation             ValidPreparation                 `json:"preparation"`
		Index                   uint32                           `json:"index"`
		Optional                bool                             `json:"optional"`
		StartTimerAutomatically bool                             `json:"startTimerAutomatically"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList []*RecipeStep

	// RecipeStepCreationRequestInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationRequestInput struct {
		_ struct{} `json:"-"`

		EstimatedTimeInSeconds  OptionalUint32Range                                  `json:"estimatedTimeInSeconds"`
		TemperatureInCelsius    OptionalFloat32Range                                 `json:"temperatureInCelsius"`
		PreparationID           string                                               `json:"preparationID"`
		Notes                   string                                               `json:"notes"`
		ConditionExpression     string                                               `json:"conditionExpression"`
		ExplicitInstructions    string                                               `json:"explicitInstructions"`
		Instruments             []*RecipeStepInstrumentCreationRequestInput          `json:"instruments"`
		Vessels                 []*RecipeStepVesselCreationRequestInput              `json:"vessels"`
		Products                []*RecipeStepProductCreationRequestInput             `json:"products"`
		Ingredients             []*RecipeStepIngredientCreationRequestInput          `json:"ingredients"`
		CompletionConditions    []*RecipeStepCompletionConditionCreationRequestInput `json:"completionConditions"`
		Index                   uint32                                               `json:"index"`
		Optional                bool                                                 `json:"optional"`
		StartTimerAutomatically bool                                                 `json:"startTimerAutomatically"`
	}

	// RecipeStepDatabaseCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		EstimatedTimeInSeconds  OptionalUint32Range
		TemperatureInCelsius    OptionalFloat32Range
		BelongsToRecipe         string
		PreparationID           string
		ID                      string
		Notes                   string
		ExplicitInstructions    string
		ConditionExpression     string
		Ingredients             []*RecipeStepIngredientDatabaseCreationInput
		Instruments             []*RecipeStepInstrumentDatabaseCreationInput
		Vessels                 []*RecipeStepVesselDatabaseCreationInput
		Products                []*RecipeStepProductDatabaseCreationInput
		CompletionConditions    []*RecipeStepCompletionConditionDatabaseCreationInput
		Index                   uint32
		Optional                bool
		StartTimerAutomatically bool
	}

	// RecipeStepUpdateRequestInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateRequestInput struct {
		_ struct{} `json:"-"`

		EstimatedTimeInSeconds  OptionalUint32Range  `json:"estimatedTimeInSeconds"`
		TemperatureInCelsius    OptionalFloat32Range `json:"temperatureInCelsius"`
		Notes                   *string              `json:"notes,omitempty"`
		Preparation             *ValidPreparation    `json:"preparation,omitempty"`
		Index                   *uint32              `json:"index,omitempty"`
		Optional                *bool                `json:"optional,omitempty"`
		ExplicitInstructions    *string              `json:"explicitInstructions,omitempty"`
		ConditionExpression     *string              `json:"conditionExpression,omitempty"`
		StartTimerAutomatically *bool                `json:"startTimerAutomatically"`
		BelongsToRecipe         string               `json:"belongsToRecipe"`
	}

	// RecipeStepSearchSubset represents the subset of values suitable to index for search.
	RecipeStepSearchSubset struct {
		_ struct{} `json:"-"`

		Preparation string    `json:"preparation,omitempty"`
		Ingredients []NamedID `json:"ingredients,omitempty"`
		Instruments []NamedID `json:"instruments,omitempty"`
		Vessels     []NamedID `json:"vessels,omitempty"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID string) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID string) (*RecipeStep, error)
		GetRecipeSteps(ctx context.Context, recipeID string, filter *QueryFilter) (*QueryFilteredResult[RecipeStep], error)
		CreateRecipeStep(ctx context.Context, input *RecipeStepDatabaseCreationInput) (*RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, updated *RecipeStep) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID string) error
	}

	// RecipeStepDataService describes a structure capable of serving traffic related to recipe steps.
	RecipeStepDataService interface {
		ListHandler(http.ResponseWriter, *http.Request)
		CreateHandler(http.ResponseWriter, *http.Request)
		ReadHandler(http.ResponseWriter, *http.Request)
		UpdateHandler(http.ResponseWriter, *http.Request)
		ArchiveHandler(http.ResponseWriter, *http.Request)
		ImageUploadHandler(http.ResponseWriter, *http.Request)
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

	if input.EstimatedTimeInSeconds.Min != nil && input.EstimatedTimeInSeconds.Min != x.EstimatedTimeInSeconds.Min {
		x.EstimatedTimeInSeconds.Min = input.EstimatedTimeInSeconds.Min
	}

	if input.EstimatedTimeInSeconds.Max != nil && input.EstimatedTimeInSeconds.Max != x.EstimatedTimeInSeconds.Max {
		x.EstimatedTimeInSeconds.Max = input.EstimatedTimeInSeconds.Max
	}

	if input.TemperatureInCelsius.Min != nil && (x.TemperatureInCelsius.Min == nil || (*input.TemperatureInCelsius.Min != 0 && *input.TemperatureInCelsius.Min != *x.TemperatureInCelsius.Min)) {
		x.TemperatureInCelsius.Min = input.TemperatureInCelsius.Min
	}

	if input.TemperatureInCelsius.Max != nil && (x.TemperatureInCelsius.Max == nil || (*input.TemperatureInCelsius.Max != 0 && *input.TemperatureInCelsius.Max != *x.TemperatureInCelsius.Max)) {
		x.TemperatureInCelsius.Max = input.TemperatureInCelsius.Max
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	// TODO: do something about products here

	if input.Optional != nil && *input.Optional != x.Optional {
		x.Optional = *input.Optional
	}

	if input.ExplicitInstructions != nil && *input.ExplicitInstructions != x.ExplicitInstructions {
		x.ExplicitInstructions = *input.ExplicitInstructions
	}

	if input.ConditionExpression != nil && *input.ConditionExpression != x.ConditionExpression {
		x.ConditionExpression = *input.ConditionExpression
	}

	if input.StartTimerAutomatically != nil && *input.StartTimerAutomatically != x.StartTimerAutomatically {
		x.StartTimerAutomatically = *input.StartTimerAutomatically
	}
}

var _ validation.ValidatableWithContext = (*RecipeStepCreationRequestInput)(nil)

// ValidateWithContext validates a RecipeStepCreationRequestInput.
func (x *RecipeStepCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	var err *multierror.Error

	if len(x.Instruments) == 0 && len(x.Vessels) == 0 {
		err = multierror.Append(err, errOneInstrumentOrVesselRequired)
	}

	validationErr := validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.Products, validation.Required),
	)

	if validationErr != nil {
		err = multierror.Append(err, validationErr)
	}

	return err.ErrorOrNil()
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

var _ validation.ValidatableWithContext = (*RecipeStepUpdateRequestInput)(nil)

// ValidateWithContext validates a RecipeStepUpdateRequestInput.
func (x *RecipeStepUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Index, validation.Required),
		validation.Field(&x.Preparation, validation.Required),
		validation.Field(&x.EstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.TemperatureInCelsius, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
