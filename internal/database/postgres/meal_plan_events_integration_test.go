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

func buildMealPlanEventForIntegrationTest(meal *types.Meal) *types.MealPlanEvent {
	exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()

	exampleMealPlanOption := buildMealPlanOptionForIntegrationTest(meal)
	exampleMealPlanOption.BelongsToMealPlanEvent = exampleMealPlanEvent.ID
	exampleMealPlanEvent.Options = []*types.MealPlanOption{
		exampleMealPlanOption,
	}

	return exampleMealPlanEvent
}

func createMealPlanEventForTest(t *testing.T, ctx context.Context, exampleMealPlanEvent *types.MealPlanEvent, dbc *Querier) *types.MealPlanEvent {
	t.Helper()

	// create
	dbInput := converters.ConvertMealPlanEventToMealPlanEventDatabaseCreationInput(exampleMealPlanEvent)

	created, err := dbc.CreateMealPlanEvent(ctx, dbInput)
	assert.NoError(t, err)
	require.NotNil(t, created)

	exampleMealPlanEvent.CreatedAt = created.CreatedAt
	for i, option := range created.Options {
		exampleMealPlanEvent.Options[i].CreatedAt = option.CreatedAt
		require.Equal(t, exampleMealPlanEvent.Options[i].Meal.ID, option.Meal.ID)
		exampleMealPlanEvent.Options[i].Meal = option.Meal
	}
	assert.Equal(t, exampleMealPlanEvent, created)

	mealPlanEvent, err := dbc.GetMealPlanEvent(ctx, created.BelongsToMealPlan, created.ID)
	require.NoError(t, err)

	exampleMealPlanEvent.CreatedAt = mealPlanEvent.CreatedAt
	exampleMealPlanEvent.StartsAt = mealPlanEvent.StartsAt
	exampleMealPlanEvent.EndsAt = mealPlanEvent.EndsAt
	for i, option := range mealPlanEvent.Options {
		exampleMealPlanEvent.Options[i].CreatedAt = option.CreatedAt
		require.Equal(t, exampleMealPlanEvent.Options[i].Meal.ID, option.Meal.ID)
		exampleMealPlanEvent.Options[i].Meal = option.Meal
		exampleMealPlanEvent.Options[i].Chosen = option.Chosen
	}
	require.Equal(t, exampleMealPlanEvent.CreatedAt, mealPlanEvent.CreatedAt)
	require.Equal(t, exampleMealPlanEvent.LastUpdatedAt, mealPlanEvent.LastUpdatedAt)
	require.Equal(t, exampleMealPlanEvent.ID, mealPlanEvent.ID)

	assert.Equal(t, exampleMealPlanEvent, mealPlanEvent)

	return mealPlanEvent
}

func TestQuerier_Integration_MealPlanEvents(t *testing.T) {
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
	newEvent := buildMealPlanEventForIntegrationTest(newMeal)
	newEvent.BelongsToMealPlan = mealPlan.ID

	exampleMealPlanEvent := fakes.BuildFakeMealPlanEvent()
	exampleMealPlanEvent.BelongsToMealPlan = mealPlan.ID
	exampleMealPlanEvent.Options = newEvent.Options
	for i := range exampleMealPlanEvent.Options {
		exampleMealPlanEvent.Options[i].BelongsToMealPlanEvent = exampleMealPlanEvent.ID
	}

	// create
	createdMealPlanEvents := []*types.MealPlanEvent{
		mealPlan.Events[0],
	}
	createdMealPlanEvents = append(createdMealPlanEvents, createMealPlanEventForTest(t, ctx, exampleMealPlanEvent, dbc))

	// fetch as list
	mealPlanEvents, err := dbc.GetMealPlanEvents(ctx, mealPlan.ID, nil)
	assert.NoError(t, err)
	assert.NotEmpty(t, mealPlanEvents)
	assert.Equal(t, len(createdMealPlanEvents), len(mealPlanEvents.Data))

	assert.NoError(t, dbc.UpdateMealPlanEvent(ctx, createdMealPlanEvents[0]))

	_, err = dbc.MealPlanEventIsEligibleForVoting(ctx, mealPlan.ID, createdMealPlanEvents[0].ID)
	assert.NoError(t, err)

	// delete
	for _, mealPlanEvent := range createdMealPlanEvents {
		assert.NoError(t, dbc.ArchiveMealPlanEvent(ctx, mealPlan.ID, mealPlanEvent.ID))

		var exists bool
		exists, err = dbc.MealPlanEventExists(ctx, mealPlanEvent.ID, householdID)
		assert.NoError(t, err)
		assert.False(t, exists)
	}
}
