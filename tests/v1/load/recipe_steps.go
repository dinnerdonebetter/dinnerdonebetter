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

// fetchRandomRecipeStep retrieves a random recipe step from the list of available recipe steps.
func fetchRandomRecipeStep(ctx context.Context, c *client.V1Client, recipeID uint64) *models.RecipeStep {
	recipeStepsRes, err := c.GetRecipeSteps(ctx, recipeID, nil)
	if err != nil || recipeStepsRes == nil || len(recipeStepsRes.RecipeSteps) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepsRes.RecipeSteps))
	return &recipeStepsRes.RecipeSteps[randIndex]
}

func buildRecipeStepActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStep": {
			Name: "CreateRecipeStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create recipe.
				exampleRecipe := fakemodels.BuildFakeRecipe()
				exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
				createdRecipe, err := c.CreateRecipe(ctx, exampleRecipeInput)
				if err != nil {
					return nil, err
				}

				recipeStepInput := fakemodels.BuildFakeRecipeStepCreationInput()
				recipeStepInput.BelongsToRecipe = createdRecipe.ID

				return c.BuildCreateRecipeStepRequest(ctx, recipeStepInput)
			},
			Weight: 100,
		},
		"GetRecipeStep": {
			Name: "GetRecipeStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeStep := fetchRandomRecipeStep(ctx, c, randomRecipe.ID)
				if randomRecipeStep == nil {
					return nil, fmt.Errorf("retrieving random recipe step: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepRequest(ctx, randomRecipe.ID, randomRecipeStep.ID)
			},
			Weight: 100,
		},
		"GetRecipeSteps": {
			Name: "GetRecipeSteps",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepsRequest(ctx, randomRecipe.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStep": {
			Name: "UpdateRecipeStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				if randomRecipeStep := fetchRandomRecipeStep(ctx, c, randomRecipe.ID); randomRecipeStep != nil {
					newRecipeStep := fakemodels.BuildFakeRecipeStepCreationInput()
					randomRecipeStep.Index = newRecipeStep.Index
					randomRecipeStep.ValidPreparationID = newRecipeStep.ValidPreparationID
					randomRecipeStep.PrerequisiteStepID = newRecipeStep.PrerequisiteStepID
					randomRecipeStep.MinEstimatedTimeInSeconds = newRecipeStep.MinEstimatedTimeInSeconds
					randomRecipeStep.MaxEstimatedTimeInSeconds = newRecipeStep.MaxEstimatedTimeInSeconds
					randomRecipeStep.YieldsProductName = newRecipeStep.YieldsProductName
					randomRecipeStep.YieldsQuantity = newRecipeStep.YieldsQuantity
					randomRecipeStep.Notes = newRecipeStep.Notes
					return c.BuildUpdateRecipeStepRequest(ctx, randomRecipeStep)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStep": {
			Name: "ArchiveRecipeStep",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeStep := fetchRandomRecipeStep(ctx, c, randomRecipe.ID)
				if randomRecipeStep == nil {
					return nil, fmt.Errorf("retrieving random recipe step: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeStepRequest(ctx, randomRecipe.ID, randomRecipeStep.ID)
			},
			Weight: 85,
		},
	}
}
