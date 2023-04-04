package splitio

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
	"github.com/prixfixeco/backend/pkg/types/fakes"

	"github.com/splitio/go-client/v6/splitio/conf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewFeatureFlagManager(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{SDKKey: t.Name()}

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			// SplitChangesDTO structure to map JSON message sent by Split servers.
			data := struct {
				Till  int64 `json:"till"`
				Since int64 `json:"since"`
			}{
				Till:  time.Now().Unix(),
				Since: time.Now().Unix(),
			}

			require.NoError(t, json.NewEncoder(res).Encode(data))
			res.WriteHeader(http.StatusOK)
		}))

		actual, err := NewFeatureFlagManager(cfg, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), ts.Client(), func(config *conf.SplitSdkConfig) *conf.SplitSdkConfig {
			config.Advanced.SdkURL = ts.URL
			return config
		})
		require.NoError(t, err)
		require.NotNil(t, actual)
	})

	T.Run("with missing http client", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{SDKKey: t.Name()}

		actual, err := NewFeatureFlagManager(cfg, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), nil)
		require.Error(t, err)
		require.Nil(t, actual)
	})

	T.Run("with nil config", func(t *testing.T) {
		t.Parallel()

		actual, err := NewFeatureFlagManager(nil, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), http.DefaultClient)
		require.Error(t, err)
		require.Nil(t, actual)
	})

	T.Run("with missing SDK key", func(t *testing.T) {
		t.Parallel()

		cfg := &Config{}

		actual, err := NewFeatureFlagManager(cfg, logging.NewNoopLogger(), tracing.NewNoopTracerProvider(), http.DefaultClient)
		require.Error(t, err)
		require.Nil(t, actual)
	})
}

type mockClient struct {
	mock.Mock
}

func (m *mockClient) Treatment(key interface{}, feature string, attributes map[string]interface{}) string {
	return m.Called(key, feature, attributes).String(0)
}

func TestFeatureFlagManager_CanUseFeature(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUsername := fakes.BuildFakeUser().Username

		fakeClient := &mockClient{}
		fakeClient.On("Treatment", exampleUsername, t.Name(), map[string]any(nil)).Return("on", nil)

		ffm := &featureFlagManager{
			splitClient: fakeClient,
		}

		actual, err := ffm.CanUseFeature(ctx, exampleUsername, t.Name())
		assert.NoError(t, err)
		assert.True(t, actual)
	})
}
