package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanOptionDataType indicates an event is related to a meal plan option.
	MealPlanOptionDataType dataType = "meal_plan_option"
)

func init() {
	gob.Register(new(MealPlanOption))
	gob.Register(new(MealPlanOptionList))
	gob.Register(new(MealPlanOptionCreationRequestInput))
	gob.Register(new(MealPlanOptionUpdateRequestInput))
}

type (
	// MealPlanOption represents a meal plan option.
	MealPlanOption struct {
		_                 struct{}
		LastUpdatedOn     *uint64 `json:"lastUpdatedOn"`
		ArchivedOn        *uint64 `json:"archivedOn"`
		RecipeID          string  `json:"recipeID"`
		Notes             string  `json:"notes"`
		ID                string  `json:"id"`
		BelongsToMealPlan string  `json:"belongsToMealPlan"`
		CreatedOn         uint64  `json:"createdOn"`
		DayOfWeek         uint8   `json:"dayOfWeek"`
	}

	// MealPlanOptionList represents a list of meal plan options.
	MealPlanOptionList struct {
		_               struct{}
		MealPlanOptions []*MealPlanOption `json:"mealPlanOptions"`
		Pagination
	}

	// MealPlanOptionCreationRequestInput represents what a user could set as input for creating meal plan options.
	MealPlanOptionCreationRequestInput struct {
		_                 struct{}
		ID                string `json:"-"`
		RecipeID          string `json:"recipeID"`
		Notes             string `json:"notes"`
		BelongsToMealPlan string `json:"-"`
		DayOfWeek         uint8  `json:"dayOfWeek"`
	}

	// MealPlanOptionDatabaseCreationInput represents what a user could set as input for creating meal plan options.
	MealPlanOptionDatabaseCreationInput struct {
		_                 struct{}
		ID                string `json:"id"`
		RecipeID          string `json:"recipeID"`
		Notes             string `json:"notes"`
		BelongsToMealPlan string `json:"belongsToMealPlan"`
		DayOfWeek         uint8  `json:"dayOfWeek"`
	}

	// MealPlanOptionUpdateRequestInput represents what a user could set as input for updating meal plan options.
	MealPlanOptionUpdateRequestInput struct {
		_                 struct{}
		RecipeID          string `json:"recipeID"`
		Notes             string `json:"notes"`
		BelongsToMealPlan string `json:"-"`
		DayOfWeek         uint8  `json:"dayOfWeek"`
	}

	// MealPlanOptionDataManager describes a structure capable of storing meal plan options permanently.
	MealPlanOptionDataManager interface {
		MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanOptionID string) (bool, error)
		GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*MealPlanOption, error)
		GetTotalMealPlanOptionCount(ctx context.Context) (uint64, error)
		GetMealPlanOptions(ctx context.Context, mealPlanID string, filter *QueryFilter) (*MealPlanOptionList, error)
		GetMealPlanOptionsWithIDs(ctx context.Context, mealPlanID string, limit uint8, ids []string) ([]*MealPlanOption, error)
		CreateMealPlanOption(ctx context.Context, input *MealPlanOptionDatabaseCreationInput) (*MealPlanOption, error)
		UpdateMealPlanOption(ctx context.Context, updated *MealPlanOption) error
		ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error
	}

	// MealPlanOptionDataService describes a structure capable of serving traffic related to meal plan options.
	MealPlanOptionDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an MealPlanOptionUpdateRequestInput with a meal plan option.
func (x *MealPlanOption) Update(input *MealPlanOptionUpdateRequestInput) {
	if input.DayOfWeek != 0 && input.DayOfWeek != x.DayOfWeek {
		x.DayOfWeek = input.DayOfWeek
	}

	if input.RecipeID != "" && input.RecipeID != x.RecipeID {
		x.RecipeID = input.RecipeID
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

var _ validation.ValidatableWithContext = (*MealPlanOptionCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionCreationRequestInput.
func (x *MealPlanOptionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.DayOfWeek, validation.Required),
		validation.Field(&x.RecipeID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanOptionDatabaseCreationInput.
func (x *MealPlanOptionDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.DayOfWeek, validation.Required),
		validation.Field(&x.RecipeID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

// MealPlanOptionDatabaseCreationInputFromMealPlanOptionCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanOptionDatabaseCreationInputFromMealPlanOptionCreationInput(input *MealPlanOptionCreationRequestInput) *MealPlanOptionDatabaseCreationInput {
	x := &MealPlanOptionDatabaseCreationInput{
		BelongsToMealPlan: input.BelongsToMealPlan,
		DayOfWeek:         input.DayOfWeek,
		RecipeID:          input.RecipeID,
		Notes:             input.Notes,
	}

	return x
}

var _ validation.ValidatableWithContext = (*MealPlanOptionUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionUpdateRequestInput.
func (x *MealPlanOptionUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.DayOfWeek, validation.Required),
		validation.Field(&x.RecipeID, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.BelongsToMealPlan, validation.Required),
	)
}
