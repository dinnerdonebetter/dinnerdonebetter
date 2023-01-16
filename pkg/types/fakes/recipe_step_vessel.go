package fakes

import (
	fake "github.com/brianvoe/gofakeit/v5"

	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"
)

// BuildFakeRecipeStepVessel builds a faked recipe step vessel.
func BuildFakeRecipeStepVessel() *types.RecipeStepVessel {
	return &types.RecipeStepVessel{
		ID:                   BuildFakeID(),
		Instrument:           BuildFakeValidInstrument(),
		Name:                 buildUniqueString(),
		RecipeStepProductID:  nil,
		Notes:                buildUniqueString(),
		CreatedAt:            BuildFakeTime(),
		BelongsToRecipeStep:  fake.UUID(),
		MinimumQuantity:      fake.Uint32(),
		MaximumQuantity:      fake.Uint32(),
		VesselPredicate:      buildUniqueString(),
		UnavailableAfterStep: fake.Bool(),
	}
}

// BuildFakeRecipeStepVesselList builds a faked RecipeStepVesselList.
func BuildFakeRecipeStepVesselList() *types.QueryFilteredResult[types.RecipeStepVessel] {
	var examples []*types.RecipeStepVessel
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeRecipeStepVessel())
	}

	return &types.QueryFilteredResult[types.RecipeStepVessel]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeRecipeStepVesselUpdateRequestInput builds a faked RecipeStepVesselUpdateRequestInput from a recipe step vessel.
func BuildFakeRecipeStepVesselUpdateRequestInput() *types.RecipeStepVesselUpdateRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepVessel()
	return converters.ConvertRecipeStepVesselToRecipeStepVesselUpdateRequestInput(recipeStepInstrument)
}

// BuildFakeRecipeStepVesselCreationRequestInput builds a faked RecipeStepVesselCreationRequestInput.
func BuildFakeRecipeStepVesselCreationRequestInput() *types.RecipeStepVesselCreationRequestInput {
	recipeStepInstrument := BuildFakeRecipeStepVessel()
	return converters.ConvertRecipeStepVesselToRecipeStepVesselCreationRequestInput(recipeStepInstrument)
}
