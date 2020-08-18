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

// fetchRandomRecipeStepInstrument retrieves a random recipe step instrument from the list of available recipe step instruments.
func fetchRandomRecipeStepInstrument(ctx context.Context, c *client.V1Client, recipeID, recipeStepID uint64) *models.RecipeStepInstrument {
	recipeStepInstrumentsRes, err := c.GetRecipeStepInstruments(ctx, recipeID, recipeStepID, nil)
	if err != nil || recipeStepInstrumentsRes == nil || len(recipeStepInstrumentsRes.RecipeStepInstruments) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepInstrumentsRes.RecipeStepInstruments))
	return &recipeStepInstrumentsRes.RecipeStepInstruments[randIndex]
}

func buildRecipeStepInstrumentActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepInstrument": {
			Name: "CreateRecipeStepInstrument",
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

				recipeStepInstrumentInput := fakemodels.BuildFakeRecipeStepInstrumentCreationInput()
				recipeStepInstrumentInput.BelongsToRecipeStep = createdRecipeStep.ID

				return c.BuildCreateRecipeStepInstrumentRequest(ctx, createdRecipe.ID, recipeStepInstrumentInput)
			},
			Weight: 100,
		},
		"GetRecipeStepInstrument": {
			Name: "GetRecipeStepInstrument",
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

				randomRecipeStepInstrument := fetchRandomRecipeStepInstrument(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepInstrument == nil {
					return nil, fmt.Errorf("retrieving random recipe step instrument: %w", ErrUnavailableYet)
				}

				return c.BuildGetRecipeStepInstrumentRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepInstrument.ID)
			},
			Weight: 100,
		},
		"GetRecipeStepInstruments": {
			Name: "GetRecipeStepInstruments",
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

				return c.BuildGetRecipeStepInstrumentsRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepInstrument": {
			Name: "UpdateRecipeStepInstrument",
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

				if randomRecipeStepInstrument := fetchRandomRecipeStepInstrument(ctx, c, randomRecipe.ID, randomRecipeStep.ID); randomRecipeStepInstrument != nil {
					newRecipeStepInstrument := fakemodels.BuildFakeRecipeStepInstrumentCreationInput()
					randomRecipeStepInstrument.InstrumentID = newRecipeStepInstrument.InstrumentID
					randomRecipeStepInstrument.RecipeStepID = newRecipeStepInstrument.RecipeStepID
					randomRecipeStepInstrument.Notes = newRecipeStepInstrument.Notes
					return c.BuildUpdateRecipeStepInstrumentRequest(ctx, randomRecipe.ID, randomRecipeStepInstrument)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepInstrument": {
			Name: "ArchiveRecipeStepInstrument",
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

				randomRecipeStepInstrument := fetchRandomRecipeStepInstrument(ctx, c, randomRecipe.ID, randomRecipeStep.ID)
				if randomRecipeStepInstrument == nil {
					return nil, fmt.Errorf("retrieving random recipe step instrument: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRecipeStepInstrumentRequest(ctx, randomRecipe.ID, randomRecipeStep.ID, randomRecipeStepInstrument.ID)
			},
			Weight: 85,
		},
	}
}
