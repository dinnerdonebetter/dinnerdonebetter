package types

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
)

func TestValidIngredientState_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientState{}
		input := &ValidIngredientStateUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}
