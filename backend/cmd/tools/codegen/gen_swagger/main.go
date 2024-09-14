package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"
	"github.com/dinnerdonebetter/backend/pkg/types"
	openapi "github.com/swaggest/openapi-go/openapi31"
	"log"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
)

var routeParamRegex = regexp.MustCompile("\\{[a-zA-Z]+\\}")

var tagReplacements = map[string]string{
	"steps":                 "recipe_steps",
	"prep_tasks":            "recipe_prep_tasks",
	"completion_conditions": "recipe_step_completion_conditions",
}

func getTypeName(input any) string {
	t := reflect.TypeOf(input)

	// Dereference pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name()
}

func main() {
	ctx := context.Background()

	rawCfg, err := os.ReadFile("environments/dev/config_files/service-config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg *config.InstanceConfig
	if err = json.Unmarshal(rawCfg, &cfg); err != nil {
		log.Fatal(err)
	}

	cfg.Neutralize()

	// build our server struct.
	srv, err := build.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	routeDefinitions := []RouteDefinition{}
	for _, route := range srv.Router().Routes() {
		if strings.Contains(route.Path, "_meta_") {
			continue
		}

		pathArgs := []string{}
		for _, pathArg := range routeParamRegex.FindAllString(route.Path, -1) {
			pathArgs = append(pathArgs, strings.TrimPrefix(strings.TrimSuffix(pathArg, "}"), "{"))
		}

		routeInfo, ok := routeInfoMap[fmt.Sprintf("%s %s", route.Method, route.Path)]
		if !ok {
			panic("unknown route")
		}

		routeDef := RouteDefinition{
			Method:        route.Method,
			Path:          route.Path,
			PathArguments: pathArgs,
			ListRoute:     routeInfo.ListRoute,
		}

		if routeInfo.ResponseType != nil {
			routeInfo.ResponseType = getTypeName(routeInfo.ResponseType)
		}

		if routeInfo.InputType != nil {
			routeDef.InputType = getTypeName(routeInfo.InputType)
		}

		for _, part := range strings.Split(route.Path, "/") {
			if strings.TrimSpace(part) != "" && !strings.HasPrefix(part, "{") && part != "api" && part != "v1" {
				if rep, ok := tagReplacements[part]; ok {
					routeDef.Tags = append(routeDef.Tags, rep)
				} else {
					routeDef.Tags = append(routeDef.Tags, part)
				}
			}
		}

		routeDefinitions = append(routeDefinitions, routeDef)
	}

	spec := baseSpec()

	paths := &openapi.Paths{
		MapOfPathItemValues: map[string]openapi.PathItem{},
	}

	for _, rd := range routeDefinitions {
		op := rd.ToOperation()

		if _, ok := paths.MapOfPathItemValues[rd.Path]; ok {
			// path already present
			item := paths.MapOfPathItemValues[rd.Path]

			switch rd.Method {
			case http.MethodGet:
				item.Get = op
			case http.MethodPut:
				item.Put = op
			case http.MethodPost:
				item.Post = op
			case http.MethodDelete:
				item.Delete = op
			}

			paths.MapOfPathItemValues[rd.Path] = item
		} else {
			// path is not yet present
			item := openapi.PathItem{}

			switch rd.Method {
			case http.MethodGet:
				item.Get = op
			case http.MethodPut:
				item.Put = op
			case http.MethodPost:
				item.Post = op
			case http.MethodDelete:
				item.Delete = op
			}

			paths.MapOfPathItemValues[rd.Path] = item
		}
	}

	spec.Paths = paths

	typeSchemas := map[string]map[string]any{
		getTypeName(&types.ResponseDetails{}): SchemaFromInstance(&types.ResponseDetails{}).RenderSchema(),
	}
	for _, v := range routeInfoMap {
		if v.InputType != nil {
			if _, ok := typeSchemas[getTypeName(v.InputType)]; ok {
				continue
			}

			typeSchemas[getTypeName(v.InputType)] = v.InputTypeSchema.RenderSchema()
		}
	}

	for _, v := range routeInfoMap {
		if v.ResponseType != nil {
			if _, ok := typeSchemas[getTypeName(v.ResponseType)]; ok {
				continue
			}

			typeSchemas[getTypeName(v.ResponseType)] = v.ResponseTypeSchema.RenderSchema()
		}
	}

	for name, schema := range typeSchemas {
		spec.Components.Schemas[name] = schema
	}

	output, err := spec.MarshalYAML()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(output))
}

func buildQueryFilterPathParams() []openapi.ParameterOrReference {
	return []openapi.ParameterOrReference{
		{
			Parameter: &openapi.Parameter{
				Name:        "page",
				In:          "query",
				Description: nil,
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
				Description: nil,
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
				Description: nil,
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
				Description: nil,
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
				Description: nil,
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
				Description: nil,
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
					"enum": []string{"1", "t", "T", "true", "TRUE", "True", "0", "f", "F", "false", "FALSE", "False"},
				},
			},
		},
		{
			Parameter: &openapi.Parameter{
				Name:        "sortBy",
				In:          "query",
				Description: nil,
				Required:    pointer.To(true),
				Schema: map[string]any{
					"type": "string",
					"enum": []string{"asc", "desc"},
				},
			},
		},
	}
}
