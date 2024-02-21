package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientMeasurementUnit_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnit{
			MaximumAllowableQuantity: pointers.Pointer(float32(3.21)),
		}
		input := &ValidIngredientMeasurementUnitUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.MaximumAllowableQuantity = pointers.Pointer(float32(1.23))

		x.Update(input)
	})
}

func TestValidIngredientMeasurementUnitCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientMeasurementUnitCreationRequestInput{
			Notes:                    t.Name(),
			ValidMeasurementUnitID:   t.Name(),
			ValidIngredientID:        t.Name(),
			MinimumAllowableQuantity: fake.Float32(),
			MaximumAllowableQuantity: pointers.Pointer(fake.Float32()),
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
			ID:                       t.Name(),
			ValidMeasurementUnitID:   t.Name(),
			ValidIngredientID:        t.Name(),
			MinimumAllowableQuantity: fake.Float32(),
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
			Notes:                    pointers.Pointer(t.Name()),
			ValidMeasurementUnitID:   pointers.Pointer(t.Name()),
			ValidIngredientID:        pointers.Pointer(t.Name()),
			MinimumAllowableQuantity: pointers.Pointer(fake.Float32()),
			MaximumAllowableQuantity: pointers.Pointer(fake.Float32()),
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
