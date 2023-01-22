package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func TestMealPlanOptionCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionCreationRequestInput{
			AssignedCook:       pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			AssignedDishwasher: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MealID:             fake.LoremIpsumSentence(exampleQuantity),
			Notes:              fake.LoremIpsumSentence(exampleQuantity),
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
			AssignedCook:           pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			AssignedDishwasher:     pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToMealPlanEvent: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MealID:                 pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:                  pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
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
