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
		_                  struct{}
		LastUpdatedOn      *uint64 `json:"lastUpdatedOn"`
		ArchivedOn         *uint64 `json:"archivedOn"`
		ID                 string  `json:"id"`
		BelongsToHousehold string  `json:"belongsToHousehold"`
		Notes              string  `json:"notes"`
		MealPlanOptionID   string  `json:"mealPlanOptionID"`
		CreatedOn          uint64  `json:"createdOn"`
		Points             int16   `json:"points"`
		Abstain            bool    `json:"abstain"`
		DayOfWeek          uint8   `json:"dayOfWeek"`
	}

	// MealPlanOptionVoteList represents a list of meal plan option votes.
	MealPlanOptionVoteList struct {
		_                   struct{}
		MealPlanOptionVotes []*MealPlanOptionVote `json:"mealPlanOptionVotes"`
		Pagination
	}

	// MealPlanOptionVoteCreationRequestInput represents what a user could set as input for creating meal plan option votes.
	MealPlanOptionVoteCreationRequestInput struct {
		_                  struct{}
		ID                 string `json:"-"`
		MealPlanOptionID   string `json:"mealPlanOptionID"`
		Notes              string `json:"notes"`
		BelongsToHousehold string `json:"-"`
		Points             int16  `json:"points"`
		DayOfWeek          uint8  `json:"dayOfWeek"`
		Abstain            bool   `json:"abstain"`
	}

	// MealPlanOptionVoteDatabaseCreationInput represents what a user could set as input for creating meal plan option votes.
	MealPlanOptionVoteDatabaseCreationInput struct {
		_                  struct{}
		ID                 string `json:"id"`
		MealPlanOptionID   string `json:"mealPlanOptionID"`
		Notes              string `json:"notes"`
		BelongsToHousehold string `json:"belongsToHousehold"`
		Points             int16  `json:"points"`
		DayOfWeek          uint8  `json:"dayOfWeek"`
		Abstain            bool   `json:"abstain"`
	}

	// MealPlanOptionVoteUpdateRequestInput represents what a user could set as input for updating meal plan option votes.
	MealPlanOptionVoteUpdateRequestInput struct {
		_                  struct{}
		MealPlanOptionID   string `json:"mealPlanOptionID"`
		Notes              string `json:"notes"`
		BelongsToHousehold string `json:"-"`
		Points             int16  `json:"points"`
		DayOfWeek          uint8  `json:"dayOfWeek"`
		Abstain            bool   `json:"abstain"`
	}

	// MealPlanOptionVoteDataManager describes a structure capable of storing meal plan option votes permanently.
	MealPlanOptionVoteDataManager interface {
		MealPlanOptionVoteExists(ctx context.Context, mealPlanOptionVoteID string) (bool, error)
		GetMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID string) (*MealPlanOptionVote, error)
		GetTotalMealPlanOptionVoteCount(ctx context.Context) (uint64, error)
		GetMealPlanOptionVotes(ctx context.Context, filter *QueryFilter) (*MealPlanOptionVoteList, error)
		GetMealPlanOptionVotesWithIDs(ctx context.Context, householdID string, limit uint8, ids []string) ([]*MealPlanOptionVote, error)
		CreateMealPlanOptionVote(ctx context.Context, input *MealPlanOptionVoteDatabaseCreationInput) (*MealPlanOptionVote, error)
		UpdateMealPlanOptionVote(ctx context.Context, updated *MealPlanOptionVote) error
		ArchiveMealPlanOptionVote(ctx context.Context, mealPlanOptionVoteID, householdID string) error
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
	if input.MealPlanOptionID != "" && input.MealPlanOptionID != x.MealPlanOptionID {
		x.MealPlanOptionID = input.MealPlanOptionID
	}

	if input.DayOfWeek != 0 && input.DayOfWeek != x.DayOfWeek {
		x.DayOfWeek = input.DayOfWeek
	}

	if input.Points != 0 && input.Points != x.Points {
		x.Points = input.Points
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
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.DayOfWeek, validation.Required),
		validation.Field(&x.Points, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteDatabaseCreationInput.
func (x *MealPlanOptionVoteDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.DayOfWeek, validation.Required),
		validation.Field(&x.Points, validation.Required),
		validation.Field(&x.Notes, validation.Required),
		validation.Field(&x.BelongsToHousehold, validation.Required),
	)
}

// MealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVoteCreationInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanOptionVoteDatabaseCreationInputFromMealPlanOptionVoteCreationInput(input *MealPlanOptionVoteCreationRequestInput) *MealPlanOptionVoteDatabaseCreationInput {
	x := &MealPlanOptionVoteDatabaseCreationInput{
		MealPlanOptionID: input.MealPlanOptionID,
		DayOfWeek:        input.DayOfWeek,
		Points:           input.Points,
		Abstain:          input.Abstain,
		Notes:            input.Notes,
	}

	return x
}

var _ validation.ValidatableWithContext = (*MealPlanOptionVoteUpdateRequestInput)(nil)

// ValidateWithContext validates a MealPlanOptionVoteUpdateRequestInput.
func (x *MealPlanOptionVoteUpdateRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.DayOfWeek, validation.Required),
		validation.Field(&x.Points, validation.Required),
		validation.Field(&x.Notes, validation.Required),
	)
}
