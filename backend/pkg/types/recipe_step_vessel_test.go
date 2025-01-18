package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepVessel_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVessel{
			Quantity: Uint16RangeWithOptionalMax{
				Max: pointer.To(uint16(1234)),
				Min: 1234,
			},
			Vessel: &ValidVessel{},
		}
		input := &RecipeStepVesselUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.UnavailableAfterStep = pointer.To(true)
		input.Quantity = Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: pointer.To(uint16(1)),
			Max: pointer.To(uint16(1)),
		}
		input.VesselID = pointer.To(t.Name())

		x.Update(input)
	})
}

func TestRecipeStepVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			Name:                t.Name(),
			RecipeStepProductID: pointer.To(t.Name()),
			Notes:               t.Name(),
			Quantity: Uint16RangeWithOptionalMax{
				Max: pointer.To(fake.Uint16()),
				Min: fake.Uint16(),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepVesselDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselDatabaseCreationInput{
			ID:                  t.Name(),
			BelongsToRecipeStep: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepVesselUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselUpdateRequestInput{
			Name:                pointer.To(t.Name()),
			BelongsToRecipeStep: pointer.To(t.Name()),
			RecipeStepProductID: pointer.To(t.Name()),
			Notes:               pointer.To(t.Name()),
			Quantity: Uint16RangeWithOptionalMaxUpdateRequestInput{
				Max: pointer.To(fake.Uint16()),
				Min: pointer.To(fake.Uint16()),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
