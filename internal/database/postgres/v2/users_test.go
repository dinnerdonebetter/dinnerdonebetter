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

	require.Equal(t, user.CreatedAt, exampleUser.CreatedAt, "CreatedAt is wrong")
	require.Equal(t, user.LastUpdatedAt, exampleUser.LastUpdatedAt, "LastUpdatedAt is wrong")
	require.Equal(t, user.ArchivedAt, exampleUser.ArchivedAt, "ArchivedAt is wrong")
	require.Equal(t, user.PasswordLastChangedAt, exampleUser.PasswordLastChangedAt, "PasswordLastChangedAt is wrong")
	require.Equal(t, user.LastAcceptedTOS, exampleUser.LastAcceptedTOS, "LastAcceptedTOS is wrong")
	require.Equal(t, user.LastAcceptedPrivacyPolicy, exampleUser.LastAcceptedPrivacyPolicy, "LastAcceptedPrivacyPolicy is wrong")
	require.Equal(t, user.TwoFactorSecretVerifiedAt, exampleUser.TwoFactorSecretVerifiedAt, "TwoFactorSecretVerifiedAt is wrong")
	require.Equal(t, user.AvatarSrc, exampleUser.AvatarSrc, "AvatarSrc is wrong")
	require.Equal(t, user.Birthday, exampleUser.Birthday, "Birthday is wrong")
	require.Equal(t, user.AccountStatusExplanation, exampleUser.AccountStatusExplanation, "AccountStatusExplanation is wrong")
	require.Equal(t, user.TwoFactorSecret, exampleUser.TwoFactorSecret, "TwoFactorSecret is wrong")
	require.Equal(t, user.HashedPassword, exampleUser.HashedPassword, "HashedPassword is wrong")
	require.Equal(t, user.ID, exampleUser.ID, "ID is wrong")
	require.Equal(t, user.AccountStatus, exampleUser.AccountStatus, "AccountStatus is wrong")
	require.Equal(t, user.Username, exampleUser.Username, "Username is wrong")
	require.Equal(t, user.FirstName, exampleUser.FirstName, "FirstName is wrong")
	require.Equal(t, user.LastName, exampleUser.LastName, "LastName is wrong")
	require.Equal(t, user.EmailAddress, exampleUser.EmailAddress, "EmailAddress is wrong")
	require.Equal(t, user.EmailAddressVerifiedAt, exampleUser.EmailAddressVerifiedAt, "EmailAddressVerifiedAt is wrong")
	require.Equal(t, user.ServiceRole, exampleUser.ServiceRole, "ServiceRole is wrong")
	require.Equal(t, user.RequiresPasswordChange, exampleUser.RequiresPasswordChange, "RequiresPasswordChange is wrong")

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

	// update
	//updatedUser := fakes.BuildFakeUser()
	//updatedUser.ID = createdUsers[0].ID
	//var x User
	//require.NoError(t, copier.Copy(&x, updatedUser))
	//require.NoError(t, dbc.UpdateUser(ctx, updatedUser))

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
