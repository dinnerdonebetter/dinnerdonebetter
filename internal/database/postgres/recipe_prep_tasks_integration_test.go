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

func createRecipePrepTaskForTest(t *testing.T, ctx context.Context, exampleRecipePrepTask *types.RecipePrepTask, dbc *Querier) *types.RecipePrepTask {
	t.Helper()

	// create
	if exampleRecipePrepTask == nil {
		exampleRecipePrepTask = fakes.BuildFakeRecipePrepTask()
	}
	dbInput := converters.ConvertRecipePrepTaskToRecipePrepTaskDatabaseCreationInput(exampleRecipePrepTask)

	created, err := dbc.CreateRecipePrepTask(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleRecipePrepTask.CreatedAt = created.CreatedAt
	assert.Equal(t, exampleRecipePrepTask, created)

	recipePrepTask, err := dbc.GetRecipePrepTask(ctx, created.BelongsToRecipe, created.ID)
	exampleRecipePrepTask.CreatedAt = recipePrepTask.CreatedAt

	assert.NoError(t, err)
	assert.Equal(t, exampleRecipePrepTask, recipePrepTask)

	return created
}

func TestQuerier_Integration_RecipePrepTasks(t *testing.T) {
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

	user := createUserForTest(t, ctx, nil, dbc)
	householdID, err := dbc.GetDefaultHouseholdIDForUser(ctx, user.ID)
	require.NoError(t, err)
	require.NotEmpty(t, householdID)

	exampleRecipe := buildRecipeForTestCreation(t, ctx, user.ID, dbc)
	createdRecipe := createRecipeForTest(t, ctx, exampleRecipe, dbc, true)

	exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()
	exampleRecipePrepTask.BelongsToRecipe = createdRecipe.ID
	exampleRecipePrepTask.TaskSteps = []*types.RecipePrepTaskStep{
		{
			ID:                      identifiers.New(),
			BelongsToRecipeStep:     exampleRecipe.Steps[0].ID,
			BelongsToRecipePrepTask: exampleRecipePrepTask.ID,
			SatisfiesRecipeStep:     true,
		},
	}
	createdRecipePrepTasks := []*types.RecipePrepTask{
		exampleRecipe.PrepTasks[0],
	}

	// create
	createdRecipePrepTasks = append(createdRecipePrepTasks, createRecipePrepTaskForTest(t, ctx, exampleRecipePrepTask, dbc))

	// fetch as list
	recipePrepTasks, err := dbc.GetRecipePrepTasksForRecipe(ctx, exampleRecipe.ID)
	assert.NoError(t, err)
	assert.NotEmpty(t, recipePrepTasks)
	assert.Equal(t, len(createdRecipePrepTasks), len(recipePrepTasks))

	// delete
	for _, recipePrepTask := range createdRecipePrepTasks {
		assert.NoError(t, dbc.ArchiveRecipePrepTask(ctx, createdRecipe.ID, recipePrepTask.ID))

		var exists bool
		exists, err = dbc.RecipePrepTaskExists(ctx, createdRecipe.ID, recipePrepTask.ID)
		assert.NoError(t, err)
		assert.False(t, exists)

		var y *types.RecipePrepTask
		y, err = dbc.GetRecipePrepTask(ctx, createdRecipe.ID, recipePrepTask.ID)
		assert.Nil(t, y)
		assert.Error(t, err)
		assert.ErrorIs(t, err, sql.ErrNoRows)
	}
}
