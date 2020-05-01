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

// fetchRandomIterationMedia retrieves a random iteration media from the list of available iteration medias.
func fetchRandomIterationMedia(ctx context.Context, c *client.V1Client, recipeID, recipeIterationID uint64) *models.IterationMedia {
	iterationMediasRes, err := c.GetIterationMedias(ctx, recipeID, recipeIterationID, nil)
	if err != nil || iterationMediasRes == nil || len(iterationMediasRes.IterationMedias) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(iterationMediasRes.IterationMedias))
	return &iterationMediasRes.IterationMedias[randIndex]
}

func buildIterationMediaActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateIterationMedia": {
			Name: "CreateIterationMedia",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create recipe.
				exampleRecipe := fakemodels.BuildFakeRecipe()
				exampleRecipeInput := fakemodels.BuildFakeRecipeCreationInputFromRecipe(exampleRecipe)
				createdRecipe, err := c.CreateRecipe(ctx, exampleRecipeInput)
				if err != nil {
					return nil, err
				}

				// Create recipe iteration.
				exampleRecipeIteration := fakemodels.BuildFakeRecipeIteration()
				exampleRecipeIteration.BelongsToRecipe = createdRecipe.ID
				exampleRecipeIterationInput := fakemodels.BuildFakeRecipeIterationCreationInputFromRecipeIteration(exampleRecipeIteration)
				createdRecipeIteration, err := c.CreateRecipeIteration(ctx, exampleRecipeIterationInput)
				if err != nil {
					return nil, err
				}

				iterationMediaInput := fakemodels.BuildFakeIterationMediaCreationInput()
				iterationMediaInput.BelongsToRecipeIteration = createdRecipeIteration.ID

				return c.BuildCreateIterationMediaRequest(ctx, createdRecipe.ID, iterationMediaInput)
			},
			Weight: 100,
		},
		"GetIterationMedia": {
			Name: "GetIterationMedia",
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

				randomIterationMedia := fetchRandomIterationMedia(ctx, c, randomRecipe.ID, randomRecipeIteration.ID)
				if randomIterationMedia == nil {
					return nil, fmt.Errorf("retrieving random iteration media: %w", ErrUnavailableYet)
				}

				return c.BuildGetIterationMediaRequest(ctx, randomRecipe.ID, randomRecipeIteration.ID, randomIterationMedia.ID)
			},
			Weight: 100,
		},
		"GetIterationMedias": {
			Name: "GetIterationMedias",
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

				return c.BuildGetIterationMediasRequest(ctx, randomRecipe.ID, randomRecipeIteration.ID, nil)
			},
			Weight: 100,
		},
		"UpdateIterationMedia": {
			Name: "UpdateIterationMedia",
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

				if randomIterationMedia := fetchRandomIterationMedia(ctx, c, randomRecipe.ID, randomRecipeIteration.ID); randomIterationMedia != nil {
					newIterationMedia := fakemodels.BuildFakeIterationMediaCreationInput()
					randomIterationMedia.Source = newIterationMedia.Source
					randomIterationMedia.Mimetype = newIterationMedia.Mimetype
					return c.BuildUpdateIterationMediaRequest(ctx, randomRecipe.ID, randomIterationMedia)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveIterationMedia": {
			Name: "ArchiveIterationMedia",
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

				randomIterationMedia := fetchRandomIterationMedia(ctx, c, randomRecipe.ID, randomRecipeIteration.ID)
				if randomIterationMedia == nil {
					return nil, fmt.Errorf("retrieving random iteration media: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveIterationMediaRequest(ctx, randomRecipe.ID, randomRecipeIteration.ID, randomIterationMedia.ID)
			},
			Weight: 85,
		},
	}
}
