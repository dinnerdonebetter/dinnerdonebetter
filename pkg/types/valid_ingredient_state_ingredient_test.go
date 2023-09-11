package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientStateIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientStateIngredient{}
		input := &ValidIngredientStateIngredientUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}

func TestValidIngredientStateIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientStateIngredientCreationRequestInput{
			Notes:                  t.Name(),
			ValidIngredientStateID: t.Name(),
			ValidIngredientID:      t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientStateIngredientCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientStateIngredientUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientStateIngredientUpdateRequestInput{
			Notes:                  pointers.Pointer(t.Name()),
			ValidIngredientStateID: pointers.Pointer(t.Name()),
			ValidIngredientID:      pointers.Pointer(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientStateIngredientUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
