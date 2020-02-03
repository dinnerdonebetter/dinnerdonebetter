package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRecipeStepEvent_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &RecipeStepEvent{}

		expected := &RecipeStepEventUpdateInput{
			EventType:         "example",
			Done:              false,
			RecipeIterationID: 1,
			RecipeStepID:      1,
		}

		i.Update(expected)
		assert.Equal(t, expected.EventType, i.EventType)
		assert.Equal(t, expected.Done, i.Done)
		assert.Equal(t, expected.RecipeIterationID, i.RecipeIterationID)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}
