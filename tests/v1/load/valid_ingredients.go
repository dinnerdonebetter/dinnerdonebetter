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

// fetchRandomValidIngredient retrieves a random valid ingredient from the list of available valid ingredients.
func fetchRandomValidIngredient(ctx context.Context, c *client.V1Client) *models.ValidIngredient {
	validIngredientsRes, err := c.GetValidIngredients(ctx, nil)
	if err != nil || validIngredientsRes == nil || len(validIngredientsRes.ValidIngredients) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(validIngredientsRes.ValidIngredients))
	return &validIngredientsRes.ValidIngredients[randIndex]
}

func buildValidIngredientActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateValidIngredient": {
			Name: "CreateValidIngredient",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				validIngredientInput := fakemodels.BuildFakeValidIngredientCreationInput()

				return c.BuildCreateValidIngredientRequest(ctx, validIngredientInput)
			},
			Weight: 100,
		},
		"GetValidIngredient": {
			Name: "GetValidIngredient",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidIngredientRequest(ctx, randomValidIngredient.ID)
			},
			Weight: 100,
		},
		"GetValidIngredients": {
			Name: "GetValidIngredients",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetValidIngredientsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateValidIngredient": {
			Name: "UpdateValidIngredient",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomValidIngredient := fetchRandomValidIngredient(ctx, c); randomValidIngredient != nil {
					newValidIngredient := fakemodels.BuildFakeValidIngredientCreationInput()
					randomValidIngredient.Name = newValidIngredient.Name
					randomValidIngredient.Variant = newValidIngredient.Variant
					randomValidIngredient.Description = newValidIngredient.Description
					randomValidIngredient.Warning = newValidIngredient.Warning
					randomValidIngredient.ContainsEgg = newValidIngredient.ContainsEgg
					randomValidIngredient.ContainsDairy = newValidIngredient.ContainsDairy
					randomValidIngredient.ContainsPeanut = newValidIngredient.ContainsPeanut
					randomValidIngredient.ContainsTreeNut = newValidIngredient.ContainsTreeNut
					randomValidIngredient.ContainsSoy = newValidIngredient.ContainsSoy
					randomValidIngredient.ContainsWheat = newValidIngredient.ContainsWheat
					randomValidIngredient.ContainsShellfish = newValidIngredient.ContainsShellfish
					randomValidIngredient.ContainsSesame = newValidIngredient.ContainsSesame
					randomValidIngredient.ContainsFish = newValidIngredient.ContainsFish
					randomValidIngredient.ContainsGluten = newValidIngredient.ContainsGluten
					randomValidIngredient.AnimalFlesh = newValidIngredient.AnimalFlesh
					randomValidIngredient.AnimalDerived = newValidIngredient.AnimalDerived
					randomValidIngredient.MeasurableByVolume = newValidIngredient.MeasurableByVolume
					randomValidIngredient.Icon = newValidIngredient.Icon
					return c.BuildUpdateValidIngredientRequest(ctx, randomValidIngredient)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveValidIngredient": {
			Name: "ArchiveValidIngredient",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveValidIngredientRequest(ctx, randomValidIngredient.ID)
			},
			Weight: 85,
		},
	}
}
