package datachangemessagehandler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"sync"
	"time"

	"github.com/dinnerdonebetter/backend/internal/domain/audit"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/webhooks"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/encoding"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	textsearch "github.com/dinnerdonebetter/backend/internal/platform/search/text"
	coreemails "github.com/dinnerdonebetter/backend/internal/services/identity/emails"
	coreindexing "github.com/dinnerdonebetter/backend/internal/services/identity/indexing"
	eatingemails "github.com/dinnerdonebetter/backend/internal/services/mealplanning/emails"
	eatingindexing "github.com/dinnerdonebetter/backend/internal/services/mealplanning/indexing"
)

func (a *AsyncDataChangeMessageHandler) DataChangesEventHandler(ctx context.Context, rawMsg []byte) error {
	ctx, span := a.tracer.StartSpan(ctx)
	defer span.End()

	start := time.Now()

	var dataChangeMessage audit.DataChangeMessage
	if err := json.NewDecoder(bytes.NewReader(rawMsg)).Decode(&dataChangeMessage); err != nil {
		return fmt.Errorf("decoding JSON body: %w", err)
	}

	if err := a.handleDataChangeMessage(ctx, &dataChangeMessage); err != nil {
		return observability.PrepareAndLogError(err, a.logger, span, "handling data change message")
	}

	a.dataChangesExecutionTimeHistogram.Record(ctx, float64(time.Since(start).Milliseconds()))

	return nil
}

func (a *AsyncDataChangeMessageHandler) handleDataChangeMessage(
	ctx context.Context,
	changeMessage *audit.DataChangeMessage,
) error {
	ctx, span := a.tracer.StartSpan(ctx)

	if changeMessage == nil {
		return errRequiredDataIsNil
	}

	logger := a.logger.WithValue("event_type", changeMessage.EventType)

	if changeMessage.UserID != "" && changeMessage.EventType != "" {
		if err := a.analyticsEventReporter.EventOccurred(ctx, changeMessage.EventType, changeMessage.UserID, changeMessage.Context); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "notifying customer data platform")
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
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
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := a.handleOutboundNotifications(ctx, changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer(s)")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := a.handleSearchIndexUpdates(ctx, changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "updating search index)")
		}
		wg.Done()
	}()

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
		recipe, parseError := parseValueFromEventContext[mealplanning.Recipe](ctx, changeMessage, a.decoder, keys.RecipeKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if recipe == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for Recipe")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     recipe.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.RecipeArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.MealCreatedServiceEventType,
		mealplanning.MealUpdatedServiceEventType,
		mealplanning.MealArchivedServiceEventType:
		meal, parseError := parseValueFromEventContext[mealplanning.Meal](ctx, changeMessage, a.decoder, keys.MealKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if meal == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for Meal")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     meal.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.MealArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientCreatedServiceEventType,
		mealplanning.ValidIngredientUpdatedServiceEventType,
		mealplanning.ValidIngredientArchivedServiceEventType:
		validIngredient, parseError := parseValueFromEventContext[mealplanning.ValidIngredient](ctx, changeMessage, a.decoder, keys.ValidIngredientKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validIngredient == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredient")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validIngredient.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidInstrumentCreatedServiceEventType,
		mealplanning.ValidInstrumentUpdatedServiceEventType,
		mealplanning.ValidInstrumentArchivedServiceEventType:
		validInstrument, parseError := parseValueFromEventContext[mealplanning.ValidInstrument](ctx, changeMessage, a.decoder, keys.ValidInstrumentKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validInstrument == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidInstrument")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validInstrument.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidInstrumentArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidMeasurementUnitCreatedServiceEventType,
		mealplanning.ValidMeasurementUnitUpdatedServiceEventType,
		mealplanning.ValidMeasurementUnitArchivedServiceEventType:
		validMeasurementUnit, parseError := parseValueFromEventContext[mealplanning.ValidMeasurementUnit](ctx, changeMessage, a.decoder, keys.ValidMeasurementUnitKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validMeasurementUnit == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidMeasurementUnit")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validMeasurementUnit.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidMeasurementUnitArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidPreparationCreatedServiceEventType,
		mealplanning.ValidPreparationUpdatedServiceEventType,
		mealplanning.ValidPreparationArchivedServiceEventType:
		validPreparation, parseError := parseValueFromEventContext[mealplanning.ValidPreparation](ctx, changeMessage, a.decoder, keys.ValidPreparationKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validPreparation == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparation")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validPreparation.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidPreparationArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientStateCreatedServiceEventType,
		mealplanning.ValidIngredientStateUpdatedServiceEventType,
		mealplanning.ValidIngredientStateArchivedServiceEventType:
		validIngredientState, parseError := parseValueFromEventContext[mealplanning.ValidIngredientState](ctx, changeMessage, a.decoder, keys.ValidIngredientStateKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validIngredientState == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientState")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validIngredientState.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientStateArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientMeasurementUnitCreatedServiceEventType,
		mealplanning.ValidIngredientMeasurementUnitUpdatedServiceEventType,
		mealplanning.ValidIngredientMeasurementUnitArchivedServiceEventType:
		validIngredientMeasurementUnit, parseError := parseValueFromEventContext[mealplanning.ValidIngredientMeasurementUnit](ctx, changeMessage, a.decoder, keys.ValidIngredientMeasurementUnitKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validIngredientMeasurementUnit == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientMeasurementUnit")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validIngredientMeasurementUnit.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidIngredientMeasurementUnitArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidPreparationInstrumentCreatedServiceEventType,
		mealplanning.ValidPreparationInstrumentUpdatedServiceEventType,
		mealplanning.ValidPreparationInstrumentArchivedServiceEventType:
		validPreparationInstrument, parseError := parseValueFromEventContext[mealplanning.ValidPreparationInstrument](ctx, changeMessage, a.decoder, keys.ValidPreparationInstrumentKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validPreparationInstrument == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparationInstrument")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validPreparationInstrument.ID,
			IndexType: eatingindexing.IndexTypeRecipes,
			Delete:    changeMessage.EventType == mealplanning.ValidPreparationInstrumentArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case mealplanning.ValidIngredientPreparationCreatedServiceEventType,
		mealplanning.ValidIngredientPreparationUpdatedServiceEventType,
		mealplanning.ValidIngredientPreparationArchivedServiceEventType:
		validIngredientPreparation, parseError := parseValueFromEventContext[mealplanning.ValidIngredientPreparation](ctx, changeMessage, a.decoder, keys.ValidIngredientPreparationKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if validIngredientPreparation == nil {
			return observability.PrepareAndLogError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientPreparation")
		}

		if err := a.searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     validIngredientPreparation.ID,
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

		evf, parseError := parseValueFromEventContext[string](ctx, changeMessage, a.decoder, keys.UserEmailVerificationTokenKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		msg, err = coreemails.BuildVerifyEmailAddressEmail(user, *evf, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		outboundEmailMessages = append(outboundEmailMessages, msg)

	case identity.UserEmailAddressVerificationEmailRequestedEventType:
		emailType = "email address verification"
		evf, parseError := parseValueFromEventContext[string](ctx, changeMessage, a.decoder, keys.UserEmailVerificationTokenKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		msg, err = coreemails.BuildVerifyEmailAddressEmail(user, *evf, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		outboundEmailMessages = append(outboundEmailMessages, msg)

	case mealplanning.MealPlanCreatedServiceEventType:
		emailType = "meal plan created"
		mealPlan, parseError := parseValueFromEventContext[mealplanning.MealPlan](ctx, changeMessage, a.decoder, keys.MealPlanKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
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
				msg, err = eatingemails.BuildMealPlanCreatedEmail(user, mealPlan, envCfg)
				if err != nil {
					return observability.PrepareAndLogError(err, logger, span, "building meal plan created email")
				}

				outboundEmailMessages = append(outboundEmailMessages, msg)
			}
		}
	case identity.PasswordResetTokenCreatedEventType:
		emailType = "password reset request"
		prt, parseError := parseValueFromEventContext[auth.PasswordResetToken](ctx, changeMessage, a.decoder, keys.PasswordResetTokenKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if prt == nil {
			return observability.PrepareError(fmt.Errorf("password reset token is nil"), span, "publishing password reset token email")
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
		accountInvite, parseError := parseValueFromEventContext[identity.AccountInvitation](ctx, changeMessage, a.decoder, keys.AccountInvitationKey)
		if parseError != nil {
			return observability.PrepareAndLogError(parseError, logger, span, "parsing email verification token")
		}

		if accountInvite == nil {
			return observability.PrepareError(fmt.Errorf("account invitation is nil"), span, "publishing password reset token redemption email")
		}

		msg, err = coreemails.BuildInviteMemberEmail(user, accountInvite, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building email message")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)
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

func parseValueFromEventContext[T any](ctx context.Context, changeMessage *audit.DataChangeMessage, decoder encoding.ServerEncoderDecoder, key string) (*T, error) {
	var x T
	if z, ok := changeMessage.Context[key]; ok {
		switch y := z.(type) {
		case string:
			if err := decoder.DecodeBytes(ctx, []byte(y), &z); err != nil {
				return nil, err
			}
		case []byte:
			z = string(y)
		}
	}

	return &x, nil
}
