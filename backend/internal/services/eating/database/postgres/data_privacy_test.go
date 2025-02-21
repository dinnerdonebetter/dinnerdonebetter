package postgres

/*

import (
	"context"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/services/eating/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Integration_DataPrivacy(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleUser := fakes.BuildFakeUser()
	exampleUser.Username = fmt.Sprintf("%d", hashStringToNumber(exampleUser.Username))
	exampleUser.TwoFactorSecretVerifiedAt = nil

	// create
	createdUser := createUserForTest(t, ctx, exampleUser, dbc)

	collection, err := dbc.AggregateUserData(ctx, createdUser.ID)
	require.NoError(t, err)
	assert.NotNil(t, collection.User)

	assert.NoError(t, dbc.DeleteUser(ctx, createdUser.ID))
}

func TestQuerier_DeleteUser(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.DeleteUser(ctx, ""))
	})
}

*/
