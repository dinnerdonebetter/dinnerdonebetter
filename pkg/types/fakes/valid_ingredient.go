package fakes

import (
	"gitlab.com/prixfixe/prixfixe/pkg/types"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidIngredient builds a faked valid ingredient.
func BuildFakeValidIngredient() *types.ValidIngredient {
	return &types.ValidIngredient{
		ID:                uint64(fake.Uint32()),
		ExternalID:        fake.UUID(),
		Name:              BuildUniqueName(),
		Variant:           fake.Word(),
		Description:       fake.Word(),
		Warning:           fake.Word(),
		ContainsEgg:       fake.Bool(),
		ContainsDairy:     fake.Bool(),
		ContainsPeanut:    fake.Bool(),
		ContainsTreeNut:   fake.Bool(),
		ContainsSoy:       fake.Bool(),
		ContainsWheat:     fake.Bool(),
		ContainsShellfish: fake.Bool(),
		ContainsSesame:    fake.Bool(),
		ContainsFish:      fake.Bool(),
		ContainsGluten:    fake.Bool(),
		AnimalFlesh:       fake.Bool(),
		AnimalDerived:     fake.Bool(),
		Volumetric:        fake.Bool(),
		IconPath:          fake.Word(),
		CreatedOn:         uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidIngredientList builds a faked ValidIngredientList.
func BuildFakeValidIngredientList() *types.ValidIngredientList {
	var examples []*types.ValidIngredient
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeValidIngredient())
	}

	return &types.ValidIngredientList{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		ValidIngredients: examples,
	}
}

// BuildFakeValidIngredientUpdateInput builds a faked ValidIngredientUpdateInput from a valid ingredient.
func BuildFakeValidIngredientUpdateInput() *types.ValidIngredientUpdateInput {
	validIngredient := BuildFakeValidIngredient()
	return &types.ValidIngredientUpdateInput{
		Name:              validIngredient.Name,
		Variant:           validIngredient.Variant,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		AnimalDerived:     validIngredient.AnimalDerived,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
	}
}

// BuildFakeValidIngredientUpdateInputFromValidIngredient builds a faked ValidIngredientUpdateInput from a valid ingredient.
func BuildFakeValidIngredientUpdateInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientUpdateInput {
	return &types.ValidIngredientUpdateInput{
		Name:              validIngredient.Name,
		Variant:           validIngredient.Variant,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		AnimalDerived:     validIngredient.AnimalDerived,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
	}
}

// BuildFakeValidIngredientCreationInput builds a faked ValidIngredientCreationInput.
func BuildFakeValidIngredientCreationInput() *types.ValidIngredientCreationInput {
	validIngredient := BuildFakeValidIngredient()
	return BuildFakeValidIngredientCreationInputFromValidIngredient(validIngredient)
}

// BuildFakeValidIngredientCreationInputFromValidIngredient builds a faked ValidIngredientCreationInput from a valid ingredient.
func BuildFakeValidIngredientCreationInputFromValidIngredient(validIngredient *types.ValidIngredient) *types.ValidIngredientCreationInput {
	return &types.ValidIngredientCreationInput{
		Name:              validIngredient.Name,
		Variant:           validIngredient.Variant,
		Description:       validIngredient.Description,
		Warning:           validIngredient.Warning,
		ContainsEgg:       validIngredient.ContainsEgg,
		ContainsDairy:     validIngredient.ContainsDairy,
		ContainsPeanut:    validIngredient.ContainsPeanut,
		ContainsTreeNut:   validIngredient.ContainsTreeNut,
		ContainsSoy:       validIngredient.ContainsSoy,
		ContainsWheat:     validIngredient.ContainsWheat,
		ContainsShellfish: validIngredient.ContainsShellfish,
		ContainsSesame:    validIngredient.ContainsSesame,
		ContainsFish:      validIngredient.ContainsFish,
		ContainsGluten:    validIngredient.ContainsGluten,
		AnimalFlesh:       validIngredient.AnimalFlesh,
		AnimalDerived:     validIngredient.AnimalDerived,
		Volumetric:        validIngredient.Volumetric,
		IconPath:          validIngredient.IconPath,
	}
}
