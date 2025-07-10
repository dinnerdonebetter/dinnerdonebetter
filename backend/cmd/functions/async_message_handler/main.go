package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	dcmh "github.com/dinnerdonebetter/backend/internal/build/functions/data_change_message_handler"
	"github.com/dinnerdonebetter/backend/internal/config"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/oauth"

	_ "go.uber.org/automaxprocs"
)

var (
	nonWebhookEventTypes = []string{
		identity.UserSignedUpServiceEventType,
		identity.UserArchivedServiceEventType,
		identity.TwoFactorSecretVerifiedServiceEventType,
		identity.TwoFactorDeactivatedServiceEventType,
		identity.TwoFactorSecretChangedServiceEventType,
		identity.PasswordResetTokenCreatedEventType,
		identity.PasswordResetTokenRedeemedEventType,
		identity.PasswordChangedEventType,
		identity.EmailAddressChangedEventType,
		identity.UsernameChangedEventType,
		identity.UserDetailsChangedEventType,
		identity.UsernameReminderRequestedEventType,
		identity.UserLoggedInServiceEventType,
		identity.UserLoggedOutServiceEventType,
		identity.UserChangedActiveAccountServiceEventType,
		identity.UserEmailAddressVerifiedEventType,
		identity.UserEmailAddressVerificationEmailRequestedEventType,
		identity.AccountMemberRemovedServiceEventType,
		identity.AccountMembershipPermissionsUpdatedServiceEventType,
		identity.AccountOwnershipTransferredServiceEventType,
		oauth.OAuth2ClientCreatedServiceEventType,
		oauth.OAuth2ClientArchivedServiceEventType,
	}
)

func main() {
	if config.ShouldCeaseOperation() {
		slog.Info("CEASE_OPERATION is set to true, exiting")
		os.Exit(0)
	}

	cfg, err := config.LoadConfigFromEnvironment[config.AsyncMessageHandlerConfig]()
	if err != nil {
		log.Fatalf("error getting config: %v", err)
	}
	cfg.Database.RunMigrations = false

	ctx := context.Background()

	dataChangeMessageHandler, err := dcmh.Build(ctx, cfg)
	if err != nil {
		log.Fatalf("error building data_change_message_handler: %v", err)
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

	dataChangeMessageHandler.SetNonWebhookEventTypes(nonWebhookEventTypes)

	if err = doTheThing(
		ctx,
		logger,
		tracerProvider,
		metricsProvider,
		cfg,
		consumerProvider,
		analyticsEventReporter,
		identityRepo,
		webhookRepo,
		mealplanningRepo,
		outboundEmailsPublisher,
		searchDataIndexPublisher,
		webhookExecutionRequestPublisher,
		emailer,
		uploadManager,
		dataChangesExecutionTimeHistogram,
		outboundEmailsExecutionTimeHistogram,
		searchIndexRequestsExecutionTimeHistogram,
		userDataAggregationExecutionTimeHistogram,
		webhookExecutionTimestampHistogram,
		decoder,
		stopChan,
		errorsChan,
	); err != nil {
		closeFunc()
		log.Fatal(err)
	}

	// os.Interrupt
	<-signalChan

	go func() {
		// os.Kill
		<-signalChan
		stopChan <- true
		closeFunc()
	}()
}
