package fakes

import (
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types/converters"

	fake "github.com/brianvoe/gofakeit/v7"
)

// BuildFakeIngredientPreference builds a faked valid preparation.
func BuildFakeIngredientPreference() *types.IngredientPreference {
	return &types.IngredientPreference{
		ID:         BuildFakeID(),
		Ingredient: *BuildFakeValidIngredient(),
		Rating:     1,
		Notes:      buildUniqueString(),
		Allergy:    fake.Bool(),
		CreatedAt:  BuildFakeTime(),
	}
}

// BuildFakeIngredientPreferencesList builds a faked IngredientPreferenceList.
func BuildFakeIngredientPreferencesList() *filtering.QueryFilteredResult[types.IngredientPreference] {
	var examples []*types.IngredientPreference
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeIngredientPreference())
	}

	return &filtering.QueryFilteredResult[types.IngredientPreference]{
		Pagination: filtering.Pagination{
			Page:          1,
			Limit:         50,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeIngredientPreferenceUpdateRequestInput builds a faked IngredientPreferenceUpdateRequestInput from a valid preparation.
func BuildFakeIngredientPreferenceUpdateRequestInput() *types.IngredientPreferenceUpdateRequestInput {
	validPreparation := BuildFakeIngredientPreference()
	return converters.ConvertIngredientPreferenceToIngredientPreferenceUpdateRequestInput(validPreparation)
}

// BuildFakeIngredientPreferenceCreationRequestInput builds a faked IngredientPreferenceCreationRequestInput.
func BuildFakeIngredientPreferenceCreationRequestInput() *types.IngredientPreferenceCreationRequestInput {
	validPreparation := BuildFakeIngredientPreference()
	return converters.ConvertIngredientPreferenceToIngredientPreferenceCreationRequestInput(validPreparation)
}
