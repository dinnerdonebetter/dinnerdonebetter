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

// fetchRandomRecipeStepPreparation retrieves a random recipe step preparation from the list of available recipe step preparations.
func fetchRandomRecipeStepPreparation(ctx context.Context, c *client.V1Client, recipeID, recipeStepID uint64) *models.RecipeStepPreparation {
	recipeStepPreparationsRes, err := c.GetRecipeStepPreparations(ctx, recipeID, recipeStepID, nil)
	if err != nil || recipeStepPreparationsRes == nil || len(recipeStepPreparationsRes.RecipeStepPreparations) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepPreparationsRes.RecipeStepPreparations))
	return &recipeStepPreparationsRes.RecipeStepPreparations[randIndex]
}

func buildRecipeStepPreparationActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepPreparation": {
			Name: "CreateRecipeStepPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create recipe.
				exampleRecipe := fakemodels.BuildFakeRecipe()
				exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
				createdRecipe, err := c.CreateRecipe(ctx, exampleRecipeInput)
				if err != nil {
					return nil, err
				}

				// Create recipe step.
				exampleRecipeStep := fakemodels.BuildFakeRecipeStep()
				exampleRecipeStep.BelongsToRecipe = createdRecipe.ID
				exampleRecipeStepInput := fakemodels.BuildFakeRecipeStepCreationInputFromRecipeStep(exampleRecipeStep)
				createdRecipeStep, err := c.CreateRecipeStep(ctx, exampleRecipeStepInput)
				if err != nil {
					return nil, err
				}

				recipeStepPreparationInput := fakemodels.BuildFakeRecipeStepPreparationCreationInput()
				recipeStepPreparationInput.BelongsToRecipeStep = createdRecipeStep.ID

				return c.BuildCreateRecipeStepPreparationRequest(ctx, createdRecipe.ID, recipeStepPreparationInput)
			},
			Weight: 100,
		},
		"GetRecipeStepPreparation": {
			Name: "GetRecipeStepPreparation",
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

				randomRecipeStepPreparation := fetchRandomRecipeStepPreparation(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepPreparation == nil {
					return nil, fmt.Errorf("retrieving random recipe step preparation: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepPreparationRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepPreparation.ID)
			},
			Weight: 100,
		},
		"GetRecipeStepPreparations": {
			Name: "GetRecipeStepPreparations",
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

				return c.BuildGetRecipeStepPreparationsRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepPreparation": {
			Name: "UpdateRecipeStepPreparation",
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

				if randomRecipeStepPreparation := fetchRandomRecipeStepPreparation(ctx, c, randomRecipe.ID, randomRecipeStep.ID); randomRecipeStepPreparation != nil {
					newRecipeStepPreparation := fakemodels.BuildFakeRecipeStepPreparationCreationInput()
					randomRecipeStepPreparation.ValidPreparationID = newRecipeStepPreparation.ValidPreparationID
					randomRecipeStepPreparation.Notes = newRecipeStepPreparation.Notes
					return c.BuildUpdateRecipeStepPreparationRequest(ctx, randomRecipe.ID, randomRecipeStepPreparation)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepPreparation": {
			Name: "ArchiveRecipeStepPreparation",
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

				randomRecipeStepPreparation := fetchRandomRecipeStepPreparation(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepPreparation == nil {
					return nil, fmt.Errorf("retrieving random recipe step preparation: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeStepPreparationRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepPreparation.ID)
			},
			Weight: 85,
		},
	}
}
