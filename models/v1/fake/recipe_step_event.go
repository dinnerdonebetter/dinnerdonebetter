package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepEvent builds a faked recipe step event.
func BuildFakeRecipeStepEvent() *models.RecipeStepEvent {
	return &models.RecipeStepEvent{
		ID:                  fake.Uint64(),
		EventType:           fake.Word(),
		Done:                fake.Bool(),
		RecipeIterationID:   uint64(fake.Uint32()),
		RecipeStepID:        uint64(fake.Uint32()),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
	}
}

// BuildFakeRecipeStepEventList builds a faked RecipeStepEventList.
func BuildFakeRecipeStepEventList() *models.RecipeStepEventList {
	exampleRecipeStepEvent1 := BuildFakeRecipeStepEvent()
	exampleRecipeStepEvent2 := BuildFakeRecipeStepEvent()
	exampleRecipeStepEvent3 := BuildFakeRecipeStepEvent()

	return &models.RecipeStepEventList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
		},
		RecipeStepEvents: []models.RecipeStepEvent{
			*exampleRecipeStepEvent1,
			*exampleRecipeStepEvent2,
			*exampleRecipeStepEvent3,
		},
	}
}

// BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent builds a faked RecipeStepEventUpdateInput from a recipe step event.
func BuildFakeRecipeStepEventUpdateInputFromRecipeStepEvent(recipeStepEvent *models.RecipeStepEvent) *models.RecipeStepEventUpdateInput {
	return &models.RecipeStepEventUpdateInput{
		EventType:           recipeStepEvent.EventType,
		Done:                recipeStepEvent.Done,
		RecipeIterationID:   recipeStepEvent.RecipeIterationID,
		RecipeStepID:        recipeStepEvent.RecipeStepID,
		BelongsToRecipeStep: recipeStepEvent.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepEventCreationInput builds a faked RecipeStepEventCreationInput.
func BuildFakeRecipeStepEventCreationInput() *models.RecipeStepEventCreationInput {
	recipeStepEvent := BuildFakeRecipeStepEvent()
	return BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(recipeStepEvent)
}

// BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent builds a faked RecipeStepEventCreationInput from a recipe step event.
func BuildFakeRecipeStepEventCreationInputFromRecipeStepEvent(recipeStepEvent *models.RecipeStepEvent) *models.RecipeStepEventCreationInput {
	return &models.RecipeStepEventCreationInput{
		EventType:           recipeStepEvent.EventType,
		Done:                recipeStepEvent.Done,
		RecipeIterationID:   recipeStepEvent.RecipeIterationID,
		RecipeStepID:        recipeStepEvent.RecipeStepID,
		BelongsToRecipeStep: recipeStepEvent.BelongsToRecipeStep,
	}
}
