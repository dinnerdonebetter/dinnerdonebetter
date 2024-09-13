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
}

func (d *RouteDefinition) ToOperation() *openapi.Operation {
	op := &openapi.Operation{
		Tags:          nil,
		Summary:       nil,
		Description:   nil,
		ExternalDocs:  nil,
		ID:            nil,
		Parameters:    nil,
		RequestBody:   nil,
		Responses:     nil,
		Callbacks:     nil,
		Deprecated:    nil,
		Security:      nil,
		Servers:       nil,
		MapOfAnything: nil,
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
			Reference: &openapi.Reference{
				Ref: fmt.Sprintf("#/components/schemas/%s", d.InputType),
			},
		}
	}

	if d.ResponseType != "" {
		baseResponseSchema := map[string]any{
			"$ref": fmt.Sprintf("#/components/schemas/%s", d.ResponseType),
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

	return op
}
