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

	"github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/lib/uploads/objectstorage"

	openapi "github.com/swaggest/openapi-go/openapi31"
)

const (
	metaReadyPath   = "/_ops_/ready"
	metaLivePath    = "/_meta/live"
	jsonContentType = "application/json"
	refKey          = "$ref"
	propertiesKey   = "properties"
)

func getTypeName(input any) (string, bool) {
	t := reflect.TypeOf(input)

	// Check if the input is an array or slice
	isArray := false
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		isArray = true
		// Get the type of the element in the array/slice
		t = t.Elem()
	}

	// Dereference pointer types
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Name(), isArray
}

var skipRoutes = map[string]bool{
	metaLivePath:                     true,
	"/auth/{auth_provider}":          true,
	"/auth/{auth_provider}/callback": true,
	"/oauth2/authorize":              true,
	"/oauth2/token":                  true,
	"/api/v1/recipes/{recipeID}/steps/{recipeStepID}":        true,
	"/api/v1/recipes/{recipeID}/images":                      true,
	"/api/v1/recipes/{recipeID}/steps/{recipeStepID}/images": true,
	"/api/v1/households/{householdID}/invitations/":          true,
	"": true,
}

func main() {
	ctx := context.Background()

	rawCfg, err := os.ReadFile("deploy/environments/localdev/config_files/api_service_config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg *config.APIServiceConfig
	if err = json.Unmarshal(rawCfg, &cfg); err != nil {
		log.Fatal(err)
	}

	neutralizeConfig(cfg)

	// build our server struct.
	srv, err := api.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	allTags := map[string]struct{}{}

	routeDefinitions := []*RouteDefinition{}
	for _, route := range srv.Router().Routes() {
		if _, ok := skipRoutes[route.Path]; ok {
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

		routeDef := &RouteDefinition{
			ID:                 routeInfo.ID,
			Method:             route.Method,
			Path:               route.Path,
			SearchRoute:        routeInfo.SearchRoute,
			PathArguments:      pathArgs,
			QueryFilteredRoute: routeInfo.ListRoute,
			Description:        routeInfo.Description,
			Authless:           routeInfo.Authless,
			OAuth2Scopes:       routeInfo.OAuth2Scopes,
		}

		if routeInfo.ResponseType != nil {
			routeDef.ResponseType, routeDef.ReturnsList = getTypeName(routeInfo.ResponseType)
		}

		if routeInfo.InputType != nil {
			routeDef.InputType, _ = getTypeName(routeInfo.InputType)
		}

		if routeDef.Path == metaLivePath || routeDef.Path == metaReadyPath {
			routeDef.ResponseType = "empty"
		}

		pathParts := strings.Split(route.Path, "/")
		for i, part := range pathParts {
			if strings.TrimSpace(part) != "" && !strings.HasPrefix(part, "{") && part != "api" && part != "v1" {
				if i != len(pathParts)-1 {
					if rep, ok1 := tagReplacements[part]; ok1 {
						if _, ok2 := tagDescriptions[rep]; !ok2 {
							continue
						}
						routeDef.Tags = append(routeDef.Tags, strings.ReplaceAll(rep, "_", " "))
						allTags[rep] = struct{}{}
					} else {
						if _, ok2 := tagDescriptions[part]; !ok2 {
							continue
						}
						routeDef.Tags = append(routeDef.Tags, strings.ReplaceAll(part, "_", " "))
						allTags[part] = struct{}{}
					}
				}
			}
		}

		routeDefinitions = append(routeDefinitions, routeDef)
	}

	spec := baseSpec()

	tags := []openapi.Tag{
		{Name: "meta"},
	}
	for tag := range allTags {
		rawDescription := tagDescriptions[tag]

		var description *string
		if rawDescription != "" {
			description = &rawDescription
		}

		tags = append(tags, openapi.Tag{
			Name:        strings.ReplaceAll(tag, "_", " "),
			Description: description,
		})
	}

	slices.SortFunc(tags, func(a, b openapi.Tag) int {
		switch {
		case a.Name < b.Name:
			return -1
		case a.Name == b.Name:
			return 0
		default:
			return 1
		}
	})
	spec.Tags = tags

	paths := &openapi.Paths{
		MapOfPathItemValues: map[string]openapi.PathItem{},
	}

	for _, rd := range routeDefinitions {
		op := rd.ToOperation()
		path := strings.TrimSuffix(rd.Path, "/")

		var item openapi.PathItem
		if _, ok := paths.MapOfPathItemValues[path]; ok {
			item = paths.MapOfPathItemValues[path]
		} else {
			item = openapi.PathItem{}
		}

		switch rd.Method {
		case http.MethodGet:
			item.Get = op
		case http.MethodPut:
			item.Put = op
		case http.MethodPost:
			item.Post = op
		case http.MethodPatch:
			item.Patch = op
		case http.MethodDelete:
			item.Delete = op
		}

		paths.MapOfPathItemValues[path] = item
	}

	spec.Paths = paths

	schemas, err := parseTypes("pkg/types", "internal/lib/database/filtering")
	if err != nil {
		log.Fatal(err)
	}

	slices.SortFunc(schemas, func(a, b *openapiSchema) int {
		switch {
		case a.name < b.name:
			return -1
		case a.name == b.name:
			return 0
		default:
			return 1
		}
	})

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

		if len(propertiesMap) > 0 {
			tcm[propertiesKey] = propertiesMap
		}

		if len(schema.Enum) > 0 {
			tcm["enum"] = schema.Enum
		}

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

const (
	dataChangesTopicName = "dataChangesTopicName"
)

func neutralizeConfig(cfg *config.APIServiceConfig) {
	if err := os.Setenv("GOOGLE_CLOUD_PROJECT_ID", "something"); err != nil {
		panic(err)
	}

	cfg.Observability.Logging.Provider = "noop"

	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = "BLAHBLAHBLAHBLAHBLAHBLAHBLAHBLAH"
	cfg.Services.Auth.SSO.Google.ClientID = "blah blah blah blah"
	cfg.Services.Auth.SSO.Google.ClientSecret = "blah blah blah blah"
	cfg.Analytics.Provider = ""

	cfg.Services.Recipes.Uploads.Storage.GCP = nil
	cfg.Services.Recipes.Uploads.Storage.Provider = objectstorage.FilesystemProvider
	cfg.Services.Recipes.Uploads.Storage.FilesystemConfig = &objectstorage.FilesystemConfig{RootDirectory: "/tmp"}
	cfg.Services.DataPrivacy.Uploads.Storage.GCP = nil
	cfg.Services.DataPrivacy.Uploads.Storage.Provider = objectstorage.FilesystemProvider
	cfg.Services.DataPrivacy.Uploads.Storage.FilesystemConfig = &objectstorage.FilesystemConfig{RootDirectory: "/tmp"}

	cfg.Queues.DataChangesTopicName = dataChangesTopicName
}
