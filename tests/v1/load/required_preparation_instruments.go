package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRequiredPreparationInstrument retrieves a random required preparation instrument from the list of available required preparation instruments
func fetchRandomRequiredPreparationInstrument(c *client.V1Client) *models.RequiredPreparationInstrument {
	requiredPreparationInstrumentsRes, err := c.GetRequiredPreparationInstruments(context.Background(), nil)
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
				return c.BuildCreateRequiredPreparationInstrumentRequest(context.Background(), randmodel.RandomRequiredPreparationInstrumentCreationInput())
			},
			Weight: 100,
		},
		"GetRequiredPreparationInstrument": {
			Name: "GetRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				if randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(c); randomRequiredPreparationInstrument != nil {
					return c.BuildGetRequiredPreparationInstrumentRequest(context.Background(), randomRequiredPreparationInstrument.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRequiredPreparationInstruments": {
			Name: "GetRequiredPreparationInstruments",
			Action: func() (*http.Request, error) {
				return c.BuildGetRequiredPreparationInstrumentsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRequiredPreparationInstrument": {
			Name: "UpdateRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				if randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(c); randomRequiredPreparationInstrument != nil {
					randomRequiredPreparationInstrument.InstrumentID = randmodel.RandomRequiredPreparationInstrumentCreationInput().InstrumentID
					randomRequiredPreparationInstrument.PreparationID = randmodel.RandomRequiredPreparationInstrumentCreationInput().PreparationID
					randomRequiredPreparationInstrument.Notes = randmodel.RandomRequiredPreparationInstrumentCreationInput().Notes
					return c.BuildUpdateRequiredPreparationInstrumentRequest(context.Background(), randomRequiredPreparationInstrument)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRequiredPreparationInstrument": {
			Name: "ArchiveRequiredPreparationInstrument",
			Action: func() (*http.Request, error) {
				if randomRequiredPreparationInstrument := fetchRandomRequiredPreparationInstrument(c); randomRequiredPreparationInstrument != nil {
					return c.BuildArchiveRequiredPreparationInstrumentRequest(context.Background(), randomRequiredPreparationInstrument.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
