package mailjet

import (
	"context"
	"github.com/mailjet/mailjet-apiv3-go/v4"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/prixfixeco/backend/internal/email"
	"github.com/prixfixeco/backend/internal/observability/logging"
	"github.com/prixfixeco/backend/internal/observability/tracing"
)

func TestNewMailjetEmailer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		config := &Config{SecretKey: t.Name(), PublicKey: t.Name()}

		client, err := NewMailjetEmailer(config, logger, tracing.NewNoopTracerProvider(), &http.Client{})
		require.NotNil(t, client)
		require.NoError(t, err)
	})

	T.Run("with missing config", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		client, err := NewMailjetEmailer(nil, logger, tracing.NewNoopTracerProvider(), &http.Client{})
		require.Nil(t, client)
		require.Error(t, err)
	})

	T.Run("with missing config secret key", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		config := &Config{PublicKey: t.Name()}

		client, err := NewMailjetEmailer(config, logger, tracing.NewNoopTracerProvider(), &http.Client{})
		require.Nil(t, client)
		require.Error(t, err)
	})

	T.Run("with missing config public key", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		config := &Config{SecretKey: t.Name()}

		client, err := NewMailjetEmailer(config, logger, tracing.NewNoopTracerProvider(), &http.Client{})
		require.Nil(t, client)
		require.Error(t, err)
	})

	T.Run("with missing HTTP client", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		config := &Config{SecretKey: t.Name(), PublicKey: t.Name()}

		client, err := NewMailjetEmailer(config, logger, tracing.NewNoopTracerProvider(), nil)
		require.Nil(t, client)
		require.Error(t, err)
	})
}

func TestMailjetEmailer_SendEmail(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusOK)
		}))

		config := &Config{SecretKey: t.Name(), PublicKey: t.Name()}

		c, err := NewMailjetEmailer(config, logger, tracing.NewNoopTracerProvider(), ts.Client())
		require.NotNil(t, c)
		require.NoError(t, err)

		c.client.(*mailjet.Client).SetBaseURL(ts.URL + "/")

		ctx := context.Background()
		details := &email.OutboundEmailMessage{
			ToAddress:   t.Name(),
			ToName:      t.Name(),
			FromAddress: t.Name(),
			FromName:    t.Name(),
			Subject:     t.Name(),
			HTMLContent: t.Name(),
		}

		require.NoError(t, c.SendEmail(ctx, details))
	})

	T.Run("with error executing request", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Hour)
		}))

		config := &Config{SecretKey: t.Name(), PublicKey: t.Name()}
		client := ts.Client()

		c, err := NewMailjetEmailer(config, logger, tracing.NewNoopTracerProvider(), client)
		require.NotNil(t, c)
		require.NoError(t, err)

		c.client.(*mailjet.Client).SetBaseURL(ts.URL + "/")
		client.Timeout = time.Millisecond

		ctx := context.Background()
		details := &email.OutboundEmailMessage{
			ToAddress:   t.Name(),
			ToName:      t.Name(),
			FromAddress: t.Name(),
			FromName:    t.Name(),
			Subject:     t.Name(),
			HTMLContent: t.Name(),
		}

		require.Error(t, c.SendEmail(ctx, details))
	})
}
