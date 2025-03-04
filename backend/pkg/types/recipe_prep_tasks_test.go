package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestRecipePrepTask_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTask{}
		input := &RecipePrepTaskUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.Optional = pointer.To(true)
		input.BelongsToRecipe = pointer.To(t.Name())

		x.Update(input)
	})
}

func TestRecipePrepTaskCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskCreationRequestInput{
			BelongsToRecipe: t.Name(),
			Name:            t.Name(),
			StorageType:     t.Name(),
			StorageTemperatureInCelsius: OptionalFloat32Range{
				Max: pointer.To(fake.Float32()),
				Min: pointer.To(fake.Float32()),
			},
			TimeBufferBeforeRecipeInSeconds: Uint32RangeWithOptionalMax{
				Max: pointer.To(fake.Uint32()),
				Min: fake.Uint32(),
			},
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipePrepTaskDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskDatabaseCreationInput{
			ID:              t.Name(),
			BelongsToRecipe: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipePrepTaskUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskUpdateRequestInput{
			BelongsToRecipe: pointer.To(t.Name()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.NoError(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipePrepTaskUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
