package datachangemessagehandler

import (
	"context"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	authkeys "github.com/dinnerdonebetter/backend/internal/domain/auth/keys"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	identitykeys "github.com/dinnerdonebetter/backend/internal/domain/identity/keys"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	eatingemails "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/emails"
	mealplanningkeys "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/keys"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	coreemails "github.com/dinnerdonebetter/backend/internal/services/identity/emails"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

func (a *AsyncDataChangeMessageHandler) DataChangesEventHandler(topicName string) func(ctx context.Context, rawMsg []byte) error {
	return func(ctx context.Context, rawMsg []byte) error {
		ctx, span := a.tracer.StartSpan(ctx)
		defer span.End()

		start := time.Now()
		status := statusSuccess
		eventType := unknownValue

		defer func() {
			a.dataChangesExecutionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()),
				metric.WithAttributes(
					attribute.String("status", status),
					attribute.String("event_type", eventType),
				))
			a.recordMessagesProcessed(ctx, topicDataChanges, status)
		}()

		var dataChangeMessage audit.DataChangeMessage
		if err := a.decoder.DecodeBytes(ctx, rawMsg, &dataChangeMessage); err != nil {
			a.messageDecodeErrorsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("topic", topicDataChanges)))
			status = statusFailure
			return fmt.Errorf("decoding message body: %w", err)
		}

		eventType = dataChangeMessage.EventType

		if err := a.handleDataChangeMessage(ctx, &dataChangeMessage, topicName); err != nil {
			a.handlerErrorsCounter.Add(ctx, 1, metric.WithAttributes(attribute.String("topic", topicDataChanges)))
			status = statusFailure
			return observability.PrepareAndLogError(err, a.logger, span, "handling data change message")
		}

		return nil
	}
}

func (a *AsyncDataChangeMessageHandler) handleDataChangeMessage(
	ctx context.Context,
	changeMessage *audit.DataChangeMessage,
	topicName string,
) error {
	ctx, span := a.tracer.StartSpan(ctx)

	if changeMessage == nil {
		return errRequiredDataIsNil
	}

	logger := a.logger.WithValue("event_type", changeMessage.EventType)

	// Non-empty TestID triggers queue test behavior (acknowledge and skip business logic)
	testID := changeMessage.TestID
	if testID == "" && changeMessage.Context != nil {
		if v, ok := changeMessage.Context["test_id"].(string); ok {
			testID = v
		}
	}
	if testID != "" {
		return a.handleQueueTestMessage(ctx, logger, span, testID, topicName)
	}

	if changeMessage.UserID != "" && changeMessage.EventType != "" {
		if err := a.analyticsEventReporter.EventOccurred(ctx, changeMessage.EventType, changeMessage.UserID, changeMessage.Context); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "notifying customer data platform")
		}
	}

	var wg sync.WaitGroup

	wg.Go(func() {
		a.nonWebhookEventTypesHat.RLock()
		eventTypeIsValid := slices.Contains(a.nonWebhookEventTypes, changeMessage.EventType)
		a.nonWebhookEventTypesHat.RUnlock()

		if changeMessage.AccountID != "" && !eventTypeIsValid {
			relevantWebhooks, err := a.webhookRepo.GetWebhooksForAccountAndEvent(ctx, changeMessage.AccountID, changeMessage.EventType)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "getting webhooks")
				return
			}

			for _, webhook := range relevantWebhooks {
				if err = a.webhookExecutionRequestPublisher.Publish(ctx, &webhooks.WebhookExecutionRequest{
					WebhookID: webhook.ID,
					AccountID: changeMessage.AccountID,
					Payload:   changeMessage,
				}); err != nil {
					observability.AcknowledgeError(err, logger, span, "publishing webhook execution request")
				}
			}
		}
	})

	wg.Go(func() {
		if err := a.handleOutboundNotifications(ctx, changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer(s)")
		}
	})

	wg.Go(func() {
		if err := a.handleSearchIndexUpdates(ctx, changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "updating search index)")
		}
	})

	wg.Wait()

	return nil
}

func (a *AsyncDataChangeMessageHandler) handleSearchIndexUpdates(
	ctx context.Context,
	changeMessage *audit.DataChangeMessage,
) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	logger := a.logger.WithValue("event_type", changeMessage.EventType)

	switch changeMessage.EventType {
	case identity.UserSignedUpServiceEventType,
		identity.UserArchivedServiceEventType,
		identity.EmailAddressChangedEventType,
		identity.UsernameChangedEventType,
		identity.UserDetailsChangedEventType,
		identity.UserEmailAddressVerifiedEventType:
		if changeMessage.UserID == "" {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for User")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.UserID,
			IndexType: coreindexing.IndexTypeUsers,
			Delete:    changeMessage.EventType == identity.UserArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.RecipeCreatedServiceEventType,
		mealplanning.RecipeUpdatedServiceEventType,
		mealplanning.RecipeArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.RecipeIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for Recipe")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.RecipeArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.MealCreatedServiceEventType,
		mealplanning.MealUpdatedServiceEventType,
		mealplanning.MealArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.MealIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for Meal")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeMeals,
			Delete:    changeMessage.EventType == mealplanning.MealArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientCreatedServiceEventType,
		mealplanning.ValidIngredientUpdatedServiceEventType,
		mealplanning.ValidIngredientArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidIngredientIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredient")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeValidIngredients,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidInstrumentCreatedServiceEventType,
		mealplanning.ValidInstrumentUpdatedServiceEventType,
		mealplanning.ValidInstrumentArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidInstrumentIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidInstrument")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeValidInstruments,
			Delete:    changeMessage.EventType == mealplanning.ValidInstrumentArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidMeasurementUnitCreatedServiceEventType,
		mealplanning.ValidMeasurementUnitUpdatedServiceEventType,
		mealplanning.ValidMeasurementUnitArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidMeasurementUnitIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidMeasurementUnit")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeValidMeasurementUnits,
			Delete:    changeMessage.EventType == mealplanning.ValidMeasurementUnitArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidPreparationCreatedServiceEventType,
		mealplanning.ValidPreparationUpdatedServiceEventType,
		mealplanning.ValidPreparationArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidPreparationIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparation")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeValidPreparations,
			Delete:    changeMessage.EventType == mealplanning.ValidPreparationArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientStateCreatedServiceEventType,
		mealplanning.ValidIngredientStateUpdatedServiceEventType,
		mealplanning.ValidIngredientStateArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidIngredientStateIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientState")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeValidIngredientStates,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientStateArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientMeasurementUnitCreatedServiceEventType,
		mealplanning.ValidIngredientMeasurementUnitUpdatedServiceEventType,
		mealplanning.ValidIngredientMeasurementUnitArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidIngredientMeasurementUnitIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientMeasurementUnit")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientMeasurementUnitArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidPreparationInstrumentCreatedServiceEventType,
		mealplanning.ValidPreparationInstrumentUpdatedServiceEventType,
		mealplanning.ValidPreparationInstrumentArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidPreparationInstrumentIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparationInstrument")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidPreparationInstrumentArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientPreparationCreatedServiceEventType,
		mealplanning.ValidIngredientPreparationUpdatedServiceEventType,
		mealplanning.ValidIngredientPreparationArchivedServiceEventType:
		rowID := rowIDFromEventContext(changeMessage, mealplanningkeys.ValidIngredientPreparationIDKey)
		if rowID == "" {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientPreparation")
		}
		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     rowID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientPreparationArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	default:
		logger.Debug("event type not handled for search indexing")
		return nil
	}
}

func (a *AsyncDataChangeMessageHandler) handleOutboundNotifications(
	ctx context.Context,
	changeMessage *audit.DataChangeMessage,
) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	if changeMessage == nil {
		return errors.New("nil data change message")
	}

	envCfg := email.GetConfigForEnvironment(os.Getenv("DINNER_DONE_BETTER_SERVICE_ENVIRONMENT"))
	if envCfg == nil {
		return observability.PrepareAndLogError(email.ErrMissingEnvCfg, a.logger, span, "getting environment queuesConfig")
	}

	logger := a.logger.WithValue("event_type", changeMessage.EventType)

	// Events from background jobs (e.g. meal plan grocery list initializer) may have no UserID; skip notifications.
	if changeMessage.UserID == "" {
		return nil
	}

	user, err := a.identityRepo.GetUser(ctx, changeMessage.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting user")
	}

	var (
		emailType             string
		msg                   *email.OutboundEmailMessage
		outboundEmailMessages []*email.OutboundEmailMessage
	)

	switch changeMessage.EventType {
	case identity.UserSignedUpServiceEventType:
		emailType = "user signup"
		if err = a.analyticsEventReporter.AddUser(ctx, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}

		emailVerificationToken := stringFromEventContext(changeMessage, identitykeys.UserEmailVerificationTokenKey)
		if emailVerificationToken == "" {
			return observability.PrepareError(fmt.Errorf("email verification token required"), span, "building address verification email")
		}

		msg, err = coreemails.BuildVerifyEmailAddressEmail(user, emailVerificationToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.UserEmailAddressVerificationEmailRequestedEventType:
		emailType = "email address verification"
		emailVerificationToken := stringFromEventContext(changeMessage, identitykeys.UserEmailVerificationTokenKey)
		if emailVerificationToken == "" {
			return observability.PrepareError(fmt.Errorf("email verification token required"), span, "building address verification email")
		}

		msg, err = coreemails.BuildVerifyEmailAddressEmail(user, emailVerificationToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		outboundEmailMessages = append(outboundEmailMessages, msg)

	case mealplanning.MealPlanCreatedServiceEventType:
		emailType = "meal plan created"
		mealPlanID, ok := changeMessage.Context[mealplanningkeys.MealPlanIDKey].(string)
		if !ok {
			mealPlanID = ""
		}
		if mealPlanID == "" || changeMessage.AccountID == "" {
			return observability.PrepareError(fmt.Errorf("meal plan created event requires meal_plan.id and accountID in context"), span, "publishing meal plan created email")
		}

		var mealPlan *mealplanning.MealPlan
		mealPlan, err = a.mealPlanRepo.GetMealPlan(ctx, mealPlanID, changeMessage.AccountID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting meal plan for created email")
		}
		if mealPlan == nil {
			return observability.PrepareError(fmt.Errorf("meal plan is nil"), span, "publishing meal plan created email")
		}

		var account *identity.Account
		account, err = a.identityRepo.GetAccount(ctx, mealPlan.BelongsToAccount)
		if err != nil {
			return observability.PrepareError(err, span, "getting account")
		}

		for _, member := range account.Members {
			if member.BelongsToUser.EmailAddressVerifiedAt != nil {
				msg, err = eatingemails.BuildMealPlanCreatedEmail(member.BelongsToUser, mealPlan, envCfg)
				if err != nil {
					return observability.PrepareAndLogError(err, logger, span, "building meal plan created email")
				}

				outboundEmailMessages = append(outboundEmailMessages, msg)
			}
		}
	case identity.PasswordResetTokenCreatedEventType:
		emailType = "password reset request"
		tokenID := stringFromEventContext(changeMessage, authkeys.PasswordResetTokenIDKey)
		if tokenID == "" {
			return observability.PrepareError(fmt.Errorf("password reset token created event requires password_reset_token.id in context"), span, "building password reset email")
		}

		var prt *auth.PasswordResetToken
		prt, err = a.passwordResetTokenDataManager.GetPasswordResetTokenByID(ctx, tokenID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting password reset token")
		}
		if prt == nil {
			return observability.PrepareError(fmt.Errorf("password reset token not found"), span, "building password reset email")
		}

		msg, err = coreemails.BuildGeneratedPasswordResetTokenEmail(user, prt, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token created email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.UsernameReminderRequestedEventType:
		emailType = "username reminder"
		msg, err = coreemails.BuildUsernameReminderEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building username reminder email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.PasswordResetTokenRedeemedEventType:
		emailType = "password reset token redeemed"
		msg, err = coreemails.BuildPasswordResetTokenRedeemedEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token redemption email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.PasswordChangedEventType:
		emailType = "password reset token redeemed"
		msg, err = coreemails.BuildPasswordChangedEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.AccountInvitationCreatedServiceEventType:
		emailType = "account invitation created"
		invitationID := stringFromEventContext(changeMessage, identitykeys.AccountInvitationIDKey)
		destinationAccountID, ok := changeMessage.Context[identitykeys.DestinationAccountIDKey].(string)
		if !ok {
			destinationAccountID = ""
		}
		if invitationID == "" || destinationAccountID == "" {
			return observability.PrepareError(fmt.Errorf("account invitation created event requires %s and %s in context", identitykeys.AccountInvitationIDKey, identitykeys.DestinationAccountIDKey), span, "building invite member email")
		}

		var accountInvite *identity.AccountInvitation
		accountInvite, err = a.identityRepo.GetAccountInvitationByAccountAndID(ctx, destinationAccountID, invitationID)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting account invitation")
		}
		if accountInvite == nil {
			return observability.PrepareError(fmt.Errorf("account invitation not found"), span, "building invite member email")
		}

		msg, err = coreemails.BuildInviteMemberEmail(user, accountInvite, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building email message")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.AccountInvitationAcceptedServiceEventType:
		destinationAccountID, ok := changeMessage.Context[identitykeys.DestinationAccountIDKey].(string)
		if !ok || destinationAccountID == "" {
			logger.Debug(fmt.Sprintf("account invitation accepted: missing %s in context, skipping mobile notification", identitykeys.DestinationAccountIDKey))
			return nil
		}
		acceptedUserID := changeMessage.UserID

		var usersResult *filtering.QueryFilteredResult[identity.User]
		usersResult, err = a.identityRepo.GetUsersForAccount(ctx, destinationAccountID, filtering.DefaultQueryFilter())
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "getting users for account")
		}

		var recipientUserIDs []string
		for _, u := range usersResult.Data {
			if u != nil && u.ID != "" && u.ID != acceptedUserID {
				recipientUserIDs = append(recipientUserIDs, u.ID)
			}
		}
		if len(recipientUserIDs) == 0 {
			return nil
		}

		displayName := "Someone"
		if user != nil {
			if user.FirstName != "" || user.LastName != "" {
				displayName = strings.TrimSpace(user.FirstName + " " + user.LastName)
			} else if user.Username != "" {
				displayName = user.Username
			}
		}

		mobileReq := &notifications.MobileNotificationRequest{
			RequestType:      notifications.MobileNotificationRequestTypeHouseholdInvitationAccepted,
			RecipientUserIDs: recipientUserIDs,
			Title:            "Someone joined your household",
			Body:             fmt.Sprintf("%s joined your household", displayName),
			Context: map[string]string{
				notifications.ExcludedUserIDContextKey: acceptedUserID,
			},
		}
		if err = a.mobileNotificationsPublisher.Publish(ctx, mobileReq); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing household invitation accepted mobile notification")
		}
	}

	if len(outboundEmailMessages) == 0 {
		logger.WithValue("email_type", emailType).WithValue("outbound_emails_to_send", len(outboundEmailMessages)).Info("publishing email requests")
	}

	for _, oem := range outboundEmailMessages {
		if err = a.outboundEmailsPublisher.Publish(ctx, oem); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing %s request email", emailType)
		}
	}

	return nil
}

// stringFromEventContext returns a string value from the data change message context.
// The value may be a string or []byte depending on message serialization.
func stringFromEventContext(changeMessage *audit.DataChangeMessage, key string) string {
	if changeMessage == nil || changeMessage.Context == nil {
		return ""
	}

	v, ok := changeMessage.Context[key]
	if !ok {
		return ""
	}

	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	default:
		return ""
	}
}

// rowIDFromEventContext returns the row ID string from the data change message context.
// Producers publish only the ID (e.g. RecipeIDKey -> recipe.ID), not the full entity.
func rowIDFromEventContext(changeMessage *audit.DataChangeMessage, idKey string) string {
	return stringFromEventContext(changeMessage, idKey)
}
