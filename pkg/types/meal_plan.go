package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanDataType indicates an event is related to a meal plan.
	MealPlanDataType dataType = "meal_plan"
)

func init() {
	gob.Register(new(MealPlan))
	gob.Register(new(MealPlanList))
	gob.Register(new(MealPlanCreationRequestInput))
	gob.Register(new(MealPlanUpdateRequestInput))
}

type (
	// MealPlan represents a meal plan.
	MealPlan struct {
		_                struct{}
		ArchivedOn       *uint64 `json:"archivedOn"`
		LastUpdatedOn    *uint64 `json:"lastUpdatedOn"`
		State            string  `json:"state"`
		ID               string  `json:"id"`
		BelongsToAccount string  `json:"belongsToAccount"`
		StartsAt         uint64  `json:"startsAt"`
		EndsAt           uint64  `json:"endsAt"`
		CreatedOn        uint64  `json:"createdOn"`
	}

	// MealPlanList represents a list of meal plans.
	MealPlanList struct {
		_         struct{}
		MealPlans []*MealPlan `json:"mealPlans"`
		Pagination
	}

	// MealPlanCreationRequestInput represents what a user could set as input for creating meal plans.
	MealPlanCreationRequestInput struct {
		_                struct{}
		ID               string `json:"-"`
		State            string `json:"state"`
		BelongsToAccount string `json:"-"`
		StartsAt         uint64 `json:"startsAt"`
		EndsAt           uint64 `json:"endsAt"`
	}

	// MealPlanDatabaseCreationInput represents what a user could set as input for creating meal plans.
	MealPlanDatabaseCreationInput struct {
		_                struct{}
		ID               string `json:"id"`
		State            string `json:"state"`
		BelongsToAccount string `json:"belongsToAccount"`
		StartsAt         uint64 `json:"startsAt"`
		EndsAt           uint64 `json:"endsAt"`
	}

	// MealPlanUpdateRequestInput represents what a user could set as input for updating meal plans.
	MealPlanUpdateRequestInput struct {
		_                struct{}
		State            string `json:"state"`
		BelongsToAccount string `json:"-"`
		StartsAt         uint64 `json:"startsAt"`
		EndsAt           uint64 `json:"endsAt"`
	}

	// MealPlanDataManager describes a structure capable of storing meal plans permanently.
	MealPlanDataManager interface {
		MealPlanExists(ctx context.Context, mealPlanID string) (bool, error)
		GetMealPlan(ctx context.Context, mealPlanID string) (*MealPlan, error)
		GetTotalMealPlanCount(ctx context.Context) (uint64, error)
		GetMealPlans(ctx context.Context, filter *QueryFilter) (*MealPlanList, error)
		GetMealPlansWithIDs(ctx context.Context, accountID string, limit uint8, ids []string) ([]*MealPlan, error)
		CreateMealPlan(ctx context.Context, input *MealPlanDatabaseCreationInput) (*MealPlan, error)
		UpdateMealPlan(ctx context.Context, updated *MealPlan) error
		ArchiveMealPlan(ctx context.Context, mealPlanID, accountID string) error
	}

	// MealPlanDataService describes a structure capable of serving traffic related to meal plans.
	MealPlanDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an MealPlanUpdateRequestInput with a meal plan.
func (x *MealPlan) Update(input *MealPlanUpdateRequestInput) {
	if input.State != "" && input.State != x.State {
		x.State = input.State
	}

	if input.StartsAt != 0 && input.StartsAt != x.StartsAt {
		x.StartsAt = input.StartsAt
	}

	if input.EndsAt != 0 && input.EndsAt != x.EndsAt {
		x.EndsAt = input.EndsAt
	}
}

var _ validation.ValidatableWithContext = (*MealPlanCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanCreationRequestInput.
func (x *MealPlanCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.State, validation.Required),
		validation.Field(&x.StartsAt, validation.Required),
		validation.Field(&x.EndsAt, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanDatabaseCreationInput.
func (x *MealPlanDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.State, validation.Required),
		validation.Field(&x.StartsAt, validation.Required),
		validation.Field(&x.EndsAt, validation.Required),
		validation.Field(&x.BelongsToAccount, validation.Required),
	)
}

// MealPlanDatabaseCreationInputFromMealPlanCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanDatabaseCreationInputFromMealPlanCreationInput(input *MealPlanCreationRequestInput) *MealPlanDatabaseCreationInput {
	x := &MealPlanDatabaseCreationInput{
		State:    input.State,
		StartsAt: input.StartsAt,
		EndsAt:   input.EndsAt,
	}

	return x
}

var _ validation.ValidatableWithContext = (*MealPlanUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanUpdateRequestInput.
func (x *MealPlanUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.State, validation.Required),
		validation.Field(&x.StartsAt, validation.Required),
		validation.Field(&x.EndsAt, validation.Required),
	)
}
