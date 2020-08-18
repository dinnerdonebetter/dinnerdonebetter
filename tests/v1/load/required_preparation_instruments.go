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
func fetchRandomRequiredPreparationInstrument(ctx context.Context, c *client.V1Client) *models.RequiredPreparationInstrument {
	requiredPreparationInstrumentsRes, err := c.GetRequiredPreparationInstruments(ctx, nil)
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

				requiredPreparationInstrumentInput := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInput()

				return c.BuildCreateRequiredPreparationInstrumentRequest(ctx, requiredPreparationInstrumentInput)
			},
			Weight: 100,
		},
		"GetRequiredPreparationInstrument": {
			Name: "GetRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(ctx, c)
				if randomRequiredPreparationInstrument == nil {
					return nil, fmt.Errorf("retrieving random required preparation instrument: %w", ErrUnavailableYet)
				}

				return c.BuildGetRequiredPreparationInstrumentRequest(ctx, randomRequiredPreparationInstrument.ID)
			},
			Weight: 100,
		},
		"GetRequiredPreparationInstruments": {
			Name: "GetRequiredPreparationInstruments",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetRequiredPreparationInstrumentsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateRequiredPreparationInstrument": {
			Name: "UpdateRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(ctx, c); randomRequiredPreparationInstrument != nil {
					newRequiredPreparationInstrument := fakemodels.BuildFakeRequiredPreparationInstrumentCreationInput()
					randomRequiredPreparationInstrument.InstrumentID = newRequiredPreparationInstrument.InstrumentID
					randomRequiredPreparationInstrument.PreparationID = newRequiredPreparationInstrument.PreparationID
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

				randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(ctx, c)
				if randomRequiredPreparationInstrument == nil {
					return nil, fmt.Errorf("retrieving random required preparation instrument: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveRequiredPreparationInstrumentRequest(ctx, randomRequiredPreparationInstrument.ID)
			},
			Weight: 85,
		},
	}
}
