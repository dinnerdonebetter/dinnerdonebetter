package main

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

func baseSpec() *openapi.Spec {
	spec := &openapi.Spec{}

	spec = spec.WithInfo(openapi.Info{
		Title:          "Dinner Done Better API",
		Summary:        nil,
		Description:    pointer.To(`This is the spec for the Dinner Done Better API.`),
		TermsOfService: nil,
		Contact: &openapi.Contact{
			Name:  nil,
			URL:   nil,
			Email: pointer.To("support@dinnerdonebetter.dev"),
		},
		License: nil,
		Version: "1.0.0",
		MapOfAnything: map[string]interface{}{
			"servers": []map[string]string{
				{
					"url": "https://api.dinnerdonebetter.dev",
				},
			},
			"tags": []map[string]string{
				{
					"name":        "recipes",
					"description": "Recipe-oriented routes",
				},
			},
		},
	})

	return spec
}
