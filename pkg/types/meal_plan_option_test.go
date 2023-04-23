package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanOptionCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionCreationRequestInput{
			AssignedCook:       pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			AssignedDishwasher: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
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
			AssignedCook:           pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			AssignedDishwasher:     pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			BelongsToMealPlanEvent: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MealID:                 pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Notes:                  pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
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
