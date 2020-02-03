package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterationMedia_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &IterationMedia{}

		expected := &IterationMediaUpdateInput{
			Path:              "example",
			Mimetype:          "example",
			RecipeIterationID: 1,
			RecipeStepID:      func(x uint64) *uint64 { return &x }(1),
		}

		i.Update(expected)
		assert.Equal(t, expected.Path, i.Path)
		assert.Equal(t, expected.Mimetype, i.Mimetype)
		assert.Equal(t, expected.RecipeIterationID, i.RecipeIterationID)
		assert.Equal(t, expected.RecipeStepID, i.RecipeStepID)
	})
}
