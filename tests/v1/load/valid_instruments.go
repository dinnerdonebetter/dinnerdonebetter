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

// fetchRandomValidInstrument retrieves a random valid instrument from the list of available valid instruments.
func fetchRandomValidInstrument(ctx context.Context, c *client.V1Client) *models.ValidInstrument {
	validInstrumentsRes, err := c.GetValidInstruments(ctx, nil)
	if err != nil || validInstrumentsRes == nil || len(validInstrumentsRes.ValidInstruments) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(validInstrumentsRes.ValidInstruments))
	return &validInstrumentsRes.ValidInstruments[randIndex]
}

func buildValidInstrumentActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateValidInstrument": {
			Name: "CreateValidInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				validInstrumentInput := fakemodels.BuildFakeValidInstrumentCreationInput()

				return c.BuildCreateValidInstrumentRequest(ctx, validInstrumentInput)
			},
			Weight: 100,
		},
		"GetValidInstrument": {
			Name: "GetValidInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidInstrument := fetchRandomValidInstrument(ctx, c)
				if randomValidInstrument == nil {
					return nil, fmt.Errorf("retrieving random valid instrument: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidInstrumentRequest(ctx, randomValidInstrument.ID)
			},
			Weight: 100,
		},
		"GetValidInstruments": {
			Name: "GetValidInstruments",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetValidInstrumentsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateValidInstrument": {
			Name: "UpdateValidInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomValidInstrument := fetchRandomValidInstrument(ctx, c); randomValidInstrument != nil {
					newValidInstrument := fakemodels.BuildFakeValidInstrumentCreationInput()
					randomValidInstrument.Name = newValidInstrument.Name
					randomValidInstrument.Variant = newValidInstrument.Variant
					randomValidInstrument.Description = newValidInstrument.Description
					randomValidInstrument.Icon = newValidInstrument.Icon
					return c.BuildUpdateValidInstrumentRequest(ctx, randomValidInstrument)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveValidInstrument": {
			Name: "ArchiveValidInstrument",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidInstrument := fetchRandomValidInstrument(ctx, c)
				if randomValidInstrument == nil {
					return nil, fmt.Errorf("retrieving random valid instrument: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveValidInstrumentRequest(ctx, randomValidInstrument.ID)
			},
			Weight: 85,
		},
	}
}
