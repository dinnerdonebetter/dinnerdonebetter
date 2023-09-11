package types

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
)

func TestRecipeMedia_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeMedia{}
		input := &RecipeMediaUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}
