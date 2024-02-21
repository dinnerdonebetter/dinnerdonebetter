package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepVessel_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVessel{
			MinimumQuantity: 1234,
			MaximumQuantity: pointers.Pointer(uint32(1234)),
			Vessel:          &ValidVessel{},
		}
		input := &RecipeStepVesselUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.UnavailableAfterStep = pointers.Pointer(true)
		input.MinimumQuantity = pointers.Pointer(uint32(1))
		input.MaximumQuantity = pointers.Pointer(uint32(1))
		input.VesselID = pointers.Pointer(t.Name())

		x.Update(input)
	})
}

func TestRecipeStepVesselCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepVesselCreationRequestInput{
			Name:                t.Name(),
			RecipeStepProductID: pointers.Pointer(t.Name()),
			Notes:               t.Name(),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
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
			Name:                pointers.Pointer(t.Name()),
			BelongsToRecipeStep: pointers.Pointer(t.Name()),
			RecipeStepProductID: pointers.Pointer(t.Name()),
			Notes:               pointers.Pointer(t.Name()),
			MinimumQuantity:     pointers.Pointer(fake.Uint32()),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
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
