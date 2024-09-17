package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

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

	schemas, err := parseTypes("pkg/types")
	if err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(schemas, func(a, b *openapiSchema) int {
		if a.name < b.name {
			return -1
		} else if a.name == b.name {
			return 0
		} else {
			return 1
		}
	})

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

	allTags := map[string]struct{}{}

	routeDefinitions := []RouteDefinition{}
	for _, route := range srv.Router().Routes() {
		if strings.Contains(route.Path, "_meta_") {
			continue
		}

		if route.Path == "/auth/{auth_provider}" || route.Path == "/auth/{auth_provider}/callback" {
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
			OAuth2Scopes:  routeInfo.OAuth2Scopes,
		}

		if routeInfo.ResponseType != nil {
			routeDef.ResponseType = getTypeName(routeInfo.ResponseType)
		}

		if routeInfo.InputType != nil {
			routeDef.InputType = getTypeName(routeInfo.InputType)
		}

		pathParts := strings.Split(route.Path, "/")
		for i, part := range pathParts {
			if strings.TrimSpace(part) != "" && !strings.HasPrefix(part, "{") && part != "api" && part != "v1" {
				if i != len(pathParts)-1 {
					if rep, ok := tagReplacements[part]; ok {
						if _, ok2 := tagDescriptions[rep]; !ok2 {
							continue
						}
						routeDef.Tags = append(routeDef.Tags, rep)
						allTags[rep] = struct{}{}
					} else {
						if _, ok2 := tagDescriptions[part]; !ok2 {
							continue
						}
						routeDef.Tags = append(routeDef.Tags, part)
						allTags[part] = struct{}{}
					}
				}
			}
		}

		routeDefinitions = append(routeDefinitions, routeDef)
	}

	spec := baseSpec()

	tags := []openapi.Tag{}
	for tag, _ := range allTags {
		rawDescription := tagDescriptions[tag]

		var description *string
		if rawDescription != "" {
			description = &rawDescription
		}

		tags = append(tags, openapi.Tag{
			Name:        tag,
			Description: description,
		})
	}

	slices.SortFunc(tags, func(a, b openapi.Tag) int {
		if a.Name < b.Name {
			return -1
		} else if a.Name == b.Name {
			return 0
		} else {
			return 1
		}
	})
	spec.Tags = tags

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

	// TODO: type schema generation goes here

	convertedMap := map[string]map[string]any{}

	for _, schema := range schemas {
		if schema.name == "" {
			continue
		}

		tcm := map[string]any{
			"type": schema.Type,
		}

		propertiesMap := map[string]any{}
		for k, v := range schema.Properties {
			propertiesMap[k] = v
		}
		tcm["properties"] = propertiesMap

		convertedMap[schema.name] = tcm
	}
	spec.Components.Schemas = convertedMap

	output, err := spec.MarshalYAML()
	if err != nil {
		log.Fatal(err)
	}

	if err = os.Remove("../openapi_spec.yaml"); err != nil {
		log.Fatal(err)
	}

	if err = os.WriteFile("../openapi_spec.yaml", output, 0o600); err != nil {
		log.Fatal(err)
	}
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
					"enum": []string{"true", "false"},
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
