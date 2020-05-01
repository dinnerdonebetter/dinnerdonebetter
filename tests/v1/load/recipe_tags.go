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

// fetchRandomRecipeTag retrieves a random recipe tag from the list of available recipe tags.
func fetchRandomRecipeTag(ctx context.Context, c *client.V1Client, recipeID uint64) *models.RecipeTag {
	recipeTagsRes, err := c.GetRecipeTags(ctx, recipeID, nil)
	if err != nil || recipeTagsRes == nil || len(recipeTagsRes.RecipeTags) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeTagsRes.RecipeTags))
	return &recipeTagsRes.RecipeTags[randIndex]
}

func buildRecipeTagActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeTag": {
			Name: "CreateRecipeTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create recipe.
				exampleRecipe := fakemodels.BuildFakeRecipe()
				exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
				createdRecipe, err := c.CreateRecipe(ctx, exampleRecipeInput)
				if err != nil {
					return nil, err
				}

				recipeTagInput := fakemodels.BuildFakeRecipeTagCreationInput()
				recipeTagInput.BelongsToRecipe = createdRecipe.ID

				return c.BuildCreateRecipeTagRequest(ctx, recipeTagInput)
			},
			Weight: 100,
		},
		"GetRecipeTag": {
			Name: "GetRecipeTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeTag := fetchRandomRecipeTag(ctx, c, randomRecipe.ID)
				if randomRecipeTag == nil {
					return nil, fmt.Errorf("retrieving random recipe tag: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeTagRequest(ctx, randomRecipe.ID, randomRecipeTag.ID)
			},
			Weight: 100,
		},
		"GetRecipeTags": {
			Name: "GetRecipeTags",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeTagsRequest(ctx, randomRecipe.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeTag": {
			Name: "UpdateRecipeTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				if randomRecipeTag := fetchRandomRecipeTag(ctx, c, randomRecipe.ID); randomRecipeTag != nil {
					newRecipeTag := fakemodels.BuildFakeRecipeTagCreationInput()
					randomRecipeTag.Name = newRecipeTag.Name
					return c.BuildUpdateRecipeTagRequest(ctx, randomRecipeTag)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeTag": {
			Name: "ArchiveRecipeTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRecipe := fetchRandomRecipe(ctx, c)
				if randomRecipe == nil {
					return nil, fmt.Errorf("retrieving random recipe: %w", ErrUnavailableYet)
				}

				randomRecipeTag := fetchRandomRecipeTag(ctx, c, randomRecipe.ID)
				if randomRecipeTag == nil {
					return nil, fmt.Errorf("retrieving random recipe tag: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeTagRequest(ctx, randomRecipe.ID, randomRecipeTag.ID)
			},
			Weight: 85,
		},
	}
}
