package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomIterationMedia retrieves a random iteration media from the list of available iteration medias
func fetchRandomIterationMedia(c *client.V1Client) *models.IterationMedia {
	iterationMediasRes, err := c.GetIterationMedias(context.Background(), nil)
	if err != nil || iterationMediasRes == nil || len(iterationMediasRes.IterationMedias) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(iterationMediasRes.IterationMedias))
	return &iterationMediasRes.IterationMedias[randIndex]
}

func buildIterationMediaActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateIterationMedia": {
			Name: "CreateIterationMedia",
			Action: func() (*http.Request, error) {
				return c.BuildCreateIterationMediaRequest(context.Background(), randmodel.RandomIterationMediaCreationInput())
			},
			Weight: 100,
		},
		"GetIterationMedia": {
			Name: "GetIterationMedia",
			Action: func() (*http.Request, error) {
				if randomIterationMedia := fetchRandomIterationMedia(c); randomIterationMedia != nil {
					return c.BuildGetIterationMediaRequest(context.Background(), randomIterationMedia.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetIterationMedias": {
			Name: "GetIterationMedias",
			Action: func() (*http.Request, error) {
				return c.BuildGetIterationMediasRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateIterationMedia": {
			Name: "UpdateIterationMedia",
			Action: func() (*http.Request, error) {
				if randomIterationMedia := fetchRandomIterationMedia(c); randomIterationMedia != nil {
					randomIterationMedia.Path = randmodel.RandomIterationMediaCreationInput().Path
					randomIterationMedia.Mimetype = randmodel.RandomIterationMediaCreationInput().Mimetype
					randomIterationMedia.RecipeIterationID = randmodel.RandomIterationMediaCreationInput().RecipeIterationID
					randomIterationMedia.RecipeStepID = randmodel.RandomIterationMediaCreationInput().RecipeStepID
					return c.BuildUpdateIterationMediaRequest(context.Background(), randomIterationMedia)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveIterationMedia": {
			Name: "ArchiveIterationMedia",
			Action: func() (*http.Request, error) {
				if randomIterationMedia := fetchRandomIterationMedia(c); randomIterationMedia != nil {
					return c.BuildArchiveIterationMediaRequest(context.Background(), randomIterationMedia.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
