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
			MeasurementUnitID:                  pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:                    pointers.Float32(fake.Float32()),
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
			Name:                               pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Type:                               pointers.String(RecipeStepProductIngredientType),
			MeasurementUnitID:                  pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			MinimumQuantity:                    pointers.Float32(fake.Float32()),
			MaximumQuantity:                    pointers.Float32(fake.Float32()),
			QuantityNotes:                      pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Compostable:                        pointers.Bool(fake.Bool()),
			MaximumStorageDurationInSeconds:    pointers.Uint32(fake.Uint32()),
			MinimumStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
			MaximumStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
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
