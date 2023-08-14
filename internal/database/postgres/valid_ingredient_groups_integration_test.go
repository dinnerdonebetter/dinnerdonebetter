package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/identifiers"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createValidIngredientGroupForTest(t *testing.T, ctx context.Context, exampleValidIngredientGroup *types.ValidIngredientGroup, dbc *Querier) *types.ValidIngredientGroup {
	t.Helper()

	// create
	if exampleValidIngredientGroup == nil {
		exampleValidIngredientGroup = fakes.BuildFakeValidIngredientGroup()
	}
	dbInput := converters.ConvertValidIngredientGroupToValidIngredientGroupDatabaseCreationInput(exampleValidIngredientGroup)

	created, err := dbc.CreateValidIngredientGroup(ctx, dbInput)
	require.NoError(t, err)
	require.NotNil(t, created)
	exampleValidIngredientGroup.CreatedAt = created.CreatedAt
	for i := range exampleValidIngredientGroup.Members {
		exampleValidIngredientGroup.Members[i].CreatedAt = created.Members[i].CreatedAt
		exampleValidIngredientGroup.Members[i].ValidIngredient = created.Members[i].ValidIngredient
	}
	assert.Equal(t, exampleValidIngredientGroup, created)

	validIngredientGroup, err := dbc.GetValidIngredientGroup(ctx, created.ID)
	exampleValidIngredientGroup.CreatedAt = validIngredientGroup.CreatedAt
	for i := range exampleValidIngredientGroup.Members {
		exampleValidIngredientGroup.Members[i].CreatedAt = validIngredientGroup.Members[i].CreatedAt
		exampleValidIngredientGroup.Members[i].ValidIngredient = validIngredientGroup.Members[i].ValidIngredient
	}

	assert.NoError(t, err)
	assert.Equal(t, validIngredientGroup, exampleValidIngredientGroup)

	return created
}

func TestQuerier_Integration_ValidIngredientGroups(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
	exampleValidIngredientGroup.Members = []*types.ValidIngredientGroupMember{}
	for i := 0; i < exampleQuantity; i++ {
		exampleValidIngredientGroup.Members = append(exampleValidIngredientGroup.Members, &types.ValidIngredientGroupMember{
			ID:              identifiers.New(),
			Notes:           "",
			BelongsToGroup:  exampleValidIngredientGroup.ID,
			ValidIngredient: *createValidIngredientForTest(t, ctx, nil, dbc),
		})
	}
	createdValidIngredientGroups := []*types.ValidIngredientGroup{}

	// create
	createdValidIngredientGroups = append(createdValidIngredientGroups, createValidIngredientGroupForTest(t, ctx, exampleValidIngredientGroup, dbc))

	// update
	updatedValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
	updatedValidIngredientGroup.ID = createdValidIngredientGroups[0].ID
	assert.NoError(t, dbc.UpdateValidIngredientGroup(ctx, updatedValidIngredientGroup))

	// fetch as list
	validIngredientGroups, err := dbc.GetValidIngredientGroups(ctx, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, validIngredientGroups.Data)
	assert.Equal(t, len(createdValidIngredientGroups), len(validIngredientGroups.Data))

	// fetch via name search
	byName, err := dbc.SearchForValidIngredientGroups(ctx, updatedValidIngredientGroup.Name, nil)
	assert.NoError(t, err)
	assert.Equal(t, validIngredientGroups.Data, byName)

	// delete
	for _, validIngredientGroup := range createdValidIngredientGroups {
		assert.NoError(t, dbc.ArchiveValidIngredientGroup(ctx, validIngredientGroup.ID))

		var exists bool
		exists, err = dbc.ValidIngredientGroupExists(ctx, validIngredientGroup.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.ValidIngredientGroup
		y, err = dbc.GetValidIngredientGroup(ctx, validIngredientGroup.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
