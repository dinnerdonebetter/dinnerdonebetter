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
	"github.com/stretchr/testify/require"
)

func createHouseholdForTest(t *testing.T, ctx context.Context, exampleHousehold *types.Household, dbc *Querier) *types.Household {
	t.Helper()

	// create
	if exampleHousehold == nil {
		exampleUser := createUserForTest(t, ctx, nil, dbc)
		exampleHousehold = fakes.BuildFakeHousehold()
		exampleHousehold.BelongsToUser = exampleUser.ID
	}
	exampleHousehold.PaymentProcessorCustomerID = ""
	exampleHousehold.Members = nil
	dbInput := converters.ConvertHouseholdToHouseholdDatabaseCreationInput(exampleHousehold)

	created, err := dbc.CreateHousehold(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)
	exampleHousehold.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleHousehold, created)

	household, err := dbc.GetHousehold(ctx, created.ID)
	exampleHousehold.CreatedAt = household.CreatedAt
	exampleHousehold.Members = household.Members

	assert.NoError(t, err)
	assert.Equal(t, household, exampleHousehold)

	return created
}

func TestQuerier_Integration_Households(t *testing.T) {
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

	exampleUser := createUserForTest(t, ctx, nil, dbc)

	exampleHousehold := fakes.BuildFakeHousehold()
	exampleHousehold.Members = nil
	exampleHousehold.BelongsToUser = exampleUser.ID
	exampleHousehold.PaymentProcessorCustomerID = ""
	createdHouseholds := []*types.Household{}

	// create
	createdHouseholds = append(createdHouseholds, createHouseholdForTest(t, ctx, exampleHousehold, dbc))

	// update
	updatedHousehold := fakes.BuildFakeHousehold()
	updatedHousehold.ID = createdHouseholds[0].ID
	updatedHousehold.BelongsToUser = createdHouseholds[0].BelongsToUser
	assert.NoError(t, dbc.UpdateHousehold(ctx, updatedHousehold))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeHousehold()
		input.BelongsToUser = exampleUser.ID
		input.Name = fmt.Sprintf("%s %d", updatedHousehold.Name, i)
		createdHouseholds = append(createdHouseholds, createHouseholdForTest(t, ctx, input, dbc))
	}

	// fetch as list
	households, err := dbc.GetHouseholds(ctx, exampleUser.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, households.Data)
	assert.GreaterOrEqual(t, len(households.Data), len(createdHouseholds))

	// delete
	for _, household := range createdHouseholds {
		assert.NoError(t, dbc.ArchiveHousehold(ctx, household.ID, exampleUser.ID))

		var y *types.Household
		y, err = dbc.GetHousehold(ctx, household.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}

	assert.NoError(t, dbc.ArchiveUser(ctx, exampleUser.ID))
}
