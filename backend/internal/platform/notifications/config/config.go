package config

import (
	"context"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/platform/notifications"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications/apns"
	"github.com/dinnerdonebetter/backend/internal/platform/notifications/fcm"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

const (
	// ProviderAPNsFCM represents the real APNs + FCM implementation.
	ProviderAPNsFCM = "apns_fcm"
	// ProviderNoop represents the no-op implementation.
	ProviderNoop = "noop"
)

type (
	// APNsConfig configures APNs for iOS push notifications.
	APNsConfig struct {
		AuthKeyPath string `env:"AUTH_KEY_PATH" json:"authKeyPath"`
		KeyID       string `env:"KEY_ID"        json:"keyID"`
		TeamID      string `env:"TEAM_ID"       json:"teamID"`
		BundleID    string `env:"BUNDLE_ID"     json:"bundleID"`
		Production  bool   `env:"PRODUCTION"    json:"production"`
	}

	// FCMConfig configures FCM for Android push notifications.
	FCMConfig struct {
		// CredentialsPath is the path to the Firebase service account JSON file.
		// If empty, Application Default Credentials (ADC) are used.
		CredentialsPath string `env:"CREDENTIALS_PATH" json:"credentialsPath"`
	}

	// Config is the push notifications configuration.
	Config struct {
		APNs     *APNsConfig `env:"init"     envPrefix:"APNS_" json:"apns"`
		FCM      *FCMConfig  `env:"init"     envPrefix:"FCM_"  json:"fcm"`
		Provider string      `env:"PROVIDER" json:"provider"`
	}
)

var _ validation.ValidatableWithContext = (*Config)(nil)

// ValidateWithContext validates the Config.
func (cfg *Config) ValidateWithContext(ctx context.Context) error {
	return validation.ValidateStructWithContext(
		ctx,
		cfg,
		validation.Field(&cfg.APNs, validation.When(
			strings.ToLower(strings.TrimSpace(cfg.Provider)) == ProviderAPNsFCM,
			validation.Required,
		)),
		validation.Field(&cfg.FCM, validation.When(
			strings.ToLower(strings.TrimSpace(cfg.Provider)) == ProviderAPNsFCM,
			validation.Required,
		)),
	)
}

// ProvidePushSender returns a PushNotificationSender based on config.
// When provider is "apns_fcm" and config is valid, returns MultiPlatformPushSender.
func (cfg *Config) ProvidePushSender(
	ctx context.Context,
	logger logging.Logger,
	tracerProvider tracing.TracerProvider,
) (notifications.PushNotificationSender, error) {
	switch strings.ToLower(strings.TrimSpace(cfg.Provider)) {
	case ProviderAPNsFCM:
		if cfg.APNs == nil || cfg.FCM == nil {
			logger.Debug("push notifications: apns_fcm requested but config incomplete, using noop")
			return &notifications.NoopPushNotificationSender{}, nil
		}

		apnsCfg := &apns.Config{
			AuthKeyPath: cfg.APNs.AuthKeyPath,
			KeyID:       cfg.APNs.KeyID,
			TeamID:      cfg.APNs.TeamID,
			BundleID:    cfg.APNs.BundleID,
			Production:  cfg.APNs.Production,
		}

		apnsSender, err := apns.NewSender(apnsCfg, tracerProvider, logger)
		if err != nil {
			logger.WithValue("error", err).Error("push notifications: failed to create APNs sender, using noop", err)
			return &notifications.NoopPushNotificationSender{}, nil
		}
		fcmCfg := &fcm.Config{CredentialsPath: cfg.FCM.CredentialsPath}

		fcmSender, err := fcm.NewSender(ctx, fcmCfg, tracerProvider, logger)
		if err != nil {
			logger.WithValue("error", err).Error("push notifications: failed to create FCM sender, using noop", err)
			return &notifications.NoopPushNotificationSender{}, nil
		}

		return notifications.NewMultiPlatformPushSender(apnsSender, fcmSender, logger, tracerProvider), nil
	default:
		logger.Debug("push notifications: using noop sender")
		return &notifications.NoopPushNotificationSender{}, nil
	}
}
