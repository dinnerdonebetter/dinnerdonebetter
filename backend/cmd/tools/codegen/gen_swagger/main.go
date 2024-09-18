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
	"github.com/dinnerdonebetter/backend/internal/uploads/objectstorage"

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
		switch {
		case a.name < b.name:
			return -1
		case a.name == b.name:
			return 0
		default:
			return 1
		}
	})

	rawCfg, err := os.ReadFile("environments/local/config_files/service-config.json")
	if err != nil {
		log.Fatal(err)
	}

	var cfg *config.InstanceConfig
	if err = json.Unmarshal(rawCfg, &cfg); err != nil {
		log.Fatal(err)
	}

	neutralizeConfig(cfg)

	// build our server struct.
	srv, err := build.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	allTags := map[string]struct{}{}

	routeDefinitions := []*RouteDefinition{}
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

		routeDef := &RouteDefinition{
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
					if rep, ok1 := tagReplacements[part]; ok1 {
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
	for tag := range allTags {
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

const (
	/* #nosec G101 */
	debugCookieHashKey = "HEREISA32CHARSECRETWHICHISMADEUP"
	/* #nosec G101 */
	debugCookieBlockKey  = "DIFFERENT32CHARSECRETTHATIMADEUP"
	dataChangesTopicName = "dataChangesTopicName"
)

func neutralizeConfig(cfg *config.InstanceConfig) {
	if err := os.Setenv("GOOGLE_CLOUD_PROJECT_ID", "something"); err != nil {
		panic(err)
	}

	cfg.Database.RunMigrations = false
	cfg.Database.OAuth2TokenEncryptionKey = "BLAHBLAHBLAHBLAHBLAHBLAHBLAHBLAH"
	cfg.Services.Auth.Cookies.HashKey = debugCookieHashKey
	cfg.Services.Auth.Cookies.BlockKey = debugCookieBlockKey
	cfg.Services.Auth.SSO.Google.ClientID = "blah blah blah blah"
	cfg.Services.Auth.SSO.Google.ClientSecret = "blah blah blah blah"
	cfg.Analytics.Provider = ""
	cfg.Services.Recipes.Uploads.Storage.GCPConfig = nil
	cfg.Services.Recipes.Uploads.Storage.Provider = objectstorage.FilesystemProvider
	cfg.Services.Recipes.Uploads.Storage.FilesystemConfig = &objectstorage.FilesystemConfig{RootDirectory: "/tmp"}
	cfg.Services.RecipeSteps.Uploads.Storage.GCPConfig = nil
	cfg.Services.RecipeSteps.Uploads.Storage.Provider = objectstorage.FilesystemProvider
	cfg.Services.RecipeSteps.Uploads.Storage.FilesystemConfig = &objectstorage.FilesystemConfig{RootDirectory: "/tmp"}
	cfg.Services.ValidMeasurementUnits.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientPreparations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparationInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidInstrumentMeasurementUnits.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Recipes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeSteps.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepProducts.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepInstruments.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Meals.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlans.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanEvents.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanOptionVotes.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanTasks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Households.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInvitations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Users.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientGroups.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Webhooks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Auth.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipePrepTasks.DataChangesTopicName = dataChangesTopicName
	cfg.Services.MealPlanGroceryListItems.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidMeasurementUnitConversions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStates.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepCompletionConditions.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidIngredientStateIngredients.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeStepVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ServiceSettings.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ServiceSettingConfigurations.DataChangesTopicName = dataChangesTopicName
	cfg.Services.UserIngredientPreferences.DataChangesTopicName = dataChangesTopicName
	cfg.Services.RecipeRatings.DataChangesTopicName = dataChangesTopicName
	cfg.Services.HouseholdInstrumentOwnerships.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.ValidPreparationVessels.DataChangesTopicName = dataChangesTopicName
	cfg.Services.Workers.DataChangesTopicName = dataChangesTopicName
	cfg.Services.UserNotifications.DataChangesTopicName = dataChangesTopicName
}
