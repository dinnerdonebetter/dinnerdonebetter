package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

type RouteDefinition struct {
	Method           string
	Summary          string
	ID               string
	Path             string
	ResponseType     string
	RequestBody      string
	Description      string
	InputType        string
	PathArguments    []string
	Tags             []string
	OAuth2Scopes     []string
	MainResponseCode int
	ListRoute        bool
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
	for _, part := range strings.Split(d.Path, "/") {
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

	if _, ok := routesWithoutAuth[d.Path]; !ok {
		op.Security = []map[string][]string{
			{"cookieAuth": []string{}},
		}

		if len(d.OAuth2Scopes) > 0 {
			op.Security = append(op.Security, map[string][]string{"oauth2": d.OAuth2Scopes})
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

		var secondResponseSchema map[string]any
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
				"400": {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							"application/json": {
								Schema: map[string]any{
									"$ref": "#/components/schemas/APIResponseWithError",
								},
							},
							"application/xml": {
								Schema: map[string]any{
									"$ref": "#/components/schemas/APIResponseWithError",
								},
							},
						},
					},
				},
				"401": {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							"application/json": {
								Schema: map[string]any{
									"$ref": "#/components/schemas/APIResponseWithError",
								},
							},
							"application/xml": {
								Schema: map[string]any{
									"$ref": "#/components/schemas/APIResponseWithError",
								},
							},
						},
					},
				},
				"500": {
					Response: &openapi.Response{
						Content: map[string]openapi.MediaType{
							"application/json": {
								Schema: map[string]any{
									"$ref": "#/components/schemas/APIResponseWithError",
								},
							},
							"application/xml": {
								Schema: map[string]any{
									"$ref": "#/components/schemas/APIResponseWithError",
								},
							},
						},
					},
				},
			},
		}
	}

	op.Tags = append(op.Tags, d.Tags...)

	return op
}
