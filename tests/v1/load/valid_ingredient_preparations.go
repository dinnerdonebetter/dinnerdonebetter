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

// fetchRandomValidIngredientPreparation retrieves a random valid ingredient preparation from the list of available valid ingredient preparations.
func fetchRandomValidIngredientPreparation(ctx context.Context, c *client.V1Client, validIngredientID uint64) *models.ValidIngredientPreparation {
	validIngredientPreparationsRes, err := c.GetValidIngredientPreparations(ctx, validIngredientID, nil)
	if err != nil || validIngredientPreparationsRes == nil || len(validIngredientPreparationsRes.ValidIngredientPreparations) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(validIngredientPreparationsRes.ValidIngredientPreparations))
	return &validIngredientPreparationsRes.ValidIngredientPreparations[randIndex]
}

func buildValidIngredientPreparationActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateValidIngredientPreparation": {
			Name: "CreateValidIngredientPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create valid ingredient.
				exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, err := c.CreateValidIngredient(ctx, exampleValidIngredientInput)
				if err != nil {
					return nil, err
				}

				validIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInput()
				validIngredientPreparationInput.BelongsToValidIngredient = createdValidIngredient.ID

				return c.BuildCreateValidIngredientPreparationRequest(ctx, validIngredientPreparationInput)
			},
			Weight: 100,
		},
		"GetValidIngredientPreparation": {
			Name: "GetValidIngredientPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				randomValidIngredientPreparation := fetchRandomValidIngredientPreparation(ctx, c, randomValidIngredient.ID)
				if randomValidIngredientPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient preparation: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidIngredientPreparationRequest(ctx, randomValidIngredient.ID, randomValidIngredientPreparation.ID)
			},
			Weight: 100,
		},
		"GetValidIngredientPreparations": {
			Name: "GetValidIngredientPreparations",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidIngredientPreparationsRequest(ctx, randomValidIngredient.ID, nil)
			},
			Weight: 100,
		},
		"UpdateValidIngredientPreparation": {
			Name: "UpdateValidIngredientPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				if randomValidIngredientPreparation := fetchRandomValidIngredientPreparation(ctx, c, randomValidIngredient.ID); randomValidIngredientPreparation != nil {
					newValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparationCreationInput()
					randomValidIngredientPreparation.Notes = newValidIngredientPreparation.Notes
					return c.BuildUpdateValidIngredientPreparationRequest(ctx, randomValidIngredientPreparation)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveValidIngredientPreparation": {
			Name: "ArchiveValidIngredientPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				randomValidIngredientPreparation := fetchRandomValidIngredientPreparation(ctx, c, randomValidIngredient.ID)
				if randomValidIngredientPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient preparation: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveValidIngredientPreparationRequest(ctx, randomValidIngredient.ID, randomValidIngredientPreparation.ID)
			},
			Weight: 85,
		},
	}
}
