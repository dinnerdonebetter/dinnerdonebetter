package main

import (
	"encoding/json"
	"net/http"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
)

func mustnt(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	//ctx := context.Background()

	spec := openapi3.T{
		OpenAPI: "3.0.0",
		Info: &openapi3.Info{
			Title:          "Dinner Done Better API",
			Description:    "Dinner Done Better API",
			TermsOfService: "",
			Contact: &openapi3.Contact{
				Name:  "",
				URL:   "",
				Email: "",
			},
			License: &openapi3.License{
				Name: "nunya",
				URL:  "",
			},
			Version: "0.0.1",
		},
		Paths: openapi3.Paths{},
		Components: &openapi3.Components{
			Extensions:      map[string]any{},
			Schemas:         map[string]*openapi3.SchemaRef{},
			Parameters:      map[string]*openapi3.ParameterRef{},
			Headers:         map[string]*openapi3.HeaderRef{},
			RequestBodies:   map[string]*openapi3.RequestBodyRef{},
			Responses:       map[string]*openapi3.ResponseRef{},
			SecuritySchemes: map[string]*openapi3.SecuritySchemeRef{},
			Examples:        map[string]*openapi3.ExampleRef{},
			Links:           map[string]*openapi3.LinkRef{},
			Callbacks:       map[string]*openapi3.CallbackRef{},
		},
	}

	validIngredientSchemaRef, err := openapi3gen.NewSchemaRefForValue(&types.ValidIngredient{}, spec.Components.Schemas, openapi3gen.UseAllExportedFields(), openapi3gen.ThrowErrorOnCycle())
	mustnt(err)
	spec.Components.Schemas["ValidIngredient"] = validIngredientSchemaRef

	validIngredientCreationRequestInputSchemaRef, err := openapi3gen.NewSchemaRefForValue(&types.ValidIngredientCreationRequestInput{}, spec.Components.Schemas, openapi3gen.UseAllExportedFields(), openapi3gen.ThrowErrorOnCycle())
	mustnt(err)
	spec.Components.Schemas["ValidIngredientCreationRequestInput"] = validIngredientCreationRequestInputSchemaRef

	spec.AddOperation("/api/v1/valid_ingredients", http.MethodPost, &openapi3.Operation{
		OperationID: "createRandomValidIngredient",
		Parameters:  openapi3.Parameters{},
		RequestBody: &openapi3.RequestBodyRef{
			Value: &openapi3.RequestBody{
				Required: true,
				Content: openapi3.Content{
					"application/json": &openapi3.MediaType{
						Schema: &openapi3.SchemaRef{
							Ref: "#/components/schemas/ValidIngredientCreationRequestInput",
						},
					},
				},
			},
		},
		Responses: openapi3.Responses{
			"201": &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: pointers.Pointer(" "),
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/ValidIngredient",
							},
						},
					},
				},
			},
		},
	})

	spec.AddOperation("/api/v1/valid_ingredients", http.MethodGet, &openapi3.Operation{
		OperationID: "getValidIngredients",
		Parameters: openapi3.Parameters{
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
		},
		Responses: openapi3.Responses{
			"200": &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: pointers.Pointer(" "),
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Value: &openapi3.Schema{
									Type: "array",
									Items: &openapi3.SchemaRef{
										Ref: "#/components/schemas/ValidIngredient",
									},
								},
							},
						},
					},
				},
			},
		},
	})

	spec.AddOperation("/api/v1/valid_ingredients/random", http.MethodGet, &openapi3.Operation{
		OperationID: "getRandomValidIngredient",
		Parameters:  openapi3.Parameters{},
		Responses: openapi3.Responses{
			"200": &openapi3.ResponseRef{
				Value: &openapi3.Response{
					Description: pointers.Pointer(" "),
					Content: openapi3.Content{
						"application/json": &openapi3.MediaType{
							Schema: &openapi3.SchemaRef{
								Ref: "#/components/schemas/ValidIngredient",
							},
						},
					},
				},
			},
		},
	})

	//mustnt(spec.Validate(ctx))

	marshalledSpec, err := json.MarshalIndent(spec, "", "  ")
	mustnt(err)

	println(string(marshalledSpec))
}
