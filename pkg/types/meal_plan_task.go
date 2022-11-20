package types

import (
	"context"
	"encoding/gob"
	"net/http"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// MealPlanTaskDataType indicates an event is related to a meal plan task.
	MealPlanTaskDataType dataType = "meal_plan_task"

	// MealPlanTaskCreatedCustomerEventType indicates a meal plan task was created.
	MealPlanTaskCreatedCustomerEventType CustomerEventType = "meal_plan_task_created"
	// MealPlanTaskStatusChangedCustomerEventType indicates a meal plan task was created.
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
}

type (
	// MealPlanTask represents a meal plan task.
	MealPlanTask struct {
		_                   struct{}
		CreatedAt           time.Time      `json:"createdAt"`
		LastUpdatedAt       *time.Time     `json:"lastUpdatedAt"`
		CompletedAt         *time.Time     `json:"completedAt"`
		AssignedToUser      *string        `json:"assignedToUser"`
		ID                  string         `json:"id"`
		Status              string         `json:"status"`
		CreationExplanation string         `json:"creationExplanation"`
		StatusExplanation   string         `json:"statusExplanation"`
		MealPlanOption      MealPlanOption `json:"mealPlanOption"`
		RecipePrepTask      RecipePrepTask `json:"recipePrepTask"`
	}

	// MealPlanTaskCreationRequestInput represents a meal plan task.
	MealPlanTaskCreationRequestInput struct {
		_                   struct{}
		AssignedToUser      *string `json:"assignedToUser"`
		Status              string  `json:"status"`
		CreationExplanation string  `json:"creationExplanation"`
		StatusExplanation   string  `json:"statusExplanation"`
		MealPlanOptionID    string  `json:"mealPlanOptionID"`
		RecipePrepTaskID    string  `json:"recipePrepTaskID"`
	}

	// MealPlanTaskDatabaseCreationInput represents what a user could set as input for creating meal plan tasks.
	MealPlanTaskDatabaseCreationInput struct {
		_                   struct{}
		AssignedToUser      *string
		CreationExplanation string
		StatusExplanation   string
		MealPlanOptionID    string
		RecipePrepTaskID    string
		ID                  string
	}

	// MealPlanTaskStatusChangeRequestInput represents what a user could set as input for updating meal plan tasks.
	MealPlanTaskStatusChangeRequestInput struct {
		_                 struct{}
		Status            *string `json:"status"`
		StatusExplanation *string `json:"statusExplanation"`
		AssignedToUser    *string `json:"assignedToUser"`
		ID                string  `json:"-"`
	}

	// MealPlanTaskDatabaseCreationEstimate represents what a user could set as input for creating meal plan tasks.
	MealPlanTaskDatabaseCreationEstimate struct {
		_                   struct{}
		CreationExplanation string `json:"creationExplanation"`
	}

	// MealPlanTaskDataManager describes a structure capable of storing meal plan tasks permanently.
	MealPlanTaskDataManager interface {
		MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error)
		CreateMealPlanTask(ctx context.Context, input *MealPlanTaskDatabaseCreationInput) (*MealPlanTask, error)
		GetMealPlanTask(ctx context.Context, mealPlanTaskID string) (*MealPlanTask, error)
		GetMealPlanTasksForMealPlan(ctx context.Context, mealPlanID string) ([]*MealPlanTask, error)
		CreateMealPlanTasksForMealPlanOption(ctx context.Context, inputs []*MealPlanTaskDatabaseCreationInput) ([]*MealPlanTask, error)
		ChangeMealPlanTaskStatus(ctx context.Context, input *MealPlanTaskStatusChangeRequestInput) error
		MarkMealPlanAsHavingTasksCreated(ctx context.Context, mealPlanID string) error
	}

	// MealPlanTaskDataService describes a structure capable of serving traffic related to meal plan tasks.
	MealPlanTaskDataService interface {
		ListByMealPlanHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		StatusChangeHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an MealPlanTaskStatusChangeRequestInput with a meal plan task.
func (x *MealPlanTask) Update(input *MealPlanTaskStatusChangeRequestInput) {
	if input.StatusExplanation != nil && *input.StatusExplanation != x.StatusExplanation {
		x.StatusExplanation = *input.StatusExplanation
	}

	if input.Status != nil && *input.Status != x.Status {
		x.Status = *input.Status
	}

	if input.AssignedToUser != nil && *input.AssignedToUser != "" {
		if x.AssignedToUser == nil || *x.AssignedToUser != *input.AssignedToUser {
			x.AssignedToUser = input.AssignedToUser
		}
	}

	x.AssignedToUser = input.AssignedToUser
}

var _ validation.ValidatableWithContext = (*MealPlanTaskCreationRequestInput)(nil)

// ValidateWithContext validates a MealPlanTaskDatabaseCreationInput.
func (x *MealPlanTaskCreationRequestInput) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		x,
		validation.Field(&x.MealPlanOptionID, validation.Required),
		validation.Field(&x.RecipePrepTaskID, validation.Required),
		validation.Field(&x.Status, validation.In(
			MealPlanTaskStatusUnfinished,
			MealPlanTaskStatusDelayed,
			MealPlanTaskStatusIgnored,
			MealPlanTaskStatusCanceled,
			MealPlanTaskStatusFinished,
		)),
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
		validation.Field(&x.RecipePrepTaskID, validation.Required),
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
