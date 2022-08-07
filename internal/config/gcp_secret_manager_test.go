package config

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestGetAPIServerConfigFromGoogleCloudRunEnvironment(T *testing.T) {
	T.Parallel()

	// TODO: actually write tests for me, coward

	T.Run("standard", func(t *testing.T) {
		t.SkipNow()

		require.NoError(t, os.Setenv(gcpConfigFilePathEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpPortEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpDatabaseSocketDirEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpDatabaseUserEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpDatabaseUserPasswordEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpDatabaseNameEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpDatabaseInstanceConnNameEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpCookieHashKeyEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpCookieBlockKeyEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpPASETOLocalKeyEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpSendgridTokenEnvVarKey, ""))
		require.NoError(t, os.Setenv(gcpSegmentTokenEnvVarKey, ""))

		ctx := context.Background()
		cfg, err := GetAPIServerConfigFromGoogleCloudRunEnvironment(ctx)
		assert.NotNil(t, cfg)
		assert.NoError(t, err)

		require.NoError(t, os.Unsetenv(gcpConfigFilePathEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpPortEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseSocketDirEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseUserEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseUserPasswordEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseNameEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpDatabaseInstanceConnNameEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpCookieHashKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpCookieBlockKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpPASETOLocalKeyEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpSendgridTokenEnvVarKey))
		require.NoError(t, os.Unsetenv(gcpSegmentTokenEnvVarKey))
	})
}
