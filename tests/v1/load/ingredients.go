package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomIngredient retrieves a random ingredient from the list of available ingredients
func fetchRandomIngredient(c *client.V1Client) *models.Ingredient {
	ingredientsRes, err := c.GetIngredients(context.Background(), nil)
	if err != nil || ingredientsRes == nil || len(ingredientsRes.Ingredients) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(ingredientsRes.Ingredients))
	return &ingredientsRes.Ingredients[randIndex]
}

func buildIngredientActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateIngredient": {
			Name: "CreateIngredient",
			Action: func() (*http.Request, error) {
				return c.BuildCreateIngredientRequest(context.Background(), randmodel.RandomIngredientCreationInput())
			},
			Weight: 100,
		},
		"GetIngredient": {
			Name: "GetIngredient",
			Action: func() (*http.Request, error) {
				if randomIngredient := fetchRandomIngredient(c); randomIngredient != nil {
					return c.BuildGetIngredientRequest(context.Background(), randomIngredient.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetIngredients": {
			Name: "GetIngredients",
			Action: func() (*http.Request, error) {
				return c.BuildGetIngredientsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateIngredient": {
			Name: "UpdateIngredient",
			Action: func() (*http.Request, error) {
				if randomIngredient := fetchRandomIngredient(c); randomIngredient != nil {
					randomIngredient.Name = randmodel.RandomIngredientCreationInput().Name
					randomIngredient.Variant = randmodel.RandomIngredientCreationInput().Variant
					randomIngredient.Description = randmodel.RandomIngredientCreationInput().Description
					randomIngredient.Warning = randmodel.RandomIngredientCreationInput().Warning
					randomIngredient.ContainsEgg = randmodel.RandomIngredientCreationInput().ContainsEgg
					randomIngredient.ContainsDairy = randmodel.RandomIngredientCreationInput().ContainsDairy
					randomIngredient.ContainsPeanut = randmodel.RandomIngredientCreationInput().ContainsPeanut
					randomIngredient.ContainsTreeNut = randmodel.RandomIngredientCreationInput().ContainsTreeNut
					randomIngredient.ContainsSoy = randmodel.RandomIngredientCreationInput().ContainsSoy
					randomIngredient.ContainsWheat = randmodel.RandomIngredientCreationInput().ContainsWheat
					randomIngredient.ContainsShellfish = randmodel.RandomIngredientCreationInput().ContainsShellfish
					randomIngredient.ContainsSesame = randmodel.RandomIngredientCreationInput().ContainsSesame
					randomIngredient.ContainsFish = randmodel.RandomIngredientCreationInput().ContainsFish
					randomIngredient.ContainsGluten = randmodel.RandomIngredientCreationInput().ContainsGluten
					randomIngredient.AnimalFlesh = randmodel.RandomIngredientCreationInput().AnimalFlesh
					randomIngredient.AnimalDerived = randmodel.RandomIngredientCreationInput().AnimalDerived
					randomIngredient.ConsideredStaple = randmodel.RandomIngredientCreationInput().ConsideredStaple
					randomIngredient.Icon = randmodel.RandomIngredientCreationInput().Icon
					return c.BuildUpdateIngredientRequest(context.Background(), randomIngredient)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveIngredient": {
			Name: "ArchiveIngredient",
			Action: func() (*http.Request, error) {
				if randomIngredient := fetchRandomIngredient(c); randomIngredient != nil {
					return c.BuildArchiveIngredientRequest(context.Background(), randomIngredient.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
