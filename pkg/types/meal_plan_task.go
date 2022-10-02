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
	gob.Register(new(MealPlanTaskList))
}

type (
	// MealPlanTask represents a meal plan task.
	MealPlanTask struct {
		_                    struct{}
		CannotCompleteBefore time.Time      `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time      `json:"cannotCompleteAfter"`
		CreatedAt            time.Time      `json:"createdAt"`
		CompletedAt          *time.Time     `json:"completedAt"`
		AssignedToUser       *string        `json:"assignedToUser"`
		ID                   string         `json:"id"`
		Status               string         `json:"status"`
		CreationExplanation  string         `json:"creationExplanation"`
		StatusExplanation    string         `json:"statusExplanation"`
		RecipeSteps          []*RecipeStep  `json:"recipeSteps"`
		MealPlanOption       MealPlanOption `json:"mealPlanOption"`
	}

	// MealPlanTaskRecipeStep represents a meal plan task's recipe step.
	MealPlanTaskRecipeStep struct {
		_              struct{}
		MealPlanTaskID string
		RecipeStepID   string
		ID             string
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
		AssignedToUser       *string   `json:"assignedToUser"`
		Status               string    `json:"status"`
		CreationExplanation  string    `json:"creationExplanation"`
		StatusExplanation    string    `json:"statusExplanation"`
		MealPlanOptionID     string    `json:"mealPlanOptionID"`
		RecipeStepIDs        []string  `json:"recipeStepIDs"`
	}

	// MealPlanTaskDatabaseCreationInput represents what a user could set as input for creating meal plan tasks.
	MealPlanTaskDatabaseCreationInput struct {
		_                    struct{}
		CannotCompleteBefore time.Time
		CannotCompleteAfter  time.Time
		AssignedToUser       *string
		CreationExplanation  string
		StatusExplanation    string
		MealPlanOptionID     string
		ID                   string
		RecipeSteps          []*MealPlanTaskRecipeStepDatabaseCreationInput
	}

	// MealPlanTaskRecipeStepDatabaseCreationInput represents what a user could set as input for creating meal plan tasks.
	MealPlanTaskRecipeStepDatabaseCreationInput struct {
		_                     struct{}
		BelongsToMealPlanTask string
		SatisfiesRecipeStep   string
		ID                    string
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
		CreationExplanation string   `json:"creationExplanation"`
		RecipeStepIDs       []string `json:"recipeStepIDs"`
	}

	// MealPlanTaskDataManager describes a structure capable of storing meal plan tasks permanently.
	MealPlanTaskDataManager interface {
		MealPlanTaskExists(ctx context.Context, mealPlanID, mealPlanTaskID string) (bool, error)
		CreateMealPlanTask(ctx context.Context, input *MealPlanTaskDatabaseCreationInput) (*MealPlanTask, error)
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
		validation.Field(&x.CannotCompleteBefore, validation.Required),
		validation.Field(&x.CannotCompleteAfter, validation.Required),
		validation.Field(&x.Status, validation.In(
			MealPlanTaskStatusUnfinished,
			MealPlanTaskStatusDelayed,
			// MealPlanTaskStatusIgnored,
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
		validation.Field(&x.RecipeSteps, validation.Required),
		validation.Field(&x.CannotCompleteBefore, validation.Required),
		validation.Field(&x.CannotCompleteAfter, validation.Required),
	)
}

// MealPlanTaskDatabaseCreationInputFromMealPlanTaskCreationRequestInput creates a DatabaseCreationInput from a CreationInput.
func MealPlanTaskDatabaseCreationInputFromMealPlanTaskCreationRequestInput(input *MealPlanTaskCreationRequestInput) *MealPlanTaskDatabaseCreationInput {
	x := &MealPlanTaskDatabaseCreationInput{
		CannotCompleteBefore: input.CannotCompleteBefore,
		CannotCompleteAfter:  input.CannotCompleteAfter,
		AssignedToUser:       input.AssignedToUser,
		CreationExplanation:  input.CreationExplanation,
		StatusExplanation:    input.StatusExplanation,
		MealPlanOptionID:     input.MealPlanOptionID,
		RecipeSteps:          []*MealPlanTaskRecipeStepDatabaseCreationInput{},
	}

	for _, recipeStepID := range input.RecipeStepIDs {
		x.RecipeSteps = append(x.RecipeSteps, &MealPlanTaskRecipeStepDatabaseCreationInput{SatisfiesRecipeStep: recipeStepID})
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
			// MealPlanTaskStatusIgnored,
			MealPlanTaskStatusCanceled,
			MealPlanTaskStatusFinished,
		)),
	)
}