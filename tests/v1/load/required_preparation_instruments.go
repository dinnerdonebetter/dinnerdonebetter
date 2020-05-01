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

// fetchRandomRequiredPreparationInstrument retrieves a random required preparation instrument from the list of available required preparation instruments.
func fetchRandomRequiredPreparationInstrument(ctx context.Context, c *client.V1Client, validPreparationID uint64) *models.RequiredPreparationInstrument {
	requiredPreparationInstrumentsRes, err := c.GetRequiredPreparationInstruments(ctx, validPreparationID, nil)
	if err != nil || requiredPreparationInstrumentsRes == nil || len(requiredPreparationInstrumentsRes.RequiredPreparationInstruments) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(requiredPreparationInstrumentsRes.RequiredPreparationInstruments))
	return &requiredPreparationInstrumentsRes.RequiredPreparationInstruments[randIndex]
}

func buildRequiredPreparationInstrumentActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRequiredPreparationInstrument": {
			Name: "CreateRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create valid preparation.
				exampleValidPreparation := fakemodels.BuildFakeValidPreparation()
				exampleValidPreparationInput := fakemodels.BuildFakeValidPreparationCreationInputFromValidPreparation(exampleValidPreparation)
				createdValidPreparation, err := c.CreateValidPreparation(ctx, exampleValidPreparationInput)
				if err != nil {
					return nil, err
				}

				requiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInput()
				requiredPreparationInstrumentInput.BelongsToValidPreparation = createdValidPreparation.ID

				return c.BuildCreateRequiredPreparationInstrumentRequest(ctx, requiredPreparationInstrumentInput)
			},
			Weight: 100,
		},
		"GetRequiredPreparationInstrument": {
			Name: "GetRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidPreparation := fetchRandomValidPreparation(ctx, c)
				if randomValidPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid preparation: %w", ErrUnavailableYet)
				}

				randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(ctx, c, randomValidPreparation.ID)
				if randomRequiredPreparationInstrument == nil {
					return nil, fmt.Errorf("retrieving random required preparation instrument: %w", ErrUnavailableYet)
				}

				return c.BuildGetRequiredPreparationInstrumentRequest(ctx, randomValidPreparation.ID, randomRequiredPreparationInstrument.ID)
			},
			Weight: 100,
		},
		"GetRequiredPreparationInstruments": {
			Name: "GetRequiredPreparationInstruments",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidPreparation := fetchRandomValidPreparation(ctx, c)
				if randomValidPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid preparation: %w", ErrUnavailableYet)
				}

				return c.BuildGetRequiredPreparationInstrumentsRequest(ctx, randomValidPreparation.ID, nil)
			},
			Weight: 100,
		},
		"UpdateRequiredPreparationInstrument": {
			Name: "UpdateRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidPreparation := fetchRandomValidPreparation(ctx, c)
				if randomValidPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid preparation: %w", ErrUnavailableYet)
				}

				if randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(ctx, c, randomValidPreparation.ID); randomRequiredPreparationInstrument != nil {
					newRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInput()
					randomRequiredPreparationInstrument.ValidInstrumentID = newRequiredPreparationInstrument.ValidInstrumentID
					randomRequiredPreparationInstrument.Notes = newRequiredPreparationInstrument.Notes
					return c.BuildUpdateRequiredPreparationInstrumentRequest(ctx, randomRequiredPreparationInstrument)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRequiredPreparationInstrument": {
			Name: "ArchiveRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidPreparation := fetchRandomValidPreparation(ctx, c)
				if randomValidPreparation == nil {
					return nil, fmt.Errorf("retrieving random valid preparation: %w", ErrUnavailableYet)
				}

				randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(ctx, c, randomValidPreparation.ID)
				if randomRequiredPreparationInstrument == nil {
					return nil, fmt.Errorf("retrieving random required preparation instrument: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRequiredPreparationInstrumentRequest(ctx, randomValidPreparation.ID, randomRequiredPreparationInstrument.ID)
			},
			Weight: 85,
		},
	}
}
