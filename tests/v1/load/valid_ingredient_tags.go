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

// fetchRandomValidIngredientTag retrieves a random valid ingredient tag from the list of available valid ingredient tags.
func fetchRandomValidIngredientTag(ctx context.Context, c *client.V1Client) *models.ValidIngredientTag {
	validIngredientTagsRes, err := c.GetValidIngredientTags(ctx, nil)
	if err != nil || validIngredientTagsRes == nil || len(validIngredientTagsRes.ValidIngredientTags) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(validIngredientTagsRes.ValidIngredientTags))
	return &validIngredientTagsRes.ValidIngredientTags[randIndex]
}

func buildValidIngredientTagActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateValidIngredientTag": {
			Name: "CreateValidIngredientTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				validIngredientTagInput := fakemodels.BuildFakeValidIngredientTagCreationInput()

				return c.BuildCreateValidIngredientTagRequest(ctx, validIngredientTagInput)
			},
			Weight: 100,
		},
		"GetValidIngredientTag": {
			Name: "GetValidIngredientTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredientTag := fetchRandomValidIngredientTag(ctx, c)
				if randomValidIngredientTag == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient tag: %w", ErrUnavailableYet)
				}

				return c.BuildGetValidIngredientTagRequest(ctx, randomValidIngredientTag.ID)
			},
			Weight: 100,
		},
		"GetValidIngredientTags": {
			Name: "GetValidIngredientTags",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				return c.BuildGetValidIngredientTagsRequest(ctx, nil)
			},
			Weight: 100,
		},
		"UpdateValidIngredientTag": {
			Name: "UpdateValidIngredientTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				if randomValidIngredientTag := fetchRandomValidIngredientTag(ctx, c); randomValidIngredientTag != nil {
					newValidIngredientTag := fakemodels.BuildFakeValidIngredientTagCreationInput()
					randomValidIngredientTag.Name = newValidIngredientTag.Name
					return c.BuildUpdateValidIngredientTagRequest(ctx, randomValidIngredientTag)
				}

				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveValidIngredientTag": {
			Name: "ArchiveValidIngredientTag",
			Action: func() (*http.Request, error) {
				ctx := context.Background()

				randomValidIngredientTag := fetchRandomValidIngredientTag(ctx, c)
				if randomValidIngredientTag == nil {
					return nil, fmt.Errorf("retrieving random valid ingredient tag: %w", ErrUnavailableYet)
				}

				return c.BuildArchiveValidIngredientTagRequest(ctx, randomValidIngredientTag.ID)
			},
			Weight: 85,
		},
	}
}
