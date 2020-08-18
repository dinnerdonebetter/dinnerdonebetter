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

// fetchRandomRecipeStepEvent retrieves a random recipe step event from the list of available recipe step events.
func fetchRandomRecipeStepEvent(ctx context.Context, c *client.V1Client, recipeID, recipeStepID uint64) *models.RecipeStepEvent {
	recipeStepEventsRes, err := c.GetRecipeStepEvents(ctx, recipeID, recipeStepID, nil)
	if err != nil || recipeStepEventsRes == nil || len(recipeStepEventsRes.RecipeStepEvents) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepEventsRes.RecipeStepEvents))
	return &recipeStepEventsRes.RecipeStepEvents[randIndex]
}

func buildRecipeStepEventActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepEvent": {
			Name: "CreateRecipeStepEvent",
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

				recipeStepEventInput := fakemodels.BuildFakeRecipeStepEventCreationInput()
				recipeStepEventInput.BelongsToRecipeStep = createdRecipeStep.ID

				return c.BuildCreateRecipeStepEventRequest(ctx, createdRecipe.ID, recipeStepEventInput)
			},
			Weight: 100,
		},
		"GetRecipeStepEvent": {
			Name: "GetRecipeStepEvent",
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

				randomRecipeStepEvent := fetchRandomRecipeStepEvent(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepEvent == nil {
					return nil, fmt.Errorf("retrieving random recipe step event: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepEventRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepEvent.ID)
			},
			Weight: 100,
		},
		"GetRecipeStepEvents": {
			Name: "GetRecipeStepEvents",
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

				return c.BuildGetRecipeStepEventsRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepEvent": {
			Name: "UpdateRecipeStepEvent",
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

				if randomRecipeStepEvent := fetchRandomRecipeStepEvent(ctx, c, randomRecipe.ID, randomRecipeStep.ID); randomRecipeStepEvent != nil {
					newRecipeStepEvent := fakemodels.BuildFakeRecipeStepEventCreationInput()
					randomRecipeStepEvent.EventType = newRecipeStepEvent.EventType
					randomRecipeStepEvent.Done = newRecipeStepEvent.Done
					randomRecipeStepEvent.RecipeIterationID = newRecipeStepEvent.RecipeIterationID
					randomRecipeStepEvent.RecipeStepID = newRecipeStepEvent.RecipeStepID
					return c.BuildUpdateRecipeStepEventRequest(ctx, randomRecipe.ID, randomRecipeStepEvent)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepEvent": {
			Name: "ArchiveRecipeStepEvent",
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

				randomRecipeStepEvent := fetchRandomRecipeStepEvent(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepEvent == nil {
					return nil, fmt.Errorf("retrieving random recipe step event: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeStepEventRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepEvent.ID)
			},
			Weight: 85,
		},
	}
}
