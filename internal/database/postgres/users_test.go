package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_GetUser(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUser(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetUserWithUnverifiedTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserWithUnverifiedTwoFactorSecret(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetUserByEmail(T *testing.T) {
	T.Parallel()

	T.Run("with invalid email", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetUserByEmail(ctx, "")
		assert.Error(t, err)
		assert.Empty(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("with invalid username", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserByUsername(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetAdminUserByUsername(T *testing.T) {
	T.Parallel()

	T.Run("with empty username", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.GetAdminUserByUsername(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_SearchForUsersByUsername(T *testing.T) {
	T.Parallel()

	T.Run("with invalid username", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForUsersByUsername(ctx, "")
		assert.Error(t, err)
		assert.NotNil(t, actual)
		assert.Empty(t, actual)
	})
}

func TestQuerier_GetUserThatNeedSearchIndexing(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		_, db := buildTestClient(t)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_MarkUserAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserAsIndexed(ctx, ""))
	})
}

func TestQuerier_CreateUser(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateUser(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateUsername(T *testing.T) {
	T.Parallel()

	T.Run("with empty user", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserUsername(ctx, "", t.Name()))
	})
}

func TestQuerier_UpdateUserDetails(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.UpdateUserDetails(ctx, exampleUser.ID, nil))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserAvatar(T *testing.T) {
	T.Parallel()

	T.Run("with empty input", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.UpdateUserAvatar(ctx, exampleUser.ID, ""))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_UpdateUserPassword(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleHashedPassword := "$argon2i$v=19$m=64,t=10,p=4$RjFtMmRmU2lGYU9CMk1mMw$cuGR9AhTczPR6xDOSAMW+SvEYFyLEIS+7nlRdC9f6ys"

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserPassword(ctx, "", exampleHashedPassword))
	})

	T.Run("with invalid new hash", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserPassword(ctx, exampleUser.ID, ""))
	})
}

func TestQuerier_UpdateUserTwoFactorSecret(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserTwoFactorSecret(ctx, "", exampleUser.TwoFactorSecret))
	})

	T.Run("with invalid new secret", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserTwoFactorSecret(ctx, exampleUser.ID, ""))
	})
}

func TestQuerier_MarkUserTwoFactorSecretAsVerified(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserTwoFactorSecretAsVerified(ctx, ""))
	})
}

func TestQuerier_MarkUserTwoFactorSecretAsUnverified(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleSecret := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserTwoFactorSecretAsUnverified(ctx, "", exampleSecret))
	})

	T.Run("with invalid secret", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkUserTwoFactorSecretAsUnverified(ctx, exampleUserID, ""))
	})
}

func TestQuerier_ArchiveUser(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveUser(ctx, ""))
	})

	T.Run("with error beginning transaction", func(t *testing.T) {
		t.Parallel()

		exampleUserID := fakes.BuildFakeID()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin().WillReturnError(errors.New("blah"))

		assert.Error(t, c.ArchiveUser(ctx, exampleUserID))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetUserByEmailAddressVerificationToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserByEmailAddressVerificationToken(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_MarkUserEmailAddressAsVerified(T *testing.T) {
	T.Parallel()

	T.Run("with missing user ID", func(t *testing.T) {
		t.Parallel()

		exampleInput := fakes.BuildFakeEmailAddressVerificationRequestInput()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.MarkUserEmailAddressAsVerified(ctx, "", exampleInput.Token)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		exampleUser := fakes.BuildFakeUser()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.MarkUserEmailAddressAsVerified(ctx, exampleUser.ID, "")
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}
