package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

type RouteDefinition struct {
	RequestBody        string
	InputType          string
	Method             string
	ID                 string
	Path               string
	ResponseType       string
	Summary            string
	Description        string
	PathArguments      []string
	Tags               []string
	OAuth2Scopes       []string
	MainResponseCode   int
	Authless           bool
	ReturnsList        bool
	QueryFilteredRoute bool
	SearchRoute        bool
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
	operationParts := []string{d.Method}
	for part := range strings.SplitSeq(d.Path, "/") {
		if part == "api" || part == "v1" || strings.TrimSpace(part) == "" {
			continue
		}
		operationParts = append(operationParts, strings.TrimSuffix(strings.TrimPrefix(part, "{"), "}"))
	}

	var description string
	if d.Description != "" {
		description = d.Description
	} else {
		description = "Operation for "
		switch d.Method {
		case http.MethodGet:
			description += "fetching"
		case http.MethodPut, http.MethodPatch:
			description += "updating"
		case http.MethodPost:
			description += "creating"
		case http.MethodDelete:
			description += "archiving"
		}

		if d.ResponseType != "" {
			description += fmt.Sprintf(" %s", d.ResponseType)
		}
	}

	opID := strings.Join(operationParts, "_")
	if d.ID != "" {
		opID = d.ID
	}

	op := &openapi.Operation{
		ID:          pointer.To(opID),
		Tags:        []string{},
		Summary:     nil,
		Description: pointer.To(description),
		Parameters:  []openapi.ParameterOrReference{},
	}

	if _, ok := routesWithoutAuth[d.Path]; !ok && !d.Authless {
		if len(d.OAuth2Scopes) > 0 {
			op.Security = append(op.Security, map[string][]string{"oauth2": d.OAuth2Scopes})
		}
	}

	if strings.HasSuffix(d.Path, "search") || d.SearchRoute {
		d.QueryFilteredRoute = true
		op.Parameters = append(op.Parameters, openapi.ParameterOrReference{
			Parameter: &openapi.Parameter{
				Name:        "q",
				In:          "query",
				Description: pointer.To("the search query parameter"),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
				},
			},
		})
	}

	if d.QueryFilteredRoute {
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
					jsonContentType: {
						Schema: map[string]any{
							refKey: fmt.Sprintf("#/components/schemas/%s", d.InputType),
						},
					},
					"application/xml": {
						Schema: map[string]any{
							refKey: fmt.Sprintf("#/components/schemas/%s", d.InputType),
						},
					},
				},
				Required:      pointer.To(true),
				MapOfAnything: nil,
			},
		}
	}

	if d.ResponseType != "" {
		_, responseTypeIsNative := nativeTypesMap[d.ResponseType]
		baseResponseSchema := map[string]any{
			refKey: "#/components/schemas/APIResponse",
		}

		var secondResponseSchema map[string]any

		if !responseTypeIsNative {
			if d.QueryFilteredRoute || d.ReturnsList {
				secondResponseSchema = map[string]any{
					"type": "object",
					propertiesKey: map[string]any{
						"data": map[string]any{
							"type": "array",
							"items": map[string]any{
								refKey: fmt.Sprintf("#/components/schemas/%s", d.ResponseType),
							},
						},
					},
				}
			} else {
				secondResponseSchema = map[string]any{
					"type": "object",
					propertiesKey: map[string]any{
						"data": map[string]any{
							"type": "object",
							refKey: fmt.Sprintf("#/components/schemas/%s", d.ResponseType),
						},
					},
				}
			}
		} else {
			if d.QueryFilteredRoute {
				secondResponseSchema = map[string]any{
					"type": "object",
					propertiesKey: map[string]any{
						"data": map[string]any{
							"type": "array",
							"items": map[string]any{
								"type": d.ResponseType,
							},
						},
					},
				}
			} else {
				secondResponseSchema = map[string]any{
					"type": "object",
					propertiesKey: map[string]any{
						"data": map[string]any{
							"type": d.ResponseType,
						},
					},
				}
			}
		}

		var statusCode int
		switch {
		case d.Method == http.MethodPost:
			statusCode = http.StatusCreated
		case d.Method == http.MethodDelete:
			statusCode = http.StatusAccepted
		default:
			statusCode = http.StatusOK
		}

		op.Responses = &openapi.Responses{
			MapOfResponseOrReferenceValues: map[string]openapi.ResponseOrReference{
				fmt.Sprintf("%d", statusCode): {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							jsonContentType: {
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
				"400": {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							jsonContentType: {
								Schema: map[string]any{
									refKey: "#/components/schemas/APIResponseWithError",
								},
							},
							"application/xml": {
								Schema: map[string]any{
									refKey: "#/components/schemas/APIResponseWithError",
								},
							},
						},
					},
				},
				"401": {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							jsonContentType: {
								Schema: map[string]any{
									refKey: "#/components/schemas/APIResponseWithError",
								},
							},
							"application/xml": {
								Schema: map[string]any{
									refKey: "#/components/schemas/APIResponseWithError",
								},
							},
						},
					},
				},
				"500": {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							jsonContentType: {
								Schema: map[string]any{
									refKey: "#/components/schemas/APIResponseWithError",
								},
							},
							"application/xml": {
								Schema: map[string]any{
									refKey: "#/components/schemas/APIResponseWithError",
								},
							},
						},
					},
				},
			},
		}
	}

	if d.Path == metaReadyPath || d.Path == metaLivePath {
		var desc string
		switch d.Path {
		case metaReadyPath:
			desc = "The service is ready to handle requests"
		case metaLivePath:
			desc = "The service is live."
		}

		op.Responses = &openapi.Responses{
			MapOfResponseOrReferenceValues: map[string]openapi.ResponseOrReference{
				fmt.Sprintf("%d", http.StatusOK): {
					Response: &openapi.Response{
						Description: desc,
					},
				},
			},
		}
	}

	if d.Path == "/api/v1/recipes/{recipeID}/dag" {
		op.Responses = &openapi.Responses{
			MapOfResponseOrReferenceValues: map[string]openapi.ResponseOrReference{
				fmt.Sprintf("%d", http.StatusOK): {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							"image/png": {
								Schema: map[string]any{
									"format": "binary",
								},
							},
						},
					},
				},
			},
		}
	}

	if d.Path == metaReadyPath || d.Path == metaLivePath {
		d.Tags = []string{"meta"}
	}

	op.Tags = append(op.Tags, d.Tags...)

	return op
}

func buildQueryFilterPathParams() []openapi.ParameterOrReference {
	return []openapi.ParameterOrReference{
		{
			Parameter: &openapi.Parameter{
				Name:        "limit",
				In:          "query",
				Description: pointer.To(fmt.Sprintf("How many results should appear in output, max is %d.", filtering.MaxQueryFilterLimit)),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "integer",
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "page",
				In:          "query",
				Description: pointer.To("What page of results should appear in output."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "integer",
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "createdBefore",
				In:          "query",
				Description: pointer.To("The latest CreatedAt date that should appear in output."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "createdAfter",
				In:          "query",
				Description: pointer.To("The earliest CreatedAt date that should appear in output."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "updatedBefore",
				In:          "query",
				Description: pointer.To("The latest UpdatedAt date that should appear in output."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "updatedAfter",
				In:          "query",
				Description: pointer.To("The earliest UpdatedAt date that should appear in output."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "includeArchived",
				In:          "query",
				Description: pointer.To("Whether or not to include archived results in output, limited to service admins."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
					"enum": []string{"true", "false"},
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "sortBy",
				In:          "query",
				Description: pointer.To("The direction in which results should be sorted."),
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
					"enum": []string{"asc", "desc"},
				},
			},
		},
	}
}
