package recipeenums

import (
	"context"
	"database/sql"
	"testing"

	types "github.com/dinnerdonebetter/backend/internal/domain/recipeenums"
	"github.com/dinnerdonebetter/backend/internal/domain/recipeenums/fakes"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	pgtesting "github.com/dinnerdonebetter/backend/internal/platform/database/postgres/testing"
	"github.com/dinnerdonebetter/backend/internal/platform/identifiers"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestQuerier_Integration_ValidIngredientGroups(t *testing.T) {
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

	exampleValidIngredientGroup := fakes.BuildFakeValidIngredientGroup()
	exampleValidIngredientGroup.Members = []*types.ValidIngredientGroupMember{}
	for i := 0; i < exampleQuantity; i++ {
		exampleValidIngredientGroup.Members = append(exampleValidIngredientGroup.Members, &types.ValidIngredientGroupMember{
			ID:              identifiers.New(),
			BelongsToGroup:  exampleValidIngredientGroup.ID,
			ValidIngredient: *CreateValidIngredientForTest(t, ctx, nil, dbc),
		})
	}
	createdValidIngredientGroups := []*types.ValidIngredientGroup{}

	// create
	createdValidIngredientGroups = append(createdValidIngredientGroups, CreateValidIngredientGroupForTest(t, ctx, exampleValidIngredientGroup, dbc))

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

func TestQuerier_ValidIngredientGroupExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c := buildInertClientForTest(t)

		actual, err := c.ValidIngredientGroupExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.GetValidIngredientGroup(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredientGroups(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)
		filter := filtering.DefaultQueryFilter()

		actual, err := c.SearchForValidIngredientGroups(ctx, "", filter)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		actual, err := c.CreateValidIngredientGroup(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.UpdateValidIngredientGroup(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientGroup(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c := buildInertClientForTest(t)

		assert.Error(t, c.ArchiveValidIngredientGroup(ctx, ""))
	})
}
