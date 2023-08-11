package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

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

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeUser()
		input.Username = fmt.Sprintf("%s_%d", exampleUser.Username, i)
		createdUsers = append(createdUsers, createUserForTest(t, ctx, input, dbc))
	}

	// fetch as list
	users, err := dbc.GetUsers(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, users.Data)

	firstUser := createdUsers[0]

	u, err := dbc.GetUserByUsername(ctx, firstUser.Username)
	assert.NoError(t, err)
	assert.Equal(t, firstUser.ID, u.ID)
	firstUser = u

	u, err = dbc.GetUserByEmail(ctx, firstUser.EmailAddress)
	assert.NoError(t, err)
	assert.Equal(t, firstUser, u)

	foundForUsername, err := dbc.SearchForUsersByUsername(ctx, firstUser.Username)
	assert.NoError(t, err)
	assert.NotEmpty(t, foundForUsername)

	// update first user's username
	firstUser.Username = fmt.Sprintf("%s_new", firstUser.Username)
	assert.NoError(t, dbc.UpdateUserUsername(ctx, firstUser.ID, firstUser.Username))

	// update first user's details
	assert.NoError(t, dbc.UpdateUserDetails(ctx, firstUser.ID, &types.UserDetailsDatabaseUpdateInput{
		FirstName: firstUser.FirstName,
		LastName:  firstUser.LastName,
		Birthday:  time.Now(),
	}))

	// update first user's avatar
	newAvatarSource := fakes.BuildFakeID()
	assert.NoError(t, dbc.UpdateUserAvatar(ctx, firstUser.ID, newAvatarSource))
	firstUser.AvatarSrc = &newAvatarSource

	// update first user's email address
	newEmailAddress := fakes.BuildFakeID()
	assert.NoError(t, dbc.UpdateUserEmailAddress(ctx, firstUser.ID, newEmailAddress))
	firstUser.EmailAddress = newEmailAddress

	// update first user's password
	newPassword := fakes.BuildFakeID()
	assert.NoError(t, dbc.UpdateUserPassword(ctx, firstUser.ID, newPassword))
	firstUser.HashedPassword = newPassword

	// update first user's two factor secret
	new2FASecret := fakes.BuildFakeID()
	assert.NoError(t, dbc.UpdateUserTwoFactorSecret(ctx, firstUser.ID, new2FASecret))
	firstUser.TwoFactorSecret = new2FASecret

	assert.NoError(t, dbc.MarkUserTwoFactorSecretAsVerified(ctx, firstUser.ID))
	assert.NoError(t, dbc.MarkUserTwoFactorSecretAsUnverified(ctx, firstUser.ID, fakes.BuildFakeID()))

	u, err = dbc.GetUserWithUnverifiedTwoFactorSecret(ctx, firstUser.ID)
	assert.NoError(t, err)
	firstUser.LastUpdatedAt = u.LastUpdatedAt                         // we've been changing a bunch of stuff
	firstUser.PasswordLastChangedAt = u.PasswordLastChangedAt         // from the UpdateUserPassword call above
	firstUser.TwoFactorSecretVerifiedAt = u.TwoFactorSecretVerifiedAt // from the two calls above
	firstUser.Birthday = u.Birthday
	firstUser.AccountStatus = u.AccountStatus
	firstUser.TwoFactorSecret = u.TwoFactorSecret
	assert.Equal(t, firstUser, u)

	userIDs, err := dbc.GetUserIDsThatNeedSearchIndexing(ctx)
	assert.NotEmpty(t, userIDs)
	assert.NoError(t, err)

	assert.NoError(t, dbc.MarkUserAsIndexed(ctx, firstUser.ID))

	token, err := dbc.GetEmailAddressVerificationTokenForUser(ctx, firstUser.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	u, err = dbc.GetUserByEmailAddressVerificationToken(ctx, token)
	assert.NoError(t, err)
	assert.Equal(t, firstUser, u)

	assert.NoError(t, dbc.MarkUserEmailAddressAsVerified(ctx, firstUser.ID, token))

	res, err := dbc.db.Exec("UPDATE users SET service_role = $1, two_factor_secret_verified_at = NOW() WHERE id = $2", "service_admin", firstUser.ID)
	assert.NoError(t, err)
	rowsAffected, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, rowsAffected, int64(1))

	u, err = dbc.GetAdminUserByUsername(ctx, firstUser.Username)
	assert.NoError(t, err)
	firstUser.AccountStatus = u.AccountStatus
	firstUser.AccountStatusExplanation = u.AccountStatusExplanation
	firstUser.ServiceRole = u.ServiceRole
	firstUser.EmailAddressVerifiedAt = u.EmailAddressVerifiedAt
	firstUser.TwoFactorSecretVerifiedAt = u.TwoFactorSecretVerifiedAt
	firstUser.LastUpdatedAt = u.LastUpdatedAt
	assert.Equal(t, firstUser, u)

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
