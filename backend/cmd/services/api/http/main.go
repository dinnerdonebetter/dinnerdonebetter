package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/build/api"
	"github.com/dinnerdonebetter/backend/internal/config"

	"github.com/KimMachineGun/automemlimit/memlimit"
	_ "go.uber.org/automaxprocs"
)

func init() {
	if _, err := memlimit.SetGoMemLimitWithOpts(
		memlimit.WithRatio(0.9),
		memlimit.WithProvider(
			memlimit.ApplyFallback(
				memlimit.FromCgroup,
				memlimit.FromSystem,
			),
		),
		memlimit.WithLogger(nil),
	); err != nil {
		slog.Error("failed to set go mem limit provider")
	}
}

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
	srv, err := api.Build(buildCtx, cfg)
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

	// Run server
	go srv.Serve()

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
	if err = srv.Shutdown(cancelCtx); err != nil {
		logger.Error("shutting down server", err)
	}
}
