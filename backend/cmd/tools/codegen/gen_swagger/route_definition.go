package main

import (
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

type RouteDefinition struct {
	Method           string
	Summary          string
	Path             string
	PathArguments    []string
	ListRoute        bool
	Tags             []string
	ResponseType     string
	MainResponseCode int
	RequestBody      string
	InputType        string
	OAuth2Scopes     []string
}

var routesWithoutAuth = map[string]struct{}{
	"/users/login":                 {},
	"/users/login/admin":           {},
	"/users/username/reminder":     {},
	"/users/password/reset":        {},
	"/users/password/reset/redeem": {},
	"/users/email_address/verify":  {},
	"/users/totp_secret/verify":    {},
	"oauth2/authorize":             {},
	"oauth2/token":                 {},
}

func (d *RouteDefinition) ToOperation() *openapi.Operation {
	op := &openapi.Operation{
		Tags:        []string{},
		Summary:     nil,
		Description: nil,
		Parameters:  []openapi.ParameterOrReference{},
	}

	if _, ok := routesWithoutAuth[d.Path]; !ok {
		op.Security = []map[string][]string{
			{"oauth2": d.OAuth2Scopes},
			{"cookieAuth": []string{}},
		}
	}

	if d.ListRoute {
		op.Parameters = append(op.Parameters, buildQueryFilterPathParams()...)
	}

	for _, arg := range d.PathArguments {
		op.Parameters = append(op.Parameters, openapi.ParameterOrReference{
			Parameter: &openapi.Parameter{
				Name:        arg,
				In:          "path",
				Description: nil,
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
				},
			},
		})
	}

	if d.InputType != "" {
		op.RequestBody = &openapi.RequestBodyOrReference{
			RequestBody: &openapi.RequestBody{
				Description: nil,
				Content: map[string]openapi.MediaType{
					"application/json": {
						Schema: map[string]any{
							"$ref": fmt.Sprintf("#/components/schemas/%s", d.InputType),
						},
					},
					"application/xml": {
						Schema: map[string]any{
							"$ref": fmt.Sprintf("#/components/schemas/%s", d.InputType),
						},
					},
				},
				Required:      pointer.To(true),
				MapOfAnything: nil,
			},
		}
	}

	if d.ResponseType != "" {
		baseResponseSchema := map[string]any{
			"$ref": "#/components/schemas/APIResponse",
		}

		secondResponseSchema := map[string]any{}
		if d.ListRoute {
			secondResponseSchema = map[string]any{
				"type": "object",
				"properties": map[string]any{
					"data": map[string]any{
						"type": "array",
						"items": map[string]any{
							"$ref": fmt.Sprintf("#/components/schemas/%s", d.ResponseType),
						},
					},
				},
			}
		} else {
			secondResponseSchema = map[string]any{
				"type": "object",
				"properties": map[string]any{
					"data": map[string]any{
						"type": "object",
						"$ref": fmt.Sprintf("#/components/schemas/%s", d.ResponseType),
					},
				},
			}
		}

		op.Responses = &openapi.Responses{
			Default: &openapi.ResponseOrReference{
				Response: &openapi.Response{
					Content: map[string]openapi.MediaType{
						"application/json": {
							Schema: map[string]any{
								"allOf": []map[string]any{
									baseResponseSchema,
									secondResponseSchema,
								},
							},
						},
						"application/xml": {
							Schema: map[string]any{
								"allOf": []map[string]any{
									baseResponseSchema,
									secondResponseSchema,
								},
							},
						},
					},
				},
			},
		}
	}

	for _, tag := range d.Tags {
		op.Tags = append(op.Tags, tag)
	}

	return op
}
