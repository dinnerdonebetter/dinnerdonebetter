package fakemodels

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeValidIngredient builds a faked valid ingredient.
func BuildFakeValidIngredient() *models.ValidIngredient {
	return &models.ValidIngredient{
		ID:                 fake.Uint64(),
		Name:               fake.Word(),
		Variant:            fake.Word(),
		Description:        fake.Word(),
		Warning:            fake.Word(),
		ContainsEgg:        fake.Bool(),
		ContainsDairy:      fake.Bool(),
		ContainsPeanut:     fake.Bool(),
		ContainsTreeNut:    fake.Bool(),
		ContainsSoy:        fake.Bool(),
		ContainsWheat:      fake.Bool(),
		ContainsShellfish:  fake.Bool(),
		ContainsSesame:     fake.Bool(),
		ContainsFish:       fake.Bool(),
		ContainsGluten:     fake.Bool(),
		AnimalFlesh:        fake.Bool(),
		AnimalDerived:      fake.Bool(),
		MeasurableByVolume: fake.Bool(),
		Icon:               fake.Word(),
		CreatedOn:          uint64(uint32(fake.Date().Unix())),
	}
}

// BuildFakeValidIngredientList builds a faked ValidIngredientList.
func BuildFakeValidIngredientList() *models.ValidIngredientList {
	exampleValidIngredient1 := BuildFakeValidIngredient()
	exampleValidIngredient2 := BuildFakeValidIngredient()
	exampleValidIngredient3 := BuildFakeValidIngredient()

	return &models.ValidIngredientList{
		Pagination: models.Pagination{
			Page:       1,
			Limit:      20,
			TotalCount: 3,
		},
		ValidIngredients: []models.ValidIngredient{
			*exampleValidIngredient1,
			*exampleValidIngredient2,
			*exampleValidIngredient3,
		},
	}
}

// BuildFakeValidIngredientUpdateInputFromValidIngredient builds a faked ValidIngredientUpdateInput from a valid ingredient.
func BuildFakeValidIngredientUpdateInputFromValidIngredient(validIngredient *models.ValidIngredient) *models.ValidIngredientUpdateInput {
	return &models.ValidIngredientUpdateInput{
		Name:               validIngredient.Name,
		Variant:            validIngredient.Variant,
		Description:        validIngredient.Description,
		Warning:            validIngredient.Warning,
		ContainsEgg:        validIngredient.ContainsEgg,
		ContainsDairy:      validIngredient.ContainsDairy,
		ContainsPeanut:     validIngredient.ContainsPeanut,
		ContainsTreeNut:    validIngredient.ContainsTreeNut,
		ContainsSoy:        validIngredient.ContainsSoy,
		ContainsWheat:      validIngredient.ContainsWheat,
		ContainsShellfish:  validIngredient.ContainsShellfish,
		ContainsSesame:     validIngredient.ContainsSesame,
		ContainsFish:       validIngredient.ContainsFish,
		ContainsGluten:     validIngredient.ContainsGluten,
		AnimalFlesh:        validIngredient.AnimalFlesh,
		AnimalDerived:      validIngredient.AnimalDerived,
		MeasurableByVolume: validIngredient.MeasurableByVolume,
		Icon:               validIngredient.Icon,
	}
}

// BuildFakeValidIngredientCreationInput builds a faked ValidIngredientCreationInput.
func BuildFakeValidIngredientCreationInput() *models.ValidIngredientCreationInput {
	validIngredient := BuildFakeValidIngredient()
	return BuildFakeValidIngredientCreationInputFromValidIngredient(validIngredient)
}

// BuildFakeValidIngredientCreationInputFromValidIngredient builds a faked ValidIngredientCreationInput from a valid ingredient.
func BuildFakeValidIngredientCreationInputFromValidIngredient(validIngredient *models.ValidIngredient) *models.ValidIngredientCreationInput {
	return &models.ValidIngredientCreationInput{
		Name:               validIngredient.Name,
		Variant:            validIngredient.Variant,
		Description:        validIngredient.Description,
		Warning:            validIngredient.Warning,
		ContainsEgg:        validIngredient.ContainsEgg,
		ContainsDairy:      validIngredient.ContainsDairy,
		ContainsPeanut:     validIngredient.ContainsPeanut,
		ContainsTreeNut:    validIngredient.ContainsTreeNut,
		ContainsSoy:        validIngredient.ContainsSoy,
		ContainsWheat:      validIngredient.ContainsWheat,
		ContainsShellfish:  validIngredient.ContainsShellfish,
		ContainsSesame:     validIngredient.ContainsSesame,
		ContainsFish:       validIngredient.ContainsFish,
		ContainsGluten:     validIngredient.ContainsGluten,
		AnimalFlesh:        validIngredient.AnimalFlesh,
		AnimalDerived:      validIngredient.AnimalDerived,
		MeasurableByVolume: validIngredient.MeasurableByVolume,
		Icon:               validIngredient.Icon,
	}
}
