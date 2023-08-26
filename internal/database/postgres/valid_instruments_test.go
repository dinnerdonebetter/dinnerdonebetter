package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidInstruments(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidInstrument(ctx, nil))
	})
}

func TestQuerier_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidInstrument(ctx, ""))
	})
}

func TestQuerier_MarkValidInstrumentAsIndexed(T *testing.T) {
	T.Parallel()

	T.Run("with invalid ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.MarkValidInstrumentAsIndexed(ctx, ""))
	})
}
