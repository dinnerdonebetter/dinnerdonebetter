package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// AdvancedPrepStepDataType indicates an event is related to a valid preparation.
	AdvancedPrepStepDataType dataType = "advanced_prep_step"

	// AdvancedPrepStepCreatedCustomerEventType indicates an advanced prep step was created.
	AdvancedPrepStepCreatedCustomerEventType CustomerEventType = "advanced_prep_step_created"
	// AdvancedPrepStepUpdatedCustomerEventType indicates an advanced prep step was updated.
	AdvancedPrepStepUpdatedCustomerEventType CustomerEventType = "advanced_prep_step_updated"
	// AdvancedPrepStepArchivedCustomerEventType indicates an advanced prep step was archived.
	AdvancedPrepStepArchivedCustomerEventType CustomerEventType = "advanced_prep_steparchived"

	// AdvancedPrepStepStatusUnfinished represents the unfinished enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusUnfinished = "unfinished"
	// AdvancedPrepStepStatusPostponed represents the postponed enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusPostponed = "postponed"
	// AdvancedPrepStepStatusCanceled represents the canceled enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusCanceled = "canceled"
	// AdvancedPrepStepStatusFinished represents the finished enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusFinished = "finished"
)

func init() {
	gob.Register(new(AdvancedPrepStep))
	gob.Register(new(AdvancedPrepStepList))
}

type (
	// AdvancedPrepStep represents a valid preparation.
	AdvancedPrepStep struct {
		_                    struct{}
		CannotCompleteBefore time.Time      `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time      `json:"cannotCompleteAfter"`
		CreatedAt            time.Time      `json:"createdAt"`
		SettledAt            *uint64        `json:"completedAt"`
		ID                   string         `json:"id"`
		Status               string         `json:"status"`
		CreationExplanation  string         `json:"creationExplanation"`
		StatusExplanation    string         `json:"statusExplanation"`
		MealPlanOption       MealPlanOption `json:"mealPlanOption"`
		RecipeStep           RecipeStep     `json:"recipeStep"`
	}

	// AdvancedPrepStepList represents a list of valid preparations.
	AdvancedPrepStepList struct {
		_                 struct{}
		AdvancedPrepSteps []*AdvancedPrepStep `json:"data"`
		Pagination
	}

	// AdvancedPrepStepDatabaseCreationInput represents what a user could set as input for creating valid preparations.
	AdvancedPrepStepDatabaseCreationInput struct {
		_                    struct{}
		CreatedAt            time.Time `json:"createdAt"`
		CannotCompleteBefore time.Time `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time `json:"cannotCompleteAfter"`
		CompletedAt          *uint64   `json:"completedAt"`
		Status               string    `json:"status"`
		CreationExplanation  string    `json:"creationExplanation"`
		StatusExplanation    string    `json:"statusExplanation"`
		MealPlanOptionID     string    `json:"mealPlanOptionID"`
		RecipeStepID         string    `json:"recipeStepID"`
		ID                   string    `json:"id"`
	}

	// AdvancedPrepStepDataManager describes a structure capable of storing valid preparations permanently.
	AdvancedPrepStepDataManager interface {
		GetAdvancedPrepStep(ctx context.Context, advancedPrepStepID string) (*AdvancedPrepStep, error)
		GetAdvancedPrepSteps(ctx context.Context, filter *QueryFilter) (*AdvancedPrepStepList, error)
		CreateAdvancedPrepStep(ctx context.Context, input *AdvancedPrepStepDatabaseCreationInput) (*AdvancedPrepStep, error)
		MarkAdvancedPrepStepAsComplete(ctx context.Context, advancedPrepStepID string) error
	}

	// AdvancedPrepStepDataService describes a structure capable of serving traffic related to valid preparations.
	AdvancedPrepStepDataService interface {
		ListHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		CompletionHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*AdvancedPrepStepDatabaseCreationInput)(nil)

// ValidateWithContext validates a AdvancedPrepStepDatabaseCreationInput.
func (x *AdvancedPrepStepDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
		validation.Field(&x.CannotCompleteBefore, validation.Required),
		validation.Field(&x.Status, validation.In(
			AdvancedPrepStepStatusUnfinished,
			AdvancedPrepStepStatusPostponed,
			AdvancedPrepStepStatusCanceled,
			AdvancedPrepStepStatusFinished,
		)),
		validation.Field(&x.CannotCompleteAfter, validation.Required),
	)
}
