package types

import (
	"context"
	"testing"
	"time"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		now := time.Now()
		inTenMinutes := time.Now().Add(time.Minute * 10)
		inOneWeek := time.Now().Add((time.Hour * 24) * 7)

		x := &MealPlanCreationRequestInput{
			VotingDeadline: uint64(now.Unix()),
			StartsAt:       uint64(inTenMinutes.Unix()),
			EndsAt:         uint64(inOneWeek.Unix()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestMealPlanUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		status := AwaitingVotesMealPlanStatus

		x := &MealPlanUpdateRequestInput{
			Status:         &status,
			VotingDeadline: uint64Pointer(uint64(time.Now().Unix())),
			StartsAt:       uint64Pointer(uint64(fake.Uint32())),
			EndsAt:         uint64Pointer(uint64(fake.Uint32())),
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
