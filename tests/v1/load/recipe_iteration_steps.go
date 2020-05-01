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

// fetchRandomRecipeIterationStep retrieves a random recipe iteration step from the list of available recipe iteration steps.
func fetchRandomRecipeIterationStep(ctx context.Context, c *client.V1Client, recipeID uint64) *models.RecipeIterationStep {
	recipeIterationStepsRes, err := c.GetRecipeIterationSteps(ctx, recipeID, nil)
	if err != nil || recipeIterationStepsRes == nil || len(recipeIterationStepsRes.RecipeIterationSteps) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeIterationStepsRes.RecipeIterationSteps))
	return &recipeIterationStepsRes.RecipeIterationSteps[randIndex]
}

func buildRecipeIterationStepActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeIterationStep": {
			Name: "CreateRecipeIterationStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create recipe.
				exampleRecipe := fakemodels.BuildFakeRecipe()
				exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
				createdRecipe, err := c.CreateRecipe(ctx, exampleRecipeInput)
				if err != nil {
					return nil, err
				}

				recipeIterationStepInput := fakemodels.BuildFakeRecipeIterationStepCreationInput()
				recipeIterationStepInput.BelongsToRecipe = createdRecipe.ID

				return c.BuildCreateRecipeIterationStepRequest(ctx, recipeIterationStepInput)
			},
			Weight: 100,
		},
		"GetRecipeIterationStep": {
			Name: "GetRecipeIterationStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeIterationStep := fetchRandomRecipeIterationStep(ctx, c, randomRecipe.ID)
				if randomRecipeIterationStep == nil {
					return nil, fmt.Errorf("retrieving random recipe iteration step: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeIterationStepRequest(ctx, randomRecipe.ID, randomRecipeIterationStep.ID)
			},
			Weight: 100,
		},
		"GetRecipeIterationSteps": {
			Name: "GetRecipeIterationSteps",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeIterationStepsRequest(ctx, randomRecipe.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeIterationStep": {
			Name: "UpdateRecipeIterationStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				if randomRecipeIterationStep := fetchRandomRecipeIterationStep(ctx, c, randomRecipe.ID); randomRecipeIterationStep != nil {
					newRecipeIterationStep := fakemodels.BuildFakeRecipeIterationStepCreationInput()
					randomRecipeIterationStep.StartedOn = newRecipeIterationStep.StartedOn
					randomRecipeIterationStep.EndedOn = newRecipeIterationStep.EndedOn
					randomRecipeIterationStep.State = newRecipeIterationStep.State
					return c.BuildUpdateRecipeIterationStepRequest(ctx, randomRecipeIterationStep)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeIterationStep": {
			Name: "ArchiveRecipeIterationStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeIterationStep := fetchRandomRecipeIterationStep(ctx, c, randomRecipe.ID)
				if randomRecipeIterationStep == nil {
					return nil, fmt.Errorf("retrieving random recipe iteration step: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeIterationStepRequest(ctx, randomRecipe.ID, randomRecipeIterationStep.ID)
			},
			Weight: 85,
		},
	}
}
