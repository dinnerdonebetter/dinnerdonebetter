package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepEvent_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepEvent{}

		expected := &RecipeStepEventUpdateInput{
			EventType:         fake.Word(),
			Done:              fake.Bool(),
			RecipeIterationID: fake.Uint64(),
			RecipeStepID:      fake.Uint64(),
		}

		i.Update(expected)
		assert.Equal(t, expected.EventType, i.EventType)
		assert.Equal(t, expected.Done, i.Done)
		assert.Equal(t, expected.RecipeIterationID, i.RecipeIterationID)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}

func TestRecipeStepEvent_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		recipeStepEvent := &RecipeStepEvent{
			EventType:         fake.Word(),
			Done:              fake.Bool(),
			RecipeIterationID: uint64(fake.Uint32()),
			RecipeStepID:      uint64(fake.Uint32()),
		}

		expected := &RecipeStepEventUpdateInput{
			EventType:         recipeStepEvent.EventType,
			Done:              recipeStepEvent.Done,
			RecipeIterationID: recipeStepEvent.RecipeIterationID,
			RecipeStepID:      recipeStepEvent.RecipeStepID,
		}
		actual := recipeStepEvent.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
