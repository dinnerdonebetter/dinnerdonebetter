package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	"github.com/stretchr/testify/assert"
)

func TestMealRatingCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealRatingCreationRequestInput{
			MealID:     t.Name(),
			Difficulty: 1.0,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealRatingDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealRatingDatabaseCreationInput{
			ID:         t.Name(),
			MealID:     t.Name(),
			ByUser:     t.Name(),
			Difficulty: 1.0,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealRatingUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealRatingUpdateRequestInput{
			ByUser:     pointers.Pointer(t.Name()),
			MealID:     pointers.Pointer(t.Name()),
			Difficulty: pointers.Pointer[float32](1.0),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
