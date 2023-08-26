package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func buildMealPlanOptionForIntegrationTest(meal *types.Meal) *types.MealPlanOption {
	exampleMealPlanOption := fakes.BuildFakeMealPlanOption()
	exampleMealPlanOption.Meal = *meal
	exampleMealPlanOption.MealScale = 1
	exampleMealPlanOption.AssignedCook = nil
	exampleMealPlanOption.AssignedDishwasher = nil
	exampleMealPlanOption.Votes = []*types.MealPlanOptionVote{}

	return exampleMealPlanOption
}

func createMealPlanOptionForTest(t *testing.T, ctx context.Context, mealPlanID string, exampleMealPlanOption *types.MealPlanOption, dbc *Querier) *types.MealPlanOption {
	t.Helper()

	// create
	dbInput := converters.ConvertMealPlanOptionToMealPlanOptionDatabaseCreationInput(exampleMealPlanOption)

	created, err := dbc.CreateMealPlanOption(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleMealPlanOption.CreatedAt = created.CreatedAt
	require.Equal(t, exampleMealPlanOption.Meal.ID, created.Meal.ID)
	exampleMealPlanOption.Meal = created.Meal
	require.Equal(t, len(exampleMealPlanOption.Votes), len(created.Votes))
	exampleMealPlanOption.Votes = created.Votes
	assert.Equal(t, exampleMealPlanOption, created)

	mealPlanOption, err := dbc.GetMealPlanOption(ctx, mealPlanID, created.BelongsToMealPlanEvent, created.ID)
	require.NoError(t, err)

	exampleMealPlanOption.CreatedAt = mealPlanOption.CreatedAt
	require.Equal(t, exampleMealPlanOption.Meal.ID, mealPlanOption.Meal.ID)
	exampleMealPlanOption.Meal = mealPlanOption.Meal
	require.Equal(t, len(exampleMealPlanOption.Votes), len(mealPlanOption.Votes))
	exampleMealPlanOption.Votes = mealPlanOption.Votes
	require.Equal(t, exampleMealPlanOption.CreatedAt, mealPlanOption.CreatedAt)
	require.Equal(t, exampleMealPlanOption.LastUpdatedAt, mealPlanOption.LastUpdatedAt)
	require.Equal(t, exampleMealPlanOption.ID, mealPlanOption.ID)

	assert.Equal(t, exampleMealPlanOption, mealPlanOption)

	return mealPlanOption
}

func TestQuerier_Integration_MealPlanOptions(t *testing.T) {
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

	recipe := createRecipeForTest(t, ctx, nil, dbc, true)
	buildMealForIntegrationTest(user.ID, recipe)
	meal := createMealForTest(t, ctx, nil, dbc)

	exampleMealPlan := buildMealPlanForIntegrationTest(user.ID, meal)
	exampleMealPlan.BelongsToHousehold = householdID
	mealPlan := createMealPlanForTest(t, ctx, exampleMealPlan, dbc)

	newMeal := createMealForTest(t, ctx, nil, dbc)
	exampleMealPlanOption := buildMealPlanOptionForIntegrationTest(newMeal)
	exampleMealPlanOption.BelongsToMealPlanEvent = mealPlan.Events[0].ID

	// create
	createdMealPlanOptions := []*types.MealPlanOption{
		mealPlan.Events[0].Options[0],
	}
	createdMealPlanOptions = append(createdMealPlanOptions, createMealPlanOptionForTest(t, ctx, mealPlan.ID, exampleMealPlanOption, dbc))

	// fetch as list
	mealPlanOptions, err := dbc.GetMealPlanOptions(ctx, mealPlan.ID, exampleMealPlanOption.BelongsToMealPlanEvent, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, mealPlanOptions)
	assert.Equal(t, len(createdMealPlanOptions), len(mealPlanOptions.Data))

	assert.NoError(t, dbc.UpdateMealPlanOption(ctx, createdMealPlanOptions[0]))

	byID, err := dbc.getMealPlanOptionByID(ctx, createdMealPlanOptions[0].ID)
	assert.NoError(t, err)
	assert.Equal(t, createdMealPlanOptions[0].ID, byID.ID)

	_, err = dbc.FinalizeMealPlanOption(ctx, mealPlan.ID, mealPlan.Events[0].ID, createdMealPlanOptions[0].ID, householdID)
	assert.NoError(t, err)

	// delete
	for _, mealPlanOption := range createdMealPlanOptions {
		assert.NoError(t, dbc.ArchiveMealPlanOption(ctx, mealPlan.ID, exampleMealPlanOption.BelongsToMealPlanEvent, mealPlanOption.ID))

		var exists bool
		exists, err = dbc.MealPlanOptionExists(ctx, mealPlan.ID, exampleMealPlanOption.BelongsToMealPlanEvent, mealPlanOption.ID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}
