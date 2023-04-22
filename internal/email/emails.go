package email

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"sync"

	"github.com/prixfixeco/backend/pkg/types"
)

const (
	companyName = "PrixFixe"

	defaultEnv = "testing"

	// SentEventType indicates a recipe step was created.
	SentEventType types.CustomerEventType = "email_sent"

	// TemplateTypeInvite is used to indicate the invite template.
	TemplateTypeInvite = "invite"
	// TemplateTypeUsernameReminder is used to indicate the username_reminder template.
	TemplateTypeUsernameReminder = "username_reminder"
	// TemplateTypePasswordReset is used to indicate the password_reset template.
	TemplateTypePasswordReset = "password_reset"
	// TemplateTypePasswordResetTokenRedeemed is used to indicate the password_reset_token_redeemed template.
	TemplateTypePasswordResetTokenRedeemed = "password_reset_token_redeemed"
)

type (
	// DeliveryRequest is the type to use when requesting emails within the service.
	DeliveryRequest struct {
		_                  struct{}
		TemplateParams     map[string]any             `json:"templateParams"`
		Invitation         *types.HouseholdInvitation `json:"invitation,omitempty"`
		PasswordResetToken *types.PasswordResetToken  `json:"passwordResetToken,omitempty"`
		UserID             string                     `json:"forUserId"`
		Template           string                     `json:"template"`
	}

	// EnvironmentConfig is the configuration for a given environment.
	EnvironmentConfig struct {
		baseURL,
		outboundInvitesEmailAddress,
		passwordResetCreationEmailAddress,
		passwordResetRedemptionEmailAddress string
	}
)

var (
	ErrMissingEnvCfg = errors.New("missing environment configuration")

	//go:embed templates/invite.tmpl
	outgoingInviteTemplate string
	//go:embed templates/username_reminder.tmpl
	usernameReminderTemplate string
	//go:embed templates/password_reset.tmpl
	passwordResetTemplate string
	//go:embed templates/password_reset_token_redeemed.tmpl
	passwordResetTokenRedeemedTemplate string

	envConfigsMapHat sync.Mutex
	envConfigsMap    = map[string]*EnvironmentConfig{
		"dev": {
			baseURL:                             "https://www.prixfixe.dev",
			outboundInvitesEmailAddress:         "invites@prixfixe.dev",
			passwordResetCreationEmailAddress:   "noreply.auth@prixfixe.dev",
			passwordResetRedemptionEmailAddress: "noreply.auth@prixfixe.dev",
		},
		defaultEnv: {
			baseURL:                             "https://not.real.lol",
			outboundInvitesEmailAddress:         "not@real.lol",
			passwordResetCreationEmailAddress:   "not@real.lol",
			passwordResetRedemptionEmailAddress: "not@real.lol",
		},
	}
)

// BaseURL returns the BaseURL field.
func (c *EnvironmentConfig) BaseURL() string {
	return c.baseURL
}

// OutboundInvitesEmailAddress returns the OutboundInvitesEmailAddress field.
func (c *EnvironmentConfig) OutboundInvitesEmailAddress() string {
	return c.outboundInvitesEmailAddress
}

// PasswordResetCreationEmailAddress returns the PasswordResetCreationEmailAddress field.
func (c *EnvironmentConfig) PasswordResetCreationEmailAddress() string {
	return c.passwordResetCreationEmailAddress
}

// PasswordResetRedemptionEmailAddress returns the PasswordResetRedemptionEmailAddress field.
func (c *EnvironmentConfig) PasswordResetRedemptionEmailAddress() string {
	return c.passwordResetRedemptionEmailAddress
}

func GetConfigForEnvironment(env string) *EnvironmentConfig {
	envConfigsMapHat.Lock()
	defer envConfigsMapHat.Unlock()

	return envConfigsMap[env]
}

type inviteContent struct {
	WebAppURL    string
	Token        string
	InvitationID string
	Note         string
}

// BuildInviteMemberEmail builds an email notifying a user that they've been invited to join a household.
func BuildInviteMemberEmail(householdInvitation *types.HouseholdInvitation, envCfg *EnvironmentConfig) (*OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, ErrMissingEnvCfg
	}

	content := &inviteContent{
		WebAppURL:    envCfg.baseURL,
		Token:        householdInvitation.Token,
		InvitationID: householdInvitation.ID,
		Note:         householdInvitation.Note,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]any{}).Parse(outgoingInviteTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   householdInvitation.ToEmail,
		ToName:      "",
		FromAddress: envCfg.outboundInvitesEmailAddress,
		FromName:    companyName,
		Subject:     fmt.Sprintf("You've been invited to join a household on %s!", companyName),
		HTMLContent: b.String(),
	}

	return msg, nil
}

type resetContent struct {
	WebAppURL string
	Token     string
}

// BuildGeneratedPasswordResetTokenEmail builds an email notifying a user that they've been invited to join a household.
func BuildGeneratedPasswordResetTokenEmail(toEmail string, passwordResetToken *types.PasswordResetToken, envCfg *EnvironmentConfig) (*OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, ErrMissingEnvCfg
	}

	content := &resetContent{
		WebAppURL: envCfg.BaseURL(),
		Token:     passwordResetToken.Token,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]any{}).Parse(passwordResetTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   toEmail,
		ToName:      "",
		FromAddress: envCfg.passwordResetCreationEmailAddress,
		FromName:    companyName,
		Subject:     fmt.Sprintf("A password reset link was requested for your %s account", companyName),
		HTMLContent: b.String(),
	}

	return msg, nil
}

type usernameReminderContent struct {
	WebAppURL string
	Username  string
}

// BuildUsernameReminderEmail builds an email notifying a user that they've been invited to join a household.
func BuildUsernameReminderEmail(toEmail, username string, envCfg *EnvironmentConfig) (*OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, ErrMissingEnvCfg
	}

	content := &usernameReminderContent{
		WebAppURL: envCfg.baseURL,
		Username:  username,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]any{}).Parse(usernameReminderTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   toEmail,
		FromAddress: envCfg.passwordResetCreationEmailAddress,
		FromName:    companyName,
		Subject:     fmt.Sprintf("A password reset link was requested for your %s account", companyName),
		HTMLContent: b.String(),
	}

	return msg, nil
}

type redemptionContent struct {
	WebAppURL string
}

// BuildPasswordResetTokenRedeemedEmail builds an email notifying a user that they've been invited to join a household.
func BuildPasswordResetTokenRedeemedEmail(toEmail string, envCfg *EnvironmentConfig) (*OutboundEmailMessage, error) {
	if envCfg == nil {
		return nil, ErrMissingEnvCfg
	}

	content := &redemptionContent{
		WebAppURL: envCfg.baseURL,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]any{}).Parse(passwordResetTokenRedeemedTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   toEmail,
		FromAddress: envCfg.passwordResetRedemptionEmailAddress,
		FromName:    companyName,
		Subject:     fmt.Sprintf("Your %s account password has been changed.", companyName),
		HTMLContent: b.String(),
	}

	return msg, nil
}
