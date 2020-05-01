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

// fetchRandomRecipeIteration retrieves a random recipe iteration from the list of available recipe iterations.
func fetchRandomRecipeIteration(ctx context.Context, c *client.V1Client, recipeID uint64) *models.RecipeIteration {
	recipeIterationsRes, err := c.GetRecipeIterations(ctx, recipeID, nil)
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
				ctx := context.Background()

				// Create recipe.
				exampleRecipe := fakemodels.BuildFakeRecipe()
				exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
				createdRecipe, err := c.CreateRecipe(ctx, exampleRecipeInput)
				if err != nil {
					return nil, err
				}

				recipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInput()
				recipeIterationInput.BelongsToRecipe = createdRecipe.ID

				return c.BuildCreateRecipeIterationRequest(ctx, recipeIterationInput)
			},
			Weight: 100,
		},
		"GetRecipeIteration": {
			Name: "GetRecipeIteration",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeIteration := fetchRandomRecipeIteration(ctx, c, randomRecipe.ID)
				if randomRecipeIteration == nil {
					return nil, fmt.Errorf("retrieving random recipe iteration: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeIterationRequest(ctx, randomRecipe.ID, randomRecipeIteration.ID)
			},
			Weight: 100,
		},
		"GetRecipeIterations": {
			Name: "GetRecipeIterations",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeIterationsRequest(ctx, randomRecipe.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeIteration": {
			Name: "UpdateRecipeIteration",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				if randomRecipeIteration := fetchRandomRecipeIteration(ctx, c, randomRecipe.ID); randomRecipeIteration != nil {
					newRecipeIteration := fakemodels.BuildFakeRecipeIterationCreationInput()
					randomRecipeIteration.EndDifficultyRating = newRecipeIteration.EndDifficultyRating
					randomRecipeIteration.EndComplexityRating = newRecipeIteration.EndComplexityRating
					randomRecipeIteration.EndTasteRating = newRecipeIteration.EndTasteRating
					randomRecipeIteration.EndOverallRating = newRecipeIteration.EndOverallRating
					return c.BuildUpdateRecipeIterationRequest(ctx, randomRecipeIteration)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeIteration": {
			Name: "ArchiveRecipeIteration",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeIteration := fetchRandomRecipeIteration(ctx, c, randomRecipe.ID)
				if randomRecipeIteration == nil {
					return nil, fmt.Errorf("retrieving random recipe iteration: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeIterationRequest(ctx, randomRecipe.ID, randomRecipeIteration.ID)
			},
			Weight: 85,
		},
	}
}
