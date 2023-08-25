package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQuerier_RecipePrepTaskExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, _ := buildTestClient(t)

		actual, err := c.RecipePrepTaskExists(ctx, "", exampleRecipePrepTask.ID)
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipePrepTaskID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.GetRecipePrepTask(ctx, "", exampleRecipePrepTaskID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateRecipePrepTask(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_GetRecipePrepTasksForRecipe(T *testing.T) {
	T.Parallel()

	T.Run("with invalid recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetRecipePrepTasksForRecipe(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		assert.Error(t, c.UpdateRecipePrepTask(ctx, nil))

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_ArchiveRecipePrepTask(T *testing.T) {
	T.Parallel()

	T.Run("with missing recipe ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleRecipePrepTask := fakes.BuildFakeRecipePrepTask()

		c, db := buildTestClient(t)

		assert.Error(t, c.ArchiveRecipePrepTask(ctx, "", exampleRecipePrepTask.ID))

		mock.AssertExpectationsForObjects(t, db)
	})
}
