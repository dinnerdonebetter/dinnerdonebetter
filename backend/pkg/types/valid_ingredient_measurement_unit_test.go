package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientMeasurementUnit_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnit{
			AllowableQuantity: Float32RangeWithOptionalMax{Max: pointer.To(float32(3.21))},
		}
		input := &ValidIngredientMeasurementUnitUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.AllowableQuantity.Max = pointer.To(float32(1.23))

		x.Update(input)
	})
}

func TestValidIngredientMeasurementUnitCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{
			Notes:                  t.Name(),
			ValidMeasurementUnitID: t.Name(),
			ValidIngredientID:      t.Name(),
			AllowableQuantity: Float32RangeWithOptionalMax{
				Max: pointer.To(fake.Float32()),
				Min: fake.Float32(),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientMeasurementUnitDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitDatabaseCreationInput{
			ID:                     t.Name(),
			ValidMeasurementUnitID: t.Name(),
			ValidIngredientID:      t.Name(),
			AllowableQuantity: Float32RangeWithOptionalMax{
				Min: fake.Float32(),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientMeasurementUnitUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitUpdateRequestInput{
			Notes:                  pointer.To(t.Name()),
			ValidMeasurementUnitID: pointer.To(t.Name()),
			ValidIngredientID:      pointer.To(t.Name()),
			AllowableQuantity: Float32RangeWithOptionalMaxUpdateRequestInput{
				Min: pointer.To(fake.Float32()),
				Max: pointer.To(fake.Float32()),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
