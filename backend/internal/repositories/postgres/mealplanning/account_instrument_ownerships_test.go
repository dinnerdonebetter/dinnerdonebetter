package mealplanning

import (
	"context"
	"database/sql"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/converters"
	"github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAccountInstrumentOwnershipForTest(t *testing.T, ctx context.Context, exampleAccountInstrumentOwnership *types.AccountInstrumentOwnership, dbc *repository) *types.AccountInstrumentOwnership {
	t.Helper()

	// create
	if exampleAccountInstrumentOwnership == nil {
		exampleAccountInstrumentOwnership = fakes.BuildFakeAccountInstrumentOwnership()
	}
	dbInput := converters.ConvertAccountInstrumentOwnershipToAccountInstrumentOwnershipDatabaseCreationInput(exampleAccountInstrumentOwnership)

	created, err := dbc.CreateAccountInstrumentOwnership(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleAccountInstrumentOwnership.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleAccountInstrumentOwnership.Instrument.ID, created.Instrument.ID)
	exampleAccountInstrumentOwnership.Instrument = created.Instrument
	assert.Equal(t, exampleAccountInstrumentOwnership, created)

	accountInstrumentOwnership, err := dbc.GetAccountInstrumentOwnership(ctx, created.ID, created.BelongsToAccount)
	exampleAccountInstrumentOwnership.CreatedAt = accountInstrumentOwnership.CreatedAt
	assert.Equal(t, exampleAccountInstrumentOwnership.Instrument.ID, accountInstrumentOwnership.Instrument.ID)
	exampleAccountInstrumentOwnership.Instrument = accountInstrumentOwnership.Instrument

	assert.NoError(t, err)
	assert.Equal(t, accountInstrumentOwnership, exampleAccountInstrumentOwnership)

	return created
}

func TestQuerier_Integration_AccountInstrumentOwnerships(t *testing.T) {
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

	user := pgtesting.CreateUserForTest(t, nil, dbc.db)
	account := pgtesting.CreateAccountForTest(t, nil, user.ID, dbc.db)

	instrument := createValidInstrumentForTest(t, ctx, nil, dbc)

	exampleAccountInstrumentOwnership := fakes.BuildFakeAccountInstrumentOwnership()
	exampleAccountInstrumentOwnership.BelongsToAccount = account.ID
	exampleAccountInstrumentOwnership.Instrument = *instrument
	createdAccountInstrumentOwnerships := []*types.AccountInstrumentOwnership{}

	// create
	createdAccountInstrumentOwnerships = append(createdAccountInstrumentOwnerships, createAccountInstrumentOwnershipForTest(t, ctx, exampleAccountInstrumentOwnership, dbc))

	// update
	assert.NoError(t, dbc.UpdateAccountInstrumentOwnership(ctx, createdAccountInstrumentOwnerships[0]))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		newInstrument := createValidInstrumentForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeAccountInstrumentOwnership()
		input.BelongsToAccount = account.ID
		input.Instrument = *newInstrument
		createdAccountInstrumentOwnerships = append(createdAccountInstrumentOwnerships, createAccountInstrumentOwnershipForTest(t, ctx, input, dbc))
	}

	// fetch as list
	accountInstrumentOwnerships, err := dbc.GetAccountInstrumentOwnerships(ctx, account.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, accountInstrumentOwnerships.Data)
	assert.Equal(t, len(createdAccountInstrumentOwnerships), len(accountInstrumentOwnerships.Data))

	// delete
	for _, accountInstrumentOwnership := range createdAccountInstrumentOwnerships {
		assert.NoError(t, dbc.ArchiveAccountInstrumentOwnership(ctx, accountInstrumentOwnership.ID, account.ID))

		var exists bool
		exists, err = dbc.AccountInstrumentOwnershipExists(ctx, accountInstrumentOwnership.ID, account.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.AccountInstrumentOwnership
		y, err = dbc.GetAccountInstrumentOwnership(ctx, accountInstrumentOwnership.ID, account.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}

func TestQuerier_AccountInstrumentOwnershipExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid account instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := fakes.BuildFakeID()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.AccountInstrumentOwnershipExists(ctx, "", exampleAccountID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid account instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := fakes.BuildFakeID()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetAccountInstrumentOwnership(ctx, "", exampleAccountID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreateAccountInstrumentOwnership(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.UpdateAccountInstrumentOwnership(ctx, nil))
	})
}

func TestQuerier_ArchiveAccountInstrumentOwnership(T *testing.T) {
	T.Parallel()

	T.Run("with invalid account instrument ownership ID", func(t *testing.T) {
		t.Parallel()

		exampleAccountID := fakes.BuildFakeID()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveAccountInstrumentOwnership(ctx, "", exampleAccountID))
	})
}
