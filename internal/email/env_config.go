package email

import (
	"errors"
	"sync"
)

type (

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
