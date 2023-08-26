package postgres

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/pkg/types/fakes"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_UserIngredientPreferenceExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIngredientPreferenceExists(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with invalid user ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserIngredientPreferenceID := fakes.BuildFakeID()

		c, _ := buildTestClient(t)

		actual, err := c.UserIngredientPreferenceExists(ctx, exampleUserIngredientPreferenceID, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		actual, err := c.GetUserIngredientPreference(ctx, "", exampleUserID)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateUserIngredientPreference(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateUserIngredientPreference(ctx, nil))
	})
}

func TestQuerier_ArchiveUserIngredientPreference(T *testing.T) {
	T.Parallel()

	T.Run("with invalid user ingredient preference ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleUserID := fakes.BuildFakeID()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveUserIngredientPreference(ctx, "", exampleUserID))
	})
}
