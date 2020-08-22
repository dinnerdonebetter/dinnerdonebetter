package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
	fakemodels "gitlab.com/prixfixe/prixfixe/models/v1/fake"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestV1Client_BuildValidInstrumentExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		actual, err := c.BuildValidInstrumentExistsRequest(ctx, exampleValidInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidInstrument.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidInstrument.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_instruments/%d", exampleValidInstrument.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetValidInstrumentRequest(ctx, exampleValidInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidInstrument.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidInstrument.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_instruments/%d", exampleValidInstrument.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidInstrument))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidInstrument, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleValidInstrument.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_instruments/%d", exampleValidInstrument.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetValidInstrumentsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetValidInstrumentsRequest(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		expectedPath := "/api/v1/valid_instruments"

		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidInstrumentList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidInstruments(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidInstrumentList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetValidInstruments(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		expectedPath := "/api/v1/valid_instruments"

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetValidInstruments(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildSearchValidInstrumentsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		limit := models.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildSearchValidInstrumentsRequest(ctx, exampleQuery, limit)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_SearchValidInstruments(T *testing.T) {
	T.Parallel()

	const expectedPath = "/api/v1/valid_instruments/search"

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		limit := models.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		exampleValidInstrumentList := fakemodels.BuildFakeValidInstrumentList().ValidInstruments

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.URL.Query().Get(models.SearchQueryKey), exampleQuery, "expected and actual search query param do not match")
					assert.Equal(t, req.URL.Query().Get(models.LimitQueryKey), strconv.FormatUint(uint64(limit), 10), "expected and actual limit query param do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleValidInstrumentList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.SearchValidInstruments(ctx, exampleQuery, limit)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidInstrumentList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		limit := models.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.SearchValidInstruments(ctx, exampleQuery, limit)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		limit := models.DefaultQueryFilter().Limit
		exampleQuery := "whatever"

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.URL.Query().Get(models.SearchQueryKey), exampleQuery, "expected and actual search query param do not match")
					assert.Equal(t, req.URL.Query().Get(models.LimitQueryKey), strconv.FormatUint(uint64(limit), 10), "expected and actual limit query param do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.SearchValidInstruments(ctx, exampleQuery, limit)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		exampleInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateValidInstrumentRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		exampleInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		expectedPath := "/api/v1/valid_instruments"

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.ValidInstrumentCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleValidInstrument))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateValidInstrument(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleValidInstrument, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		exampleInput := fakemodels.BuildFakeValidInstrumentCreationInputFromValidInstrument(exampleValidInstrument)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateValidInstrument(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateValidInstrumentRequest(ctx, exampleValidInstrument)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_instruments/%d", exampleValidInstrument.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleValidInstrument))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateValidInstrument(ctx, exampleValidInstrument)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		err := buildTestClientWithInvalidURL(t).UpdateValidInstrument(ctx, exampleValidInstrument)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveValidInstrumentRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveValidInstrumentRequest(ctx, exampleValidInstrument.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleValidInstrument.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/valid_instruments/%d", exampleValidInstrument.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleValidInstrument := fakemodels.BuildFakeValidInstrument()

		err := buildTestClientWithInvalidURL(t).ArchiveValidInstrument(ctx, exampleValidInstrument.ID)
		assert.Error(t, err, "error should be returned")
	})
}
