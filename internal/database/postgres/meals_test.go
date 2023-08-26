package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/converters"
	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_MealExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.MealExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetMeal(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetMeal(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateMeal(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		actual, err := c.CreateMeal(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_CreateMealRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with missing meal ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		exampleInput := converters.ConvertMealComponentToMealComponentDatabaseCreationInput(exampleMeal.Components[0])

		err := c.CreateMealComponent(ctx, c.db, "", exampleInput)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})

	T.Run("with missing input", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, db := buildTestClient(t)

		err := c.CreateMealComponent(ctx, c.db, exampleMeal.ID, nil)
		assert.Error(t, err)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveMeal(T *testing.T) {
	T.Parallel()

	T.Run("with invalid meal ID", func(t *testing.T) {
		t.Parallel()

		exampleHouseholdID := fakes.BuildFakeID()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMeal(ctx, "", exampleHouseholdID))
	})

	T.Run("with invalid household ID", func(t *testing.T) {
		t.Parallel()

		exampleMeal := fakes.BuildFakeMeal()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveMeal(ctx, exampleMeal.ID, ""))
	})
}

func TestQuerier_MarkMealAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkMealAsIndexed(ctx, ""))
	})
}
