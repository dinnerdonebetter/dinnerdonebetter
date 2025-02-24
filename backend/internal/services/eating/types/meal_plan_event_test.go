package types

import (
	"testing"
	"time"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanEventCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealPlanEventCreationRequestInput{
			MealName: SecondBreakfastMealName,
			StartsAt: time.Now(),
			EndsAt:   time.Now(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})

	T.Run("with invalid time", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		tt := time.Now()
		x := &MealPlanEventCreationRequestInput{
			MealName: SecondBreakfastMealName,
			StartsAt: tt,
			EndsAt:   tt,
		}

		assert.Error(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanEventDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealPlanEventDatabaseCreationInput{
			ID:                t.Name(),
			BelongsToMealPlan: t.Name(),
			MealName:          SecondBreakfastMealName,
			StartsAt:          time.Now(),
			EndsAt:            time.Now(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanEventUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealPlanEventUpdateRequestInput{
			MealName: pointer.To(SecondBreakfastMealName),
			StartsAt: pointer.To(time.Now()),
			EndsAt:   pointer.To(time.Now()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanEvent_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanEvent{}
		input := &MealPlanEventUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.StartsAt = pointer.To(time.Now())
		input.EndsAt = pointer.To(time.Now())

		x.Update(input)
	})
}
