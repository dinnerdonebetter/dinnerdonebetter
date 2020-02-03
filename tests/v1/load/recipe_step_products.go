package main

import (
	"context"
	"math/rand"
	"net/http"

	client "gitlab.com/prixfixe/prixfixe/client/v1/http"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
	randmodel "gitlab.com/prixfixe/prixfixe/tests/v1/testutil/rand/model"
)

// fetchRandomRecipeStepProduct retrieves a random recipe step product from the list of available recipe step products
func fetchRandomRecipeStepProduct(c *client.V1Client) *models.RecipeStepProduct {
	recipeStepProductsRes, err := c.GetRecipeStepProducts(context.Background(), nil)
	if err != nil || recipeStepProductsRes == nil || len(recipeStepProductsRes.RecipeStepProducts) == 0 {
		return nil
	}

	randIndex := rand.Intn(len(recipeStepProductsRes.RecipeStepProducts))
	return &recipeStepProductsRes.RecipeStepProducts[randIndex]
}

func buildRecipeStepProductActions(c *client.V1Client) map[string]*Action {
	return map[string]*Action{
		"CreateRecipeStepProduct": {
			Name: "CreateRecipeStepProduct",
			Action: func() (*http.Request, error) {
				return c.BuildCreateRecipeStepProductRequest(context.Background(), randmodel.RandomRecipeStepProductCreationInput())
			},
			Weight: 100,
		},
		"GetRecipeStepProduct": {
			Name: "GetRecipeStepProduct",
			Action: func() (*http.Request, error) {
				if randomRecipeStepProduct := fetchRandomRecipeStepProduct(c); randomRecipeStepProduct != nil {
					return c.BuildGetRecipeStepProductRequest(context.Background(), randomRecipeStepProduct.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"GetRecipeStepProducts": {
			Name: "GetRecipeStepProducts",
			Action: func() (*http.Request, error) {
				return c.BuildGetRecipeStepProductsRequest(context.Background(), nil)
			},
			Weight: 100,
		},
		"UpdateRecipeStepProduct": {
			Name: "UpdateRecipeStepProduct",
			Action: func() (*http.Request, error) {
				if randomRecipeStepProduct := fetchRandomRecipeStepProduct(c); randomRecipeStepProduct != nil {
					randomRecipeStepProduct.Name = randmodel.RandomRecipeStepProductCreationInput().Name
					randomRecipeStepProduct.RecipeStepID = randmodel.RandomRecipeStepProductCreationInput().RecipeStepID
					return c.BuildUpdateRecipeStepProductRequest(context.Background(), randomRecipeStepProduct)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 100,
		},
		"ArchiveRecipeStepProduct": {
			Name: "ArchiveRecipeStepProduct",
			Action: func() (*http.Request, error) {
				if randomRecipeStepProduct := fetchRandomRecipeStepProduct(c); randomRecipeStepProduct != nil {
					return c.BuildArchiveRecipeStepProductRequest(context.Background(), randomRecipeStepProduct.ID)
				}
				return nil, ErrUnavailableYet
			},
			Weight: 85,
		},
	}
}
