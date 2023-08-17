package postgres

import (
	"context"
	"database/sql"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createRecipeStepVesselForTest(t *testing.T, ctx context.Context, recipeID string, exampleRecipeStepVessel *types.RecipeStepVessel, dbc *Querier) *types.RecipeStepVessel {
	t.Helper()

	// create
	if exampleRecipeStepVessel == nil {
		user := createUserForTest(t, ctx, nil, dbc)
		exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
		createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
		exampleRecipeStep := createdRecipe.Steps[0]

		exampleRecipeStepVessel = fakes.BuildFakeRecipeStepVessel()
		exampleRecipeStepVessel.BelongsToRecipeStep = exampleRecipeStep.ID
	}
	dbInput := converters.ConvertRecipeStepVesselToRecipeStepVesselDatabaseCreationInput(exampleRecipeStepVessel)

	created, err := dbc.CreateRecipeStepVessel(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipeStepVessel.CreatedAt = created.CreatedAt
	exampleRecipeStepVessel.Vessel.CreatedAt = created.Vessel.CreatedAt
	exampleRecipeStepVessel.Vessel = created.Vessel
	assert.Equal(t, exampleRecipeStepVessel, created)

	recipeStepVessel, err := dbc.GetRecipeStepVessel(ctx, recipeID, exampleRecipeStepVessel.BelongsToRecipeStep, exampleRecipeStepVessel.ID)

	exampleRecipeStepVessel.CreatedAt = recipeStepVessel.CreatedAt
	exampleRecipeStepVessel.Vessel.CreatedAt = recipeStepVessel.Vessel.CreatedAt
	exampleRecipeStepVessel.Vessel = recipeStepVessel.Vessel

	require.Equal(t, exampleRecipeStepVessel, recipeStepVessel)

	assert.NoError(t, err)
	assert.Equal(t, recipeStepVessel, exampleRecipeStepVessel)

	return created
}

func TestQuerier_Integration_RecipeStepVessels(t *testing.T) {
	if !runningContainerTests {
		t.SkipNow()
	}

	ctx := context.Background()
	dbc, container := buildDatabaseClientForTest(t, ctx)

	defer func(t *testing.T) {
		t.Helper()
		assert.NoError(t, container.Terminate(ctx))
	}(t)

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)
	exampleRecipeStep := createdRecipe.Steps[0]

	validVessel := createValidVesselForTest(t, ctx, nil, dbc)
	exampleRecipeStepVessel := fakes.BuildFakeRecipeStepVessel()
	exampleRecipeStepVessel.Vessel = validVessel
	exampleRecipeStepVessel.BelongsToRecipeStep = exampleRecipeStep.ID
	createdRecipeStepVessels := []*types.RecipeStepVessel{
		exampleRecipeStep.Vessels[0],
	}

	// create
	createdRecipeStepVessels = append(createdRecipeStepVessels, createRecipeStepVesselForTest(t, ctx, exampleRecipe.ID, exampleRecipeStepVessel, dbc))

	// create more
	for i := 0; i < exampleQuantity; i++ {
		validVessel = createValidVesselForTest(t, ctx, nil, dbc)
		input := fakes.BuildFakeRecipeStepVessel()
		input.Vessel = validVessel
		input.BelongsToRecipeStep = exampleRecipeStep.ID
		createdRecipeStepVessels = append(createdRecipeStepVessels, createRecipeStepVesselForTest(t, ctx, exampleRecipe.ID, input, dbc))
	}

	// fetch as list
	recipeStepVessels, err := dbc.GetRecipeStepVessels(ctx, exampleRecipe.ID, createdRecipeStepVessels[0].BelongsToRecipeStep, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipeStepVessels.Data)
	assert.Equal(t, len(createdRecipeStepVessels), len(recipeStepVessels.Data))

	// delete
	for _, recipeStepVessel := range createdRecipeStepVessels {
		assert.NoError(t, dbc.ArchiveRecipeStepVessel(ctx, exampleRecipeStep.ID, recipeStepVessel.ID))

		var exists bool
		exists, err = dbc.RecipeStepVesselExists(ctx, exampleRecipe.ID, recipeStepVessel.BelongsToRecipeStep, recipeStepVessel.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipeStepVessel
		y, err = dbc.GetRecipeStepVessel(ctx, exampleRecipe.ID, recipeStepVessel.BelongsToRecipeStep, recipeStepVessel.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
