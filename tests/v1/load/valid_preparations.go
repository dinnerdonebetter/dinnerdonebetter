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

// fetchRandomValidPreparation retrieves a random valid preparation from the list of available valid preparations.
func fetchRandomValidPreparation(ctx context.Context, c *client.V1Client) *models.ValidPreparation {
	validPreparationsRes, err := c.GetValidPreparations(ctx, nil)
	if err != nil || validPreparationsRes == nil || len(validPreparationsRes.ValidPreparations) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(validPreparationsRes.ValidPreparations))
	return &validPreparationsRes.ValidPreparations[randIndex]
}

func buildValidPreparationActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateValidPreparation": {
			Name: "CreateValidPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				validPreparationInput := fakemodels.BuildFakeValidPreparationCreationInput()

				return c.BuildCreateValidPreparationRequest(ctx, validPreparationInput)
			},
			Weight: 100,
		},
		"GetValidPreparation": {
			Name: "GetValidPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidPreparation := fetchRandomValidPreparation(ctx, c)
				if randomValidPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid preparation: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidPreparationRequest(ctx, randomValidPreparation.ID)
			},
			Weight: 100,
		},
		"GetValidPreparations": {
			Name: "GetValidPreparations",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetValidPreparationsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateValidPreparation": {
			Name: "UpdateValidPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomValidPreparation := fetchRandomValidPreparation(ctx, c); randomValidPreparation != nil {
					newValidPreparation := fakemodels.BuildFakeValidPreparationCreationInput()
					randomValidPreparation.Name = newValidPreparation.Name
					randomValidPreparation.Description = newValidPreparation.Description
					randomValidPreparation.Icon = newValidPreparation.Icon
					return c.BuildUpdateValidPreparationRequest(ctx, randomValidPreparation)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveValidPreparation": {
			Name: "ArchiveValidPreparation",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidPreparation := fetchRandomValidPreparation(ctx, c)
				if randomValidPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid preparation: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveValidPreparationRequest(ctx, randomValidPreparation.ID)
			},
			Weight: 85,
		},
	}
}
