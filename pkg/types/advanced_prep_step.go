package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanTaskDataType indicates an event is related to an advanced prep step.
	MealPlanTaskDataType dataType = "meal_plan_task"

	// MealPlanTaskCreatedCustomerEventType indicates an advanced prep step was created.
	MealPlanTaskCreatedCustomerEventType CustomerEventType = "meal_plan_task_created"
	// MealPlanTaskStatusChangedCustomerEventType indicates an advanced prep step was created.
	MealPlanTaskStatusChangedCustomerEventType CustomerEventType = "meal_plan_task_status_changed"

	// MealPlanTaskStatusUnfinished represents the unfinished enum member for advanced prep step status in the DB.
	MealPlanTaskStatusUnfinished = "unfinished"
	// MealPlanTaskStatusDelayed represents the delayed enum member for advanced prep step status in the DB.
	MealPlanTaskStatusDelayed = "delayed"
	// MealPlanTaskStatusIgnored represents the ignored enum member for advanced prep step status in the DB.
	MealPlanTaskStatusIgnored = "ignored"
	// MealPlanTaskStatusCanceled represents the canceled enum member for advanced prep step status in the DB.
	MealPlanTaskStatusCanceled = "canceled"
	// MealPlanTaskStatusFinished represents the finished enum member for advanced prep step status in the DB.
	MealPlanTaskStatusFinished = "finished"
)

func init() {
	gob.Register(new(MealPlanTask))
	gob.Register(new(MealPlanTaskList))
}

type (
	// MealPlanTask represents a advanced prep step.
	MealPlanTask struct {
		_                    struct{}
		CannotCompleteBefore time.Time      `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time      `json:"cannotCompleteAfter"`
		CreatedAt            time.Time      `json:"createdAt"`
		CompletedAt          *time.Time     `json:"completedAt"`
		RecipeStep           RecipeStep     `json:"recipeStep"`
		ID                   string         `json:"id"`
		AssignedToUser       *string        `json:"assignedToUser"`
		Status               string         `json:"status"`
		CreationExplanation  string         `json:"creationExplanation"`
		StatusExplanation    string         `json:"statusExplanation"`
		MealPlanOption       MealPlanOption `json:"mealPlanOption"`
	}

	// MealPlanTaskList represents a list of advanced prep steps.
	MealPlanTaskList struct {
		_             struct{}
		MealPlanTasks []*MealPlanTask `json:"data"`
		Pagination
	}

	// MealPlanTaskDatabaseCreationInput represents what a user could set as input for creating advanced prep steps.
	MealPlanTaskDatabaseCreationInput struct {
		_                    struct{}
		CannotCompleteBefore time.Time
		CannotCompleteAfter  time.Time
		CompletedAt          *time.Time
		AssignedToUser       *string
		Status               string
		CreationExplanation  string
		StatusExplanation    string
		MealPlanOptionID     string
		RecipeStepID         string
		ID                   string
	}

	// MealPlanTaskStatusChangeRequestInput represents what a user could set as input for updating advanced prep steps.
	MealPlanTaskStatusChangeRequestInput struct {
		_                 struct{}
		Status            string  `json:"status"`
		StatusExplanation string  `json:"statusExplanation"`
		AssignedToUser    *string `json:"assignedToUser"`
		ID                string  `json:"-"`
	}

	// MealPlanTaskDatabaseCreationEstimate represents what a user could set as input for creating advanced prep steps.
	MealPlanTaskDatabaseCreationEstimate struct {
		_                   struct{}
		CreationExplanation string `json:"creationExplanation"`
		RecipeStepID        string `json:"recipeStepID"`
	}

	// MealPlanTaskDataManager describes a structure capable of storing advanced prep steps permanently.
	MealPlanTaskDataManager interface {
		MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error)
		GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*MealPlanTask, error)
		GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*MealPlanTask, error)
		CreateMealPlanTasksForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*MealPlanTaskDatabaseCreationInput) ([]*MealPlanTask, error)
		ChangeMealPlanTaskStatus(ctx context.Context, input *MealPlanTaskStatusChangeRequestInput) error
	}

	// MealPlanTaskDataService describes a structure capable of serving traffic related to advanced prep steps.
	MealPlanTaskDataService interface {
		ListByMealPlanHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		StatusChangeHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*MealPlanTaskDatabaseCreationInput)(nil)

// ValidateWithContext validates a MealPlanTaskDatabaseCreationInput.
func (x *MealPlanTaskDatabaseCreationInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.RecipeStepID, validation.Required),
		validation.Field(&x.CannotCompleteBefore, validation.Required),
		validation.Field(&x.Status, validation.In(
			MealPlanTaskStatusUnfinished,
			MealPlanTaskStatusDelayed,
			MealPlanTaskStatusIgnored,
			MealPlanTaskStatusCanceled,
			MealPlanTaskStatusFinished,
		)),
		validation.Field(&x.CannotCompleteAfter, validation.Required),
	)
}

// ValidateWithContext validates a MealPlanTaskStatusChangeRequestInput.
func (x *MealPlanTaskStatusChangeRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.ID, validation.Required),
		validation.Field(&x.Status, validation.In(
			MealPlanTaskStatusUnfinished,
			MealPlanTaskStatusDelayed,
			MealPlanTaskStatusIgnored,
			MealPlanTaskStatusCanceled,
			MealPlanTaskStatusFinished,
		)),
	)
}
