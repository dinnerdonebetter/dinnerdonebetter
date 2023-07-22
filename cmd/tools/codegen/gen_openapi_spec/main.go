package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"

	"github.com/dinnerdonebetter/backend/cmd/tools/codegen"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/objectstorage"
	"github.com/dinnerdonebetter/backend/internal/routing"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"

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

func fetchServiceRouter() (routing.Router, error) {
	var cfg *config.InstanceConfig
	configBytes, err := os.ReadFile("environments/local/config_files/service-config.json")
	if err != nil {
		return nil, fmt.Errorf("reading local config file: %w", err)
	}

	if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
		return nil, fmt.Errorf("decoding config file contents: %w", err)
	}

	cfg.Database.RunMigrations = false
	cfg.Services.Users.Uploads.Storage.Provider = objectstorage.MemoryProvider
	cfg.Services.Users.Uploads.Storage.FilesystemConfig = nil
	cfg.Services.Recipes.Uploads.Storage.Provider = objectstorage.MemoryProvider
	cfg.Services.Recipes.Uploads.Storage.FilesystemConfig = nil
	cfg.Services.RecipeSteps.Uploads.Storage.Provider = objectstorage.MemoryProvider
	cfg.Services.RecipeSteps.Uploads.Storage.FilesystemConfig = nil

	// build our server struct.
	srv, err := build.Build(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("building server: %w", err)
	}

	return srv.Router(), nil
}

func main() {
	if _, err := fetchServiceRouter(); err != nil {
		panic(err)
	}

	spec := openapi3.T{
		Servers: openapi3.Servers{
			&openapi3.Server{
				URL:         "https://api.dinnerdonebetter.dev",
				Description: "Development server",
			},
		},
		OpenAPI: "3.0.3",
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

	routeSpecs := []*routeSpec{
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
			routeParams: []*routeParam{
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
		addRoute(&spec, rs)
	}

	marshalledSpec, err := yaml.Marshal(spec)
	mustnt(err)

	const outputFilepath = "./openapi-schema.yaml"
	if err = os.Remove(outputFilepath); err != nil {
		panic(err)
	}

	if err = os.WriteFile(outputFilepath, marshalledSpec, 0o600); err != nil {
		panic(err)
	}
}
