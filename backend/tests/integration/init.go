package integration

import (
	"context"
	"fmt"
	grpcapi "github.com/dinnerdonebetter/backend/internal/build/services/api/grpc"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"log"
	"os"

	"github.com/dinnerdonebetter/backend/internal/config"
)

const (
	apiConfigurationFilepath = "../../deploy/environments/testing/config_files/integration-tests-config.json"
)

func deriveServerConfig() (*config.APIServiceConfig, error) {
	wd, _ := os.Getwd()
	fmt.Println(wd)

	content, err := os.ReadFile(apiConfigurationFilepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read api configuration file: %w", err)
	}

	decoder := encoding.ProvideServerEncoderDecoder(logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), encoding.ContentTypeJSON)

	var x *config.APIServiceConfig
	if err = decoder.DecodeBytes(context.Background(), content, &x); err != nil {
		return nil, fmt.Errorf("failed to decode api configuration file: %w", err)
	}

	return x, nil
}

func init() {
	ctx := context.Background()

	cfg, err := deriveServerConfig()
	if err != nil {
		log.Fatal(err)
	}

	_, _, dbCfg, err := pgtesting.BuildDatabaseContainer(ctx, "integration_testing")
	if err != nil {
		log.Fatal(err)
	}
	cfg.Database = *dbCfg

	grpcServer, err := grpcapi.Build(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	go grpcServer.Serve()
}
