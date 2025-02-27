package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeRecipeStepVessel builds a faked recipe step vessel.
func BuildFakeRecipeStepVessel() *types.RecipeStepVessel {
	return &types.RecipeStepVessel{
		ID:                   BuildFakeID(),
		Vessel:               BuildFakeValidVessel(),
		Name:                 buildUniqueString(),
		RecipeStepProductID:  nil,
		Notes:                buildUniqueString(),
		CreatedAt:            BuildFakeTime(),
		BelongsToRecipeStep:  fake.UUID(),
		Quantity:             BuildFakeUint16RangeWithOptionalMax(),
		VesselPreposition:    buildUniqueString(),
		UnavailableAfterStep: fake.Bool(),
	}
}

// BuildFakeRecipeStepVesselsList builds a faked RecipeStepVesselList.
func BuildFakeRecipeStepVesselsList() *filtering.QueryFilteredResult[types.RecipeStepVessel] {
	var examples []*types.RecipeStepVessel
	for range exampleQuantity {
		examples = append(examples, BuildFakeRecipeStepVessel())
	}

	return &filtering.QueryFilteredResult[types.RecipeStepVessel]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
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
