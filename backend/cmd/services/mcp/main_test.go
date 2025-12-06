package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_schemaFromType(T *testing.T) {
	T.Parallel()

	T.Run("SearchValidIngredientsInvocation", func(t *testing.T) {
		t.Parallel()

		expected := validIngredientSearchInputSchema()
		actual := schemaForType(SearchValidIngredientsInvocation{})

		assert.Equal(t, expected, actual)
	})

	T.Run("SearchValidIngredientsResult", func(t *testing.T) {
		t.Parallel()

		expected := validIngredientOutputSchema()
		actual := schemaForType(SearchValidIngredientsResult{})

		assert.Equal(t, expected, actual)
	})
}
