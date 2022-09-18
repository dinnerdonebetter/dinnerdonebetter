package types

import (
	"context"
	"encoding/gob"
	"errors"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanDataType indicates an event is related to a meal plan.
	MealPlanDataType dataType = "meal_plan"

	// MealPlanCreatedCustomerEventType indicates a meal plan was created.
	MealPlanCreatedCustomerEventType CustomerEventType = "meal_plan_created"
	// MealPlanUpdatedCustomerEventType indicates a meal plan was updated.
	MealPlanUpdatedCustomerEventType CustomerEventType = "meal_plan_updated"
	// MealPlanArchivedCustomerEventType indicates a meal plan was archived.
	MealPlanArchivedCustomerEventType CustomerEventType = "meal_plan_archived"
	// MealPlanFinalizedCustomerEventType indicates a meal plan was finalized.
	MealPlanFinalizedCustomerEventType CustomerEventType = "meal_plan_finalized"

	// AwaitingVotesMealPlanStatus indicates a household invitation is pending.
	AwaitingVotesMealPlanStatus MealPlanStatus = "awaiting_votes"
	// FinalizedMealPlanStatus indicates a household invitation was accepted.
	FinalizedMealPlanStatus MealPlanStatus = "finalized"
)

func init() {
	gob.Register(new(MealPlan))
	gob.Register(new(MealPlanList))
	gob.Register(new(MealPlanCreationRequestInput))
	gob.Register(new(MealPlanUpdateRequestInput))
}

type (
	// MealPlanStatus is the type to use/compare against when checking meal plan status.
	MealPlanStatus string

	// MealPlan represents a meal plan.
	MealPlan struct {
		_                  struct{}
		CreatedAt          time.Time        `json:"createdAt"`
		VotingDeadline     time.Time        `json:"votingDeadline"`
		ArchivedAt         *time.Time       `json:"archivedAt"`
		LastUpdatedAt      *time.Time       `json:"lastUpdatedAt"`
		Notes              string           `json:"notes"`
		BelongsToHousehold string           `json:"belongsToHousehold"`
		Status             MealPlanStatus   `json:"status"`
		ID                 string           `json:"id"`
		Events             []*MealPlanEvent `json:"events"`
	}

	// MealPlanList represents a list of meal plans.
	MealPlanList struct {
		_         struct{}
		MealPlans []*MealPlan `json:"data"`
		Pagination
	}

	// MealPlanCreationRequestInput represents what a user could set as input for creating meal plans.
	MealPlanCreationRequestInput struct {
		_                  struct{}
		VotingDeadline     time.Time                            `json:"votingDeadline"`
		BelongsToHousehold string                               `json:"-"`
		Notes              string                               `json:"notes"`
		ID                 string                               `json:"-"`
		Events             []*MealPlanEventCreationRequestInput `json:"events"`
	}

	// MealPlanDatabaseCreationInput represents what a user could set as input for creating meal plans.
	MealPlanDatabaseCreationInput struct {
		_                  struct{}
		VotingDeadline     time.Time                             `json:"votingDeadline"`
		BelongsToHousehold string                                `json:"belongsToHousehold"`
		Notes              string                                `json:"notes"`
		ID                 string                                `json:"id"`
		Events             []*MealPlanEventDatabaseCreationInput `json:"events"`
	}

	// MealPlanUpdateRequestInput represents what a user could set as input for updating meal plans.
	MealPlanUpdateRequestInput struct {
		_                  struct{}
		BelongsToHousehold *string    `json:"-"`
		Notes              *string    `json:"notes"`
		VotingDeadline     *time.Time `json:"votingDeadline"`
	}

	// FinalizedMealPlanDatabaseResult represents what is returned by the above query.
	FinalizedMealPlanDatabaseResult struct {
		MealPlanID       string
		MealPlanEventID  string
		MealPlanOptionID string
		MealID           string
		RecipeIDs        []string
	}

	// MealPlanDataManager describes a structure capable of storing meal plans permanently.
	MealPlanDataManager interface {
		MealPlanExists(ctx context.Context, mealPlanID, householdID string) (bool, error)
		GetMealPlan(ctx context.Context, mealPlanID, householdID string) (*MealPlan, error)
		GetMealPlans(ctx context.Context, householdID string, filter *QueryFilter) (*MealPlanList, error)
		CreateMealPlan(ctx context.Context, input *MealPlanDatabaseCreationInput) (*MealPlan, error)
		UpdateMealPlan(ctx context.Context, updated *MealPlan) error
		ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error
		AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (changed bool, err error)
		GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*FinalizedMealPlanDatabaseResult, error)
		GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*MealPlan, error)
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
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var errTooFewUniqueMeals = errors.New("too many instances of the same meal")
var errInvalidVotingDeadline = errors.New("invalid voting deadline")

var _ validation.ValidatableWithContext = (*MealPlanCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanCreationRequestInput.
func (x *MealPlanCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	if time.Now().After(x.VotingDeadline) {
		return errInvalidVotingDeadline
	}

	return nil
}

var _ validation.ValidatableWithContext = (*MealPlanDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanDatabaseCreationInput.
func (x *MealPlanDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.VotingDeadline, validation.Required),
		validation.Field(&x.BelongsToHousehold, validation.Required),
	)
}

// MealPlanUpdateRequestInputFromMealPlan creates a DatabaseCreationInput from a CreationInput.
func MealPlanUpdateRequestInputFromMealPlan(input *MealPlan) *MealPlanUpdateRequestInput {
	x := &MealPlanUpdateRequestInput{
		BelongsToHousehold: &input.BelongsToHousehold,
		Notes:              &input.Notes,
		VotingDeadline:     &input.VotingDeadline,
	}

	return x
}

// MealPlanDatabaseCreationInputFromMealPlanCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanDatabaseCreationInputFromMealPlanCreationInput(input *MealPlanCreationRequestInput) *MealPlanDatabaseCreationInput {
	events := []*MealPlanEventDatabaseCreationInput{}
	for _, e := range input.Events {
		events = append(events, MealPlanEventDatabaseCreationInputFromMealPlanEventCreationRequestInput(e))
	}

	x := &MealPlanDatabaseCreationInput{
		Notes:          input.Notes,
		VotingDeadline: input.VotingDeadline,
		Events:         events,
	}

	return x
}

var _ validation.ValidatableWithContext = (*MealPlanUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanUpdateRequestInput.
func (x *MealPlanUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.VotingDeadline, validation.Required),
	)
}
