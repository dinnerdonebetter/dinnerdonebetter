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

// fetchRandomIngredientTagMapping retrieves a random ingredient tag mapping from the list of available ingredient tag mappings.
func fetchRandomIngredientTagMapping(ctx context.Context, c *client.V1Client, validIngredientID uint64) *models.IngredientTagMapping {
	ingredientTagMappingsRes, err := c.GetIngredientTagMappings(ctx, validIngredientID, nil)
	if err != nil || ingredientTagMappingsRes == nil || len(ingredientTagMappingsRes.IngredientTagMappings) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(ingredientTagMappingsRes.IngredientTagMappings))
	return &ingredientTagMappingsRes.IngredientTagMappings[randIndex]
}

func buildIngredientTagMappingActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateIngredientTagMapping": {
			Name: "CreateIngredientTagMapping",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				// Create valid ingredient.
				exampleValidIngredient := fakemodels.BuildFakeValidIngredient()
				exampleValidIngredientInput := fakemodels.BuildFakeValidIngredientCreationInputFromValidIngredient(exampleValidIngredient)
				createdValidIngredient, err := c.CreateValidIngredient(ctx, exampleValidIngredientInput)
				if err != nil {
					return nil, err
				}

				ingredientTagMappingInput := fakemodels.BuildFakeIngredientTagMappingCreationInput()
				ingredientTagMappingInput.BelongsToValidIngredient = createdValidIngredient.ID

				return c.BuildCreateIngredientTagMappingRequest(ctx, ingredientTagMappingInput)
			},
			Weight: 100,
		},
		"GetIngredientTagMapping": {
			Name: "GetIngredientTagMapping",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				randomIngredientTagMapping := fetchRandomIngredientTagMapping(ctx, c, randomValidIngredient.ID)
				if randomIngredientTagMapping == nil {
					return nil, fmt.Errorf("retrieving random ingredient tag mapping: %w", ErrUnavailableYet)
				}

				return c.BuildGetIngredientTagMappingRequest(ctx, randomValidIngredient.ID, randomIngredientTagMapping.ID)
			},
			Weight: 100,
		},
		"GetIngredientTagMappings": {
			Name: "GetIngredientTagMappings",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				return c.BuildGetIngredientTagMappingsRequest(ctx, randomValidIngredient.ID, nil)
			},
			Weight: 100,
		},
		"UpdateIngredientTagMapping": {
			Name: "UpdateIngredientTagMapping",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				if randomIngredientTagMapping := fetchRandomIngredientTagMapping(ctx, c, randomValidIngredient.ID); randomIngredientTagMapping != nil {
					newIngredientTagMapping := fakemodels.BuildFakeIngredientTagMappingCreationInput()
					randomIngredientTagMapping.ValidIngredientTagID = newIngredientTagMapping.ValidIngredientTagID
					return c.BuildUpdateIngredientTagMappingRequest(ctx, randomIngredientTagMapping)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveIngredientTagMapping": {
			Name: "ArchiveIngredientTagMapping",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredient := fetchRandomValidIngredient(ctx, c)
				if randomValidIngredient == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient: %w", ErrUnavailableYet)
				}

				randomIngredientTagMapping := fetchRandomIngredientTagMapping(ctx, c, randomValidIngredient.ID)
				if randomIngredientTagMapping == nil {
					return nil, fmt.Errorf("retrieving random ingredient tag mapping: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveIngredientTagMappingRequest(ctx, randomValidIngredient.ID, randomIngredientTagMapping.ID)
			},
			Weight: 85,
		},
	}
}
