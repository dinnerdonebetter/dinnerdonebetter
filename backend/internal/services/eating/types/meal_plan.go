package types

import (
	"context"
	"net/http"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanElectionMethodSchulze is used to denote the Schulze election method.
	MealPlanElectionMethodSchulze = "schulze"
	// MealPlanElectionMethodInstantRunoff is used to denote the Instant Runoff election method.
	MealPlanElectionMethodInstantRunoff = "instant-runoff"

	// MealPlanStatusAwaitingVotes indicates a household invitation is pending.
	MealPlanStatusAwaitingVotes MealPlanStatus = "awaiting_votes"
	// MealPlanStatusFinalized indicates a household invitation was accepted.
	MealPlanStatusFinalized MealPlanStatus = "finalized"
)

type (
	// MealPlanStatus is the type to use/compare against when checking meal plan status.
	MealPlanStatus string

	// MealPlan represents a meal plan.
	MealPlan struct {
		_ struct{} `json:"-"`

		CreatedAt              time.Time        `json:"createdAt"`
		VotingDeadline         time.Time        `json:"votingDeadline"`
		ArchivedAt             *time.Time       `json:"archivedAt"`
		LastUpdatedAt          *time.Time       `json:"lastUpdatedAt"`
		ID                     string           `json:"id"`
		Status                 string           `json:"status"`
		Notes                  string           `json:"notes"`
		ElectionMethod         string           `json:"electionMethod"`
		BelongsToHousehold     string           `json:"belongsToHousehold"`
		CreatedByUser          string           `json:"createdBy"`
		Events                 []*MealPlanEvent `json:"events"`
		GroceryListInitialized bool             `json:"groceryListInitialized"`
		TasksCreated           bool             `json:"tasksCreated"`
	}

	// MealPlanCreationRequestInput represents what a user could set as input for creating meal plans.
	MealPlanCreationRequestInput struct {
		_ struct{} `json:"-"`

		VotingDeadline time.Time                            `json:"votingDeadline"`
		Notes          string                               `json:"notes"`
		ElectionMethod string                               `json:"electionMethod"`
		Events         []*MealPlanEventCreationRequestInput `json:"events"`
	}

	// MealPlanDatabaseCreationInput represents what a user could set as input for creating meal plans.
	MealPlanDatabaseCreationInput struct {
		_ struct{} `json:"-"`

		VotingDeadline     time.Time                             `json:"-"`
		BelongsToHousehold string                                `json:"-"`
		Notes              string                                `json:"-"`
		ID                 string                                `json:"-"`
		ElectionMethod     string                                `json:"-"`
		CreatedByUser      string                                `json:"-"`
		Events             []*MealPlanEventDatabaseCreationInput `json:"-"`
	}

	// MealPlanUpdateRequestInput represents what a user could set as input for updating meal plans.
	MealPlanUpdateRequestInput struct {
		_ struct{} `json:"-"`

		BelongsToHousehold *string    `json:"-"`
		Notes              *string    `json:"notes,omitempty"`
		VotingDeadline     *time.Time `json:"votingDeadline,omitempty"`
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
		GetMealPlansForHousehold(ctx context.Context, householdID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[MealPlan], error)
		CreateMealPlan(ctx context.Context, input *MealPlanDatabaseCreationInput) (*MealPlan, error)
		UpdateMealPlan(ctx context.Context, updated *MealPlan) error
		ArchiveMealPlan(ctx context.Context, mealPlanID, householdID string) error
		AttemptToFinalizeMealPlan(ctx context.Context, mealPlanID, householdID string) (bool, error)
		GetFinalizedMealPlanIDsForTheNextWeek(ctx context.Context) ([]*FinalizedMealPlanDatabaseResult, error)
		GetUnfinalizedMealPlansWithExpiredVotingPeriods(ctx context.Context) ([]*MealPlan, error)
		GetFinalizedMealPlansWithUninitializedGroceryLists(ctx context.Context) ([]*MealPlan, error)
	}

	// MealPlanDataService describes a structure capable of serving traffic related to meal plans.
	MealPlanDataService interface {
		ListMealPlanHandler(http.ResponseWriter, *http.Request)
		CreateMealPlanHandler(http.ResponseWriter, *http.Request)
		ReadMealPlanHandler(http.ResponseWriter, *http.Request)
		UpdateMealPlanHandler(http.ResponseWriter, *http.Request)
		ArchiveMealPlanHandler(http.ResponseWriter, *http.Request)
		FinalizeMealPlanHandler(http.ResponseWriter, *http.Request)
	}
)

// Update merges an MealPlanUpdateRequestInput with a meal plan.
func (x *MealPlan) Update(input *MealPlanUpdateRequestInput) {
	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
	}
}

var _ validation.ValidatableWithContext = (*MealPlanCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanCreationRequestInput.
func (x *MealPlanCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	if time.Now().After(x.VotingDeadline) {
		return ErrInvalidVotingDeadline
	}

	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.VotingDeadline, validation.Required),
		validation.Field(&x.Events, validation.Required),
		validation.Field(&x.ElectionMethod, validation.In(MealPlanElectionMethodSchulze, MealPlanElectionMethodInstantRunoff)),
	)
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
		validation.Field(&x.CreatedByUser, validation.Required),
	)
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
