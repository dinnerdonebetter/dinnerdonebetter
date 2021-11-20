package types

import (
	"context"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanOptionCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanOptionCreationRequestInput{
			BelongsToMealPlan: fake.LoremIpsumSentence(exampleQuantity),
			Day: time.Weekday(fake.RandomInt([]int{
				int(time.Monday),
				int(time.Tuesday),
				int(time.Wednesday),
				int(time.Thursday),
				int(time.Friday),
				int(time.Saturday),
				int(time.Sunday),
			})),
			MealName: MealName(fake.RandomString([]string{
				string(BreakfastMealName),
				string(SecondBreakfastMealName),
				string(BrunchMealName),
				string(LunchMealName),
				string(SupperMealName),
				string(DinnerMealName),
			})),
			RecipeID: fake.LoremIpsumSentence(exampleQuantity),
			Notes:    fake.LoremIpsumSentence(exampleQuantity),
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
			BelongsToMealPlan: fake.LoremIpsumSentence(exampleQuantity),
			Day: time.Weekday(fake.RandomInt([]int{
				int(time.Monday),
				int(time.Tuesday),
				int(time.Wednesday),
				int(time.Thursday),
				int(time.Friday),
				int(time.Saturday),
				int(time.Sunday),
			})),
			MealName: MealName(fake.RandomString([]string{
				string(BreakfastMealName),
				string(SecondBreakfastMealName),
				string(BrunchMealName),
				string(LunchMealName),
				string(SupperMealName),
				string(DinnerMealName),
			})),
			RecipeID: fake.LoremIpsumSentence(exampleQuantity),
			Notes:    fake.LoremIpsumSentence(exampleQuantity),
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
