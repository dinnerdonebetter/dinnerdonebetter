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

	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/branding"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/config/envvars"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/domain/oauth"
	"github.com/dinnerdonebetter/dinnerdonebetter/backend/pkg/client"

	"github.com/verygoodsoftwarenotvirus/platform/v4/encoding"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/logging"
	"github.com/verygoodsoftwarenotvirus/platform/v4/observability/tracing"
	"github.com/verygoodsoftwarenotvirus/platform/v4/version"

	"github.com/modelcontextprotocol/go-sdk/auth"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/spf13/pflag"
)

const (
	defaultMcpServerConfigurationFilepath = "deploy/environments/localdev/config_files/mcp_server_config.json"

	transportStdio = "stdio"
	transportSSE   = "sse"
	transportHTTP  = "http"

	defaultPort    = 8888
	defaultBaseURL = "http://localhost:8888"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse command-line flags
	transport := pflag.String("transport", transportHTTP, fmt.Sprintf("Transport method: %s, %s, or %s", transportStdio, transportSSE, transportHTTP))
	baseURL := pflag.String("base-url", defaultBaseURL, "Public base URL of the MCP server (used for OAuth2 metadata)")
	pflag.Parse()

	// Validate transport flag
	validTransports := map[string]bool{transportStdio: true, transportSSE: true, transportHTTP: true}
	if !validTransports[*transport] {
		log.Fatalf("Invalid transport method: %s. Allowed values are: %s, %s, %s", *transport, transportStdio, transportSSE, transportHTTP)
	}

	// Allow override via env var
	if envBase := os.Getenv("MCP_BASE_URL"); envBase != "" {
		*baseURL = envBase
	}

	ctx := context.Background()

	configFilepath := os.Getenv(config.ConfigurationFilePathEnvVarKey)
	if configFilepath == "" {
		configFilepath = defaultMcpServerConfigurationFilepath
	}

	// When running locally (not in Kubernetes), override with localhost values
	if os.Getenv(config.RunningInKubernetesEnvVarKey) != "true" {
		must(os.Setenv(envvars.APIServiceHTTPAPIServerURLEnvVarKey, "http://localhost:8000"))
		must(os.Setenv(envvars.APIServiceGrpcAPIServerURLEnvVarKey, ":8001"))
		must(os.Setenv(envvars.APIServiceOauth2APIClientIDEnvVarKey, strings.Repeat("A", oauth.ClientIDSize)))
		must(os.Setenv(envvars.APIServiceOauth2APIClientSecretEnvVarKey, strings.Repeat("A", oauth.ClientSecretSize)))
	}

	cfg, err := config.LoadConfigFromPath[config.MCPServiceConfig](ctx, configFilepath)
	if err != nil {
		log.Fatal(err)
	}

	pillars, err := cfg.Observability.ProvidePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}
	logger := pillars.Logger

	if err = cfg.ValidateWithContext(ctx); err != nil {
		log.Fatal(err)
	}

	// gRPC address for backend
	grpcAddr := cfg.APIServiceConnection.GRPCAPIServerURL
	if grpcAddr == "" {
		grpcAddr = ":8001"
	}

	// Build unauthenticated client for LoginForToken calls during OAuth2 flow.
	unauthedClient, err := client.BuildUnauthenticatedGRPCClient(grpcAddr)
	if err != nil {
		log.Fatalf("failed to build unauthenticated gRPC client: %v", err)
	}

	// Create token store for per-user auth.
	tokens := newTokenStore(
		grpcAddr,
		cfg.APIServiceConnection.OAuth2APIClientID,
		cfg.APIServiceConnection.OAuth2APIClientSecret,
		cfg.APIServiceConnection.HTTPAPIServerURL,
	)
	tokens.startCleanupLoop(ctx)

	helper := &mcpToolManager{tokens: tokens}
	server := helper.setupServer()

	log.Printf("serving now with transport: %s", *transport)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	go func() {
		<-signalChan
		os.Exit(0)
	}()

	switch *transport {
	case transportStdio:
		if err = server.Run(ctx, &mcp.StdioTransport{}); err != nil {
			logger.Error("serving MCP server via stdio", err)
			log.Fatal(err)
		}
	case transportSSE:
		sseHandler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			return server
		}, &mcp.SSEOptions{})

		mux := buildHTTPMux(sseHandler, tokens, pillars.Logger, pillars.TracerProvider, *baseURL, unauthedClient)

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", defaultPort),
			Handler:           mux,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err = srv.ListenAndServe(); err != nil {
			logger.Error("starting MCP server via SSE", err)
		}
	case transportHTTP:
		handlerOpts := &mcp.StreamableHTTPOptions{
			Stateless:      true,
			JSONResponse:   true,
			Logger:         slog.New(&slog.JSONHandler{}),
			EventStore:     mcp.NewMemoryEventStore(nil),
			SessionTimeout: 0,
		}
		streamHandler := mcp.NewStreamableHTTPHandler(func(request *http.Request) *mcp.Server {
			return server
		}, handlerOpts)

		mux := buildHTTPMux(streamHandler, tokens, pillars.Logger, pillars.TracerProvider, *baseURL, unauthedClient)

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", defaultPort),
			Handler:           mux,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err = srv.ListenAndServe(); err != nil {
			logger.Error("starting MCP server via HTTP", err)
		}
	}
}

// buildHTTPMux creates an HTTP mux with OAuth2 routes (unauthenticated) and the MCP handler (authenticated).
func buildHTTPMux(mcpHandler http.Handler, tokens *tokenStore, logger logging.Logger, tracerProvider tracing.TracerProvider, baseURL string, unauthedClient client.Client) *http.ServeMux {
	mux := http.NewServeMux()

	encoder := encoding.ProvideServerEncoderDecoder(logger, tracerProvider, encoding.ContentTypeJSON)
	// Ops routes (unauthenticated).
	mux.HandleFunc("GET /_ops_/live", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /_ops_/ready", func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("GET /_ops_/version", func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")
		encoder.EncodeResponseWithStatus(req.Context(), res, version.Get(), http.StatusOK)
	})

	// Register OAuth2 routes (no auth middleware — these handle authentication themselves).
	registerOAuth2Routes(mux, tokens, baseURL, unauthedClient)

	// Wrap the MCP handler with bearer token auth middleware.
	authMiddleware := auth.RequireBearerToken(tokens.verifyToken, &auth.RequireBearerTokenOptions{
		ResourceMetadataURL: baseURL + "/.well-known/oauth-protected-resource",
	})
	mux.Handle("/mcp", authMiddleware(mcpHandler))

	return mux
}

type mcpToolManager struct {
	tokens *tokenStore
}

// clientFromRequest resolves the per-user gRPC client from the MCP request's auth token.
func (h *mcpToolManager) clientFromRequest(req *mcp.CallToolRequest) (client.Client, error) {
	if req.Extra == nil || req.Extra.TokenInfo == nil {
		return nil, fmt.Errorf("not authenticated")
	}
	rawToken, ok := req.Extra.TokenInfo.Extra["raw_token"].(string)
	if !ok || rawToken == "" {
		return nil, fmt.Errorf("bearer token not found in request")
	}
	return h.tokens.clientForToken(rawToken)
}

func (h *mcpToolManager) setupServer() *mcp.Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: fmt.Sprintf("%s-mcp", branding.CompanyNameSlug), Version: "v1.0.0"}, nil)

	mcp.AddTool(mcpServer, getValidIngredientTool, h.GetValidIngredient())
	mcp.AddTool(mcpServer, searchForValidIngredientsTool, h.SearchForValidIngredients())
	mcp.AddTool(mcpServer, validIngredientCreationTool, h.CreateValidIngredient())
	mcp.AddTool(mcpServer, validIngredientUpdateTool, h.UpdateValidIngredient())

	mcp.AddTool(mcpServer, getValidPreparationTool, h.GetValidPreparation())
	mcp.AddTool(mcpServer, searchForValidPreparationsTool, h.SearchForValidPreparations())
	mcp.AddTool(mcpServer, validPreparationCreationTool, h.CreateValidPreparation())
	mcp.AddTool(mcpServer, validPreparationUpdateTool, h.UpdateValidPreparation())

	mcp.AddTool(mcpServer, getValidMeasurementUnitTool, h.GetValidMeasurementUnit())
	mcp.AddTool(mcpServer, searchForValidMeasurementUnitsTool, h.SearchForValidMeasurementUnits())
	mcp.AddTool(mcpServer, validMeasurementUnitCreationTool, h.CreateValidMeasurementUnit())
	mcp.AddTool(mcpServer, validMeasurementUnitUpdateTool, h.UpdateValidMeasurementUnit())

	mcp.AddTool(mcpServer, getValidIngredientPreparationTool, h.GetValidIngredientPreparation())
	mcp.AddTool(mcpServer, getValidIngredientPreparationsTool, h.GetValidIngredientPreparations())
	mcp.AddTool(mcpServer, validIngredientPreparationCreationTool, h.CreateValidIngredientPreparation())
	mcp.AddTool(mcpServer, validIngredientPreparationUpdateTool, h.UpdateValidIngredientPreparation())

	mcp.AddTool(mcpServer, getValidPrepTaskConfigTool, h.GetValidPrepTaskConfig())
	mcp.AddTool(mcpServer, getValidPrepTaskConfigsTool, h.GetValidPrepTaskConfigs())
	mcp.AddTool(mcpServer, getValidPrepTaskConfigsByIngredientTool, h.GetValidPrepTaskConfigsByIngredient())
	mcp.AddTool(mcpServer, getValidPrepTaskConfigsByPreparationTool, h.GetValidPrepTaskConfigsByPreparation())
	mcp.AddTool(mcpServer, getValidPrepTaskConfigsByIngredientAndPreparationTool, h.GetValidPrepTaskConfigsByIngredientAndPreparation())
	mcp.AddTool(mcpServer, validPrepTaskConfigCreationTool, h.CreateValidPrepTaskConfig())
	mcp.AddTool(mcpServer, validPrepTaskConfigUpdateTool, h.UpdateValidPrepTaskConfig())

	mcp.AddTool(mcpServer, getValidIngredientMeasurementUnitTool, h.GetValidIngredientMeasurementUnit())
	mcp.AddTool(mcpServer, getValidIngredientMeasurementUnitsTool, h.GetValidIngredientMeasurementUnits())
	mcp.AddTool(mcpServer, validIngredientMeasurementUnitCreationTool, h.CreateValidIngredientMeasurementUnit())
	mcp.AddTool(mcpServer, validIngredientMeasurementUnitUpdateTool, h.UpdateValidIngredientMeasurementUnit())

	mcp.AddTool(mcpServer, getValidVesselTool, h.GetValidVessel())
	mcp.AddTool(mcpServer, searchForValidVesselsTool, h.SearchForValidVessels())
	mcp.AddTool(mcpServer, validVesselCreationTool, h.CreateValidVessel())
	mcp.AddTool(mcpServer, validVesselUpdateTool, h.UpdateValidVessel())

	mcp.AddTool(mcpServer, getValidMeasurementUnitConversionTool, h.GetValidMeasurementUnitConversion())
	mcp.AddTool(mcpServer, getValidMeasurementUnitConversionsForUnitTool, h.GetValidMeasurementUnitConversionsForUnit())
	mcp.AddTool(mcpServer, getValidMeasurementUnitConversionsForIngredientsTool, h.GetValidMeasurementUnitConversionsForIngredients())
	mcp.AddTool(mcpServer, validMeasurementUnitConversionCreationTool, h.CreateValidMeasurementUnitConversion())
	mcp.AddTool(mcpServer, validMeasurementUnitConversionUpdateTool, h.UpdateValidMeasurementUnitConversion())

	mcp.AddTool(mcpServer, getValidIngredientStateTool, h.GetValidIngredientState())
	mcp.AddTool(mcpServer, searchForValidIngredientStatesTool, h.SearchForValidIngredientStates())
	mcp.AddTool(mcpServer, validIngredientStateCreationTool, h.CreateValidIngredientState())
	mcp.AddTool(mcpServer, validIngredientStateUpdateTool, h.UpdateValidIngredientState())

	mcp.AddTool(mcpServer, getValidIngredientStateIngredientTool, h.GetValidIngredientStateIngredient())
	mcp.AddTool(mcpServer, getValidIngredientStateIngredientsTool, h.GetValidIngredientStateIngredients())
	mcp.AddTool(mcpServer, validIngredientStateIngredientCreationTool, h.CreateValidIngredientStateIngredient())
	mcp.AddTool(mcpServer, validIngredientStateIngredientUpdateTool, h.UpdateValidIngredientStateIngredient())

	mcp.AddTool(mcpServer, getValidInstrumentTool, h.GetValidInstrument())
	mcp.AddTool(mcpServer, searchForValidInstrumentsTool, h.SearchForValidInstruments())
	mcp.AddTool(mcpServer, validInstrumentCreationTool, h.CreateValidInstrument())
	mcp.AddTool(mcpServer, validInstrumentUpdateTool, h.UpdateValidInstrument())

	mcp.AddTool(mcpServer, getValidPreparationInstrumentTool, h.GetValidPreparationInstrument())
	mcp.AddTool(mcpServer, getValidPreparationInstrumentsTool, h.GetValidPreparationInstruments())
	mcp.AddTool(mcpServer, validPreparationInstrumentCreationTool, h.CreateValidPreparationInstrument())
	mcp.AddTool(mcpServer, validPreparationInstrumentUpdateTool, h.UpdateValidPreparationInstrument())

	mcp.AddTool(mcpServer, getValidPreparationVesselTool, h.GetValidPreparationVessel())
	mcp.AddTool(mcpServer, getValidPreparationVesselsTool, h.GetValidPreparationVessels())
	mcp.AddTool(mcpServer, validPreparationVesselCreationTool, h.CreateValidPreparationVessel())
	mcp.AddTool(mcpServer, validPreparationVesselUpdateTool, h.UpdateValidPreparationVessel())

	mcp.AddTool(mcpServer, getRecipeStepInstrumentTool, h.GetRecipeStepInstrument())
	mcp.AddTool(mcpServer, getRecipeStepInstrumentsTool, h.GetRecipeStepInstruments())
	mcp.AddTool(mcpServer, recipeStepInstrumentCreationTool, h.CreateRecipeStepInstrument())
	mcp.AddTool(mcpServer, recipeStepInstrumentUpdateTool, h.UpdateRecipeStepInstrument())

	mcp.AddTool(mcpServer, getRecipeStepProductTool, h.GetRecipeStepProduct())
	mcp.AddTool(mcpServer, getRecipeStepProductsTool, h.GetRecipeStepProducts())
	mcp.AddTool(mcpServer, recipeStepProductCreationTool, h.CreateRecipeStepProduct())
	mcp.AddTool(mcpServer, recipeStepProductUpdateTool, h.UpdateRecipeStepProduct())

	mcp.AddTool(mcpServer, getRecipeStepIngredientTool, h.GetRecipeStepIngredient())
	mcp.AddTool(mcpServer, getRecipeStepIngredientsTool, h.GetRecipeStepIngredients())
	mcp.AddTool(mcpServer, recipeStepIngredientCreationTool, h.CreateRecipeStepIngredient())
	mcp.AddTool(mcpServer, recipeStepIngredientUpdateTool, h.UpdateRecipeStepIngredient())

	mcp.AddTool(mcpServer, getRecipePrepTaskTool, h.GetRecipePrepTask())
	mcp.AddTool(mcpServer, getRecipePrepTasksTool, h.GetRecipePrepTasks())
	mcp.AddTool(mcpServer, recipePrepTaskCreationTool, h.CreateRecipePrepTask())
	mcp.AddTool(mcpServer, recipePrepTaskUpdateTool, h.UpdateRecipePrepTask())

	mcp.AddTool(mcpServer, getRecipeStepVesselTool, h.GetRecipeStepVessel())
	mcp.AddTool(mcpServer, getRecipeStepVesselsTool, h.GetRecipeStepVessels())
	mcp.AddTool(mcpServer, recipeStepVesselCreationTool, h.CreateRecipeStepVessel())
	mcp.AddTool(mcpServer, recipeStepVesselUpdateTool, h.UpdateRecipeStepVessel())

	mcp.AddTool(mcpServer, getRecipeStepCompletionConditionTool, h.GetRecipeStepCompletionCondition())
	mcp.AddTool(mcpServer, getRecipeStepCompletionConditionsTool, h.GetRecipeStepCompletionConditions())
	mcp.AddTool(mcpServer, recipeStepCompletionConditionCreationTool, h.CreateRecipeStepCompletionCondition())
	mcp.AddTool(mcpServer, recipeStepCompletionConditionUpdateTool, h.UpdateRecipeStepCompletionCondition())

	mcp.AddTool(mcpServer, getRecipeStepTool, h.GetRecipeStep())
	mcp.AddTool(mcpServer, getRecipeStepsTool, h.GetRecipeSteps())
	mcp.AddTool(mcpServer, recipeStepCreationTool, h.CreateRecipeStep())
	mcp.AddTool(mcpServer, recipeStepUpdateTool, h.UpdateRecipeStep())

	mcp.AddTool(mcpServer, getRecipeTool, h.GetRecipe())
	mcp.AddTool(mcpServer, getRecipesTool, h.GetRecipes())
	mcp.AddTool(mcpServer, searchForRecipesTool, h.SearchForRecipes())
	mcp.AddTool(mcpServer, recipeCreationTool, h.CreateRecipe())
	mcp.AddTool(mcpServer, recipeUpdateTool, h.UpdateRecipe())

	// Issue Reports
	mcp.AddTool(mcpServer, getIssueReportTool, h.GetIssueReport())
	mcp.AddTool(mcpServer, getIssueReportsTool, h.GetIssueReports())
	mcp.AddTool(mcpServer, getIssueReportsForAccountTool, h.GetIssueReportsForAccount())
	mcp.AddTool(mcpServer, createIssueReportTool, h.CreateIssueReport())
	mcp.AddTool(mcpServer, updateIssueReportTool, h.UpdateIssueReport())
	mcp.AddTool(mcpServer, archiveIssueReportTool, h.ArchiveIssueReport())

	// Users / Identity
	mcp.AddTool(mcpServer, getUserTool, h.GetUser())
	mcp.AddTool(mcpServer, getUsersTool, h.GetUsers())
	mcp.AddTool(mcpServer, searchForUsersTool, h.SearchForUsers())
	mcp.AddTool(mcpServer, getAccountTool, h.GetAccount())
	mcp.AddTool(mcpServer, getAccountsForUserTool, h.GetAccountsForUser())
	mcp.AddTool(mcpServer, updateUserDetailsTool, h.UpdateUserDetails())

	// Webhooks
	mcp.AddTool(mcpServer, getWebhookTool, h.GetWebhook())
	mcp.AddTool(mcpServer, getWebhooksTool, h.GetWebhooks())
	mcp.AddTool(mcpServer, createWebhookTool, h.CreateWebhook())
	mcp.AddTool(mcpServer, archiveWebhookTool, h.ArchiveWebhook())
	mcp.AddTool(mcpServer, getWebhookTriggerEventsTool, h.GetWebhookTriggerEvents())
	mcp.AddTool(mcpServer, createWebhookTriggerEventTool, h.CreateWebhookTriggerEvent())
	mcp.AddTool(mcpServer, addWebhookTriggerConfigTool, h.AddWebhookTriggerConfig())
	mcp.AddTool(mcpServer, archiveWebhookTriggerConfigTool, h.ArchiveWebhookTriggerConfig())

	// Waitlists
	mcp.AddTool(mcpServer, getWaitlistTool, h.GetWaitlist())
	mcp.AddTool(mcpServer, getWaitlistsTool, h.GetWaitlists())
	mcp.AddTool(mcpServer, getActiveWaitlistsTool, h.GetActiveWaitlists())
	mcp.AddTool(mcpServer, createWaitlistTool, h.CreateWaitlist())
	mcp.AddTool(mcpServer, archiveWaitlistTool, h.ArchiveWaitlist())
	mcp.AddTool(mcpServer, getWaitlistSignupTool, h.GetWaitlistSignup())
	mcp.AddTool(mcpServer, getWaitlistSignupsForWaitlistTool, h.GetWaitlistSignupsForWaitlist())
	mcp.AddTool(mcpServer, createWaitlistSignupTool, h.CreateWaitlistSignup())
	mcp.AddTool(mcpServer, archiveWaitlistSignupTool, h.ArchiveWaitlistSignup())

	return mcpServer
}
