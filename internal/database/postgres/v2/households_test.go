package v2

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/jinzhu/copier"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createHouseholdForTest(t *testing.T, ctx context.Context, exampleHousehold *types.Household, dbc *DatabaseClient) *types.Household {
	t.Helper()

	// create
	if exampleHousehold == nil {
		exampleHousehold = fakes.BuildFakeHousehold()
	}

	user := createUserForTest(t, ctx, nil, dbc)
	exampleHousehold.BelongsToUser = user.ID

	exampleHousehold.Members = nil

	var x Household
	require.NoError(t, copier.Copy(&x, exampleHousehold))

	created, err := dbc.CreateHousehold(ctx, &x)
	assert.NoError(t, err)
	assert.Equal(t, exampleHousehold, created)

	household, err := dbc.GetHousehold(ctx, created.ID)
	exampleHousehold.CreatedAt = household.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, household, exampleHousehold)

	return created
}

func TestDatabaseClient_Households(t *testing.T) {
	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleHousehold := fakes.BuildFakeHousehold()
	createdHouseholds := []*types.Household{}

	// create
	createdHouseholds = append(createdHouseholds, createHouseholdForTest(t, ctx, exampleHousehold, dbc))

	// update
	updatedHousehold := fakes.BuildFakeHousehold()
	updatedHousehold.ID = createdHouseholds[0].ID
	var x Household
	require.NoError(t, copier.Copy(&x, updatedHousehold))
	assert.NoError(t, dbc.UpdateHousehold(ctx, updatedHousehold))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		input := fakes.BuildFakeHousehold()
		input.Name = fmt.Sprintf("%s %d", exampleHousehold.Name, i)
		createdHouseholds = append(createdHouseholds, createHouseholdForTest(t, ctx, input, dbc))
	}

	// fetch as list
	households, err := dbc.GetHouseholds(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, households.Data)
	assert.Equal(t, len(createdHouseholds), len(households.Data))

	// fetch as list of IDs
	householdIDs := []string{}
	for _, household := range createdHouseholds {
		householdIDs = append(householdIDs, household.ID)
	}

	// delete
	for _, household := range createdHouseholds {
		assert.NoError(t, dbc.ArchiveHousehold(ctx, household.ID))

		var y *types.Household
		y, err = dbc.GetHousehold(ctx, household.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
