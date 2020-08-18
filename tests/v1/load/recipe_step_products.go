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

// fetchRandomRecipeStepProduct retrieves a random recipe step product from the list of available recipe step products.
func fetchRandomRecipeStepProduct(ctx context.Context, c *client.V1Client, recipeID, recipeStepID uint64) *models.RecipeStepProduct {
	recipeStepProductsRes, err := c.GetRecipeStepProducts(ctx, recipeID, recipeStepID, nil)
	if err != nil || recipeStepProductsRes == nil || len(recipeStepProductsRes.RecipeStepProducts) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepProductsRes.RecipeStepProducts))
	return &recipeStepProductsRes.RecipeStepProducts[randIndex]
}

func buildRecipeStepProductActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepProduct": {
			Name: "CreateRecipeStepProduct",
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

				recipeStepProductInput := fakemodels.BuildFakeRecipeStepProductCreationInput()
				recipeStepProductInput.BelongsToRecipeStep = createdRecipeStep.ID

				return c.BuildCreateRecipeStepProductRequest(ctx, createdRecipe.ID, recipeStepProductInput)
			},
			Weight: 100,
		},
		"GetRecipeStepProduct": {
			Name: "GetRecipeStepProduct",
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

				randomRecipeStepProduct := fetchRandomRecipeStepProduct(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepProduct == nil {
					return nil, fmt.Errorf("retrieving random recipe step product: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepProductRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepProduct.ID)
			},
			Weight: 100,
		},
		"GetRecipeStepProducts": {
			Name: "GetRecipeStepProducts",
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

				return c.BuildGetRecipeStepProductsRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepProduct": {
			Name: "UpdateRecipeStepProduct",
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

				if randomRecipeStepProduct := fetchRandomRecipeStepProduct(ctx, c, randomRecipe.ID, randomRecipeStep.ID); randomRecipeStepProduct != nil {
					newRecipeStepProduct := fakemodels.BuildFakeRecipeStepProductCreationInput()
					randomRecipeStepProduct.Name = newRecipeStepProduct.Name
					randomRecipeStepProduct.RecipeStepID = newRecipeStepProduct.RecipeStepID
					return c.BuildUpdateRecipeStepProductRequest(ctx, randomRecipe.ID, randomRecipeStepProduct)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepProduct": {
			Name: "ArchiveRecipeStepProduct",
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

				randomRecipeStepProduct := fetchRandomRecipeStepProduct(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepProduct == nil {
					return nil, fmt.Errorf("retrieving random recipe step product: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeStepProductRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepProduct.ID)
			},
			Weight: 85,
		},
	}
}
