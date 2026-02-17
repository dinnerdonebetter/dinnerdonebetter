package mealplanning

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/types"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepVessel_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVessel{
			Quantity: types.Uint16RangeWithOptionalMax{
				Max: new(uint16(1234)),
				Min: 1234,
			},
			Vessel: &ValidVessel{},
		}
		input := &RecipeStepVesselUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.UnavailableAfterStep = new(true)
		input.Quantity = types.Uint16RangeWithOptionalMaxUpdateRequestInput{
			Min: new(uint16(1)),
			Max: new(uint16(1)),
		}
		input.VesselID = new(t.Name())

		x.Update(input)
	})
}

func TestRecipeStepVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			ValidPreparationVesselID: new(t.Name()),
			Name:                     t.Name(),
			RecipeStepProductID:      new(t.Name()),
			Notes:                    t.Name(),
			Quantity: types.Uint16RangeWithOptionalMax{
				Max: new(fake.Uint16()),
				Min: fake.Uint16(),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("recipe step product does not require bridge IDs", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			ProductOfRecipeStepIndex:        new(uint64(0)),
			ProductOfRecipeStepProductIndex: new(uint64(0)),
			Name:                            t.Name(),
			Notes:                           t.Name(),
			Quantity: types.Uint16RangeWithOptionalMax{
				Max: new(fake.Uint16()),
				Min: fake.Uint16(),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{}

		actual := x.ValidateWithContext(t.Context())
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

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselDatabaseCreationInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}

func TestRecipeStepVesselUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselUpdateRequestInput{
			Name:                new(t.Name()),
			BelongsToRecipeStep: new(t.Name()),
			RecipeStepProductID: new(t.Name()),
			Notes:               new(t.Name()),
			Quantity: types.Uint16RangeWithOptionalMaxUpdateRequestInput{
				Max: new(fake.Uint16()),
				Min: new(fake.Uint16()),
			},
		}

		actual := x.ValidateWithContext(t.Context())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselUpdateRequestInput{}

		actual := x.ValidateWithContext(t.Context())
		assert.Error(t, actual)
	})
}
