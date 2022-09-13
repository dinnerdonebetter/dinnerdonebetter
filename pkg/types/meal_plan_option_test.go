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
			AssignedCook:           stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			AssignedDishwasher:     stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToMealPlanEvent: fake.LoremIpsumSentence(exampleQuantity),
			MealID:                 fake.LoremIpsumSentence(exampleQuantity),
			Notes:                  fake.LoremIpsumSentence(exampleQuantity),
			PrepStepsCreated:       false,
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
			AssignedCook:           stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			AssignedDishwasher:     stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToMealPlanEvent: stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MealID:                 stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:                  stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			PrepStepsCreated:       boolPointer(false),
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
