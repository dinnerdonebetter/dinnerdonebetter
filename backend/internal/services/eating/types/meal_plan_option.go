package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type (
	// MealPlanOption represents a meal plan option.
	MealPlanOption struct {
		_ struct{} `json:"-"`

		CreatedAt              time.Time             `json:"createdAt"`
		LastUpdatedAt          *time.Time            `json:"lastUpdatedAt"`
		AssignedCook           *string               `json:"assignedCook"`
		ArchivedAt             *time.Time            `json:"archivedAt"`
		AssignedDishwasher     *string               `json:"assignedDishwasher"`
		Notes                  string                `json:"notes"`
		BelongsToMealPlanEvent string                `json:"belongsToMealPlanEvent"`
		ID                     string                `json:"id"`
		Votes                  []*MealPlanOptionVote `json:"votes"`
		Meal                   Meal                  `json:"meal"`
		MealScale              float32               `json:"mealScale"`
		Chosen                 bool                  `json:"chosen"`
		TieBroken              bool                  `json:"tieBroken"`
	}

	// MealPlanOptionCreationRequestInput represents what a user could set as input for creating meal plan options.
	MealPlanOptionCreationRequestInput struct {
		_ struct{} `json:"-"`

		AssignedCook       *string `json:"assignedCook"`
		AssignedDishwasher *string `json:"assignedDishwasher"`
		MealID             string  `json:"mealID"`
		Notes              string  `json:"notes"`
		MealScale          float32 `json:"mealScale"`
	}

	// MealPlanOptionDatabaseCreationInput represents what a user could set as input for creating meal plan options.
	MealPlanOptionDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		ID                     string  `json:"-"`
		MealID                 string  `json:"-"`
		Notes                  string  `json:"-"`
		AssignedCook           *string `json:"-"`
		AssignedDishwasher     *string `json:"-"`
		BelongsToMealPlanEvent string  `json:"-"`
		MealScale              float32 `json:"-"`
	}

	// MealPlanOptionUpdateRequestInput represents what a user could set as input for updating meal plan options.
	MealPlanOptionUpdateRequestInput struct {
		_ struct{} `json:"-"`

		MealID                 *string  `json:"mealID,omitempty"`
		Notes                  *string  `json:"notes,omitempty"`
		AssignedCook           *string  `json:"assignedCook,omitempty"`
		AssignedDishwasher     *string  `json:"assignedDishwasher,omitempty"`
		MealScale              *float32 `json:"mealScale,omitempty"`
		BelongsToMealPlanEvent *string  `json:"-"`
	}

	// MealPlanOptionDataManager describes a structure capable of storing meal plan options permanently.
	MealPlanOptionDataManager interface {
		MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (bool, error)
		GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) (*MealPlanOption, error)
		GetMealPlanOptions(ctx context.Context, mealPlanID, mealPlanEventID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[MealPlanOption], error)
		CreateMealPlanOption(ctx context.Context, input *MealPlanOptionDatabaseCreationInput) (*MealPlanOption, error)
		UpdateMealPlanOption(ctx context.Context, updated *MealPlanOption) error
		ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID string) error
		FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanEventID, mealPlanOptionID, householdID string) (changed bool, err error)
	}

	// MealPlanOptionDataService describes a structure capable of serving traffic related to meal plan options.
	MealPlanOptionDataService interface {
		ListMealPlanOptionHandler(http.ResponseWriter, *http.Request)
		CreateMealPlanOptionHandler(http.ResponseWriter, *http.Request)
		ReadMealPlanOptionHandler(http.ResponseWriter, *http.Request)
		UpdateMealPlanOptionHandler(http.ResponseWriter, *http.Request)
		ArchiveMealPlanOptionHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealPlanOptionUpdateRequestInput with a meal plan option.
func (x *MealPlanOption) Update(input *MealPlanOptionUpdateRequestInput) {
	if input.MealID != nil && *input.MealID != x.Meal.ID {
		// we should do something better here
		x.Meal = Meal{ID: *input.MealID}
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.AssignedCook != nil && input.AssignedCook != x.AssignedCook {
		x.AssignedCook = input.AssignedCook
	}

	if input.AssignedDishwasher != nil && input.AssignedDishwasher != x.AssignedDishwasher {
		x.AssignedDishwasher = input.AssignedDishwasher
	}

	if input.MealScale != nil && *input.MealScale != x.MealScale {
		x.MealScale = *input.MealScale
	}
}

var _ validation.ValidatableWithContext = (*MealPlanOptionCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionCreationRequestInput.
func (x *MealPlanOptionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanOptionDatabaseCreationInput.
func (x *MealPlanOptionDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToMealPlanEvent, validation.Required),
		validation.Field(&x.MealID, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionUpdateRequestInput.
func (x *MealPlanOptionUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealID, validation.Required),
		validation.Field(&x.BelongsToMealPlanEvent, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.BelongsToMealPlanEvent, validation.Required),
	)
}
