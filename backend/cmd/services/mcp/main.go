package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/config/envvars"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	grpcconverters "github.com/dinnerdonebetter/backend/internal/grpc/converters"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mealplanningconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"
	"github.com/invopop/jsonschema"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pquerna/otp/totp"
)

const (
	jsonSchemaVersion = "https://json-schema.org/draft/2020-12/schema"
	objType           = "object"
	arrType           = "array"
	strType           = "string"
	boolType          = "boolean"
	intType           = "integer"
	dtFmt             = "date-time"
)

const (
	adminServerConfigurationFilepath = "deploy/environments/localdev/config_files/admin_webapp_config.json"

	// TODO: get these values another way
	tempUsername     = "admin_user"
	tempPassword     = "admin_pass"
	tempTOTPTokenKey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ctx := context.Background()

	// We don't yet have a way to write these values into the AdminWebappConfig (because they're not present in the root APIConfig struct).
	// This approach is an atrocious hack that I have to employ because I wasn't smart enough to design a better config generation system.
	must(os.Setenv(envvars.APIServiceHTTPAPIServerURLEnvVarKey, "http://localhost:8000"))
	must(os.Setenv(envvars.APIServiceGrpcAPIServerURLEnvVarKey, ":8001"))
	must(os.Setenv(envvars.APIServiceOauth2APIClientIDEnvVarKey, strings.Repeat("A", oauth.ClientIDSize)))
	must(os.Setenv(envvars.APIServiceOauth2APIClientSecretEnvVarKey, strings.Repeat("A", oauth.ClientSecretSize)))

	cfg, err := config.LoadConfigFromPath[config.AdminWebappConfig](ctx, adminServerConfigurationFilepath)
	if err != nil {
		log.Fatal(err)
	}

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = logger, tracerProvider

	if err = cfg.ValidateWithContext(ctx); err != nil {
		log.Fatal(err)
	}

	// Build gRPC client
	grpcAddr := cfg.APIServiceConnection.GRPCAPIServerURL
	if grpcAddr == "" {
		grpcAddr = ":8001" // fallback to default
	}

	unauthedClient, err := client.BuildUnauthenticatedGRPCClient(grpcAddr)
	if err != nil {
		log.Fatalf("failed to build gRPC client: %v", err)
	}

	totpToken, err := totp.GenerateCode(strings.ToUpper(tempTOTPTokenKey), time.Now().UTC())
	if err != nil {
		log.Fatalf("failed to build generate TOTP code: %v", err)
	}

	tokenRes, err := unauthedClient.AdminLoginForToken(ctx, &authsvc.AdminLoginForTokenRequest{
		Input: &authsvc.UserLoginInput{
			Username:  tempUsername,
			Password:  tempPassword,
			TOTPToken: totpToken,
		},
	})
	if err != nil {
		log.Fatalf("failed to get access token: %v", err)
	}

	c, err := localdev.BuildInsecureOAuthedGRPCClient(
		ctx,
		cfg.APIServiceConnection.OAuth2APIClientID,
		cfg.APIServiceConnection.OAuth2APIClientSecret,
		cfg.APIServiceConnection.HTTPAPIServerURL,
		grpcAddr,
		tokenRes.Result.AccessToken,
	)

	helper := &mcpToolHelper{
		logger: logger,
		tracer: tracing.NewTracer(tracerProvider.Tracer("mcp-server")),
		client: c,
	}
	server := helper.setupServer()

	log.Println("serving now")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	handlerOpts := &mcp.StreamableHTTPOptions{
		Stateless:      true,
		JSONResponse:   true,
		Logger:         slog.New(&slog.JSONHandler{}),
		EventStore:     mcp.NewMemoryEventStore(nil),
		SessionTimeout: 0,
	}
	handler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
		return server
	}, handlerOpts)

	go func() {
		if err = http.ListenAndServe(fmt.Sprintf(":%d", 8888), handler); err != nil {
			logger.Error("starting MCP server", err)
		}
	}()

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
	}()
}

// queryFilterSchema returns the JSON schema for a QueryFilter object
func queryFilterSchema() map[string]any {
	return map[string]any{
		"type": objType,
		"properties": map[string]any{
			"SortBy": map[string]any{
				"type":        strType,
				"description": "Field to sort by",
			},
			"CreatedAfter": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results created after this timestamp (ISO 8601)",
			},
			"CreatedBefore": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results created before this timestamp (ISO 8601)",
			},
			"UpdatedAfter": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results updated after this timestamp (ISO 8601)",
			},
			"UpdatedBefore": map[string]any{
				"type":        strType,
				"format":      dtFmt,
				"description": "Filter results updated before this timestamp (ISO 8601)",
			},
			"PageSize": map[string]any{
				"type":        intType,
				"description": "Maximum number of results to return",
			},
			"IncludeArchived": map[string]any{
				"type":        boolType,
				"description": "Whether to include archived items",
			},
			"Cursor": map[string]any{
				"type":        strType,
				"description": "Pagination cursor for fetching next page",
			},
		},
	}
}

type (
	SearchValidIngredientsInvocation struct {
		Query            string
		Filter           *filtering.QueryFilter
		UseSearchService bool
	}

	SearchValidIngredientsResult struct {
		Results []*mealplanning.ValidIngredient
	}
)

type mcpToolHelper struct {
	logger logging.Logger
	tracer tracing.Tracer
	client client.Client
}

func (h *mcpToolHelper) setupServer() *mcp.Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "dinnerdonebetter-mcp", Version: "v1.0.0"}, nil)

	validIngredientSearchTool, validIngredientSearchFunc := h.SearchForValidIngredients()
	mcp.AddTool(mcpServer, validIngredientSearchTool, validIngredientSearchFunc)

	return mcpServer
}

func validIngredientSearchInputSchema() map[string]any {
	return map[string]any{
		"$schema": jsonSchemaVersion,
		"type":    objType,
		"properties": map[string]any{
			"Filter": queryFilterSchema(),
			"Query": map[string]any{
				"type":        strType,
				"description": "Search query string to match ingredient names or descriptions",
			},
			"UseSearchService": map[string]any{
				"type":        boolType,
				"description": "Whether to use the search service for more advanced search capabilities",
			},
		},
	}
}

func validIngredientOutputSchema() map[string]any {
	return map[string]any{
		"$schema": jsonSchemaVersion,
		"type":    objType,
		"items": map[string]any{
			"type": arrType,
			"name": "Results",
			"properties": map[string]any{
				"createdAt": map[string]any{
					"type":   strType,
					"format": dtFmt,
				},
				"lastUpdatedAt": map[string]any{
					"type":   []any{strType, "null"},
					"format": dtFmt,
				},
				"archivedAt": map[string]any{
					"type":   []any{strType, "null"},
					"format": dtFmt,
				},
				"storageTemperatureInCelsius": map[string]any{
					"type": objType,
					"properties": map[string]any{
						"min": map[string]any{
							"type": []any{"number", "null"},
						},
						"max": map[string]any{
							"type": []any{"number", "null"},
						},
					},
				},
				"iconPath": map[string]any{
					"type": strType,
				},
				"warning": map[string]any{
					"type": strType,
				},
				"pluralName": map[string]any{
					"type": strType,
				},
				"storageInstructions": map[string]any{
					"type": strType,
				},
				"name": map[string]any{
					"type": strType,
				},
				"id": map[string]any{
					"type": strType,
				},
				"description": map[string]any{
					"type": strType,
				},
				"slug": map[string]any{
					"type": strType,
				},
				"shoppingSuggestions": map[string]any{
					"type": strType,
				},
				"containsShellfish": map[string]any{
					"type": boolType,
				},
				"isLiquid": map[string]any{
					"type": boolType,
				},
				"containsPeanut": map[string]any{
					"type": boolType,
				},
				"containsTreeNut": map[string]any{
					"type": boolType,
				},
				"containsEgg": map[string]any{
					"type": boolType,
				},
				"containsWheat": map[string]any{
					"type": boolType,
				},
				"containsSoy": map[string]any{
					"type": boolType,
				},
				"animalDerived": map[string]any{
					"type": boolType,
				},
				"restrictToPreparations": map[string]any{
					"type": boolType,
				},
				"containsSesame": map[string]any{
					"type": boolType,
				},
				"containsFish": map[string]any{
					"type": boolType,
				},
				"containsGluten": map[string]any{
					"type": boolType,
				},
				"containsDairy": map[string]any{
					"type": boolType,
				},
				"containsAlcohol": map[string]any{
					"type": boolType,
				},
				"animalFlesh": map[string]any{
					"type": boolType,
				},
				"isStarch": map[string]any{
					"type": boolType,
				},
				"isProtein": map[string]any{
					"type": boolType,
				},
				"isGrain": map[string]any{
					"type": boolType,
				},
				"isFruit": map[string]any{
					"type": boolType,
				},
				"isSalt": map[string]any{
					"type": boolType,
				},
				"isFat": map[string]any{
					"type": boolType,
				},
				"isAcid": map[string]any{
					"type": boolType,
				},
				"isHeat": map[string]any{
					"type": boolType,
				},
			},
		},
	}
}

func schemaForType(x any) map[string]any {
	var y map[string]any
	encoding.MustDecodeJSON(encoding.MustEncodeJSON(jsonschema.Reflect(x)), &y)

	// Transform the schema to match the expected format
	result := transformSchema(y)
	return result
}

func transformSchema(schema map[string]any) map[string]any {
	// Keep $defs for reference resolution
	defs, hasDefs := schema["$defs"].(map[string]any)

	// If schema has $ref, resolve it from $defs
	resolvedSchema := schema
	if ref, ok := schema["$ref"].(string); ok && hasDefs {
		if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
			if def, ok := defs[refName].(map[string]any); ok {
				resolvedSchema = def
			}
		}
	}

	// Build the result with $schema and type
	result := map[string]any{
		"$schema": jsonSchemaVersion,
		"type":    objType,
	}

	// Extract properties
	props, ok := resolvedSchema["properties"].(map[string]any)
	if !ok {
		return result
	}

	// Special case: if there's only a "Results" or "results" property that's an array, transform to items format
	if len(props) == 1 {
		var resultsProp map[string]any
		var found bool
		// Check both camelCase and PascalCase
		if resultsProp, found = props["results"].(map[string]any); !found {
			resultsProp, found = props["Results"].(map[string]any)
		}
		if found {
			if propType, ok := resultsProp["type"].(string); ok && propType == "array" {
				// This is the SearchValidIngredientsResult case
				itemsSchema := transformResultsArraySchema(resultsProp, defs, hasDefs)
				result["items"] = itemsSchema
				return result
			}
		}
	}

	// Transform properties
	transformedProps := make(map[string]any)
	for key, value := range props {
		prop, ok := value.(map[string]any)
		if !ok {
			transformedProps[key] = value
			continue
		}

		// Handle $ref in properties (for nested types like QueryFilter)
		if ref, ok := prop["$ref"].(string); ok && hasDefs {
			if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
				if def, ok := defs[refName].(map[string]any); ok {
					prop = transformQueryFilterSchema(def)
				}
			}
		}

		// Convert camelCase to PascalCase for field names
		pascalKey := toPascalCase(key)

		// Add description if not present
		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}
		if _, hasDesc := propCopy["description"]; !hasDesc {
			if desc := getFieldDescription(pascalKey); desc != "" {
				propCopy["description"] = desc
			}
		}

		transformedProps[pascalKey] = propCopy
	}

	result["properties"] = transformedProps
	return result
}

func transformQueryFilterSchema(schema map[string]any) map[string]any {
	result := map[string]any{
		"type": objType,
	}

	props, ok := schema["properties"].(map[string]any)
	if !ok {
		return result
	}

	transformedProps := make(map[string]any)
	for key, value := range props {
		prop, ok := value.(map[string]any)
		if !ok {
			continue
		}

		// Convert camelCase to PascalCase
		pascalKey := toPascalCase(key)

		// Special case: limit -> PageSize
		if key == "limit" {
			pascalKey = "PageSize"
		}

		// Add descriptions based on field name
		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}

		// Add description if not present
		if _, hasDesc := propCopy["description"]; !hasDesc {
			propCopy["description"] = getFieldDescription(pascalKey)
		}

		transformedProps[pascalKey] = propCopy
	}

	result["properties"] = transformedProps
	return result
}

func transformResultsArraySchema(arrayProp map[string]any, defs map[string]any, hasDefs bool) map[string]any {
	result := map[string]any{
		"name": "Results",
		"type": arrType,
	}

	// Get the items schema
	items, ok := arrayProp["items"].(map[string]any)
	if !ok {
		return result
	}

	// Resolve $ref if present
	var itemSchema map[string]any
	if ref, ok := items["$ref"].(string); ok && hasDefs {
		if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
			if def, ok := defs[refName].(map[string]any); ok {
				itemSchema = def
			}
		}
	} else {
		itemSchema = items
	}

	if itemSchema == nil {
		return result
	}

	// Extract properties from the item schema
	props, ok := itemSchema["properties"].(map[string]any)
	if !ok {
		return result
	}

	// Transform properties (keep camelCase as expected)
	transformedProps := make(map[string]any)
	for key, value := range props {
		prop, ok := value.(map[string]any)
		if !ok {
			transformedProps[key] = value
			continue
		}

		propCopy := make(map[string]any)
		for k, v := range prop {
			propCopy[k] = v
		}

		// Handle $ref in properties (for nested types like OptionalFloat32Range)
		if ref, ok := propCopy["$ref"].(string); ok && hasDefs {
			if refName := strings.TrimPrefix(ref, "#/$defs/"); refName != ref {
				if def, ok := defs[refName].(map[string]any); ok {
					// Replace $ref with the actual schema
					for k, v := range def {
						propCopy[k] = v
					}
					delete(propCopy, "$ref")
				}
			}
		}

		// Handle nullable types - check if the field is a pointer type
		// For nullable fields, type should be ["string", "null"] or ["number", "null"]
		if propType, ok := propCopy["type"].(string); ok {
			// Check if this is a nullable field based on the expected schema
			nullableFields := map[string]bool{
				"archivedAt":    true,
				"lastUpdatedAt": true,
			}
			if nullableFields[key] {
				propCopy["type"] = []any{propType, "null"}
			}
		}

		// Handle nested objects with nullable number fields
		if propType, ok := propCopy["type"].(string); ok && propType == objType {
			// Remove additionalProperties if present
			delete(propCopy, "additionalProperties")
			if nestedProps, ok := propCopy["properties"].(map[string]any); ok {
				transformedNestedProps := make(map[string]any)
				for nestedKey, nestedValue := range nestedProps {
					nestedProp, ok := nestedValue.(map[string]any)
					if !ok {
						transformedNestedProps[nestedKey] = nestedValue
						continue
					}
					nestedPropCopy := make(map[string]any)
					for k, v := range nestedProp {
						nestedPropCopy[k] = v
					}
					// Make min and max nullable
					if nestedKey == "min" || nestedKey == "max" {
						if nestedType, ok := nestedPropCopy["type"].(string); ok && nestedType == "number" {
							nestedPropCopy["type"] = []any{"number", "null"}
						}
					}
					transformedNestedProps[nestedKey] = nestedPropCopy
				}
				propCopy["properties"] = transformedNestedProps
			}
		}

		transformedProps[key] = propCopy
	}

	result["properties"] = transformedProps
	return result
}

func toPascalCase(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func getFieldDescription(fieldName string) string {
	descriptions := map[string]string{
		"SortBy":           "Field to sort by",
		"CreatedAfter":     "Filter results created after this timestamp (ISO 8601)",
		"CreatedBefore":    "Filter results created before this timestamp (ISO 8601)",
		"UpdatedAfter":     "Filter results updated after this timestamp (ISO 8601)",
		"UpdatedBefore":    "Filter results updated before this timestamp (ISO 8601)",
		"PageSize":         "Maximum number of results to return",
		"IncludeArchived":  "Whether to include archived items",
		"Cursor":           "Pagination cursor for fetching next page",
		"Query":            "Search query string to match ingredient names or descriptions",
		"UseSearchService": "Whether to use the search service for more advanced search capabilities",
	}
	if desc, ok := descriptions[fieldName]; ok {
		return desc
	}
	return ""
}

func (h *mcpToolHelper) SearchForValidIngredients() (*mcp.Tool, mcp.ToolHandlerFor[*SearchValidIngredientsInvocation, *SearchValidIngredientsResult]) {
	tool := &mcp.Tool{
		Name:         "SearchForValidIngredients",
		Description:  "Search for valid ingredients with optional filtering and query string",
		InputSchema:  validIngredientSearchInputSchema(),
		OutputSchema: validIngredientOutputSchema(), //encoding.MustEncodeJSON(jsonschema.Reflect(SearchValidIngredientsResult{})),
	}

	return tool, func(ctx context.Context, _ *mcp.CallToolRequest, x *SearchValidIngredientsInvocation) (*mcp.CallToolResult, *SearchValidIngredientsResult, error) {
		results, err := h.client.SearchForValidIngredients(ctx, &mealplanninggrpc.SearchForValidIngredientsRequest{
			Filter:           grpcconverters.ConvertQueryFilterToGRPCQueryFilter(x.Filter, filtering.Pagination{}),
			Query:            x.Query,
			UseSearchService: x.UseSearchService,
		})
		if err != nil {
			return nil, nil, err
		}

		out := &SearchValidIngredientsResult{}
		for _, result := range results.Results {
			out.Results = append(out.Results, mealplanningconverters.ConvertGRPCValidIngredientToValidIngredient(result))
		}

		return nil, out, nil
	}
}
