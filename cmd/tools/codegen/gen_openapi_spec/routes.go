package main

import (
	"fmt"
	"log"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	"github.com/getkin/kin-openapi/openapi3"
)

type routeParam struct {
	name        string
	description string
	typ         string
	inQuery     bool
}

type routeSpec struct {
	path               string
	method             string
	description        string
	operationID        string
	returnTypeName     string
	routeParams        []*routeParam
	tags               []string
	returnCode         int
	returnsContent     bool
	returnsArray       bool
	acceptsQueryFilter bool
}

func (r *routeSpec) Parameters() openapi3.Parameters {
	output := []*openapi3.ParameterRef{}

	if r.acceptsQueryFilter {
		output = append(
			output,
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "limit",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "number",
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "page",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "number",
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "sortBy",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "string",
							Enum: []any{
								"asc",
								"desc",
							},
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "createdAfter",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:   "string",
							Format: "date-time",
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "createdBefore",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:   "string",
							Format: "date-time",
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "updatedAfter",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:   "string",
							Format: "date-time",
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "updatedBefore",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type:   "string",
							Format: "date-time",
						},
					},
				},
			},
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            "includeArchived",
					In:              "query",
					Description:     "",
					Required:        false,
					Deprecated:      false,
					AllowEmptyValue: true,
					Style:           "form",
					Explode:         pointers.Pointer(false),
					AllowReserved:   false,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "string",
							Enum: []any{
								"true",
							},
						},
					},
				},
			},
		)
	}

	for _, p := range r.routeParams {
		in := "path"
		style := "simple"
		if p.inQuery {
			in = "query"
			style = "form"
		}

		if p.typ == "" {
			log.Panicf("missing type for route %s parameter %s", r.operationID, p.name)
		}

		output = append(
			output,
			&openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					Name:            p.name,
					In:              in,
					Description:     p.description,
					Required:        true,
					AllowEmptyValue: false,
					Style:           style,
					AllowReserved:   !p.inQuery,
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: p.typ,
						},
					},
				},
			},
		)
	}

	return output
}

func addRoute(spec *openapi3.T, r *routeSpec) {
	contentReturn := openapi3.Content{}
	if r.returnsContent {
		if r.returnsArray {
			contentReturn = openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Value: &openapi3.Schema{
							Type: "array",
							Items: &openapi3.SchemaRef{
								Ref: fmt.Sprintf("#/components/schemas/%s", r.returnTypeName),
							},
						},
					},
				},
			}
		} else {
			contentReturn = openapi3.Content{
				"application/json": &openapi3.MediaType{
					Schema: &openapi3.SchemaRef{
						Ref: fmt.Sprintf("#/components/schemas/%s", r.returnTypeName),
					},
				},
			}
		}
	}

	spec.AddOperation(
		r.path,
		r.method,
		&openapi3.Operation{
			OperationID: r.operationID,
			Parameters:  r.Parameters(),
			Tags:        r.tags,
			Responses: openapi3.Responses{
				fmt.Sprintf("%d", r.returnCode): &openapi3.ResponseRef{
					Value: &openapi3.Response{
						Description: &r.description,
						Content:     contentReturn,
					},
				},
			},
		},
	)
}
