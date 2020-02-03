package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipeStepEvent retrieves a random recipe step event from the list of available recipe step events
func fetchRandomRecipeStepEvent(c *client.V1Client) *models.RecipeStepEvent {
	recipeStepEventsRes, err := c.GetRecipeStepEvents(context.Background(), nil)
	if err != nil || recipeStepEventsRes == nil || len(recipeStepEventsRes.RecipeStepEvents) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepEventsRes.RecipeStepEvents))
	return &recipeStepEventsRes.RecipeStepEvents[randIndex]
}

func buildRecipeStepEventActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepEvent": {
			Name: "CreateRecipeStepEvent",
			Action: func() (*http.Request, error) {
				return c.BuildCreateRecipeStepEventRequest(context.Background(), randmodel.RandomRecipeStepEventCreationInput())
			},
			Weight: 100,
		},
		"GetRecipeStepEvent": {
			Name: "GetRecipeStepEvent",
			Action: func() (*http.Request, error) {
				if randomRecipeStepEvent := fetchRandomRecipeStepEvent(c); randomRecipeStepEvent != nil {
					return c.BuildGetRecipeStepEventRequest(context.Background(), randomRecipeStepEvent.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipeStepEvents": {
			Name: "GetRecipeStepEvents",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipeStepEventsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepEvent": {
			Name: "UpdateRecipeStepEvent",
			Action: func() (*http.Request, error) {
				if randomRecipeStepEvent := fetchRandomRecipeStepEvent(c); randomRecipeStepEvent != nil {
					randomRecipeStepEvent.EventType = randmodel.RandomRecipeStepEventCreationInput().EventType
					randomRecipeStepEvent.Done = randmodel.RandomRecipeStepEventCreationInput().Done
					randomRecipeStepEvent.RecipeIterationID = randmodel.RandomRecipeStepEventCreationInput().RecipeIterationID
					randomRecipeStepEvent.RecipeStepID = randmodel.RandomRecipeStepEventCreationInput().RecipeStepID
					return c.BuildUpdateRecipeStepEventRequest(context.Background(), randomRecipeStepEvent)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepEvent": {
			Name: "ArchiveRecipeStepEvent",
			Action: func() (*http.Request, error) {
				if randomRecipeStepEvent := fetchRandomRecipeStepEvent(c); randomRecipeStepEvent != nil {
					return c.BuildArchiveRecipeStepEventRequest(context.Background(), randomRecipeStepEvent.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
