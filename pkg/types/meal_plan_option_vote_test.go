package types

import (
	"context"
	"math"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanOptionVoteCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionVoteCreationRequestInput{
			Abstain: fake.Bool(),
			Rank:    uint8(fake.Number(1, math.MaxUint8)),
			Notes:   fake.LoremIpsumSentence(exampleQuantity),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionVoteCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestMealPlanOptionVoteUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionVoteUpdateRequestInput{
			Rank:    uint8(fake.Number(1, math.MaxUint8)),
			Abstain: fake.Bool(),
			Notes:   fake.LoremIpsumSentence(exampleQuantity),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionVoteUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
