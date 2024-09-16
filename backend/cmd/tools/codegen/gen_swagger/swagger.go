package main

import (
	"github.com/dinnerdonebetter/backend/internal/authorization"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"strings"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

func gatherScopes() map[string]string {
	allScopes := map[string]string{}
	for _, scope := range authorization.ServiceAdminPermissions {
		rawPerm := scope.ID()

		rawPermParts := strings.Split(rawPerm, ".")
		moddedPermParts := make([]string, len(rawPermParts))
		for i, part := range rawPermParts {
			moddedPermParts[i] = strings.ReplaceAll(part, "_", " ")
		}

		allScopes[scope.ID()] = strings.Join(moddedPermParts, " ")
	}

	for _, scope := range authorization.HouseholdAdminPermissions {
		rawPerm := scope.ID()

		rawPermParts := strings.Split(rawPerm, ".")
		moddedPermParts := make([]string, len(rawPermParts))
		for i, part := range rawPermParts {
			moddedPermParts[i] = strings.ReplaceAll(part, "_", " ")
		}

		allScopes[scope.ID()] = strings.Join(moddedPermParts, " ")
	}

	for _, scope := range authorization.HouseholdMemberPermissions {
		rawPerm := scope.ID()

		rawPermParts := strings.Split(rawPerm, ".")
		moddedPermParts := make([]string, len(rawPermParts))
		for i, part := range rawPermParts {
			moddedPermParts[i] = strings.ReplaceAll(part, "_", " ")
		}

		allScopes[scope.ID()] = strings.Join(moddedPermParts, " ")
	}

	return allScopes
}

func baseSpec() *openapi.Spec {
	allScopes := gatherScopes()

	spec := &openapi.Spec{
		Openapi: "3.1.0",
		Components: &openapi.Components{
			SecuritySchemes: map[string]openapi.SecuritySchemeOrReference{
				"cookieAuth": {
					SecurityScheme: &openapi.SecurityScheme{
						Description: nil,
						MapOfAnything: map[string]any{
							"type": "apiKey",
							"in":   "cookie",
							"name": "ddb_api_cookie",
						},
					},
				},
				"oauth2": {
					SecurityScheme: &openapi.SecurityScheme{
						Description: nil,
						Oauth2: &openapi.SecuritySchemeOauth2{Flows: openapi.OauthFlows{
							Implicit: &openapi.OauthFlowsDefsImplicit{
								AuthorizationURL: "/oauth2/authorize",
								Scopes:           map[string]string{},
							},
						}},
						MapOfAnything: map[string]any{
							"type": "oauth2",
							"flows": map[string]any{
								"authorizationCode": map[string]any{
									"authorizationUrl": "/oauth2/authorize",
									"scopes":           allScopes,
								},
							},
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
