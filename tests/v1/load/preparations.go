package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomPreparation retrieves a random preparation from the list of available preparations
func fetchRandomPreparation(c *client.V1Client) *models.Preparation {
	preparationsRes, err := c.GetPreparations(context.Background(), nil)
	if err != nil || preparationsRes == nil || len(preparationsRes.Preparations) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(preparationsRes.Preparations))
	return &preparationsRes.Preparations[randIndex]
}

func buildPreparationActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreatePreparation": {
			Name: "CreatePreparation",
			Action: func() (*http.Request, error) {
				return c.BuildCreatePreparationRequest(context.Background(), randmodel.RandomPreparationCreationInput())
			},
			Weight: 100,
		},
		"GetPreparation": {
			Name: "GetPreparation",
			Action: func() (*http.Request, error) {
				if randomPreparation := fetchRandomPreparation(c); randomPreparation != nil {
					return c.BuildGetPreparationRequest(context.Background(), randomPreparation.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetPreparations": {
			Name: "GetPreparations",
			Action: func() (*http.Request, error) {
				return c.BuildGetPreparationsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdatePreparation": {
			Name: "UpdatePreparation",
			Action: func() (*http.Request, error) {
				if randomPreparation := fetchRandomPreparation(c); randomPreparation != nil {
					randomPreparation.Name = randmodel.RandomPreparationCreationInput().Name
					randomPreparation.Variant = randmodel.RandomPreparationCreationInput().Variant
					randomPreparation.Description = randmodel.RandomPreparationCreationInput().Description
					randomPreparation.AllergyWarning = randmodel.RandomPreparationCreationInput().AllergyWarning
					randomPreparation.Icon = randmodel.RandomPreparationCreationInput().Icon
					return c.BuildUpdatePreparationRequest(context.Background(), randomPreparation)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchivePreparation": {
			Name: "ArchivePreparation",
			Action: func() (*http.Request, error) {
				if randomPreparation := fetchRandomPreparation(c); randomPreparation != nil {
					return c.BuildArchivePreparationRequest(context.Background(), randomPreparation.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
