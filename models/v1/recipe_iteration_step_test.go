package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeIterationStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeIterationStep{}

		expected := &RecipeIterationStepUpdateInput{
			StartedOn: func(x uint64) *uint64 { return &x }(fake.Uint64()),
			EndedOn:   func(x uint64) *uint64 { return &x }(fake.Uint64()),
			State:     fake.Word(),
		}

		i.Update(expected)
		assert.Equal(t, expected.StartedOn, i.StartedOn)
		assert.Equal(t, expected.EndedOn, i.EndedOn)
		assert.Equal(t, expected.State, i.State)
	})
}

func TestRecipeIterationStep_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeIterationStep := &RecipeIterationStep{
			StartedOn: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			EndedOn:   func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			State:     fake.Word(),
		}

		expected := &RecipeIterationStepUpdateInput{
			StartedOn: recipeIterationStep.StartedOn,
			EndedOn:   recipeIterationStep.EndedOn,
			State:     recipeIterationStep.State,
		}
		actual := recipeIterationStep.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
