package v2

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/require"
)

func createUserForTest(t *testing.T, ctx context.Context, exampleUser *types.User, dbc *DatabaseClient) *types.User {
	t.Helper()

	// create
	if exampleUser == nil {
		exampleUser = fakes.BuildFakeUser()
	}
	var x User
	require.NoError(t, copier.Copy(&x, exampleUser))

	created, err := dbc.CreateUser(ctx, &x)
	require.NoError(t, err)
	require.Equal(t, exampleUser, created)

	user, err := dbc.GetUser(ctx, created.ID)
	exampleUser.CreatedAt = user.CreatedAt
	exampleUser.TwoFactorSecretVerifiedAt = user.TwoFactorSecretVerifiedAt
	exampleUser.Birthday = user.Birthday

	require.NoError(t, err)
	require.Equal(t, user, exampleUser)

	return created
}

func TestDatabaseClient_Users(t *testing.T) {
	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		require.NoError(t, container.Terminate(ctx))
	}(t)

	exampleUser := fakes.BuildFakeUser()
	createdUsers := []*types.User{}

	// create
	createdUsers = append(createdUsers, createUserForTest(t, ctx, exampleUser, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeUser()
		input.Username = fmt.Sprintf("%s %d", exampleUser.Username, i)
		createdUsers = append(createdUsers, createUserForTest(t, ctx, input, dbc))
	}

	// fetch as list
	users, err := dbc.GetUsers(ctx, nil)
	require.NoError(t, err)
	require.NotEmpty(t, users.Data)
	require.Equal(t, len(createdUsers), len(users.Data))

	// fetch as list of IDs
	userIDs := []string{}
	for _, user := range createdUsers {
		userIDs = append(userIDs, user.ID)
	}

	byIDs, err := dbc.GetUsersWithIDs(ctx, userIDs)
	require.NoError(t, err)
	require.Equal(t, users.Data, byIDs)

	// fetch via name search
	byName, err := dbc.SearchForUsers(ctx, exampleUser.Username, nil)
	require.NoError(t, err)
	require.Equal(t, users, byName)

	// TODO: update

	// delete
	for _, user := range createdUsers {
		require.NoError(t, dbc.ArchiveUser(ctx, user.ID))

		var y *types.User
		y, err = dbc.GetUser(ctx, user.ID)
		require.Nil(t, y)
		require.Error(t, err)
		require.ErrorIs(t, err, sql.ErrNoRows)
	}
}
