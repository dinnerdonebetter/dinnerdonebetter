package main

import (
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

func baseSpec() *openapi.Spec {
	spec := &openapi.Spec{
		Openapi: "3.1.0",
		Components: &openapi.Components{
			SecuritySchemes: map[string]openapi.SecuritySchemeOrReference{
				"cookieAuth": {
					SecurityScheme: &openapi.SecurityScheme{
						Description: nil,
						APIKey:      nil,
						HTTP:        nil,
						HTTPBearer:  nil,
						Oauth2: &openapi.SecuritySchemeOauth2{Flows: openapi.OauthFlows{
							Implicit: &openapi.OauthFlowsDefsImplicit{
								AuthorizationURL: "/oauth2/authorize",
								RefreshURL:       nil,
								Scopes:           nil,
								MapOfAnything:    nil,
							},
						}},
						Oidc:      nil,
						MutualTLS: nil,
						MapOfAnything: map[string]any{
							"type": "apiKey",
							"in":   "cookie",
							"name": "ddb_api_cookie",
						},
					},
				},
			},
			Schemas: map[string]map[string]any{},
		},
	}

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
