package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// AdvancedPrepStepDataType indicates an event is related to an advanced prep step.
	AdvancedPrepStepDataType dataType = "advanced_prep_step"

	// AdvancedPrepStepCreatedCustomerEventType indicates an advanced prep step was created.
	AdvancedPrepStepCreatedCustomerEventType CustomerEventType = "advanced_prep_step_created"
	// AdvancedPrepStepStatusChangedCustomerEventType indicates an advanced prep step was created.
	AdvancedPrepStepStatusChangedCustomerEventType CustomerEventType = "advanced_prep_step_status_changed"

	// AdvancedPrepStepStatusUnfinished represents the unfinished enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusUnfinished = "unfinished"
	// AdvancedPrepStepStatusDelayed represents the delayed enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusDelayed = "delayed"
	// AdvancedPrepStepStatusIgnored represents the ignored enum member for advanced prep step status in the DB.
	AdvancedPrepStepStatusIgnored = "ignored"
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
	// AdvancedPrepStep represents a advanced prep step.
	AdvancedPrepStep struct {
		_                    struct{}
		CannotCompleteBefore time.Time      `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time      `json:"cannotCompleteAfter"`
		CreatedAt            time.Time      `json:"createdAt"`
		CompletedAt          *time.Time     `json:"completedAt"`
		RecipeStep           RecipeStep     `json:"recipeStep"`
		ID                   string         `json:"id"`
		Status               string         `json:"status"`
		CreationExplanation  string         `json:"creationExplanation"`
		StatusExplanation    string         `json:"statusExplanation"`
		MealPlanOption       MealPlanOption `json:"mealPlanOption"`
	}

	// AdvancedPrepStepList represents a list of advanced prep steps.
	AdvancedPrepStepList struct {
		_                 struct{}
		AdvancedPrepSteps []*AdvancedPrepStep `json:"data"`
		Pagination
	}

	// AdvancedPrepStepDatabaseCreationInput represents what a user could set as input for creating advanced prep steps.
	AdvancedPrepStepDatabaseCreationInput struct {
		_                    struct{}
		CannotCompleteBefore time.Time
		CannotCompleteAfter  time.Time
		CompletedAt          *time.Time
		Status               string
		CreationExplanation  string
		StatusExplanation    string
		MealPlanOptionID     string
		RecipeStepID         string
		ID                   string
	}

	// AdvancedPrepStepStatusChangeRequestInput represents what a user could set as input for updating advanced prep steps.
	AdvancedPrepStepStatusChangeRequestInput struct {
		_                 struct{}
		Status            string `json:"status"`
		StatusExplanation string `json:"statusExplanation"`
		ID                string `json:"-"`
	}

	// AdvancedPrepStepDatabaseCreationEstimate represents what a user could set as input for creating advanced prep steps.
	AdvancedPrepStepDatabaseCreationEstimate struct {
		_                   struct{}
		CreationExplanation string `json:"creationExplanation"`
		RecipeStepID        string `json:"recipeStepID"`
	}

	// AdvancedPrepStepDataManager describes a structure capable of storing advanced prep steps permanently.
	AdvancedPrepStepDataManager interface {
		AdvancedPrepStepExists(ctx context.Context, mealPlanID, advancedPrepStepID string) (bool, error)
		GetAdvancedPrepStep(ctx context.Context, advancedPrepStepID string) (*AdvancedPrepStep, error)
		GetAdvancedPrepStepsForMealPlan(ctx context.Context, mealPlanID string) ([]*AdvancedPrepStep, error)
		CreateAdvancedPrepStepsForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*AdvancedPrepStepDatabaseCreationInput) ([]*AdvancedPrepStep, error)
		ChangeAdvancedPrepStepStatus(ctx context.Context, input *AdvancedPrepStepStatusChangeRequestInput) error
	}

	// AdvancedPrepStepDataService describes a structure capable of serving traffic related to advanced prep steps.
	AdvancedPrepStepDataService interface {
		ListByMealPlanHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		StatusChangeHandler(res http.ResponseWriter, req *http.Request)
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
			AdvancedPrepStepStatusDelayed,
			AdvancedPrepStepStatusIgnored,
			AdvancedPrepStepStatusCanceled,
			AdvancedPrepStepStatusFinished,
		)),
		validation.Field(&x.CannotCompleteAfter, validation.Required),
	)
}

// ValidateWithContext validates a AdvancedPrepStepStatusChangeRequestInput.
func (x *AdvancedPrepStepStatusChangeRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Status, validation.In(
			AdvancedPrepStepStatusUnfinished,
			AdvancedPrepStepStatusDelayed,
			AdvancedPrepStepStatusIgnored,
			AdvancedPrepStepStatusCanceled,
			AdvancedPrepStepStatusFinished,
		)),
	)
}
