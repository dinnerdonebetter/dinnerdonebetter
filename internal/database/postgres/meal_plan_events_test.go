package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestQuerier_MealPlanEventExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanEventID := fakes.BuildFakeID()

		c, db := buildTestClient(t)

		actual, err := c.MealPlanEventExists(ctx, "", exampleMealPlanEventID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with invalid meal plan event ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleMealPlanID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.MealPlanEventExists(ctx, exampleMealPlanID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanEvent(ctx, "", exampleMealPlanEventID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with invalid meal plan event ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanEvent(ctx, exampleMealPlan.ID, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_getMealPlanEventsForMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.getMealPlanEventsForMealPlan(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetMealPlanEvents(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMealPlanEvents(ctx, "", nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_createMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		db.ExpectBegin()
		tx, err := c.db.BeginTx(ctx, nil)
		require.NoError(t, err)
		require.NotNil(t, tx)

		actual, err := c.createMealPlanEvent(ctx, tx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateMealPlanEvent(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.UpdateMealPlanEvent(ctx, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMealPlanEvent(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal plan event ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlan := fakes.BuildFakeMealPlan()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanEvent(ctx, "", exampleMealPlan.ID))
	})

	T.Run("with invalid meal plan ID", func(t *testing.T) {
		t.Parallel()

		exampleMealPlanEventID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMealPlanEvent(ctx, exampleMealPlanEventID, ""))
	})
}
