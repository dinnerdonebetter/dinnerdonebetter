package sendgrid

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/prixfixeco/api_server/internal/email"
	"github.com/prixfixeco/api_server/internal/observability/logging"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
)

func TestNewSendGridEmailer(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		logger := logging.NewNoopLogger()

		client, err := NewSendGridEmailer(t.Name(), logger, tracing.NewNoopTracerProvider(), &http.Client{})
		require.NotNil(t, client)
		require.NoError(t, err)
	})
}

func TestSendGridEmailer_SendEmail(T *testing.T) {
	T.Run("standard", func(t *testing.T) {
		logger := logging.NewNoopLogger()

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusAccepted)
		}))

		c, err := NewSendGridEmailer(t.Name(), logger, tracing.NewNoopTracerProvider(), ts.Client())
		require.NotNil(t, c)
		require.NoError(t, err)

		c.client.Request.BaseURL = ts.URL

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
		logger := logging.NewNoopLogger()

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			time.Sleep(time.Hour)
		}))
		client := ts.Client()
		client.Timeout = time.Millisecond

		c, err := NewSendGridEmailer(t.Name(), logger, tracing.NewNoopTracerProvider(), client)
		require.NotNil(t, c)
		require.NoError(t, err)

		c.client.Request.BaseURL = ts.URL

		ctx := context.Background()
		details := &email.OutboundEmailMessage{
			ToAddress:   t.Name(),
			ToName:      t.Name(),
			FromAddress: t.Name(),
			FromName:    t.Name(),
			Subject:     t.Name(),
			HTMLContent: t.Name(),
		}

		err = c.SendEmail(ctx, details)
		require.Error(t, err)
	})

	T.Run("with invalid response code", func(t *testing.T) {
		logger := logging.NewNoopLogger()

		ts := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
			res.WriteHeader(http.StatusInternalServerError)
		}))

		c, err := NewSendGridEmailer(t.Name(), logger, tracing.NewNoopTracerProvider(), ts.Client())
		require.NotNil(t, c)
		require.NoError(t, err)

		c.client.Request.BaseURL = ts.URL

		ctx := context.Background()
		details := &email.OutboundEmailMessage{
			ToAddress:   t.Name(),
			ToName:      t.Name(),
			FromAddress: t.Name(),
			FromName:    t.Name(),
			Subject:     t.Name(),
			HTMLContent: t.Name(),
		}

		err = c.SendEmail(ctx, details)
		require.Error(t, err)
	})
}
