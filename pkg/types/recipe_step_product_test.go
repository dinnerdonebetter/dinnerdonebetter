package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:                               fake.LoremIpsumSentence(exampleQuantity),
			Type:                               RecipeStepProductIngredientType,
			MeasurementUnitID:                  pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:                    pointers.Pointer(fake.Float32()),
			QuantityNotes:                      fake.LoremIpsumSentence(exampleQuantity),
			Compostable:                        fake.Bool(),
			MaximumStorageDurationInSeconds:    pointers.Pointer(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{
			Name:                               pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Type:                               pointers.Pointer(RecipeStepProductIngredientType),
			MeasurementUnitID:                  pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:                    pointers.Pointer(fake.Float32()),
			MaximumQuantity:                    pointers.Pointer(fake.Float32()),
			QuantityNotes:                      pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Compostable:                        pointers.Pointer(fake.Bool()),
			MaximumStorageDurationInSeconds:    pointers.Pointer(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
