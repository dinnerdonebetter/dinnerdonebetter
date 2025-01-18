package main

import (
	"github.com/dinnerdonebetter/backend/internal/pointer"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

func baseSpec() *openapi.Spec {
	spec := &openapi.Spec{
		Openapi: "3.1.0",
		Servers: []openapi.Server{
			{
				URL:         "https://api.dinnerdonebetter.dev",
				Description: pointer.To("dev API server"),
			},
		},
		Components: &openapi.Components{
			SecuritySchemes: map[string]openapi.SecuritySchemeOrReference{
				"oauth2": {
					SecurityScheme: &openapi.SecurityScheme{
						Description: nil,
						Oauth2: &openapi.SecuritySchemeOauth2{Flows: openapi.OauthFlows{
							Implicit: &openapi.OauthFlowsDefsImplicit{
								AuthorizationURL: "/oauth2/authorize",
								Scopes: map[string]string{
									serviceAdmin:    "service-level administrator capabilities",
									householdAdmin:  "household-level administrator capabilities",
									householdMember: "household-level user capabilities",
								},
							},
						}},
						MapOfAnything: map[string]any{
							"type": "oauth2",
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
	})

	return spec
}
