package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"
)

// fetchRandomRecipe retrieves a random recipe from the list of available recipes.
func fetchRandomRecipe(ctx context.Context, c *client.V1Client) *models.Recipe {
	recipesRes, err := c.GetRecipes(ctx, nil)
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
				ctx := context.Background()

				recipeInput := fakemodels.BuildFakeRecipeCreationInput()

				return c.BuildCreateRecipeRequest(ctx, recipeInput)
			},
			Weight: 100,
		},
		"GetRecipe": {
			Name: "GetRecipe",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeRequest(ctx, randomRecipe.ID)
			},
			Weight: 100,
		},
		"GetRecipes": {
			Name: "GetRecipes",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetRecipesRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateRecipe": {
			Name: "UpdateRecipe",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomRecipe := fetchRandomRecipe(ctx, c); randomRecipe != nil {
					newRecipe := fakemodels.BuildFakeRecipeCreationInput()
					randomRecipe.Name = newRecipe.Name
					randomRecipe.Source = newRecipe.Source
					randomRecipe.Description = newRecipe.Description
					randomRecipe.InspiredByRecipeID = newRecipe.InspiredByRecipeID
					return c.BuildUpdateRecipeRequest(ctx, randomRecipe)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipe": {
			Name: "ArchiveRecipe",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeRequest(ctx, randomRecipe.ID)
			},
			Weight: 85,
		},
	}
}
