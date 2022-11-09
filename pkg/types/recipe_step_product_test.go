package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/api_server/internal/pointers"
)

func TestRecipeStepProductCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationRequestInput{
			Name:                               fake.LoremIpsumSentence(exampleQuantity),
			Type:                               RecipeStepProductIngredientType,
			MeasurementUnitID:                  fake.LoremIpsumSentence(exampleQuantity),
			MinimumQuantity:                    fake.Float32(),
			QuantityNotes:                      fake.LoremIpsumSentence(exampleQuantity),
			Compostable:                        fake.Bool(),
			MaximumStorageDurationInSeconds:    pointers.Uint32(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
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
			Name:                               stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Type:                               stringPointer(RecipeStepProductIngredientType),
			MeasurementUnitID:                  stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:                    float32Pointer(fake.Float32()),
			MaximumQuantity:                    float32Pointer(fake.Float32()),
			QuantityNotes:                      stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Compostable:                        boolPointer(fake.Bool()),
			MaximumStorageDurationInSeconds:    uint32Pointer(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: float32Pointer(fake.Float32()),
			MaximumStorageTemperatureInCelsius: float32Pointer(fake.Float32()),
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
