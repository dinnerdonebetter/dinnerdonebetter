package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KimMachineGun/automemlimit/memlimit"
	_ "go.uber.org/automaxprocs"

	"github.com/dinnerdonebetter/backend/internal/build/api"
	"github.com/dinnerdonebetter/backend/internal/config"
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

	cfg, err := config.FetchForApplication(rootCtx, config.GetAPIServiceConfigFromGoogleCloudRunEnvironment)
	if err != nil {
		log.Fatal(err)
	}

	// only allow initialization to take so long.
	buildCtx, cancel := context.WithTimeout(rootCtx, cfg.Server.StartupDeadline)

	configBytes, err := os.ReadFile(os.Getenv(config.ConfigurationFilePathEnvVarKey))
	if err != nil {
		panic(err)
	}

	loadedCfg := &config.APIServiceConfig{}
	json.Unmarshal(configBytes, &loadedCfg)

	logger := cfg.Observability.Logging.ProvideLogger()
	logger.Info("building server")

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

	log.Println("serving")
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

	log.Println("shutting down")
	// Gracefully shutdown the server by waiting on existing requests (except websockets).
	if err = srv.Shutdown(cancelCtx); err != nil {
		panic(err)
	}
}
