package types

import (
	"context"
	"math"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrument{
			RecipeStepProductID: pointers.Pointer(t.Name()),
			MaximumQuantity:     pointers.Pointer(uint32(321)),
		}
		input := &RecipeStepInstrumentUpdateRequestInput{}

		fake.Struct(&input)
		input.Optional = pointers.Pointer(true)
		input.RecipeStepProductID = pointers.Pointer("whatever")
		input.MaximumQuantity = pointers.Pointer(uint32(123))

		x.Update(input)
	})
}

func TestRecipeStepInstrumentCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{
			InstrumentID:        pointers.Pointer(t.Name()),
			Name:                t.Name(),
			RecipeStepProductID: pointers.Pointer(t.Name()),
			Notes:               t.Name(),
			PreferenceRank:      uint8(fake.Number(1, math.MaxUint8)),
			Optional:            fake.Bool(),
			MinimumQuantity:     fake.Uint32(),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepInstrumentUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentUpdateRequestInput{
			InstrumentID:        pointers.Pointer(t.Name()),
			Name:                pointers.Pointer(t.Name()),
			BelongsToRecipeStep: pointers.Pointer(t.Name()),
			RecipeStepProductID: pointers.Pointer(t.Name()),
			Notes:               pointers.Pointer(t.Name()),
			PreferenceRank:      pointers.Pointer(uint8(fake.Number(1, math.MaxUint8))),
			Optional:            pointers.Pointer(fake.Bool()),
			MinimumQuantity:     pointers.Pointer(fake.Uint32()),
			MaximumQuantity:     pointers.Pointer(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepInstrumentUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
