package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidIngredientStateExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidIngredientStateExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidIngredientState(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidIngredientStates(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidIngredientStates(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidIngredientState(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidIngredientState(ctx, nil))
	})
}

func TestQuerier_ArchiveValidIngredientState(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid ingredient state ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidIngredientState(ctx, ""))
	})
}

func TestQuerier_MarkValidIngredientStateAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidIngredientStateAsIndexed(ctx, ""))
	})
}
