package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipeStepInstrument retrieves a random recipe step instrument from the list of available recipe step instruments
func fetchRandomRecipeStepInstrument(c *client.V1Client) *models.RecipeStepInstrument {
	recipeStepInstrumentsRes, err := c.GetRecipeStepInstruments(context.Background(), nil)
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
				return c.BuildCreateRecipeStepInstrumentRequest(context.Background(), randmodel.RandomRecipeStepInstrumentCreationInput())
			},
			Weight: 100,
		},
		"GetRecipeStepInstrument": {
			Name: "GetRecipeStepInstrument",
			Action: func() (*http.Request, error) {
				if randomRecipeStepInstrument := fetchRandomRecipeStepInstrument(c); randomRecipeStepInstrument != nil {
					return c.BuildGetRecipeStepInstrumentRequest(context.Background(), randomRecipeStepInstrument.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipeStepInstruments": {
			Name: "GetRecipeStepInstruments",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipeStepInstrumentsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepInstrument": {
			Name: "UpdateRecipeStepInstrument",
			Action: func() (*http.Request, error) {
				if randomRecipeStepInstrument := fetchRandomRecipeStepInstrument(c); randomRecipeStepInstrument != nil {
					randomRecipeStepInstrument.InstrumentID = randmodel.RandomRecipeStepInstrumentCreationInput().InstrumentID
					randomRecipeStepInstrument.RecipeStepID = randmodel.RandomRecipeStepInstrumentCreationInput().RecipeStepID
					randomRecipeStepInstrument.Notes = randmodel.RandomRecipeStepInstrumentCreationInput().Notes
					return c.BuildUpdateRecipeStepInstrumentRequest(context.Background(), randomRecipeStepInstrument)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepInstrument": {
			Name: "ArchiveRecipeStepInstrument",
			Action: func() (*http.Request, error) {
				if randomRecipeStepInstrument := fetchRandomRecipeStepInstrument(c); randomRecipeStepInstrument != nil {
					return c.BuildArchiveRecipeStepInstrumentRequest(context.Background(), randomRecipeStepInstrument.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
