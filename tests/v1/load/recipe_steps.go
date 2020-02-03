package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipeStep retrieves a random recipe step from the list of available recipe steps
func fetchRandomRecipeStep(c *client.V1Client) *models.RecipeStep {
	recipeStepsRes, err := c.GetRecipeSteps(context.Background(), nil)
	if err != nil || recipeStepsRes == nil || len(recipeStepsRes.RecipeSteps) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepsRes.RecipeSteps))
	return &recipeStepsRes.RecipeSteps[randIndex]
}

func buildRecipeStepActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStep": {
			Name: "CreateRecipeStep",
			Action: func() (*http.Request, error) {
				return c.BuildCreateRecipeStepRequest(context.Background(), randmodel.RandomRecipeStepCreationInput())
			},
			Weight: 100,
		},
		"GetRecipeStep": {
			Name: "GetRecipeStep",
			Action: func() (*http.Request, error) {
				if randomRecipeStep := fetchRandomRecipeStep(c); randomRecipeStep != nil {
					return c.BuildGetRecipeStepRequest(context.Background(), randomRecipeStep.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipeSteps": {
			Name: "GetRecipeSteps",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipeStepsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStep": {
			Name: "UpdateRecipeStep",
			Action: func() (*http.Request, error) {
				if randomRecipeStep := fetchRandomRecipeStep(c); randomRecipeStep != nil {
					randomRecipeStep.Index = randmodel.RandomRecipeStepCreationInput().Index
					randomRecipeStep.PreparationID = randmodel.RandomRecipeStepCreationInput().PreparationID
					randomRecipeStep.PrerequisiteStep = randmodel.RandomRecipeStepCreationInput().PrerequisiteStep
					randomRecipeStep.MinEstimatedTimeInSeconds = randmodel.RandomRecipeStepCreationInput().MinEstimatedTimeInSeconds
					randomRecipeStep.MaxEstimatedTimeInSeconds = randmodel.RandomRecipeStepCreationInput().MaxEstimatedTimeInSeconds
					randomRecipeStep.TemperatureInCelsius = randmodel.RandomRecipeStepCreationInput().TemperatureInCelsius
					randomRecipeStep.Notes = randmodel.RandomRecipeStepCreationInput().Notes
					randomRecipeStep.RecipeID = randmodel.RandomRecipeStepCreationInput().RecipeID
					return c.BuildUpdateRecipeStepRequest(context.Background(), randomRecipeStep)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStep": {
			Name: "ArchiveRecipeStep",
			Action: func() (*http.Request, error) {
				if randomRecipeStep := fetchRandomRecipeStep(c); randomRecipeStep != nil {
					return c.BuildArchiveRecipeStepRequest(context.Background(), randomRecipeStep.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
