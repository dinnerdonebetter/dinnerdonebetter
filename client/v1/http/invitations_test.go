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

func TestV1Client_BuildInvitationExistsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodHead
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		exampleInvitation := fakemodels.BuildFakeInvitation()
		actual, err := c.BuildInvitationExistsRequest(ctx, exampleInvitation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleInvitation.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_InvitationExists(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleInvitation.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/invitations/%d", exampleInvitation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodHead)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.InvitationExists(ctx, exampleInvitation.ID)

		assert.NoError(t, err, "no error should be returned")
		assert.True(t, actual)
	})

	T.Run("with erroneous response", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.InvitationExists(ctx, exampleInvitation.ID)

		assert.Error(t, err, "error should be returned")
		assert.False(t, actual)
	})
}

func TestV1Client_BuildGetInvitationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		exampleInvitation := fakemodels.BuildFakeInvitation()

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetInvitationRequest(ctx, exampleInvitation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleInvitation.ID)))
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetInvitation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleInvitation.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/invitations/%d", exampleInvitation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleInvitation))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetInvitation(ctx, exampleInvitation.ID)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleInvitation, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetInvitation(ctx, exampleInvitation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.True(t, strings.HasSuffix(req.URL.String(), strconv.Itoa(int(exampleInvitation.ID))))
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/invitations/%d", exampleInvitation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode("BLAH"))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetInvitation(ctx, exampleInvitation.ID)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildGetInvitationsRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)
		expectedMethod := http.MethodGet
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildGetInvitationsRequest(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_GetInvitations(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		expectedPath := "/api/v1/invitations"

		exampleInvitationList := fakemodels.BuildFakeInvitationList()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodGet)
					require.NoError(t, json.NewEncoder(res).Encode(exampleInvitationList))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.GetInvitations(ctx, filter)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleInvitationList, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.GetInvitations(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})

	T.Run("with invalid response", func(t *testing.T) {
		ctx := context.Background()

		filter := (*models.QueryFilter)(nil)

		expectedPath := "/api/v1/invitations"

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
		actual, err := c.GetInvitations(ctx, filter)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildCreateInvitationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleUser := fakemodels.BuildFakeUser()
		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInvitation.BelongsToUser = exampleUser.ID
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		expectedMethod := http.MethodPost
		ts := httptest.NewTLSServer(nil)

		c := buildTestClient(t, ts)
		actual, err := c.BuildCreateInvitationRequest(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_CreateInvitation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		expectedPath := "/api/v1/invitations"

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, expectedPath, "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPost)

					var x *models.InvitationCreationInput
					require.NoError(t, json.NewDecoder(req.Body).Decode(&x))

					exampleInput.BelongsToUser = 0
					assert.Equal(t, exampleInput, x)

					require.NoError(t, json.NewEncoder(res).Encode(exampleInvitation))
				},
			),
		)

		c := buildTestClient(t, ts)
		actual, err := c.CreateInvitation(ctx, exampleInput)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, exampleInvitation, actual)
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()
		exampleInput := fakemodels.BuildFakeInvitationCreationInputFromInvitation(exampleInvitation)

		c := buildTestClientWithInvalidURL(t)
		actual, err := c.CreateInvitation(ctx, exampleInput)

		assert.Nil(t, actual)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildUpdateInvitationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()
		expectedMethod := http.MethodPut

		ts := httptest.NewTLSServer(nil)
		c := buildTestClient(t, ts)
		actual, err := c.BuildUpdateInvitationRequest(ctx, exampleInvitation)

		require.NotNil(t, actual)
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_UpdateInvitation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/invitations/%d", exampleInvitation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodPut)
					assert.NoError(t, json.NewEncoder(res).Encode(exampleInvitation))
				},
			),
		)

		err := buildTestClient(t, ts).UpdateInvitation(ctx, exampleInvitation)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		err := buildTestClientWithInvalidURL(t).UpdateInvitation(ctx, exampleInvitation)
		assert.Error(t, err, "error should be returned")
	})
}

func TestV1Client_BuildArchiveInvitationRequest(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		expectedMethod := http.MethodDelete
		ts := httptest.NewTLSServer(nil)

		exampleInvitation := fakemodels.BuildFakeInvitation()

		c := buildTestClient(t, ts)
		actual, err := c.BuildArchiveInvitationRequest(ctx, exampleInvitation.ID)

		require.NotNil(t, actual)
		require.NotNil(t, actual.URL)
		assert.True(t, strings.HasSuffix(actual.URL.String(), fmt.Sprintf("%d", exampleInvitation.ID)))
		assert.NoError(t, err, "no error should be returned")
		assert.Equal(t, actual.Method, expectedMethod, "request should be a %s request", expectedMethod)
	})
}

func TestV1Client_ArchiveInvitation(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		ts := httptest.NewTLSServer(
			http.HandlerFunc(
				func(res http.ResponseWriter, req *http.Request) {
					assert.Equal(t, req.URL.Path, fmt.Sprintf("/api/v1/invitations/%d", exampleInvitation.ID), "expected and actual paths do not match")
					assert.Equal(t, req.Method, http.MethodDelete)
					res.WriteHeader(http.StatusOK)
				},
			),
		)

		err := buildTestClient(t, ts).ArchiveInvitation(ctx, exampleInvitation.ID)
		assert.NoError(t, err, "no error should be returned")
	})

	T.Run("with invalid client URL", func(t *testing.T) {
		ctx := context.Background()

		exampleInvitation := fakemodels.BuildFakeInvitation()

		err := buildTestClientWithInvalidURL(t).ArchiveInvitation(ctx, exampleInvitation.ID)
		assert.Error(t, err, "error should be returned")
	})
}
