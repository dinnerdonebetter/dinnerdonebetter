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
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"
	authsvc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/auth"
	"github.com/dinnerdonebetter/backend/internal/localdev"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/pquerna/otp/totp"
	"github.com/spf13/pflag"
)

const (
	adminServerConfigurationFilepath = "deploy/environments/localdev/config_files/admin_webapp_config.json"

	// TODO: get these values another way.
	tempUsername     = "admin_user"
	tempPassword     = "admin_pass"
	tempTOTPTokenKey = "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="

	transportStdio = "stdio"
	transportSSE   = "sse"
	transportHTTP  = "http"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Parse command-line flags
	transport := pflag.String("transport", transportHTTP, fmt.Sprintf("Transport method: %s, %s, or %s", transportStdio, transportSSE, transportHTTP))
	pflag.Parse()

	// Validate transport flag
	validTransports := map[string]bool{transportStdio: true, transportSSE: true, transportHTTP: true}
	if !validTransports[*transport] {
		log.Fatalf("Invalid transport method: %s. Allowed values are: %s, %s, %s", *transport, transportStdio, transportSSE, transportHTTP)
	}

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
	if err != nil {
		log.Fatalf("failed to build authenticated gRPC client: %v", err)
	}

	helper := &mcpToolManager{client: c}
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
		// Serve using stdio transport
		if err = server.Run(ctx, &mcp.StdioTransport{}); err != nil {
			logger.Error("serving MCP server via stdio", err)
			log.Fatal(err)
		}
	case transportSSE:
		// Serve using SSE transport
		handler := mcp.NewSSEHandler(func(request *http.Request) *mcp.Server {
			return server
		}, &mcp.SSEOptions{})

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", 8888),
			Handler:           handler,
			ReadTimeout:       15 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
		}
		if err = srv.ListenAndServe(); err != nil {
			logger.Error("starting MCP server via SSE", err)
		}
	case transportHTTP:
		// Serve using HTTP transport
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

		srv := &http.Server{
			Addr:              fmt.Sprintf(":%d", 8888),
			Handler:           handler,
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

type mcpToolManager struct {
	client client.Client
}

func (h *mcpToolManager) setupServer() *mcp.Server {
	mcpServer := mcp.NewServer(&mcp.Implementation{Name: "dinnerdonebetter-mcp", Version: "v1.0.0"}, nil)

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

	return mcpServer
}
