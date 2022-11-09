package email

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"sync"

	"github.com/prixfixeco/backend/pkg/types"
)

const defaultEnv = "testing"

var (
	urlMapHat sync.Mutex
	urlMap    = map[string]string{
		"dev":      "https://www.prixfixe.dev",
		defaultEnv: "https://not.real.lol",
	}

	emailsMapHat sync.Mutex
	emailsMap    = map[string]struct{ outboundInvites, passwordResetCreation, passwordResetRedemption string }{
		"dev": {
			outboundInvites:         "invites@prixfixe.dev",
			passwordResetCreation:   "noreply.auth@prixfixe.dev",
			passwordResetRedemption: "noreply.auth@prixfixe.dev",
		},
		defaultEnv: {
			outboundInvites:         "not@real.lol",
			passwordResetCreation:   "not@real.lol",
			passwordResetRedemption: "not@real.lol",
		},
	}
)

var (
	//go:embed templates/invite.tmpl
	outgoingInviteTemplate string
	//go:embed templates/username_reminder.tmpl
	usernameReminderTemplate string
	//go:embed templates/password_reset.tmpl
	passwordResetTemplate string
	//go:embed templates/password_reset_token_redeemed.tmpl
	passwordResetTokenRedeemedTemplate string
)

func determineEnv() string {
	env := os.Getenv("PF_ENVIRONMENT")
	if env == "" {
		env = defaultEnv
	}

	return env
}

type inviteContent struct {
	WebAppURL    string
	Token        string
	InvitationID string
	Note         string
}

// BuildInviteMemberEmail builds an email notifying a user that they've been invited to join a household.
func BuildInviteMemberEmail(householdInvitation *types.HouseholdInvitation) (*OutboundEmailMessage, error) {
	env := determineEnv()

	urlMapHat.Lock()
	envAddr, ok := urlMap[env]
	if !ok {
		return nil, fmt.Errorf("no available URL for environment")
	}
	urlMapHat.Unlock()

	emailsMapHat.Lock()
	emails, ok := emailsMap[env]
	if !ok {
		return nil, fmt.Errorf("no available email for environment")
	}
	emailsMapHat.Unlock()

	content := &inviteContent{
		WebAppURL:    envAddr,
		Token:        householdInvitation.Token,
		InvitationID: householdInvitation.ID,
		Note:         householdInvitation.Note,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]interface{}{}).Parse(outgoingInviteTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   householdInvitation.ToEmail,
		ToName:      "",
		FromAddress: emails.outboundInvites,
		FromName:    "PrixFixe",
		Subject:     "You've been invited to join a household on PrixFixe!",
		HTMLContent: b.String(),
	}

	return msg, nil
}

type resetContent struct {
	WebAppURL string
	Token     string
}

// BuildGeneratedPasswordResetTokenEmail builds an email notifying a user that they've been invited to join a household.
func BuildGeneratedPasswordResetTokenEmail(toEmail string, passwordResetToken *types.PasswordResetToken) (*OutboundEmailMessage, error) {
	env := determineEnv()

	urlMapHat.Lock()
	envAddr, ok := urlMap[env]
	if !ok {
		return nil, fmt.Errorf("no available URL for environment")
	}
	urlMapHat.Unlock()

	emailsMapHat.Lock()
	emails, ok := emailsMap[env]
	if !ok {
		return nil, fmt.Errorf("no available email for environment")
	}
	emailsMapHat.Unlock()

	content := &resetContent{
		WebAppURL: envAddr,
		Token:     passwordResetToken.Token,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]interface{}{}).Parse(passwordResetTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   toEmail,
		ToName:      "",
		FromAddress: emails.passwordResetCreation,
		FromName:    "PrixFixe",
		Subject:     "A password reset link was requested for your PrixFixe account",
		HTMLContent: b.String(),
	}

	return msg, nil
}

type usernameReminderContent struct {
	WebAppURL string
	Username  string
}

// BuildUsernameReminderEmail builds an email notifying a user that they've been invited to join a household.
func BuildUsernameReminderEmail(toEmail, username string) (*OutboundEmailMessage, error) {
	env := determineEnv()

	urlMapHat.Lock()
	envAddr, ok := urlMap[env]
	if !ok {
		return nil, fmt.Errorf("no available URL for environment")
	}
	urlMapHat.Unlock()

	emailsMapHat.Lock()
	emails, ok := emailsMap[env]
	if !ok {
		return nil, fmt.Errorf("no available email for environment")
	}
	emailsMapHat.Unlock()

	content := &usernameReminderContent{
		WebAppURL: envAddr,
		Username:  username,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]interface{}{}).Parse(usernameReminderTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   toEmail,
		ToName:      "",
		FromAddress: emails.passwordResetCreation,
		FromName:    "PrixFixe",
		Subject:     "A password reset link was requested for your PrixFixe account",
		HTMLContent: b.String(),
	}

	return msg, nil
}

type redemptionContent struct {
	WebAppURL string
}

// BuildPasswordResetTokenRedeemedEmail builds an email notifying a user that they've been invited to join a household.
func BuildPasswordResetTokenRedeemedEmail(toEmail string) (*OutboundEmailMessage, error) {
	env := determineEnv()

	urlMapHat.Lock()
	envAddr, ok := urlMap[env]
	if !ok {
		return nil, fmt.Errorf("no available URL for environment")
	}
	urlMapHat.Unlock()

	emailsMapHat.Lock()
	emails, ok := emailsMap[env]
	if !ok {
		return nil, fmt.Errorf("no available email for environment")
	}
	emailsMapHat.Unlock()

	content := &redemptionContent{
		WebAppURL: envAddr,
	}

	tmpl := template.Must(template.New("").Funcs(map[string]interface{}{}).Parse(passwordResetTokenRedeemedTemplate))
	var b bytes.Buffer
	if err := tmpl.Execute(&b, content); err != nil {
		return nil, fmt.Errorf("error rendering email template: %w", err)
	}

	msg := &OutboundEmailMessage{
		ToAddress:   toEmail,
		ToName:      "",
		FromAddress: emails.passwordResetRedemption,
		FromName:    "PrixFixe",
		Subject:     "Your PrixFixe account password has been changed.",
		HTMLContent: b.String(),
	}

	return msg, nil
}
