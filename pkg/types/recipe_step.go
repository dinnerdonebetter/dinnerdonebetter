package types

import (
	"context"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// RecipeStep represents a recipe step.
	RecipeStep struct {
		LastUpdatedOn             *uint64                 `json:"lastUpdatedOn"`
		TemperatureInCelsius      *uint16                 `json:"temperatureInCelsius"`
		ArchivedOn                *uint64                 `json:"archivedOn"`
		ExternalID                string                  `json:"externalID"`
		Why                       string                  `json:"why"`
		Notes                     string                  `json:"notes"`
		Ingredients               []*RecipeStepIngredient `json:"ingredients"`
		PrerequisiteStep          uint64                  `json:"prerequisiteStep"`
		ID                        uint64                  `json:"id"`
		Index                     uint                    `json:"index"`
		CreatedOn                 uint64                  `json:"createdOn"`
		BelongsToRecipe           uint64                  `json:"belongsToRecipe"`
		PreparationID             uint64                  `json:"preparationID"`
		MaxEstimatedTimeInSeconds uint32                  `json:"maxEstimatedTimeInSeconds"`
		MinEstimatedTimeInSeconds uint32                  `json:"minEstimatedTimeInSeconds"`
	}

	// FullRecipeStep represents a recipe step.
	FullRecipeStep struct {
		TemperatureInCelsius      *uint16                     `json:"temperatureInCelsius"`
		ArchivedOn                *uint64                     `json:"archivedOn"`
		LastUpdatedOn             *uint64                     `json:"lastUpdatedOn"`
		ExternalID                string                      `json:"externalID"`
		Why                       string                      `json:"why"`
		Notes                     string                      `json:"notes"`
		Ingredients               []*FullRecipeStepIngredient `json:"ingredients"`
		Preparation               ValidPreparation            `json:"preparation"`
		PrerequisiteStep          uint64                      `json:"prerequisiteStep"`
		Index                     uint                        `json:"index"`
		CreatedOn                 uint64                      `json:"createdOn"`
		BelongsToRecipe           uint64                      `json:"belongsToRecipe"`
		ID                        uint64                      `json:"id"`
		MaxEstimatedTimeInSeconds uint32                      `json:"maxEstimatedTimeInSeconds"`
		MinEstimatedTimeInSeconds uint32                      `json:"minEstimatedTimeInSeconds"`
	}

	// RecipeStepList represents a list of recipe steps.
	RecipeStepList struct {
		RecipeSteps []*RecipeStep `json:"recipeSteps"`
		Pagination
	}

	// RecipeStepCreationInput represents what a user could set as input for creating recipe steps.
	RecipeStepCreationInput struct {
		TemperatureInCelsius      *uint16                              `json:"temperatureInCelsius"`
		Notes                     string                               `json:"notes"`
		Why                       string                               `json:"why"`
		Ingredients               []*RecipeStepIngredientCreationInput `json:"ingredients"`
		PrerequisiteStep          uint64                               `json:"prerequisiteStep"`
		Index                     uint                                 `json:"index"`
		PreparationID             uint64                               `json:"preparationID"`
		BelongsToRecipe           uint64                               `json:"-"`
		MaxEstimatedTimeInSeconds uint32                               `json:"maxEstimatedTimeInSeconds"`
		MinEstimatedTimeInSeconds uint32                               `json:"minEstimatedTimeInSeconds"`
	}

	// RecipeStepUpdateInput represents what a user could set as input for updating recipe steps.
	RecipeStepUpdateInput struct {
		TemperatureInCelsius      *uint16 `json:"temperatureInCelsius"`
		Notes                     string  `json:"notes"`
		Why                       string  `json:"why"`
		PrerequisiteStep          uint64  `json:"prerequisiteStep"`
		BelongsToRecipe           uint64  `json:"belongsToRecipe"`
		Index                     uint    `json:"index"`
		PreparationID             uint64  `json:"preparationID"`
		MaxEstimatedTimeInSeconds uint32  `json:"maxEstimatedTimeInSeconds"`
		MinEstimatedTimeInSeconds uint32  `json:"minEstimatedTimeInSeconds"`
	}

	// RecipeStepDataManager describes a structure capable of storing recipe steps permanently.
	RecipeStepDataManager interface {
		RecipeStepExists(ctx context.Context, recipeID, recipeStepID uint64) (bool, error)
		GetRecipeStep(ctx context.Context, recipeID, recipeStepID uint64) (*RecipeStep, error)
		GetAllRecipeStepsCount(ctx context.Context) (uint64, error)
		GetAllRecipeSteps(ctx context.Context, resultChannel chan []*RecipeStep, bucketSize uint16) error
		GetRecipeSteps(ctx context.Context, recipeID uint64, filter *QueryFilter) (*RecipeStepList, error)
		GetRecipeStepsWithIDs(ctx context.Context, recipeID uint64, limit uint8, ids []uint64) ([]*RecipeStep, error)
		CreateRecipeStep(ctx context.Context, input *RecipeStepCreationInput, createdByUser uint64) (*RecipeStep, error)
		UpdateRecipeStep(ctx context.Context, updated *RecipeStep, changedByUser uint64, changes []*FieldChangeSummary) error
		ArchiveRecipeStep(ctx context.Context, recipeID, recipeStepID, archivedBy uint64) error
		GetAuditLogEntriesForRecipeStep(ctx context.Context, recipeStepID uint64) ([]*AuditLogEntry, error)
	}

	// RecipeStepDataService describes a structure capable of serving traffic related to recipe steps.
	RecipeStepDataService interface {
		AuditEntryHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an RecipeStepUpdateInput with a recipe step.
func (x *RecipeStep) Update(input *RecipeStepUpdateInput) []*FieldChangeSummary {
	var out []*FieldChangeSummary

	if input.Index != 0 && input.Index != x.Index {
		out = append(out, &FieldChangeSummary{
			FieldName: "Index",
			OldValue:  x.Index,
			NewValue:  input.Index,
		})

		x.Index = input.Index
	}

	if input.PreparationID != 0 && input.PreparationID != x.PreparationID {
		out = append(out, &FieldChangeSummary{
			FieldName: "PreparationID",
			OldValue:  x.PreparationID,
			NewValue:  input.PreparationID,
		})

		x.PreparationID = input.PreparationID
	}

	if input.PrerequisiteStep != 0 && input.PrerequisiteStep != x.PrerequisiteStep {
		out = append(out, &FieldChangeSummary{
			FieldName: "PrerequisiteStep",
			OldValue:  x.PrerequisiteStep,
			NewValue:  input.PrerequisiteStep,
		})

		x.PrerequisiteStep = input.PrerequisiteStep
	}

	if input.MinEstimatedTimeInSeconds != 0 && input.MinEstimatedTimeInSeconds != x.MinEstimatedTimeInSeconds {
		out = append(out, &FieldChangeSummary{
			FieldName: "MinEstimatedTimeInSeconds",
			OldValue:  x.MinEstimatedTimeInSeconds,
			NewValue:  input.MinEstimatedTimeInSeconds,
		})

		x.MinEstimatedTimeInSeconds = input.MinEstimatedTimeInSeconds
	}

	if input.MaxEstimatedTimeInSeconds != 0 && input.MaxEstimatedTimeInSeconds != x.MaxEstimatedTimeInSeconds {
		out = append(out, &FieldChangeSummary{
			FieldName: "MaxEstimatedTimeInSeconds",
			OldValue:  x.MaxEstimatedTimeInSeconds,
			NewValue:  input.MaxEstimatedTimeInSeconds,
		})

		x.MaxEstimatedTimeInSeconds = input.MaxEstimatedTimeInSeconds
	}

	if input.TemperatureInCelsius != nil && (x.TemperatureInCelsius == nil || (*input.TemperatureInCelsius != 0 && *input.TemperatureInCelsius != *x.TemperatureInCelsius)) {
		out = append(out, &FieldChangeSummary{
			FieldName: "TemperatureInCelsius",
			OldValue:  x.TemperatureInCelsius,
			NewValue:  input.TemperatureInCelsius,
		})

		x.TemperatureInCelsius = input.TemperatureInCelsius
	}

	if input.Notes != x.Notes {
		out = append(out, &FieldChangeSummary{
			FieldName: "Notes",
			OldValue:  x.Notes,
			NewValue:  input.Notes,
		})

		x.Notes = input.Notes
	}

	if input.Why != x.Why {
		out = append(out, &FieldChangeSummary{
			FieldName: "Why",
			OldValue:  x.Why,
			NewValue:  input.Why,
		})

		x.Why = input.Why
	}

	return out
}

var _ validation.ValidatableWithContext = (*RecipeStepCreationInput)(nil)

// ValidateWithContext validates a RecipeStepCreationInput.
func (x *RecipeStepCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Index, validation.Required),
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.PrerequisiteStep, validation.Required),
		validation.Field(&x.MinEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.MaxEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.TemperatureInCelsius, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.Why, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*RecipeStepUpdateInput)(nil)

// ValidateWithContext validates a RecipeStepUpdateInput.
func (x *RecipeStepUpdateInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Index, validation.Required),
		validation.Field(&x.PreparationID, validation.Required),
		validation.Field(&x.PrerequisiteStep, validation.Required),
		validation.Field(&x.MinEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.MaxEstimatedTimeInSeconds, validation.Required),
		validation.Field(&x.TemperatureInCelsius, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.Why, validation.Required),
	)
}
