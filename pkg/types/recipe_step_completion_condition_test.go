package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func TestRecipeStepCompletionConditionCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionCreationRequestInput{
			IngredientStateID:   t.Name(),
			BelongsToRecipeStep: t.Name(),
			Optional:            fake.Bool(),
			Ingredients: []*RecipeStepCompletionConditionIngredientCreationRequestInput{
				{
					RecipeStepIngredient: t.Name(),
				},
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepCompletionConditionUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionUpdateRequestInput{
			IngredientStateID:   pointers.String(t.Name()),
			BelongsToRecipeStep: pointers.String(t.Name()),
			Optional:            boolPointer(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCompletionConditionUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
