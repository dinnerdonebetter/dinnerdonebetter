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
	AdvancedPrepStepDataType dataType = "valid_preparation"

	// AdvancedPrepStepCreatedCustomerEventType indicates a valid preparation was created.
	AdvancedPrepStepCreatedCustomerEventType CustomerEventType = "valid_preparation_created"
	// AdvancedPrepStepUpdatedCustomerEventType indicates a valid preparation was updated.
	AdvancedPrepStepUpdatedCustomerEventType CustomerEventType = "valid_preparation_updated"
	// AdvancedPrepStepArchivedCustomerEventType indicates a valid preparation was archived.
	AdvancedPrepStepArchivedCustomerEventType CustomerEventType = "valid_preparation_archived"
)

func init() {
	gob.Register(new(AdvancedPrepStep))
	gob.Register(new(AdvancedPrepStepList))
}

type (
	// AdvancedPrepStep represents a valid preparation.
	AdvancedPrepStep struct {
		_                    struct{}
		CreatedAt            time.Time      `json:"createdAt"`
		CannotCompleteBefore time.Time      `json:"cannotCompleteBefore"`
		CannotCompleteAfter  time.Time      `json:"cannotCompleteAfter"`
		CompletedAt          *uint64        `json:"completedAt"`
		ID                   string         `json:"id"`
		RecipeStep           RecipeStep     `json:"recipeStep"`
		MealPlanOption       MealPlanOption `json:"mealPlanOption"`
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
		MealPlanOptionID     string    `json:"mealPlanOptionID"`
		RecipeStepID         string    `json:"recipeStepID"`
		ID                   string    `json:"id"`
	}

	// AdvancedPrepStepDataManager describes a structure capable of storing valid preparations permanently.
	AdvancedPrepStepDataManager interface {
		AdvancedPrepStepExists(ctx context.Context, advancedPrepStepID string) (bool, error)
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
		validation.Field(&x.CannotCompleteAfter, validation.Required),
	)
}
