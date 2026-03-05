package emails

import (
	"errors"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/branding"
	"github.com/dinnerdonebetter/backend/internal/domain/auth"
	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/email"
	"github.com/dinnerdonebetter/backend/internal/platform/internalerrors"

	"github.com/matcornic/hermes/v2"
)

var (
	ErrUnverifiedEmailRecipient = errors.New("missing email address verification for user")
)

// BuildInviteMemberEmail builds an email notifying a user that they've been invited to join an account.
func BuildInviteMemberEmail(recipient *identity.User, accountInvitation *identity.AccountInvitation, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: accountInvitation.ToEmail,
			Intros: []string{
				fmt.Sprintf("You've been invited to join an account on %s!", branding.CompanyName),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Text: "Join the fun",
						Link: fmt.Sprintf("%s/accept_invitation?i=%s&t=%s", envCfg.BaseURL(), accountInvitation.ID, accountInvitation.Token),
					},
				},
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		UserID:      recipient.ID,
		ToAddress:   accountInvitation.ToEmail,
		ToName:      recipient.FullName(),
		FromAddress: envCfg.OutboundInvitesEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     "You've been invited!",
		HTMLContent: htmlContent,
	}

	return msg, nil
}

// BuildGeneratedPasswordResetTokenEmail builds an email notifying a user that they've been invited to join an account.
func BuildGeneratedPasswordResetTokenEmail(recipient *identity.User, passwordResetToken *auth.PasswordResetToken, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if recipient.EmailAddressVerifiedAt == nil {
		return nil, ErrUnverifiedEmailRecipient
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.Username,
			Intros: []string{
				"You have received this email because a password reset was requested.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to reset your password:",
					Button: hermes.Button{
						Text: "Reset your password",
						Link: fmt.Sprintf("%s/reset_password?t=%s", envCfg.BaseURL(), passwordResetToken.Token),
					},
				},
			},
			Outros: []string{
				"If you did not request a password reset, no further action is required on your part.",
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		UserID:      recipient.ID,
		ToAddress:   recipient.EmailAddress,
		ToName:      recipient.FullName(),
		FromAddress: envCfg.PasswordResetCreationEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     fmt.Sprintf("A password reset link was requested for your %s account", branding.CompanyName),
		HTMLContent: htmlContent,
	}

	return msg, nil
}

// BuildUsernameReminderEmail builds an email notifying a user that they've been invited to join an account.
func BuildUsernameReminderEmail(recipient *identity.User, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if recipient.EmailAddressVerifiedAt == nil {
		return nil, ErrUnverifiedEmailRecipient
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.Username,
			Intros: []string{
				fmt.Sprintf("A username reminder for your %s account was requested. Your username is <b>%s</b>.", branding.CompanyName, recipient.Username),
			},
			Outros: []string{
				"If you did not request a username reminder, no further action is required on your part.",
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		UserID:      recipient.ID,
		ToName:      recipient.FullName(),
		ToAddress:   recipient.EmailAddress,
		FromAddress: envCfg.PasswordResetCreationEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     fmt.Sprintf("A password reset link was requested for your %s account", branding.CompanyName),
		HTMLContent: htmlContent,
	}

	return msg, nil
}

// BuildPasswordResetTokenRedeemedEmail builds an email notifying a user that they've been invited to join an account.
func BuildPasswordResetTokenRedeemedEmail(recipient *identity.User, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if recipient.EmailAddressVerifiedAt == nil {
		return nil, ErrUnverifiedEmailRecipient
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.Username,
			Intros: []string{
				"This is to inform you that your password has been changed upon successful redemption of a reset token.",
			},
			Outros: []string{
				"If you did not request a password reset, please contact support.",
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		UserID:      recipient.ID,
		ToAddress:   recipient.EmailAddress,
		ToName:      recipient.FullName(),
		FromAddress: envCfg.PasswordResetRedemptionEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     fmt.Sprintf("Your %s account password has been changed.", branding.CompanyName),
		HTMLContent: htmlContent,
	}

	return msg, nil
}

// BuildPasswordChangedEmail builds an email notifying a user that they've been invited to join an account.
func BuildPasswordChangedEmail(recipient *identity.User, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if recipient.EmailAddressVerifiedAt == nil {
		return nil, ErrUnverifiedEmailRecipient
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.Username,
			Intros: []string{
				"This is to inform you that your password has been changed.",
			},
			Outros: []string{
				"If you did not perform this action, please contact support.",
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		UserID:      recipient.ID,
		ToAddress:   recipient.EmailAddress,
		ToName:      recipient.FullName(),
		FromAddress: envCfg.PasswordResetRedemptionEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     fmt.Sprintf("Your %s account password has been changed.", branding.CompanyName),
		HTMLContent: htmlContent,
	}

	return msg, nil
}

var errEmailVerificationTokenRequired = errors.New("email verification token required")

// BuildVerifyEmailAddressEmail builds an email notifying a user that they've been invited to join an account.
func BuildVerifyEmailAddressEmail(recipient *identity.User, emailVerificationToken string, envCfg *email.EnvironmentConfig) (*email.OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, internalerrors.NilConfigError("email environment config")
	}

	if emailVerificationToken == "" {
		return nil, errEmailVerificationTokenRequired
	}

	if recipient.EmailAddressVerifiedAt != nil {
		return nil, fmt.Errorf("user %s has already verified their email address", recipient.ID)
	}

	e := hermes.Email{
		Body: hermes.Body{
			Name: recipient.Username,
			Intros: []string{
				fmt.Sprintf("You recently signed up for an account on %s. Please verify your email address.", branding.CompanyName),
			},
			Actions: []hermes.Action{
				{
					Instructions: "Click the button below to verify your email address:",
					Button: hermes.Button{
						Text: "Verify my email",
						Link: fmt.Sprintf("%s/verify_email_address?t=%s", envCfg.BaseURL(), emailVerificationToken),
					},
				},
			},
			Outros: []string{
				"If you did not sign up for an account, please contact support.",
			},
		},
	}

	htmlContent, err := envCfg.BuildHermes(branding.DefaultEmailBranding()).GenerateHTML(e)
	if err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &email.OutboundEmailMessage{
		UserID:      recipient.ID,
		ToAddress:   recipient.EmailAddress,
		ToName:      recipient.FullName(),
		FromAddress: envCfg.PasswordResetCreationEmailAddress(),
		FromName:    branding.CompanyName,
		Subject:     fmt.Sprintf("Verify your email with %s", branding.CompanyName),
		HTMLContent: htmlContent,
	}

	return msg, nil
}
