package apiclient

import (
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/fake"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_SetOptions(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expectedURL, err := url.Parse("https://notarealplace.lol")
		require.NoError(t, err)

		c, _ := buildSimpleTestClient(t)
		assert.NotEqual(t, expectedURL, c.URL(), "expected and actual URLs match somehow")

		exampleOption := func(client *Client) error {
			client.url = expectedURL
			return nil
		}

		require.NoError(t, c.SetOptions(exampleOption))

		assert.Equal(t, expectedURL, c.URL(), "expected and actual URLs do not match")
	})

	T.Run("with error", func(t *testing.T) {
		t.Parallel()

		c, _ := buildSimpleTestClient(t)

		exampleOption := func(client *Client) error {
			return errors.New("blah")
		}

		assert.Error(t, c.SetOptions(exampleOption))
	})
}

func TestUsingJSON(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, err := NewClient(mustParseURL(exampleURI), tracing.NewNoopTracerProvider(), UsingJSON(tracing.NewNoopTracerProvider()))
		assert.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, "application/json", c.encoder.ContentType())
	})
}

func TestUsingXML(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, err := NewClient(mustParseURL(exampleURI), tracing.NewNoopTracerProvider(), UsingXML(tracing.NewNoopTracerProvider()))
		assert.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, "application/xml", c.encoder.ContentType())
	})
}

func TestUsingLogger(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expectedURL, err := url.Parse("https://whatever.whocares.gov")
		require.NoError(t, err)

		c, err := NewClient(expectedURL, tracing.NewNoopTracerProvider(), UsingLogger(logging.NewNoopLogger()))
		assert.NotNil(t, c)
		assert.NoError(t, err)
	})
}

func TestUsingDebug(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, err := NewClient(mustParseURL(exampleURI), tracing.NewNoopTracerProvider(), UsingDebug(true))
		assert.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, true, c.debug, "REPLACE ME")
	})
}

func TestUsingTimeout(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := time.Minute

		c, err := NewClient(mustParseURL(exampleURI), tracing.NewNoopTracerProvider(), UsingTimeout(expected))
		assert.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, expected, c.authedClient.Timeout)
	})

	T.Run("with fallback to default timeout", func(t *testing.T) {
		t.Parallel()

		c, err := NewClient(mustParseURL(exampleURI), tracing.NewNoopTracerProvider(), UsingTimeout(0))

		assert.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, defaultTimeout, c.authedClient.Timeout)
	})
}

func TestImpersonatingUser(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, _ := buildSimpleTestClient(t)
		exampleUserID := fake.BuildFakeID()

		require.NoError(t, c.SetOptions(ImpersonatingUser(exampleUserID)))
		assert.Equal(t, c.impersonatedUserID, exampleUserID)
	})
}

func TestImpersonatingHousehold(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, _ := buildSimpleTestClient(t)
		exampleHouseholdID := fake.BuildFakeID()

		require.NoError(t, c.SetOptions(ImpersonatingHousehold(exampleHouseholdID)))
		assert.Equal(t, c.impersonatedHouseholdID, exampleHouseholdID)
	})
}

func TestWithoutImpersonating(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		c, _ := buildSimpleTestClient(t)
		c.impersonatedHouseholdID = fake.BuildFakeID()
		c.impersonatedUserID = fake.BuildFakeID()

		require.NoError(t, c.SetOptions(WithoutImpersonating()))
		assert.Empty(t, c.impersonatedHouseholdID)
		assert.Empty(t, c.impersonatedUserID)
	})
}
