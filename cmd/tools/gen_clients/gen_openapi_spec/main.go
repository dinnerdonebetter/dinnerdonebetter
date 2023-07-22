package main

import (
	"net/http"
	"os"
	"reflect"

	codegen "github.com/dinnerdonebetter/backend/cmd/tools/gen_clients"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"gopkg.in/yaml.v2"
)

func mustnt(err error) {
	if err != nil {
		panic(err)
	}
}

func defaultSchemaCustomizer(name string, _ reflect.Type, _ reflect.StructTag, _ *openapi3.Schema) error {
	if name == "_" {
		return &openapi3gen.ExcludeSchemaSentinel{}
	}

	return nil
}

func main() {
	spec := openapi3.T{
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL:         "https://api.dinnerdonebetter.dev",
				Description: "Development server",
			},
		},
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

	for _, typ := range codegen.TypesWeCareAbout {
		if typ.Name() == "QueryFilter" {
			continue
		}

		schemaRef, err := openapi3gen.NewSchemaRefForValue(
			typ.Type,
			spec.Components.Schemas,
			openapi3gen.UseAllExportedFields(),
			openapi3gen.SchemaCustomizer(defaultSchemaCustomizer),
		)
		if err != nil {
			panic(err)
		}

		schemaRef.Value.Description = typ.Description

		spec.Components.Schemas[typ.Name()] = schemaRef
	}

	routeSpecs := []routeSpec{
		{
			path:               "/api/v1/valid_ingredients",
			method:             http.MethodGet,
			description:        "Fetches valid ingredients",
			operationID:        "getValidIngredients",
			returnTypeName:     "ValidIngredient",
			returnCode:         http.StatusOK,
			returnsContent:     true,
			returnsArray:       true,
			acceptsQueryFilter: true,
			tags: []string{
				"Valid Ingredients",
			},
		},
		{
			path:           "/api/v1/valid_ingredients/{id}",
			method:         http.MethodGet,
			description:    "Fetches a valid ingredient",
			operationID:    "getValidIngredient",
			returnTypeName: "ValidIngredient",
			routeParams: []routeParam{
				{
					name:        "id",
					description: "the valid ingredient's id",
					typ:         "string",
				},
			},
			returnCode:     http.StatusOK,
			returnsContent: true,
			tags: []string{
				"Valid Ingredients",
			},
		},
		{
			path:           "/api/v1/valid_ingredients/random",
			method:         http.MethodGet,
			description:    "Fetches a random valid ingredient",
			operationID:    "getRandomValidIngredient",
			returnTypeName: "ValidIngredient",
			returnCode:     http.StatusOK,
			returnsContent: true,
			tags: []string{
				"Valid Ingredients",
			},
		},
		{
			path:           "/api/v1/valid_ingredients",
			method:         http.MethodPost,
			description:    "Creates a valid ingredient",
			operationID:    "createValidIngredient",
			returnTypeName: "ValidIngredient",
			returnCode:     http.StatusCreated,
			returnsContent: true,
			tags: []string{
				"Valid Ingredients",
			},
		},
	}

	for _, rs := range routeSpecs {
		addRoute(spec, rs)
	}

	//manuallyAddRoutesToSpec(spec)

	marshalledSpec, err := yaml.Marshal(spec)
	mustnt(err)

	if err = os.Remove("./cmd/tools/gen_clients/gen_openapi_spec/schema.yaml"); err != nil {
		panic(err)
	}

	if err = os.WriteFile("./cmd/tools/gen_clients/gen_openapi_spec/schema.yaml", marshalledSpec, 0o644); err != nil {
		panic(err)
	}
}

func manuallyAddRoutesToSpec(spec openapi3.T) {
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
}
