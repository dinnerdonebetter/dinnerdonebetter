package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidPreparationCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{
			Name:                  t.Name(),
			Description:           t.Name(),
			IconPath:              t.Name(),
			PastTense:             t.Name(),
			YieldsNothing:         fake.Bool(),
			RestrictToIngredients: fake.Bool(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{
			Name:                  pointers.Pointer(t.Name()),
			Description:           pointers.Pointer(t.Name()),
			IconPath:              pointers.Pointer(t.Name()),
			PastTense:             pointers.Pointer(t.Name()),
			YieldsNothing:         pointers.Pointer(fake.Bool()),
			RestrictToIngredients: pointers.Pointer(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparation{
			MaximumInstrumentCount: pointers.Pointer(int32(0)),
			MaximumIngredientCount: pointers.Pointer(int32(0)),
			MaximumVesselCount:     pointers.Pointer(int32(0)),
		}
		input := &ValidPreparationUpdateRequestInput{}

		fake.Struct(&input)
		input.YieldsNothing = pointers.Pointer(true)
		input.RestrictToIngredients = pointers.Pointer(true)
		input.MaximumIngredientCount = pointers.Pointer(int32(1))
		input.MaximumInstrumentCount = pointers.Pointer(int32(1))
		input.MaximumVesselCount = pointers.Pointer(int32(1))
		input.TemperatureRequired = pointers.Pointer(true)
		input.TimeEstimateRequired = pointers.Pointer(true)
		input.OnlyForVessels = pointers.Pointer(true)
		input.ConsumesVessel = pointers.Pointer(true)
		input.ConditionExpressionRequired = pointers.Pointer(true)

		x.Update(input)
	})
}

func TestValidPreparationCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationCreationRequestInput{
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationDatabaseCreationInput{
			ID:   t.Name(),
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationUpdateRequestInput{
			Name: pointers.Pointer(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
