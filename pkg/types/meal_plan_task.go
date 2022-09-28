package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanTaskDataType indicates an event is related to an meal plan task.
	MealPlanTaskDataType dataType = "meal_plan_task"

	// MealPlanTaskCreatedCustomerEventType indicates an meal plan task was created.
	MealPlanTaskCreatedCustomerEventType CustomerEventType = "meal_plan_task_created"
	// MealPlanTaskStatusChangedCustomerEventType indicates an meal plan task was created.
	MealPlanTaskStatusChangedCustomerEventType CustomerEventType = "meal_plan_task_status_changed"

	// MealPlanTaskStatusUnfinished represents the unfinished enum member for meal plan task status in the DB.
	MealPlanTaskStatusUnfinished = "unfinished"
	// MealPlanTaskStatusDelayed represents the delayed enum member for meal plan task status in the DB.
	MealPlanTaskStatusDelayed = "delayed"
	// MealPlanTaskStatusIgnored represents the ignored enum member for meal plan task status in the DB.
	MealPlanTaskStatusIgnored = "ignored"
	// MealPlanTaskStatusCanceled represents the canceled enum member for meal plan task status in the DB.
	MealPlanTaskStatusCanceled = "canceled"
	// MealPlanTaskStatusFinished represents the finished enum member for meal plan task status in the DB.
	MealPlanTaskStatusFinished = "finished"
)

func init() {
	gob.Register(new(MealPlanTask))
	gob.Register(new(MealPlanTaskList))
}

type (
	// MealPlanTask represents a meal plan task.
	MealPlanTask struct {
		_                    struct{}
		CannotCompleteBefore time.Time  `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time  `json:"cannotCompleteAfter"`
		CreatedAt            time.Time  `json:"createdAt"`
		CompletedAt          *time.Time `json:"completedAt"`
		ID                   string     `json:"id"`
		AssignedToUser       *string    `json:"assignedToUser"`
		Status               string     `json:"status"`
		CreationExplanation  string     `json:"creationExplanation"`
		StatusExplanation    string     `json:"statusExplanation"`
	}

	// MealPlanTaskRecipeStep represents a meal plan task.
	MealPlanTaskRecipeStep struct {
		_                    struct{}
		CannotCompleteBefore time.Time  `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time  `json:"cannotCompleteAfter"`
		CreatedAt            time.Time  `json:"createdAt"`
		CompletedAt          *time.Time `json:"completedAt"`
		ID                   string     `json:"id"`
		AssignedToUser       *string    `json:"assignedToUser"`
		Status               string     `json:"status"`
		CreationExplanation  string     `json:"creationExplanation"`
		StatusExplanation    string     `json:"statusExplanation"`
	}

	// MealPlanTaskList represents a list of meal plan tasks.
	MealPlanTaskList struct {
		_             struct{}
		MealPlanTasks []*MealPlanTask `json:"data"`
		Pagination
	}

	// MealPlanTaskCreationRequestInput represents a meal plan task.
	MealPlanTaskCreationRequestInput struct {
		_                    struct{}
		CannotCompleteBefore time.Time `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time `json:"cannotCompleteAfter"`
		RecipeStepID         string    `json:"recipeStepID"`
		AssignedToUser       *string   `json:"assignedToUser"`
		Status               string    `json:"status"`
		CreationExplanation  string    `json:"creationExplanation"`
		StatusExplanation    string    `json:"statusExplanation"`
		MealPlanOptionID     string    `json:"mealPlanOptionID"`
	}

	// MealPlanTaskDatabaseCreationInput represents what a user could set as input for creating meal plan tasks.
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

	// MealPlanTaskStatusChangeRequestInput represents what a user could set as input for updating meal plan tasks.
	MealPlanTaskStatusChangeRequestInput struct {
		_                 struct{}
		Status            string  `json:"status"`
		StatusExplanation string  `json:"statusExplanation"`
		AssignedToUser    *string `json:"assignedToUser"`
		ID                string  `json:"-"`
	}

	// MealPlanTaskDatabaseCreationEstimate represents what a user could set as input for creating meal plan tasks.
	MealPlanTaskDatabaseCreationEstimate struct {
		_                   struct{}
		CreationExplanation string `json:"creationExplanation"`
		RecipeStepID        string `json:"recipeStepID"`
	}

	// MealPlanTaskDataManager describes a structure capable of storing meal plan tasks permanently.
	MealPlanTaskDataManager interface {
		MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error)
		CreateMealPlanTask(ctx context.Context, mealPlanID string, input *MealPlanTaskDatabaseCreationInput) (*MealPlanTask, error)
		GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*MealPlanTask, error)
		GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*MealPlanTask, error)
		CreateMealPlanTasksForMealPlanOption(ctx context.Context, mealPlanOptionID string, inputs []*MealPlanTaskDatabaseCreationInput) ([]*MealPlanTask, error)
		ChangeMealPlanTaskStatus(ctx context.Context, input *MealPlanTaskStatusChangeRequestInput) error
	}

	// MealPlanTaskDataService describes a structure capable of serving traffic related to meal plan tasks.
	MealPlanTaskDataService interface {
		ListByMealPlanHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		StatusChangeHandler(res http.ResponseWriter, req *http.Request)
	}
)

var _ validation.ValidatableWithContext = (*MealPlanTaskCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanTaskDatabaseCreationInput.
func (x *MealPlanTaskCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
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

// MealPlanTaskDatabaseCreationInputFromMealPlanTaskCreationRequestInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanTaskDatabaseCreationInputFromMealPlanTaskCreationRequestInput(input *MealPlanTaskCreationRequestInput) *MealPlanTaskDatabaseCreationInput {
	x := &MealPlanTaskDatabaseCreationInput{
		CannotCompleteBefore: input.CannotCompleteBefore,
		CannotCompleteAfter:  input.CannotCompleteAfter,
		AssignedToUser:       input.AssignedToUser,
		Status:               input.Status,
		CreationExplanation:  input.CreationExplanation,
		StatusExplanation:    input.StatusExplanation,
		MealPlanOptionID:     input.MealPlanOptionID,
		RecipeStepID:         input.RecipeStepID,
	}

	return x
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
