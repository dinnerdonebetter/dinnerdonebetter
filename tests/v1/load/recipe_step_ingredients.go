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

// fetchRandomRecipeStepIngredient retrieves a random recipe step ingredient from the list of available recipe step ingredients.
func fetchRandomRecipeStepIngredient(ctx context.Context, c *client.V1Client, recipeID, recipeStepID uint64) *models.RecipeStepIngredient {
	recipeStepIngredientsRes, err := c.GetRecipeStepIngredients(ctx, recipeID, recipeStepID, nil)
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

				recipeStepIngredientInput := fakemodels.BuildFakeRecipeStepIngredientCreationInput()
				recipeStepIngredientInput.BelongsToRecipeStep = createdRecipeStep.ID

				return c.BuildCreateRecipeStepIngredientRequest(ctx, createdRecipe.ID, recipeStepIngredientInput)
			},
			Weight: 100,
		},
		"GetRecipeStepIngredient": {
			Name: "GetRecipeStepIngredient",
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

				randomRecipeStepIngredient := fetchRandomRecipeStepIngredient(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepIngredient == nil {
					return nil, fmt.Errorf("retrieving random recipe step ingredient: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepIngredientRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepIngredient.ID)
			},
			Weight: 100,
		},
		"GetRecipeStepIngredients": {
			Name: "GetRecipeStepIngredients",
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

				return c.BuildGetRecipeStepIngredientsRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepIngredient": {
			Name: "UpdateRecipeStepIngredient",
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

				if randomRecipeStepIngredient := fetchRandomRecipeStepIngredient(ctx, c, randomRecipe.ID, randomRecipeStep.ID); randomRecipeStepIngredient != nil {
					newRecipeStepIngredient := fakemodels.BuildFakeRecipeStepIngredientCreationInput()
					randomRecipeStepIngredient.IngredientID = newRecipeStepIngredient.IngredientID
					randomRecipeStepIngredient.QuantityType = newRecipeStepIngredient.QuantityType
					randomRecipeStepIngredient.QuantityValue = newRecipeStepIngredient.QuantityValue
					randomRecipeStepIngredient.QuantityNotes = newRecipeStepIngredient.QuantityNotes
					randomRecipeStepIngredient.ProductOfRecipe = newRecipeStepIngredient.ProductOfRecipe
					randomRecipeStepIngredient.IngredientNotes = newRecipeStepIngredient.IngredientNotes
					return c.BuildUpdateRecipeStepIngredientRequest(ctx, randomRecipe.ID, randomRecipeStepIngredient)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepIngredient": {
			Name: "ArchiveRecipeStepIngredient",
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

				randomRecipeStepIngredient := fetchRandomRecipeStepIngredient(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepIngredient == nil {
					return nil, fmt.Errorf("retrieving random recipe step ingredient: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeStepIngredientRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepIngredient.ID)
			},
			Weight: 85,
		},
	}
}
