package types

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipeRating_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeRating{}
		input := &RecipeRatingUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestRecipeRatingCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &RecipeRatingCreationRequestInput{
			RecipeID:   t.Name(),
			Difficulty: 1.0,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &RecipeRatingCreationRequestInput{}

		assert.Error(t, x.ValidateWithContext(ctx))
	})
}

func TestRecipeRatingDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &RecipeRatingDatabaseCreationInput{
			ID:         t.Name(),
			RecipeID:   t.Name(),
			ByUser:     t.Name(),
			Difficulty: 1.0,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &RecipeRatingDatabaseCreationInput{}

		assert.Error(t, x.ValidateWithContext(ctx))
	})
}

func TestRecipeRatingUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &RecipeRatingUpdateRequestInput{
			ByUser:     pointer.To(t.Name()),
			RecipeID:   pointer.To(t.Name()),
			Difficulty: pointer.To[float32](1.0),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &RecipeRatingUpdateRequestInput{}

		assert.Error(t, x.ValidateWithContext(ctx))
	})
}
