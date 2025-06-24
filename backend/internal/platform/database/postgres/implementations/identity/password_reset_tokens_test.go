package identity

import (
	"context"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/identity/fakes"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPasswordResetTokenForTest(t *testing.T, ctx context.Context, examplePasswordResetToken *types.PasswordResetToken, dbc *Querier) *types.PasswordResetToken {
	t.Helper()

	// create
	if examplePasswordResetToken == nil {
		examplePasswordResetToken = fakes.BuildFakePasswordResetToken()
	}
	dbInput := converters.ConvertPasswordResetTokenToPasswordResetTokenDatabaseCreationInput(examplePasswordResetToken)

	created, err := dbc.CreatePasswordResetToken(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	examplePasswordResetToken.CreatedAt = created.CreatedAt
	assert.Equal(t, examplePasswordResetToken, created)

	passwordResetToken, err := dbc.GetPasswordResetTokenByToken(ctx, created.Token)
	examplePasswordResetToken.CreatedAt = passwordResetToken.CreatedAt
	examplePasswordResetToken.ExpiresAt = passwordResetToken.ExpiresAt

	assert.NoError(t, err)
	assert.Equal(t, passwordResetToken, examplePasswordResetToken)

	return created
}

func TestQuerier_Integration_PasswordResetTokens(t *testing.T) {
	if !pgtesting.RunContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t)

	databaseURI, err := container.ConnectionString(ctx)
	require.NoError(t, err)
	require.NotEmpty(t, databaseURI)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)

	examplePasswordResetToken := fakes.BuildFakePasswordResetToken()
	examplePasswordResetToken.BelongsToUser = user.ID
	createdPasswordResetTokens := []*types.PasswordResetToken{}

	// create
	createdPasswordResetTokens = append(createdPasswordResetTokens, createPasswordResetTokenForTest(t, ctx, examplePasswordResetToken, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakePasswordResetToken()
		input.BelongsToUser = user.ID
		createdPasswordResetTokens = append(createdPasswordResetTokens, createPasswordResetTokenForTest(t, ctx, input, dbc))
	}

	// delete
	for _, passwordResetToken := range createdPasswordResetTokens {
		assert.NoError(t, dbc.RedeemPasswordResetToken(ctx, passwordResetToken.ID))
	}
}

func TestSQLQuerier_GetPasswordResetTokenByToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing token", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetPasswordResetTokenByToken(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_CreatePasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreatePasswordResetToken(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestSQLQuerier_RedeemPasswordResetToken(T *testing.T) {
	T.Parallel()

	T.Run("with missing ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual := c.RedeemPasswordResetToken(ctx, "")
		assert.Error(t, actual)
	})
}
