package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func createUserForTest(t *testing.T, ctx context.Context, exampleUser *types.User, dbc *Querier) *types.User {
	t.Helper()

	// create
	if exampleUser == nil {
		exampleUser = fakes.BuildFakeUser()
	}
	dbInput := converters.ConvertUserToUserDatabaseCreationInput(exampleUser)

	exampleUser.TwoFactorSecretVerifiedAt = nil
	created, err := dbc.CreateUser(ctx, dbInput)
	exampleUser.CreatedAt = created.CreatedAt
	exampleUser.TwoFactorSecretVerifiedAt = created.TwoFactorSecretVerifiedAt
	assert.NoError(t, err)
	assert.Equal(t, exampleUser, created)

	user, err := dbc.GetUser(ctx, created.ID)
	exampleUser.CreatedAt = user.CreatedAt
	exampleUser.Birthday = user.Birthday

	assert.NoError(t, err)
	assert.Equal(t, user, exampleUser)

	return created
}

func TestQuerier_Integration_Users(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleUser := fakes.BuildFakeUser()
	exampleUser.Username = fmt.Sprintf("%d", hashStringToNumber(exampleUser.Username))
	exampleUser.TwoFactorSecretVerifiedAt = nil
	createdUsers := []*types.User{}

	// create
	createdUsers = append(createdUsers, createUserForTest(t, ctx, exampleUser, dbc))

	// update
	// TODO

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeUser()
		input.Username = fmt.Sprintf("%s %d", exampleUser.Username, i)
		createdUsers = append(createdUsers, createUserForTest(t, ctx, input, dbc))
	}

	// fetch as list
	users, err := dbc.GetUsers(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, users.Data)
	assert.Equal(t, len(createdUsers), len(users.Data))

	// delete
	for _, user := range createdUsers {
		assert.NoError(t, dbc.ArchiveUser(ctx, user.ID))

		var y *types.User
		y, err = dbc.GetUser(ctx, user.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
