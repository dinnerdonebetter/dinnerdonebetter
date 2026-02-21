package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	outboundemailhandler "github.com/dinnerdonebetter/backend/internal/build/functions/outbound_email_handler"
	"github.com/dinnerdonebetter/backend/internal/config"

	_ "go.uber.org/automaxprocs"
)

func main() {
	config.ConditionallyCease()

	cfg, err := config.LoadConfigFromEnvironment[config.OutboundEmailHandlerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}

	ctx := context.Background()

	handler, err := outboundemailhandler.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("error building outbound_email_handler: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(
		signalChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	stopChan := make(chan bool)
	errorsChan := make(chan error)

	if err = handler.ConsumeMessages(ctx, stopChan, errorsChan); err != nil {
		log.Fatal(err)
	}

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
		stopChan <- true
	}()
}
