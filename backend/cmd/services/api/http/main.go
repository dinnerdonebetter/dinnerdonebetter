package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/server/http/build"

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

func getConfig(ctx context.Context) (*config.InstanceConfig, error) {
	var cfg *config.InstanceConfig
	if os.Getenv(config.RunningInGCPEnvVarKey) != "" {
		c, err := config.GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx)
		if err != nil {
			return nil, fmt.Errorf("fetching config from GCP: %w", err)
		}

		cfg = c
	} else if configFilepath := os.Getenv(config.FilePathEnvVarKey); configFilepath != "" {
		configBytes, err := os.ReadFile(configFilepath)
		if err != nil {
			return nil, fmt.Errorf("reading local config file: %w", err)
		}

		if err = json.NewDecoder(bytes.NewReader(configBytes)).Decode(&cfg); err != nil || cfg == nil {
			return nil, fmt.Errorf("decoding config file contents: %w", err)
		}
	} else {
		return nil, errors.New("no config provided")
	}

	if err := cfg.ValidateWithContext(ctx, true); err != nil {
		return nil, fmt.Errorf("validating config: %w", err)
	}

	return cfg, nil
}

func main() {
	rootCtx := context.Background()

	cfg, err := getConfig(rootCtx)
	if err != nil {
		log.Fatal(err)
	}

	// only allow initialization to take so long.
	buildCtx, cancel := context.WithTimeout(rootCtx, cfg.Server.StartupDeadline)

	log.Println("building server")
	// build our server struct.
	srv, err := build.Build(buildCtx, cfg)
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
