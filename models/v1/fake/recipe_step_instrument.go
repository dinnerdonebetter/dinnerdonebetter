package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeRecipeStepInstrument builds a faked recipe step instrument.
func BuildFakeRecipeStepInstrument() *models.RecipeStepInstrument {
	return &models.RecipeStepInstrument{
		ID:                  fake.Uint64(),
		InstrumentID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		RecipeStepID:        uint64(fake.Uint32()),
		Notes:               fake.Word(),
		CreatedOn:           uint64(uint32(fake.Date().Unix())),
		BelongsToRecipeStep: fake.Uint64(),
	}
}

// BuildFakeRecipeStepInstrumentList builds a faked RecipeStepInstrumentList.
func BuildFakeRecipeStepInstrumentList() *models.RecipeStepInstrumentList {
	exampleRecipeStepInstrument1 := BuildFakeRecipeStepInstrument()
	exampleRecipeStepInstrument2 := BuildFakeRecipeStepInstrument()
	exampleRecipeStepInstrument3 := BuildFakeRecipeStepInstrument()

	return &models.RecipeStepInstrumentList{
		Pagination: models.Pagination{
			Page:  1,
			Limit: 20,
		},
		RecipeStepInstruments: []models.RecipeStepInstrument{
			*exampleRecipeStepInstrument1,
			*exampleRecipeStepInstrument2,
			*exampleRecipeStepInstrument3,
		},
	}
}

// BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentUpdateInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentUpdateInputFromRecipeStepInstrument(recipeStepInstrument *models.RecipeStepInstrument) *models.RecipeStepInstrumentUpdateInput {
	return &models.RecipeStepInstrumentUpdateInput{
		InstrumentID:        recipeStepInstrument.InstrumentID,
		RecipeStepID:        recipeStepInstrument.RecipeStepID,
		Notes:               recipeStepInstrument.Notes,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
	}
}

// BuildFakeRecipeStepInstrumentCreationInput builds a faked RecipeStepInstrumentCreationInput.
func BuildFakeRecipeStepInstrumentCreationInput() *models.RecipeStepInstrumentCreationInput {
	recipeStepInstrument := BuildFakeRecipeStepInstrument()
	return BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(recipeStepInstrument)
}

// BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument builds a faked RecipeStepInstrumentCreationInput from a recipe step instrument.
func BuildFakeRecipeStepInstrumentCreationInputFromRecipeStepInstrument(recipeStepInstrument *models.RecipeStepInstrument) *models.RecipeStepInstrumentCreationInput {
	return &models.RecipeStepInstrumentCreationInput{
		InstrumentID:        recipeStepInstrument.InstrumentID,
		RecipeStepID:        recipeStepInstrument.RecipeStepID,
		Notes:               recipeStepInstrument.Notes,
		BelongsToRecipeStep: recipeStepInstrument.BelongsToRecipeStep,
	}
}
