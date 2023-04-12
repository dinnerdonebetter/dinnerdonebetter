package fakes

import (
	"github.com/prixfixeco/backend/internal/pkg/pointers"
	"github.com/prixfixeco/backend/pkg/types"
	"github.com/prixfixeco/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidIngredient builds a faked valid ingredient.
func BuildFakeValidIngredient() *types.ValidIngredient {
	return &types.ValidIngredient{
		ID:                                      BuildFakeID(),
		Name:                                    buildUniqueString(),
		Description:                             buildUniqueString(),
		Warning:                                 buildUniqueString(),
		ContainsEgg:                             fake.Bool(),
		ContainsDairy:                           fake.Bool(),
		ContainsPeanut:                          fake.Bool(),
		ContainsTreeNut:                         fake.Bool(),
		ContainsSoy:                             fake.Bool(),
		ContainsWheat:                           fake.Bool(),
		ContainsShellfish:                       fake.Bool(),
		ContainsSesame:                          fake.Bool(),
		ContainsFish:                            fake.Bool(),
		ContainsGluten:                          fake.Bool(),
		AnimalFlesh:                             fake.Bool(),
		IsMeasuredVolumetrically:                fake.Bool(),
		IsLiquid:                                fake.Bool(),
		IconPath:                                buildUniqueString(),
		CreatedAt:                               BuildFakeTime(),
		PluralName:                              buildUniqueString(),
		AnimalDerived:                           fake.Bool(),
		RestrictToPreparations:                  fake.Bool(),
		MinimumIdealStorageTemperatureInCelsius: pointers.Float32(float32(BuildFakeNumber())),
		MaximumIdealStorageTemperatureInCelsius: pointers.Float32(float32(BuildFakeNumber())),
		StorageInstructions:                     buildUniqueString(),
		Slug:                                    buildUniqueString(),
		ContainsAlcohol:                         fake.Bool(),
		ShoppingSuggestions:                     buildUniqueString(),
	}
}

// BuildFakeValidIngredientList builds a faked ValidIngredientList.
func BuildFakeValidIngredientList() *types.QueryFilteredResult[types.ValidIngredient] {
	var examples []*types.ValidIngredient
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredient())
	}

	return &types.QueryFilteredResult[types.ValidIngredient]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeValidIngredientUpdateRequestInput builds a faked ValidIngredientUpdateRequestInput from a valid ingredient.
func BuildFakeValidIngredientUpdateRequestInput() *types.ValidIngredientUpdateRequestInput {
	validIngredient := BuildFakeValidIngredient()
	return converters.ConvertValidIngredientToValidIngredientUpdateRequestInput(validIngredient)
}

// BuildFakeValidIngredientCreationRequestInput builds a faked ValidIngredientCreationRequestInput.
func BuildFakeValidIngredientCreationRequestInput() *types.ValidIngredientCreationRequestInput {
	validIngredient := BuildFakeValidIngredient()
	return converters.ConvertValidIngredientToValidIngredientCreationRequestInput(validIngredient)
}
