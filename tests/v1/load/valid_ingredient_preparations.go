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
func fetchRandomValidIngredientPreparation(ctx context.Context, c *client.V1Client) *models.ValidIngredientPreparation {
	validIngredientPreparationsRes, err := c.GetValidIngredientPreparations(ctx, nil)
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

				validIngredientPreparationInput := fakemodels.BuildFakeValidIngredientPreparationCreationInput()

				return c.BuildCreateValidIngredientPreparationRequest(ctx, validIngredientPreparationInput)
			},
			Weight: 100,
		},
		"GetValidIngredientPreparation": {
			Name: "GetValidIngredientPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredientPreparation := fetchRandomValidIngredientPreparation(ctx, c)
				if randomValidIngredientPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient preparation: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidIngredientPreparationRequest(ctx, randomValidIngredientPreparation.ID)
			},
			Weight: 100,
		},
		"GetValidIngredientPreparations": {
			Name: "GetValidIngredientPreparations",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetValidIngredientPreparationsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateValidIngredientPreparation": {
			Name: "UpdateValidIngredientPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomValidIngredientPreparation := fetchRandomValidIngredientPreparation(ctx, c); randomValidIngredientPreparation != nil {
					newValidIngredientPreparation := fakemodels.BuildFakeValidIngredientPreparationCreationInput()
					randomValidIngredientPreparation.Notes = newValidIngredientPreparation.Notes
					randomValidIngredientPreparation.ValidPreparationID = newValidIngredientPreparation.ValidPreparationID
					randomValidIngredientPreparation.ValidIngredientID = newValidIngredientPreparation.ValidIngredientID
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

				randomValidIngredientPreparation := fetchRandomValidIngredientPreparation(ctx, c)
				if randomValidIngredientPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient preparation: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveValidIngredientPreparationRequest(ctx, randomValidIngredientPreparation.ID)
			},
			Weight: 85,
		},
	}
}
