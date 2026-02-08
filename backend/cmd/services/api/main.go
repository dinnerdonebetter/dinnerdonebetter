package main

import (
	"context"
	"log"
	"os"

	apiserver "github.com/dinnerdonebetter/backend/internal/build/services/api"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/version"
	"github.com/spf13/cobra"

	_ "go.uber.org/automaxprocs"
)

func main() {
	root := &cobra.Command{
		Use:   "server",
		Short: "API server CLI",
	}
	root.AddCommand(serveCmd())
	root.AddCommand(versionCmd())

	if err := root.Execute(); err != nil {
		os.Exit(1)
	}
}

func serveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "serve",
		Short: "Run the API server (HTTP + gRPC)",
		RunE:  runServe,
	}
}

func runServe(_ *cobra.Command, _ []string) error {
	rootCtx := context.Background()

	cfg, err := config.LoadConfigFromEnvironment[config.APIServiceConfig]()
	if err != nil {
		return err
	}

	buildCtx, cancel := context.WithTimeout(rootCtx, cfg.HTTPServer.StartupDeadline)
	defer cancel()

	logger, err := cfg.Observability.Logging.ProvideLogger(rootCtx)
	if err != nil {
		log.Fatalf("could not create logger: %v", err)
	}

	server, err := apiserver.NewServer(buildCtx, logger, cfg)
	if err != nil {
		return err
	}

	server.Run()
	return nil
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version info as JSON (commit hash and build times)",
		RunE:  runVersion,
	}
}

func runVersion(_ *cobra.Command, _ []string) error {
	return version.WriteJSON()
}
