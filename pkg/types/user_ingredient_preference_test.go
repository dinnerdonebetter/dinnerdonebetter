package types

import (
	"context"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserIngredientPreferenceCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		input := &UserIngredientPreferenceCreationRequestInput{
			ValidIngredientID: t.Name(),
			Rating:            1,
		}

		assert.NoError(t, input.ValidateWithContext(ctx))
	})

	T.Run("invalid range", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		input := &UserIngredientPreferenceCreationRequestInput{
			ValidIngredientID: t.Name(),
			Rating:            math.MaxInt8,
		}

		assert.Error(t, input.ValidateWithContext(ctx))
	})

	T.Run("with valid ingredient group ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		input := &UserIngredientPreferenceCreationRequestInput{
			ValidIngredientGroupID: t.Name(),
			Rating:                 1,
		}

		assert.NoError(t, input.ValidateWithContext(ctx))
	})
}
