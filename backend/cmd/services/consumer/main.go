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

	pillars, err := cfg.Observability.ProvidePillars(ctx)
	if err != nil {
		log.Fatal(err)
	}

	if err = cfg.ValidateWithContext(ctx); err != nil {
		log.Fatal(err)
	}

	if err = pillars.Profiler.Start(ctx); err != nil {
		log.Fatal(err)
	}

	fs, err := NewConsumerFrontendServer(
		ctx,
		pillars.Logger,
		pillars.TracerProvider,
		encoding.ProvideServerEncoderDecoder(pillars.Logger, pillars.TracerProvider, encoding.ContentTypeJSON),
		chi.NewRouteParamManager(),
		cfg,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = pillars.Profiler.Shutdown(shutdownCtx); err != nil {
			pillars.Logger.Error("failed to shutdown profiler", err)
		}
	}()

	pillars.Logger.Info("serving consumer app")

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
