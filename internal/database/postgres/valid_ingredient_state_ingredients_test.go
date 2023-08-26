package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidIngredientStateIngredientExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientStateIngredientExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientStateIngredient(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientStateIngredient(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientStateIngredient(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientStateIngredient(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient preparation ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientStateIngredient(ctx, ""))
	})
}
