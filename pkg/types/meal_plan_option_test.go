package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanOptionCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionCreationRequestInput{
			BelongsToMealPlan: fake.Word(),
			DayOfWeek:         fake.Uint8(),
			RecipeID:          fake.Word(),
			Notes:             fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestMealPlanOptionUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionUpdateRequestInput{
			BelongsToMealPlan: fake.Word(),
			DayOfWeek:         fake.Uint8(),
			RecipeID:          fake.Word(),
			Notes:             fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
