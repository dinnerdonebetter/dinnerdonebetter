package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomInstrument retrieves a random instrument from the list of available instruments
func fetchRandomInstrument(c *client.V1Client) *models.Instrument {
	instrumentsRes, err := c.GetInstruments(context.Background(), nil)
	if err != nil || instrumentsRes == nil || len(instrumentsRes.Instruments) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(instrumentsRes.Instruments))
	return &instrumentsRes.Instruments[randIndex]
}

func buildInstrumentActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateInstrument": {
			Name: "CreateInstrument",
			Action: func() (*http.Request, error) {
				return c.BuildCreateInstrumentRequest(context.Background(), randmodel.RandomInstrumentCreationInput())
			},
			Weight: 100,
		},
		"GetInstrument": {
			Name: "GetInstrument",
			Action: func() (*http.Request, error) {
				if randomInstrument := fetchRandomInstrument(c); randomInstrument != nil {
					return c.BuildGetInstrumentRequest(context.Background(), randomInstrument.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetInstruments": {
			Name: "GetInstruments",
			Action: func() (*http.Request, error) {
				return c.BuildGetInstrumentsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateInstrument": {
			Name: "UpdateInstrument",
			Action: func() (*http.Request, error) {
				if randomInstrument := fetchRandomInstrument(c); randomInstrument != nil {
					randomInstrument.Name = randmodel.RandomInstrumentCreationInput().Name
					randomInstrument.Variant = randmodel.RandomInstrumentCreationInput().Variant
					randomInstrument.Description = randmodel.RandomInstrumentCreationInput().Description
					randomInstrument.Icon = randmodel.RandomInstrumentCreationInput().Icon
					return c.BuildUpdateInstrumentRequest(context.Background(), randomInstrument)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveInstrument": {
			Name: "ArchiveInstrument",
			Action: func() (*http.Request, error) {
				if randomInstrument := fetchRandomInstrument(c); randomInstrument != nil {
					return c.BuildArchiveInstrumentRequest(context.Background(), randomInstrument.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
