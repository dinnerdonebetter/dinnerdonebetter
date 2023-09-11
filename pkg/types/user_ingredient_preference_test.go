package types

import (
	"context"
	"math"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserIngredientPreference_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &UserIngredientPreference{}
		input := &UserIngredientPreferenceUpdateRequestInput{}

		fake.Struct(&input)

		x.Update(input)
	})
}

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
