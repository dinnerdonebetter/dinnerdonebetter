package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidPreparationInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidPreparationInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidPreparationInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidPreparationInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidPreparationInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveValidPreparationInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid preparation instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidPreparationInstrument(ctx, ""))
	})
}
