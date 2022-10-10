package types

import (
	"encoding/gob"
)

const (
	// RecipePrepTaskStepDataType indicates an event is related to a recipe prep task step.
	RecipePrepTaskStepDataType dataType = "recipe_prep_step"

	// RecipePrepTaskStepCreatedCustomerEventType indicates a recipe prep task step was created.
	RecipePrepTaskStepCreatedCustomerEventType CustomerEventType = "recipe_created"
	// RecipePrepTaskStepUpdatedCustomerEventType indicates a recipe prep task step was updated.
	RecipePrepTaskStepUpdatedCustomerEventType CustomerEventType = "recipe_updated"
	// RecipePrepTaskStepArchivedCustomerEventType indicates a recipe prep task step was archived.
	RecipePrepTaskStepArchivedCustomerEventType CustomerEventType = "recipe_archived"
)

func init() {
	gob.Register(new(RecipePrepTask))
	gob.Register(new(RecipePrepTaskList))
	gob.Register(new(RecipePrepTaskCreationRequestInput))
	gob.Register(new(RecipePrepTaskUpdateRequestInput))
}

type (
	// RecipePrepTaskStep represents a recipe prep task step.
	RecipePrepTaskStep struct {
		_                       struct{}
		ID                      string `json:"id"`
		BelongsToRecipeStep     string `json:"belongsToRecipeStep"`
		BelongsToRecipePrepTask string `json:"belongsToRecipeStepTask"`
		SatisfiesRecipeStep     bool   `json:"satisfiesRecipeStep"`
	}

	// RecipePrepTaskStepCreationRequestInput represents a recipe prep task step.
	RecipePrepTaskStepCreationRequestInput struct {
		_                       struct{}
		ID                      string `json:"id"`
		BelongsToRecipeStep     string `json:"belongsToRecipeStep"`
		BelongsToRecipePrepTask string `json:"-"`
		SatisfiesRecipeStep     bool   `json:"satisfiesRecipeStep"`
	}

	// RecipePrepTaskStepDatabaseCreationInput represents a recipe prep task step.
	RecipePrepTaskStepDatabaseCreationInput struct {
		_                       struct{}
		ID                      string `json:"id"`
		BelongsToRecipeStep     string `json:"belongsToRecipeStep"`
		BelongsToRecipePrepTask string `json:"belongsToRecipeStepTask"`
		SatisfiesRecipeStep     bool   `json:"satisfiesRecipeStep"`
	}

	// RecipePrepTaskStepUpdateRequestInput represents a recipe prep task step.
	RecipePrepTaskStepUpdateRequestInput struct {
		_                       struct{}
		SatisfiesRecipeStep     *bool   `json:"satisfiesRecipeStep"`
		BelongsToRecipeStep     *string `json:"belongsToRecipeStep"`
		BelongsToRecipePrepTask *string `json:"belongsToRecipeStepTask"`
		ID                      string  `json:"-"`
	}
)
