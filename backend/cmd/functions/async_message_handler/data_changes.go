package main

import (
	"context"
	"fmt"
	"os"
	"slices"
	"sync"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	"github.com/dinnerdonebetter/backend/internal/lib/messagequeue"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	textsearch "github.com/dinnerdonebetter/backend/internal/lib/search/text"
	coreemails "github.com/dinnerdonebetter/backend/internal/services/core/emails"
	eatingemails "github.com/dinnerdonebetter/backend/internal/services/eating/emails"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

func handleDataChangeMessage(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	analyticsEventReporter analytics.EventReporter,
	webhookExecutionRequestPublisher,
	outboundEmailsPublisher,
	searchDataIndexPublisher messagequeue.Publisher,
	changeMessage *types.DataChangeMessage,
) {
	ctx, span := tracer.StartSpan(ctx)

	logger = logger.WithValue("event_type", changeMessage.EventType)

	if changeMessage.UserID != "" && changeMessage.EventType != "" {
		if err := analyticsEventReporter.EventOccurred(ctx, changeMessage.EventType, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		if changeMessage.HouseholdID != "" && !slices.Contains(nonWebhookEventTypes, changeMessage.EventType) {
			relevantWebhooks, err := dataManager.GetWebhooksForHouseholdAndEvent(ctx, changeMessage.HouseholdID, changeMessage.EventType)
			if err != nil {
				observability.AcknowledgeError(err, logger, span, "getting webhooks")
			}

			for _, webhook := range relevantWebhooks {
				if err = webhookExecutionRequestPublisher.Publish(ctx, &types.WebhookExecutionRequest{
					WebhookID:   webhook.ID,
					HouseholdID: changeMessage.HouseholdID,
					Payload:     changeMessage,
				}); err != nil {
					observability.AcknowledgeError(err, logger, span, "publishing webhook execution request")
				}
			}
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := handleOutboundNotifications(ctx, logger, tracer, dataManager, outboundEmailsPublisher, analyticsEventReporter, changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer(s)")
		}
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		if err := handleSearchIndexUpdates(ctx, logger, tracer, searchDataIndexPublisher, changeMessage); err != nil {
			observability.AcknowledgeError(err, logger, span, "updating search index)")
		}
		wg.Done()
	}()

	wg.Wait()
}

func handleSearchIndexUpdates(
	ctx context.Context,
	l logging.Logger,
	tracer tracing.Tracer,
	searchDataIndexPublisher messagequeue.Publisher,
	changeMessage *types.DataChangeMessage,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	logger := l.WithValue("event_type", changeMessage.EventType)

	switch changeMessage.EventType {
	case types.UserSignedUpServiceEventType,
		types.UserArchivedServiceEventType,
		types.EmailAddressChangedEventType,
		types.UsernameChangedEventType,
		types.UserDetailsChangedEventType,
		types.UserEmailAddressVerifiedEventType:
		if changeMessage.UserID == "" {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for User")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.UserID,
			IndexType: textsearch.IndexTypeUsers,
			Delete:    changeMessage.EventType == types.UserArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.RecipeCreatedServiceEventType,
		types.RecipeUpdatedServiceEventType,
		types.RecipeArchivedServiceEventType:
		if changeMessage.Recipe == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for Recipe")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.Recipe.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.RecipeArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.MealCreatedServiceEventType,
		types.MealUpdatedServiceEventType,
		types.MealArchivedServiceEventType:
		if changeMessage.Meal == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for Meal")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.Meal.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.MealArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientCreatedServiceEventType,
		types.ValidIngredientUpdatedServiceEventType,
		types.ValidIngredientArchivedServiceEventType:
		if changeMessage.ValidIngredient == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredient")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidIngredient.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidInstrumentCreatedServiceEventType,
		types.ValidInstrumentUpdatedServiceEventType,
		types.ValidInstrumentArchivedServiceEventType:
		if changeMessage.ValidInstrument == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidInstrument")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidInstrument.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidInstrumentArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidMeasurementUnitCreatedServiceEventType,
		types.ValidMeasurementUnitUpdatedServiceEventType,
		types.ValidMeasurementUnitArchivedServiceEventType:
		if changeMessage.ValidMeasurementUnit == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidMeasurementUnit")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidMeasurementUnit.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidMeasurementUnitArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidPreparationCreatedServiceEventType,
		types.ValidPreparationUpdatedServiceEventType,
		types.ValidPreparationArchivedServiceEventType:
		if changeMessage.ValidPreparation == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparation")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidPreparation.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidPreparationArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientStateCreatedServiceEventType,
		types.ValidIngredientStateUpdatedServiceEventType,
		types.ValidIngredientStateArchivedServiceEventType:
		if changeMessage.ValidIngredientState == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientState")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidIngredientState.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientStateArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientMeasurementUnitCreatedServiceEventType,
		types.ValidIngredientMeasurementUnitUpdatedServiceEventType,
		types.ValidIngredientMeasurementUnitArchivedServiceEventType:
		if changeMessage.ValidIngredientMeasurementUnit == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientMeasurementUnit")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidIngredientMeasurementUnit.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientMeasurementUnitArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidPreparationInstrumentCreatedServiceEventType,
		types.ValidPreparationInstrumentUpdatedServiceEventType,
		types.ValidPreparationInstrumentArchivedServiceEventType:
		if changeMessage.ValidPreparationInstrument == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidPreparationInstrument")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidPreparationInstrument.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidPreparationInstrumentArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	case types.ValidIngredientPreparationCreatedServiceEventType,
		types.ValidIngredientPreparationUpdatedServiceEventType,
		types.ValidIngredientPreparationArchivedServiceEventType:
		if changeMessage.ValidIngredientPreparation == nil {
			observability.AcknowledgeError(errRequiredDataIsNil, logger, span, "updating search index for ValidIngredientPreparation")
		}

		if err := searchDataIndexPublisher.Publish(ctx, &textsearch.IndexRequest{
			RowID:     changeMessage.ValidIngredientPreparation.ID,
			IndexType: textsearch.IndexTypeRecipes,
			Delete:    changeMessage.EventType == types.ValidIngredientPreparationArchivedServiceEventType,
		}); err != nil {
			return observability.PrepareAndLogError(err, logger, span, "publishing search index update")
		}

		return nil
	default:
		logger.Debug("event type not handled for search indexing")
		return nil
	}
}

func handleOutboundNotifications(
	ctx context.Context,
	l logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	outboundEmailsPublisher messagequeue.Publisher,
	analyticsEventReporter analytics.EventReporter,
	changeMessage *types.DataChangeMessage,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	envCfg := email.GetConfigForEnvironment(os.Getenv("DINNER_DONE_BETTER_SERVICE_ENVIRONMENT"))
	if envCfg == nil {
		return observability.PrepareAndLogError(email.ErrMissingEnvCfg, l, span, "getting environment config")
	}

	logger := l.WithValue("event_type", changeMessage.EventType)

	user, err := dataManager.GetUser(ctx, changeMessage.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting user")
	}

	var (
		emailType             string
		msg                   *email.OutboundEmailMessage
		outboundEmailMessages []*email.OutboundEmailMessage
	)

	switch changeMessage.EventType {
	case types.UserSignedUpServiceEventType:
		emailType = "user signup"
		if err = analyticsEventReporter.AddUser(ctx, changeMessage.UserID, changeMessage.Context); err != nil {
			observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
		}

		msg, err = coreemails.BuildVerifyEmailAddressEmail(user, changeMessage.EmailVerificationToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		outboundEmailMessages = append(outboundEmailMessages, msg)

	case types.UserEmailAddressVerificationEmailRequestedEventType:
		emailType = "email address verification"
		msg, err = coreemails.BuildVerifyEmailAddressEmail(user, changeMessage.EmailVerificationToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		outboundEmailMessages = append(outboundEmailMessages, msg)

	case types.MealPlanCreatedServiceEventType:
		emailType = "meal plan created"
		mealPlan := changeMessage.MealPlan
		if mealPlan == nil {
			return observability.PrepareError(fmt.Errorf("meal plan is nil"), span, "publishing meal plan created email")
		}

		var household *types.Household
		household, err = dataManager.GetHousehold(ctx, mealPlan.BelongsToHousehold)
		if err != nil {
			return observability.PrepareError(err, span, "getting household")
		}

		for _, member := range household.Members {
			if member.BelongsToUser.EmailAddressVerifiedAt != nil {
				msg, err = eatingemails.BuildMealPlanCreatedEmail(user, mealPlan, envCfg)
				if err != nil {
					return observability.PrepareAndLogError(err, logger, span, "building meal plan created email")
				}

				outboundEmailMessages = append(outboundEmailMessages, msg)
			}
		}
	case types.PasswordResetTokenCreatedEventType:
		emailType = "password reset request"
		if changeMessage.PasswordResetToken == nil {
			return observability.PrepareError(fmt.Errorf("password reset token is nil"), span, "publishing password reset token email")
		}

		msg, err = coreemails.BuildGeneratedPasswordResetTokenEmail(user, changeMessage.PasswordResetToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token created email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case types.UsernameReminderRequestedEventType:
		emailType = "username reminder"
		msg, err = coreemails.BuildUsernameReminderEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building username reminder email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case types.PasswordResetTokenRedeemedEventType:
		emailType = "password reset token redeemed"
		msg, err = coreemails.BuildPasswordResetTokenRedeemedEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token redemption email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case types.PasswordChangedEventType:
		emailType = "password reset token redeemed"
		msg, err = coreemails.BuildPasswordChangedEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token email")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)

	case types.HouseholdInvitationCreatedServiceEventType:
		emailType = "household invitation created"
		if changeMessage.HouseholdInvitation == nil {
			return observability.PrepareError(fmt.Errorf("household invitation is nil"), span, "publishing password reset token redemption email")
		}

		msg, err = coreemails.BuildInviteMemberEmail(user, changeMessage.HouseholdInvitation, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building email message")
		}

		outboundEmailMessages = append(outboundEmailMessages, msg)
	}

	if len(outboundEmailMessages) == 0 {
		logger.WithValue("email_type", emailType).WithValue("outbound_emails_to_send", len(outboundEmailMessages)).Info("publishing email requests")
	}

	for _, oem := range outboundEmailMessages {
		if err = outboundEmailsPublisher.Publish(ctx, oem); err != nil {
			observability.AcknowledgeError(err, logger, span, "publishing %s request email", emailType)
		}
	}

	return nil
}
