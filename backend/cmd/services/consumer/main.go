package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/routing/chi"
)

const (
	defaultConsumerServerConfigurationFilepath = "deploy/environments/localdev/config_files/consumer_webapp_config.json"
)

func main() {
	ctx := context.Background()

	configFilepath := os.Getenv("CONFIGURATION_FILEPATH")
	if configFilepath == "" {
		configFilepath = defaultConsumerServerConfigurationFilepath
	}

	cfg, err := config.LoadConfigFromPath[config.ConsumerWebappConfig](ctx, configFilepath)
	if err != nil {
		log.Fatal(err)
	}

	logger, tracerProvider, _, err := cfg.Observability.ProvideThreePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		log.Fatal(err)
	}

	fs, err := NewConsumerFrontendServer(
		ctx,
		logger,
		tracerProvider,
		encoding.ProvideServerEncoderDecoder(logger, tracerProvider, encoding.ContentTypeJSON),
		chi.NewRouteParamManager(),
		cfg,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("serving consumer app")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	go fs.server.Serve()

	<-signalChan

	go func() {
		<-signalChan
	}()

	cancelCtx, cancelShutdown := context.WithTimeout(ctx, 10*time.Second)
	defer cancelShutdown()

	if err = fs.server.Shutdown(cancelCtx); err != nil {
		log.Println("shutting down server", err)
	}
}
