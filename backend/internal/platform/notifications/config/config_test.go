package config

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestP8File(t *testing.T) string {
	t.Helper()

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)

	keyBytes, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: keyBytes,
	}

	dir := t.TempDir()
	path := filepath.Join(dir, "AuthKey.p8")
	require.NoError(t, os.WriteFile(path, pem.EncodeToMemory(block), 0o600))
	return path
}

func TestConfig_ValidateWithContext(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	t.Run("with noop provider", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{Provider: ProviderNoop}
		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})

	t.Run("with empty provider", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{Provider: ""}
		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})

	t.Run("with apns_fcm provider and nil APNs", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider: ProviderAPNsFCM,
			APNs:     nil,
			FCM:      &FCMConfig{},
		}
		assert.Error(t, cfg.ValidateWithContext(ctx))
	})

	t.Run("with apns_fcm provider and nil FCM", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{
			Provider: ProviderAPNsFCM,
			APNs:     &APNsConfig{AuthKeyPath: "x", KeyID: "x", TeamID: "x", BundleID: "x"},
			FCM:      nil,
		}
		assert.Error(t, cfg.ValidateWithContext(ctx))
	})

	t.Run("with apns_fcm provider and both configs", func(t *testing.T) {
		t.Parallel()

		p8Path := createTestP8File(t)
		cfg := &Config{
			Provider: ProviderAPNsFCM,
			APNs:     &APNsConfig{AuthKeyPath: p8Path, KeyID: "x", TeamID: "x", BundleID: "x"},
			FCM:      &FCMConfig{},
		}
		assert.NoError(t, cfg.ValidateWithContext(ctx))
	})
}

func TestConfig_ProvidePushSender(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	logger := logging.NewNoopLogger()
	tracer := tracing.NewNoopTracerProvider()

	t.Run("with empty provider returns noop", func(t *testing.T) {
		t.Parallel()

		cfg := Config{Provider: ""}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		// Noop returns nil on SendPush
		assert.NoError(t, sender.SendPush(ctx, "ios", "token", "title", "body"))
	})

	t.Run("with noop provider returns noop", func(t *testing.T) {
		t.Parallel()

		cfg := Config{Provider: ProviderNoop}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		assert.NoError(t, sender.SendPush(ctx, "android", "token", "title", "body"))
	})

	t.Run("with apns_fcm provider and nil APNs returns noop", func(t *testing.T) {
		t.Parallel()

		cfg := Config{
			Provider: ProviderAPNsFCM,
			APNs:     nil,
			FCM:      &FCMConfig{},
		}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		assert.NoError(t, sender.SendPush(ctx, "ios", "token", "title", "body"))
	})

	t.Run("with apns_fcm provider and nil FCM returns noop", func(t *testing.T) {
		t.Parallel()

		p8Path := createTestP8File(t)
		cfg := Config{
			Provider: ProviderAPNsFCM,
			APNs:     &APNsConfig{AuthKeyPath: p8Path, KeyID: "x", TeamID: "x", BundleID: "x"},
			FCM:      nil,
		}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		// Falls back to noop because FCM is nil
		assert.NoError(t, sender.SendPush(ctx, "ios", "token", "title", "body"))
	})

	t.Run("with apns_fcm provider and invalid APNs path returns noop", func(t *testing.T) {
		t.Parallel()

		cfg := Config{
			Provider: ProviderAPNsFCM,
			APNs:     &APNsConfig{AuthKeyPath: filepath.Join(t.TempDir(), "nonexistent.p8"), KeyID: "x", TeamID: "x", BundleID: "x"},
			FCM:      &FCMConfig{},
		}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		// Falls back to noop when APNs init fails
		assert.NoError(t, sender.SendPush(ctx, "ios", "token", "title", "body"))
	})

	t.Run("with apns_fcm provider and invalid FCM path returns noop", func(t *testing.T) {
		t.Parallel()

		p8Path := createTestP8File(t)
		cfg := Config{
			Provider: ProviderAPNsFCM,
			APNs:     &APNsConfig{AuthKeyPath: p8Path, KeyID: "x", TeamID: "x", BundleID: "x"},
			FCM:      &FCMConfig{CredentialsPath: filepath.Join(t.TempDir(), "nonexistent.json")},
		}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		// Falls back to noop when FCM init fails (credentials file not found)
		assert.NoError(t, sender.SendPush(ctx, "android", "token", "title", "body"))
	})

	t.Run("with unknown provider returns noop", func(t *testing.T) {
		t.Parallel()

		cfg := Config{Provider: "unknown"}
		sender, err := cfg.ProvidePushSender(ctx, logger, tracer)
		require.NoError(t, err)
		require.NotNil(t, sender)
		assert.NoError(t, sender.SendPush(ctx, "ios", "token", "title", "body"))
	})
}
