package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipe retrieves a random recipe from the list of available recipes
func fetchRandomRecipe(c *client.V1Client) *models.Recipe {
	recipesRes, err := c.GetRecipes(context.Background(), nil)
	if err != nil || recipesRes == nil || len(recipesRes.Recipes) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipesRes.Recipes))
	return &recipesRes.Recipes[randIndex]
}

func buildRecipeActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipe": {
			Name: "CreateRecipe",
			Action: func() (*http.Request, error) {
				return c.BuildCreateRecipeRequest(context.Background(), randmodel.RandomRecipeCreationInput())
			},
			Weight: 100,
		},
		"GetRecipe": {
			Name: "GetRecipe",
			Action: func() (*http.Request, error) {
				if randomRecipe := fetchRandomRecipe(c); randomRecipe != nil {
					return c.BuildGetRecipeRequest(context.Background(), randomRecipe.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipes": {
			Name: "GetRecipes",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipesRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipe": {
			Name: "UpdateRecipe",
			Action: func() (*http.Request, error) {
				if randomRecipe := fetchRandomRecipe(c); randomRecipe != nil {
					randomRecipe.Name = randmodel.RandomRecipeCreationInput().Name
					randomRecipe.Source = randmodel.RandomRecipeCreationInput().Source
					randomRecipe.Description = randmodel.RandomRecipeCreationInput().Description
					randomRecipe.InspiredByRecipeID = randmodel.RandomRecipeCreationInput().InspiredByRecipeID
					return c.BuildUpdateRecipeRequest(context.Background(), randomRecipe)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipe": {
			Name: "ArchiveRecipe",
			Action: func() (*http.Request, error) {
				if randomRecipe := fetchRandomRecipe(c); randomRecipe != nil {
					return c.BuildArchiveRecipeRequest(context.Background(), randomRecipe.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
