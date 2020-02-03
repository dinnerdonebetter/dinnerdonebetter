package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipeStepIngredient retrieves a random recipe step ingredient from the list of available recipe step ingredients
func fetchRandomRecipeStepIngredient(c *client.V1Client) *models.RecipeStepIngredient {
	recipeStepIngredientsRes, err := c.GetRecipeStepIngredients(context.Background(), nil)
	if err != nil || recipeStepIngredientsRes == nil || len(recipeStepIngredientsRes.RecipeStepIngredients) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepIngredientsRes.RecipeStepIngredients))
	return &recipeStepIngredientsRes.RecipeStepIngredients[randIndex]
}

func buildRecipeStepIngredientActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepIngredient": {
			Name: "CreateRecipeStepIngredient",
			Action: func() (*http.Request, error) {
				return c.BuildCreateRecipeStepIngredientRequest(context.Background(), randmodel.RandomRecipeStepIngredientCreationInput())
			},
			Weight: 100,
		},
		"GetRecipeStepIngredient": {
			Name: "GetRecipeStepIngredient",
			Action: func() (*http.Request, error) {
				if randomRecipeStepIngredient := fetchRandomRecipeStepIngredient(c); randomRecipeStepIngredient != nil {
					return c.BuildGetRecipeStepIngredientRequest(context.Background(), randomRecipeStepIngredient.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipeStepIngredients": {
			Name: "GetRecipeStepIngredients",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipeStepIngredientsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepIngredient": {
			Name: "UpdateRecipeStepIngredient",
			Action: func() (*http.Request, error) {
				if randomRecipeStepIngredient := fetchRandomRecipeStepIngredient(c); randomRecipeStepIngredient != nil {
					randomRecipeStepIngredient.IngredientID = randmodel.RandomRecipeStepIngredientCreationInput().IngredientID
					randomRecipeStepIngredient.QuantityType = randmodel.RandomRecipeStepIngredientCreationInput().QuantityType
					randomRecipeStepIngredient.QuantityValue = randmodel.RandomRecipeStepIngredientCreationInput().QuantityValue
					randomRecipeStepIngredient.QuantityNotes = randmodel.RandomRecipeStepIngredientCreationInput().QuantityNotes
					randomRecipeStepIngredient.ProductOfRecipe = randmodel.RandomRecipeStepIngredientCreationInput().ProductOfRecipe
					randomRecipeStepIngredient.IngredientNotes = randmodel.RandomRecipeStepIngredientCreationInput().IngredientNotes
					randomRecipeStepIngredient.RecipeStepID = randmodel.RandomRecipeStepIngredientCreationInput().RecipeStepID
					return c.BuildUpdateRecipeStepIngredientRequest(context.Background(), randomRecipeStepIngredient)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepIngredient": {
			Name: "ArchiveRecipeStepIngredient",
			Action: func() (*http.Request, error) {
				if randomRecipeStepIngredient := fetchRandomRecipeStepIngredient(c); randomRecipeStepIngredient != nil {
					return c.BuildArchiveRecipeStepIngredientRequest(context.Background(), randomRecipeStepIngredient.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
