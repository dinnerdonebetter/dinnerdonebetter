package types

import (
	"context"
	"database/sql"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanOptionVoteDataType indicates an event is related to a meal plan option vote.
	MealPlanOptionVoteDataType dataType = "meal_plan_option_vote"

	// MealPlanOptionVoteCreatedCustomerEventType indicates a meal plan option vote was created.
	MealPlanOptionVoteCreatedCustomerEventType CustomerEventType = "meal_plan_option_vote_created"
	// MealPlanOptionVoteUpdatedCustomerEventType indicates a meal plan option vote was updated.
	MealPlanOptionVoteUpdatedCustomerEventType CustomerEventType = "meal_plan_option_vote_updated"
	// MealPlanOptionVoteArchivedCustomerEventType indicates a meal plan option vote was archived.
	MealPlanOptionVoteArchivedCustomerEventType CustomerEventType = "meal_plan_option_vote_archived"
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
		CreatedAt               time.Time  `json:"createdAt"`
		ArchivedAt              *time.Time `json:"archivedAt"`
		LastUpdatedAt           *time.Time `json:"lastUpdatedAt"`
		ID                      string     `json:"id"`
		Notes                   string     `json:"notes"`
		BelongsToMealPlanOption string     `json:"belongsToMealPlanOption"`
		ByUser                  string     `json:"byUser"`
		Rank                    uint8      `json:"rank"`
		Abstain                 bool       `json:"abstain"`
	}

	// NullableMealPlanOptionVote represents a fully nullable meal plan option vote.
	NullableMealPlanOptionVote struct {
		_                       struct{}
		Rank                    *uint8
		ID                      *string
		Notes                   *string
		BelongsToMealPlanOption *string
		ByUser                  *string
		Abstain                 *bool
		LastUpdatedAt           sql.NullTime
		CreatedAt               sql.NullTime
		ArchivedAt              sql.NullTime
	}

	// MealPlanOptionVoteList represents a list of meal plan option votes.
	MealPlanOptionVoteList struct {
		_                   struct{}
		MealPlanOptionVotes []*MealPlanOptionVote `json:"data"`
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
		Notes                   *string `json:"notes"`
		Rank                    *uint8  `json:"rank"`
		Abstain                 *bool   `json:"abstain"`
		BelongsToMealPlanOption string  `json:"belongsToMealPlanOption"`
	}

	// MealPlanOptionVoteDataManager describes a structure capable of storing meal plan option votes permanently.
	MealPlanOptionVoteDataManager interface {
		MealPlanOptionVoteExists(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (bool, error)
		GetMealPlanOptionVote(ctx context.Context, mealPlanID, mealPlanOptionID, mealPlanOptionVoteID string) (*MealPlanOptionVote, error)
		GetMealPlanOptionVotes(ctx context.Context, mealPlanID, mealPlanOptionID string, filter *QueryFilter) (*MealPlanOptionVoteList, error)
		GetMealPlanOptionVotesForMealPlanOption(ctx context.Context, mealPlanID, mealPlanOptionID string) (x []*MealPlanOptionVote, err error)
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
	if input.Rank != nil && *input.Rank != x.Rank {
		x.Rank = *input.Rank
	}

	if input.Abstain != nil && *input.Abstain != x.Abstain {
		x.Abstain = *input.Abstain
	}

	if input.Notes != nil && *input.Notes != x.Notes {
		x.Notes = *input.Notes
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

// MealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote creates a DatabaseCreationInput from a CreationInput.
func MealPlanOptionVoteUpdateRequestInputFromMealPlanOptionVote(input *MealPlanOptionVote) *MealPlanOptionVoteUpdateRequestInput {
	x := &MealPlanOptionVoteUpdateRequestInput{
		Notes:                   &input.Notes,
		BelongsToMealPlanOption: input.BelongsToMealPlanOption,
		Rank:                    &input.Rank,
		Abstain:                 &input.Abstain,
	}

	return x
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
