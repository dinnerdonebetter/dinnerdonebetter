package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanOptionDataType indicates an event is related to a meal plan option.
	MealPlanOptionDataType dataType = "meal_plan_option"

	// MealPlanOptionCreatedCustomerEventType indicates a meal plan option was created.
	MealPlanOptionCreatedCustomerEventType CustomerEventType = "meal_plan_option_created"
	// MealPlanOptionUpdatedCustomerEventType indicates a meal plan option was updated.
	MealPlanOptionUpdatedCustomerEventType CustomerEventType = "meal_plan_option_updated"
	// MealPlanOptionArchivedCustomerEventType indicates a meal plan option was archived.
	MealPlanOptionArchivedCustomerEventType CustomerEventType = "meal_plan_option_archived"
	// MealPlanOptionFinalizedCreatedCustomerEventType indicates a meal plan option was created.
	MealPlanOptionFinalizedCreatedCustomerEventType CustomerEventType = "meal_plan_option_finalized"

	// BreakfastMealName represents breakfast.
	BreakfastMealName MealName = "breakfast"
	// SecondBreakfastMealName represents second breakfast.
	SecondBreakfastMealName MealName = "second_breakfast"
	// BrunchMealName represents brunch.
	BrunchMealName MealName = "brunch"
	// LunchMealName represents lunch.
	LunchMealName MealName = "lunch"
	// SupperMealName represents supper.
	SupperMealName MealName = "supper"
	// DinnerMealName represents dinner.
	DinnerMealName MealName = "dinner"
)

func init() {
	gob.Register(new(MealPlanOption))
	gob.Register(new(MealPlanOptionList))
	gob.Register(new(MealPlanOptionCreationRequestInput))
	gob.Register(new(MealPlanOptionUpdateRequestInput))
}

type (
	// MealName is an enumeration for meal names.
	MealName string

	// MealPlanOption represents a meal plan option.
	MealPlanOption struct {
		_                 struct{}
		ArchivedAt        *uint64               `json:"archivedAt"`
		LastUpdatedAt     *uint64               `json:"lastUpdatedAt"`
		ID                string                `json:"id"`
		BelongsToMealPlan string                `json:"belongsToMealPlan"`
		Notes             string                `json:"notes"`
		MealName          MealName              `json:"mealName"`
		AssignedCook      *string               `json:"assignedCook"`
		Votes             []*MealPlanOptionVote `json:"votes"`
		Meal              Meal                  `json:"meal"`
		CreatedAt         uint64                `json:"createdAt"`
		Day               time.Weekday          `json:"day"`
		Chosen            bool                  `json:"chosen"`
		TieBroken         bool                  `json:"tieBroken"`
	}

	// MealPlanOptionList represents a list of meal plan options.
	MealPlanOptionList struct {
		_               struct{}
		MealPlanOptions []*MealPlanOption `json:"data"`
		Pagination
	}

	// MealPlanOptionCreationRequestInput represents what a user could set as input for creating meal plan options.
	MealPlanOptionCreationRequestInput struct {
		_                 struct{}
		ID                string       `json:"-"`
		MealID            string       `json:"mealID"`
		Notes             string       `json:"notes"`
		AssignedCook      *string      `json:"assignedCook"`
		MealName          MealName     `json:"mealName"`
		BelongsToMealPlan string       `json:"-"`
		Day               time.Weekday `json:"day"`
	}

	// MealPlanOptionDatabaseCreationInput represents what a user could set as input for creating meal plan options.
	MealPlanOptionDatabaseCreationInput struct {
		_                 struct{}
		ID                string       `json:"id"`
		MealID            string       `json:"mealID"`
		Notes             string       `json:"notes"`
		AssignedCook      *string      `json:"assignedCook"`
		MealName          MealName     `json:"mealName"`
		BelongsToMealPlan string       `json:"belongsToMealPlan"`
		Day               time.Weekday `json:"day"`
	}

	// MealPlanOptionUpdateRequestInput represents what a user could set as input for updating meal plan options.
	MealPlanOptionUpdateRequestInput struct {
		_                 struct{}
		MealID            *string       `json:"mealID"`
		Notes             *string       `json:"notes"`
		AssignedCook      *string       `json:"assignedCook"`
		MealName          *MealName     `json:"mealName"`
		BelongsToMealPlan *string       `json:"-"`
		Day               *time.Weekday `json:"day"`
	}

	// MealPlanOptionDataManager describes a structure capable of storing meal plan options permanently.
	MealPlanOptionDataManager interface {
		MealPlanOptionExists(ctx context.Context, mealPlanID, mealPlanOptionID string) (bool, error)
		GetMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (*MealPlanOption, error)
		GetMealPlanOptions(ctx context.Context, mealPlanID string, filter *QueryFilter) (*MealPlanOptionList, error)
		GetMealPlanOptionsWithIDs(ctx context.Context, mealPlanID string, limit uint8, ids []string) ([]*MealPlanOption, error)
		CreateMealPlanOption(ctx context.Context, input *MealPlanOptionDatabaseCreationInput) (*MealPlanOption, error)
		UpdateMealPlanOption(ctx context.Context, updated *MealPlanOption) error
		ArchiveMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) error
		FinalizeMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID, householdID string) (changed bool, err error)
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
	if input.Day != nil && *input.Day != x.Day {
		x.Day = *input.Day
	}

	if input.MealID != nil && *input.MealID != x.Meal.ID {
		// we should do something better here
		x.Meal = Meal{ID: *input.MealID}
	}

	if input.MealName != nil && *input.MealName != x.MealName {
		// we should do something better here
		x.MealName = *input.MealName
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.AssignedCook != nil && input.AssignedCook != x.AssignedCook {
		x.AssignedCook = input.AssignedCook
	}
}

var _ validation.ValidatableWithContext = (*MealPlanOptionCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionCreationRequestInput.
func (x *MealPlanOptionCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealID, validation.Required),
		validation.Field(&x.MealName, validation.In(
			BreakfastMealName,
			SecondBreakfastMealName,
			BrunchMealName,
			LunchMealName,
			SupperMealName,
			DinnerMealName,
		)),
		validation.Field(&x.Day, validation.In(
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
			time.Saturday,
			time.Sunday,
		)),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanOptionDatabaseCreationInput.
func (x *MealPlanOptionDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.BelongsToMealPlan, validation.Required),
		validation.Field(&x.MealName, validation.Required),
		validation.Field(&x.MealID, validation.Required),
		validation.Field(&x.Day, validation.In(
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
			time.Saturday,
			time.Sunday,
		)),
	)
}

// MealPlanOptionUpdateRequestInputFromMealPlanOption creates a DatabaseCreationInput from a CreationInput.
func MealPlanOptionUpdateRequestInputFromMealPlanOption(input *MealPlanOption) *MealPlanOptionUpdateRequestInput {
	x := &MealPlanOptionUpdateRequestInput{
		MealID:            &input.Meal.ID,
		Notes:             &input.Notes,
		MealName:          &input.MealName,
		BelongsToMealPlan: &input.BelongsToMealPlan,
		Day:               &input.Day,
	}

	return x
}

// MealPlanOptionDatabaseCreationInputFromMealPlanOptionCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanOptionDatabaseCreationInputFromMealPlanOptionCreationInput(input *MealPlanOptionCreationRequestInput) *MealPlanOptionDatabaseCreationInput {
	x := &MealPlanOptionDatabaseCreationInput{
		BelongsToMealPlan: input.BelongsToMealPlan,
		Day:               input.Day,
		MealName:          input.MealName,
		MealID:            input.MealID,
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
		validation.Field(&x.MealID, validation.Required),
		validation.Field(&x.BelongsToMealPlan, validation.Required),
		validation.Field(&x.MealName, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.BelongsToMealPlan, validation.Required),
		validation.Field(&x.Day, validation.In(
			time.Monday,
			time.Tuesday,
			time.Wednesday,
			time.Thursday,
			time.Friday,
			time.Saturday,
			time.Sunday,
		)),
	)
}
