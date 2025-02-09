package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpcapi "github.com/dinnerdonebetter/backend/internal/build/services/api/grpc"
	httpapi "github.com/dinnerdonebetter/backend/internal/build/services/api/http"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func main() {
	rootCtx := context.Background()

	cfg, err := config.LoadConfigFromEnvironment[config.APIServiceConfig]()
	if err != nil {
		log.Fatal(err)
	}

	// only allow initialization to take so long.
	buildCtx, cancel := context.WithTimeout(rootCtx, cfg.HTTPServer.StartupDeadline)

	logger, err := cfg.Observability.Logging.ProvideLogger(rootCtx)
	if err != nil {
		log.Fatalf("could not create logger: %v", err)
	}

	// build our server struct.
	httpServer, err := httpapi.Build(buildCtx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer, err := grpcapi.Build(buildCtx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	cancel()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	// Run servers
	go httpServer.Serve()
	go grpcServer.Serve()

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
	}()

	cancelCtx, cancelShutdown := context.WithTimeout(rootCtx, 10*time.Second)
	defer cancelShutdown()

	logger.Info("shutting down")

	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err = httpServer.Shutdown(cancelCtx); err != nil {
		logger.Error("shutting down server", err)
	}
}
