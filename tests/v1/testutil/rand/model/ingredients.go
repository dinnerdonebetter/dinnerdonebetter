package randmodel

import (
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	fake "github.com/brianvoe/gofakeit"
)

// RandomIngredientCreationInput creates a random IngredientInput
func RandomIngredientCreationInput() *models.IngredientCreationInput {
	x := &models.IngredientCreationInput{
		Name:              fake.Word(),
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
		ConsideredStaple:  fake.Bool(),
		Icon:              fake.Word(),
	}

	return x
}
