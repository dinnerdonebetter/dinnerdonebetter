package types

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMealPlanCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanCreationRequestInput{
			VotingDeadline: time.Now().Add(24 * time.Hour),
			Events: []*MealPlanEventCreationRequestInput{
				{
					MealName: BreakfastMealName,
					Notes:    t.Name(),
					StartsAt: time.Now(),
					EndsAt:   time.Now().Add(24 * time.Hour),
				},
			},
		}

		assert.NoError(t, x.ValidateWithContext(context.Background()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanCreationRequestInput{}

		assert.Error(t, x.ValidateWithContext(context.Background()))
	})
}

func TestMealPlanUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleTime := time.Now()

		x := &MealPlanUpdateRequestInput{
			VotingDeadline: &exampleTime,
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
