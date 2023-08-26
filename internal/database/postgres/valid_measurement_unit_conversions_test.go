package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerier_ValidMeasurementUnitConversionExists(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidMeasurementUnitConversionExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})
}

func TestQuerier_GetValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidMeasurementUnitConversion(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_CreateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidMeasurementUnitConversion(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})
}

func TestQuerier_UpdateValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidMeasurementUnitConversion(ctx, nil))
	})
}

func TestQuerier_ArchiveValidMeasurementUnitConversion(T *testing.T) {
	T.Parallel()

	T.Run("with invalid valid measurement conversion ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidMeasurementUnitConversion(ctx, ""))
	})
}
