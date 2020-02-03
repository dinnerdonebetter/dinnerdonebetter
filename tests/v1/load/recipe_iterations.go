package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipeIteration retrieves a random recipe iteration from the list of available recipe iterations
func fetchRandomRecipeIteration(c *client.V1Client) *models.RecipeIteration {
	recipeIterationsRes, err := c.GetRecipeIterations(context.Background(), nil)
	if err != nil || recipeIterationsRes == nil || len(recipeIterationsRes.RecipeIterations) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeIterationsRes.RecipeIterations))
	return &recipeIterationsRes.RecipeIterations[randIndex]
}

func buildRecipeIterationActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeIteration": {
			Name: "CreateRecipeIteration",
			Action: func() (*http.Request, error) {
				return c.BuildCreateRecipeIterationRequest(context.Background(), randmodel.RandomRecipeIterationCreationInput())
			},
			Weight: 100,
		},
		"GetRecipeIteration": {
			Name: "GetRecipeIteration",
			Action: func() (*http.Request, error) {
				if randomRecipeIteration := fetchRandomRecipeIteration(c); randomRecipeIteration != nil {
					return c.BuildGetRecipeIterationRequest(context.Background(), randomRecipeIteration.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipeIterations": {
			Name: "GetRecipeIterations",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipeIterationsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipeIteration": {
			Name: "UpdateRecipeIteration",
			Action: func() (*http.Request, error) {
				if randomRecipeIteration := fetchRandomRecipeIteration(c); randomRecipeIteration != nil {
					randomRecipeIteration.RecipeID = randmodel.RandomRecipeIterationCreationInput().RecipeID
					randomRecipeIteration.EndDifficultyRating = randmodel.RandomRecipeIterationCreationInput().EndDifficultyRating
					randomRecipeIteration.EndComplexityRating = randmodel.RandomRecipeIterationCreationInput().EndComplexityRating
					randomRecipeIteration.EndTasteRating = randmodel.RandomRecipeIterationCreationInput().EndTasteRating
					randomRecipeIteration.EndOverallRating = randmodel.RandomRecipeIterationCreationInput().EndOverallRating
					return c.BuildUpdateRecipeIterationRequest(context.Background(), randomRecipeIteration)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeIteration": {
			Name: "ArchiveRecipeIteration",
			Action: func() (*http.Request, error) {
				if randomRecipeIteration := fetchRandomRecipeIteration(c); randomRecipeIteration != nil {
					return c.BuildArchiveRecipeIterationRequest(context.Background(), randomRecipeIteration.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
