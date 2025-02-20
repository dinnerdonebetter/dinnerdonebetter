package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanEventCreatedServiceEventType indicates a meal plan was created.
	MealPlanEventCreatedServiceEventType = "meal_plan_event_created"
	// MealPlanEventUpdatedServiceEventType indicates a meal plan was updated.
	MealPlanEventUpdatedServiceEventType = "meal_plan_event_updated"
	// MealPlanEventArchivedServiceEventType indicates a meal plan was archived.
	/* #nosec G101 */
	MealPlanEventArchivedServiceEventType = "meal_plan_event_archived"

	// BreakfastMealName represents breakfast.
	BreakfastMealName = "breakfast"
	// SecondBreakfastMealName represents second breakfast.
	SecondBreakfastMealName = "second_breakfast"
	// BrunchMealName represents brunch.
	BrunchMealName = "brunch"
	// LunchMealName represents lunch.
	LunchMealName = "lunch"
	// SupperMealName represents supper.
	SupperMealName = "supper"
	// DinnerMealName represents dinner.
	DinnerMealName = "dinner"
)

func init() {
	gob.Register(new(MealPlanEvent))
	gob.Register(new(MealPlanEventCreationRequestInput))
	gob.Register(new(MealPlanEventUpdateRequestInput))
}

type (
	// MealPlanEvent represents a meal plan.
	MealPlanEvent struct {
		_ struct{} `json:"-"`

		CreatedAt         time.Time         `json:"createdAt"`
		StartsAt          time.Time         `json:"startsAt"`
		EndsAt            time.Time         `json:"endsAt"`
		ArchivedAt        *time.Time        `json:"archivedAt"`
		LastUpdatedAt     *time.Time        `json:"lastUpdatedAt"`
		MealName          string            `json:"mealName"`
		Notes             string            `json:"notes"`
		BelongsToMealPlan string            `json:"belongsToMealPlan"`
		ID                string            `json:"id"`
		Options           []*MealPlanOption `json:"options"`
	}

	// MealPlanEventCreationRequestInput represents what a user could set as input for creating meal plans.
	MealPlanEventCreationRequestInput struct {
		_ struct{} `json:"-"`

		EndsAt   time.Time                             `json:"endsAt"`
		StartsAt time.Time                             `json:"startsAt"`
		Notes    string                                `json:"notes"`
		MealName string                                `json:"mealName"`
		Options  []*MealPlanOptionCreationRequestInput `json:"options"`
	}

	// MealPlanEventDatabaseCreationInput represents what a user could set as input for creating meal plans.
	MealPlanEventDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		StartsAt          time.Time                              `json:"-"`
		EndsAt            time.Time                              `json:"-"`
		BelongsToMealPlan string                                 `json:"-"`
		Notes             string                                 `json:"-"`
		MealName          string                                 `json:"-"`
		ID                string                                 `json:"-"`
		Options           []*MealPlanOptionDatabaseCreationInput `json:"-"`
	}

	// MealPlanEventUpdateRequestInput represents what a user could set as input for updating meal plans.
	MealPlanEventUpdateRequestInput struct {
		_ struct{} `json:"-"`

		Notes             *string    `json:"notes,omitempty"`
		StartsAt          *time.Time `json:"startsAt,omitempty"`
		MealName          *string    `json:"mealName,omitempty"`
		EndsAt            *time.Time `json:"endsAt,omitempty"`
		BelongsToMealPlan string     `json:"-"`
	}

	// MealPlanEventDataManager describes a structure capable of storing meal plans permanently.
	MealPlanEventDataManager interface {
		MealPlanEventExists(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error)
		GetMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) (*MealPlanEvent, error)
		GetMealPlanEvents(ctx context.Context, mealPlanID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[MealPlanEvent], error)
		MealPlanEventIsEligibleForVoting(ctx context.Context, mealPlanID, mealPlanEventID string) (bool, error)
		CreateMealPlanEvent(ctx context.Context, input *MealPlanEventDatabaseCreationInput) (*MealPlanEvent, error)
		UpdateMealPlanEvent(ctx context.Context, updated *MealPlanEvent) error
		ArchiveMealPlanEvent(ctx context.Context, mealPlanID, mealPlanEventID string) error
	}

	// MealPlanEventDataService describes a structure capable of serving traffic related to meal plans.
	MealPlanEventDataService interface {
		ListMealPlanEventHandler(http.ResponseWriter, *http.Request)
		CreateMealPlanEventHandler(http.ResponseWriter, *http.Request)
		ReadMealPlanEventHandler(http.ResponseWriter, *http.Request)
		UpdateMealPlanEventHandler(http.ResponseWriter, *http.Request)
		ArchiveMealPlanEventHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealPlanEventUpdateRequestInput with a meal plan.
func (x *MealPlanEvent) Update(input *MealPlanEventUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}

	if input.StartsAt != nil && *input.StartsAt != x.StartsAt {
		x.StartsAt = *input.StartsAt
	}

	if input.EndsAt != nil && *input.EndsAt != x.EndsAt {
		x.EndsAt = *input.EndsAt
	}

	if input.MealName != nil && *input.MealName != x.MealName {
		x.MealName = *input.MealName
	}
}

var _ validation.ValidatableWithContext = (*MealPlanEventCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanEventCreationRequestInput.
func (x *MealPlanEventCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	if x.StartsAt == x.EndsAt || x.StartsAt.After(x.EndsAt) {
		return ErrStartsAfterItEnds
	}

	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.StartsAt, validation.Required),
		validation.Field(&x.EndsAt, validation.Required),
		validation.Field(&x.MealName, validation.In(
			BreakfastMealName,
			SecondBreakfastMealName,
			BrunchMealName,
			LunchMealName,
			SupperMealName,
			DinnerMealName,
		)),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanEventDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanEventDatabaseCreationInput.
func (x *MealPlanEventDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.StartsAt, validation.Required),
		validation.Field(&x.EndsAt, validation.Required),
		validation.Field(&x.MealName, validation.In(
			BreakfastMealName,
			SecondBreakfastMealName,
			BrunchMealName,
			LunchMealName,
			SupperMealName,
			DinnerMealName,
		)),
		validation.Field(&x.BelongsToMealPlan, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanEventUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanEventUpdateRequestInput.
func (x *MealPlanEventUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.StartsAt, validation.Required),
		validation.Field(&x.EndsAt, validation.Required),
		validation.Field(&x.MealName, validation.In(
			BreakfastMealName,
			SecondBreakfastMealName,
			BrunchMealName,
			LunchMealName,
			SupperMealName,
			DinnerMealName,
		)),
	)
}
