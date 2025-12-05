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
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	mealplanningconverters "github.com/dinnerdonebetter/backend/internal/services/mealplanning/grpc/converters"
	"github.com/dinnerdonebetter/backend/pkg/client"
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

func (h *mcpToolHelper) SearchForValidIngredients() (*mcp.Tool, mcp.ToolHandlerFor[*SearchValidIngredientsInvocation, *SearchValidIngredientsResult]) {
	tool := &mcp.Tool{
		Name:        "SearchForValidIngredients",
		Description: "Search for valid ingredients with optional filtering and query string",
		InputSchema: map[string]any{
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
		},
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
