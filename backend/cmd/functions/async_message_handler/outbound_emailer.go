package main

import (
	"context"
	"os"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/lib/analytics"
	"github.com/dinnerdonebetter/backend/internal/lib/email"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	coreemails "github.com/dinnerdonebetter/backend/internal/services/core/emails"
	eatingemails "github.com/dinnerdonebetter/backend/internal/services/eating/emails"
)

func handleEmailRequest(
	ctx context.Context,
	logger logging.Logger,
	tracer tracing.Tracer,
	dataManager database.DataManager,
	emailer email.Emailer,
	analyticsEventReporter analytics.EventReporter,
	emailDeliveryRequest *email.DeliveryRequest,
) error {
	ctx, span := tracer.StartSpan(ctx)
	defer span.End()

	logger = logger.WithValue("template", emailDeliveryRequest.Template)

	user, err := dataManager.GetUser(ctx, emailDeliveryRequest.UserID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "getting user")
	}

	envCfg := email.GetConfigForEnvironment(os.Getenv("DINNER_DONE_BETTER_SERVICE_ENVIRONMENT"))
	if envCfg == nil {
		return observability.PrepareAndLogError(email.ErrMissingEnvCfg, logger, nil, "getting environment config")
	}

	var (
		mail                   *email.OutboundEmailMessage
		shouldSkipIfUnverified = true
		emailType              string
	)

	switch emailDeliveryRequest.Template {
	case email.TemplateTypeInvite:
		if emailDeliveryRequest.Invitation == nil {
			return observability.PrepareAndLogError(err, logger, span, "missing household invitation")
		}

		mail, err = coreemails.BuildInviteMemberEmail(emailDeliveryRequest.Invitation, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building email message")
		}
		emailType = "invite"
	case email.TemplateTypeUsernameReminder:
		mail, err = coreemails.BuildUsernameReminderEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building username reminder email")
		}
		emailType = "username reminder"
	case email.TemplateTypePasswordResetTokenCreated:
		if emailDeliveryRequest.PasswordResetToken == nil {
			return observability.PrepareAndLogError(err, logger, span, "missing password reset token")
		}

		mail, err = coreemails.BuildGeneratedPasswordResetTokenEmail(user, emailDeliveryRequest.PasswordResetToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token created email")
		}
		emailType = "password reset token"
	case email.TemplateTypePasswordReset:
		mail, err = coreemails.BuildPasswordChangedEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token email")
		}
		emailType = "password reset token"
	case email.TemplateTypePasswordResetTokenRedeemed:
		mail, err = coreemails.BuildPasswordResetTokenRedeemedEmail(user, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building password reset token redemption email")
		}
		emailType = "password reset token redemption"
	case email.TemplateTypeMealPlanCreated:
		mail, err = eatingemails.BuildMealPlanCreatedEmail(user, emailDeliveryRequest.MealPlan, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building meal plan created email")
		}
		emailType = "meal plan created"
	case email.TemplateTypeVerifyEmailAddress:
		shouldSkipIfUnverified = false
		mail, err = coreemails.BuildVerifyEmailAddressEmail(user, emailDeliveryRequest.EmailVerificationToken, envCfg)
		if err != nil {
			return observability.PrepareAndLogError(err, logger, span, "building address verification email")
		}
		emailType = "email address verification"
	}

	if shouldSkipIfUnverified && user.EmailAddressVerifiedAt == nil {
		logger.Info("user email address not verified, skipping email delivery")
		return nil
	}

	logger.Info("sending email")

	if err = emailer.SendEmail(ctx, mail); err != nil {
		observability.AcknowledgeError(err, logger, span, "sending %s email", emailType)
	}

	if err = analyticsEventReporter.EventOccurred(ctx, email.SentEventType, emailDeliveryRequest.UserID, emailDeliveryRequest.TemplateParams); err != nil {
		observability.AcknowledgeError(err, logger, span, "notifying customer data platform")
	}

	return nil
}
