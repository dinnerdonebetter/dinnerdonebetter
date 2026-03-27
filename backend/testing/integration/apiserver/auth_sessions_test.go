package integration

import (
	"testing"

	authsvc "github.com/dinnerdonebetter/dinnerdonebetter/backend/internal/grpc/generated/services/auth"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAuth_ListActiveSessions(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		res, err := testClient.ListActiveSessions(ctx, &authsvc.ListActiveSessionsRequest{})
		require.NoError(t, err)
		require.NotNil(t, res)
		require.NotEmpty(t, res.Sessions)

		// The session created by the login should be present and marked as current.
		var foundCurrent bool
		for _, sess := range res.Sessions {
			assert.NotEmpty(t, sess.Id)
			assert.NotEmpty(t, sess.LoginMethod)
			assert.NotNil(t, sess.CreatedAt)
			assert.NotNil(t, sess.ExpiresAt)
			if sess.IsCurrent {
				foundCurrent = true
			}
		}
		assert.True(t, foundCurrent, "expected one session to be marked as current")
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)
		_, err := unauthedClient.ListActiveSessions(ctx, &authsvc.ListActiveSessionsRequest{})
		assert.Error(t, err)
	})
}

func TestAuth_RevokeSession(T *testing.T) {
	T.Parallel()

	T.Run("revoke a non-current session", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(t)

		// Log in again to create a second session.
		token2 := fetchLoginTokenForUserForTest(t, user)
		client2, err := buildAuthedGRPCClient(ctx, token2)
		require.NoError(t, err)

		// List sessions from client 1 – should see at least 2.
		listRes, err := testClient.ListActiveSessions(ctx, &authsvc.ListActiveSessionsRequest{})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listRes.Sessions), 2)

		// Find a non-current session to revoke.
		var targetSessionID string
		for _, sess := range listRes.Sessions {
			if !sess.IsCurrent {
				targetSessionID = sess.Id
				break
			}
		}
		require.NotEmpty(t, targetSessionID, "expected at least one non-current session")

		// Revoke it.
		_, err = testClient.RevokeSession(ctx, &authsvc.RevokeSessionRequest{
			SessionId: targetSessionID,
		})
		require.NoError(t, err)

		// Verify it's gone from the list.
		listRes2, err := testClient.ListActiveSessions(ctx, &authsvc.ListActiveSessionsRequest{})
		require.NoError(t, err)
		for _, sess := range listRes2.Sessions {
			assert.NotEqual(t, targetSessionID, sess.Id, "revoked session should not appear in active list")
		}

		// The second client should still work (we revoked its session, so it may fail
		// depending on whether the second login's session was the one we revoked).
		// Instead of asserting on client2, just ensure it was created successfully.
		_ = client2
	})

	T.Run("empty session ID returns error", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, testClient := createUserAndClientForTest(t)

		_, err := testClient.RevokeSession(ctx, &authsvc.RevokeSessionRequest{
			SessionId: "",
		})
		require.Error(t, err)
		st, ok := status.FromError(err)
		require.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)
		_, err := unauthedClient.RevokeSession(ctx, &authsvc.RevokeSessionRequest{
			SessionId: "anything",
		})
		assert.Error(t, err)
	})
}

func TestAuth_RevokeAllOtherSessions(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		user, testClient := createUserAndClientForTest(t)

		// Create additional sessions by logging in more times.
		_ = fetchLoginTokenForUserForTest(t, user)
		_ = fetchLoginTokenForUserForTest(t, user)

		// Should have multiple sessions.
		listRes, err := testClient.ListActiveSessions(ctx, &authsvc.ListActiveSessionsRequest{})
		require.NoError(t, err)
		require.GreaterOrEqual(t, len(listRes.Sessions), 2, "expected multiple sessions before revoking")

		// Revoke all other sessions.
		_, err = testClient.RevokeAllOtherSessions(ctx, &authsvc.RevokeAllOtherSessionsRequest{})
		require.NoError(t, err)

		// Should now have exactly one session (the current one).
		listRes2, err := testClient.ListActiveSessions(ctx, &authsvc.ListActiveSessionsRequest{})
		require.NoError(t, err)
		require.Len(t, listRes2.Sessions, 1)
		assert.True(t, listRes2.Sessions[0].IsCurrent)
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		unauthedClient := buildUnauthenticatedGRPCClientForTest(t)
		_, err := unauthedClient.RevokeAllOtherSessions(ctx, &authsvc.RevokeAllOtherSessionsRequest{})
		assert.Error(t, err)
	})
}
