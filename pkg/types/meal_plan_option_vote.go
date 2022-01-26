package types

import (
	"context"
	"encoding/gob"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanOptionVoteDataType indicates an event is related to a meal plan option vote.
	MealPlanOptionVoteDataType dataType = "meal_plan_option_vote"
)

func init() {
	gob.Register(new(MealPlanOptionVote))
	gob.Register(new(MealPlanOptionVoteList))
	gob.Register(new(MealPlanOptionVoteCreationRequestInput))
	gob.Register(new(MealPlanOptionVoteUpdateRequestInput))
}

type (
	// MealPlanOptionVote represents a meal plan option vote.
	MealPlanOptionVote struct {
		_                       struct{}
		LastUpdatedOn           *uint64 `json:"lastUpdatedOn"`
		ArchivedOn              *uint64 `json:"archivedOn"`
		ID                      string  `json:"id"`
		Notes                   string  `json:"notes"`
		BelongsToMealPlanOption string  `json:"belongsToMealPlanOption"`
		ByUser                  string  `json:"byUser"`
		CreatedOn               uint64  `json:"createdOn"`
		Rank                    uint8   `json:"rank"`
		Abstain                 bool    `json:"abstain"`
	}

	// MealPlanOptionVoteList represents a list of meal plan option votes.
	MealPlanOptionVoteList struct {
		_                   struct{}
		MealPlanOptionVotes []*MealPlanOptionVote `json:"mealPlanOptionVotes"`
		Pagination
	}

	// MealPlanOptionVoteCreationInput represents what a user could set as input for creating meal plan option votes.
	MealPlanOptionVoteCreationInput struct {
		_                       struct{}
		ID                      string `json:"-"`
		Notes                   string `json:"notes"`
		ByUser                  string `json:"-"`
		BelongsToMealPlanOption string `json:"belongsToMealPlanOption"`
		Rank                    uint8  `json:"rank"`
		Abstain                 bool   `json:"abstain"`
	}

	// MealPlanOptionVoteCreationRequestInput is a pending container for multiple votes.
	MealPlanOptionVoteCreationRequestInput struct {
		_      struct{}
		ByUser string                             `json:"-"`
		Votes  []*MealPlanOptionVoteCreationInput `json:"votes"`
	}

	// MealPlanOptionVoteDatabaseCreationInput represents what a user could set as input for creating meal plan option votes.
	MealPlanOptionVoteDatabaseCreationInput struct {
		_      struct{}
		ByUser string                             `json:"-"`
		Votes  []*MealPlanOptionVoteCreationInput `json:"votes"`
	}

	// MealPlanOptionVoteUpdateRequestInput represents what a user could set as input for updating meal plan option votes.
	MealPlanOptionVoteUpdateRequestInput struct {
		_                       struct{}
		Notes                   string `json:"notes"`
		BelongsToMealPlanOption string `json:"belongsToMealPlanOption"`
		Rank                    uint8  `json:"rank"`
		Abstain                 bool   `json:"abstain"`
	}

	// MealPlanOptionVoteDataManager describes a structure capable of storing meal plan option votes permanently.
	MealPlanOptionVoteDataManager interface {
		MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (bool, error)
		GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (*MealPlanOptionVote, error)
		GetTotalMealPlanOptionVoteCount(ctx context.Context) (uint64, error)
		GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanOptionID string, filter *QueryFilter) (*MealPlanOptionVoteList, error)
		GetMealPlanOptionVotesWithIDs(ctx context.Context, mealPlanOptionID string, limit uint8, ids []string) ([]*MealPlanOptionVote, error)
		CreateMealPlanOptionVote(ctx context.Context, input *MealPlanOptionVoteDatabaseCreationInput) ([]*MealPlanOptionVote, error)
		UpdateMealPlanOptionVote(ctx context.Context, updated *MealPlanOptionVote) error
		ArchiveMealPlanOptionVote(ctx context.Context, mealPlanOptionID, mealPlanOptionVoteID string) error
	}

	// MealPlanOptionVoteDataService describes a structure capable of serving traffic related to meal plan option votes.
	MealPlanOptionVoteDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an MealPlanOptionVoteUpdateRequestInput with a meal plan option vote.
func (x *MealPlanOptionVote) Update(input *MealPlanOptionVoteUpdateRequestInput) {
	if input.Rank != 0 && input.Rank != x.Rank {
		x.Rank = input.Rank
	}

	if input.Abstain != x.Abstain {
		x.Abstain = input.Abstain
	}

	if input.Notes != "" && input.Notes != x.Notes {
		x.Notes = input.Notes
	}
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteCreationRequestInput.
func (x *MealPlanOptionVoteCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	// LATER: we should validate the contents of each individual vote
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Votes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteDatabaseCreationInput.
func (x *MealPlanOptionVoteDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.Votes, validation.Required),
	)
}

// MealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVoteCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVoteCreationInput(input *MealPlanOptionVoteCreationRequestInput) *MealPlanOptionVoteDatabaseCreationInput {
	x := &MealPlanOptionVoteDatabaseCreationInput{
		Votes:  input.Votes,
		ByUser: input.ByUser,
	}

	return x
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteUpdateRequestInput.
func (x *MealPlanOptionVoteUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		// we cannot validate rank because zero is a valid value.
		validation.Field(&x.BelongsToMealPlanOption, validation.Required),
	)
}
