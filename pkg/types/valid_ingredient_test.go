package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientCreationRequestInput{
			Name:                                    fake.LoremIpsumSentence(exampleQuantity),
			Description:                             fake.LoremIpsumSentence(exampleQuantity),
			Warning:                                 fake.LoremIpsumSentence(exampleQuantity),
			ContainsEgg:                             fake.Bool(),
			ContainsDairy:                           fake.Bool(),
			ContainsPeanut:                          fake.Bool(),
			ContainsTreeNut:                         fake.Bool(),
			ContainsSoy:                             fake.Bool(),
			ContainsWheat:                           fake.Bool(),
			ContainsShellfish:                       fake.Bool(),
			ContainsSesame:                          fake.Bool(),
			ContainsFish:                            fake.Bool(),
			ContainsGluten:                          fake.Bool(),
			AnimalFlesh:                             fake.Bool(),
			IsMeasuredVolumetrically:                fake.Bool(),
			IconPath:                                fake.LoremIpsumSentence(exampleQuantity),
			PluralName:                              fake.LoremIpsumSentence(exampleQuantity),
			AnimalDerived:                           fake.Bool(),
			RestrictToPreparations:                  fake.Bool(),
			MinimumIdealStorageTemperatureInCelsius: fake.Float32(),
			MaximumIdealStorageTemperatureInCelsius: fake.Float32(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientUpdateRequestInput{
			Name:                                    stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Description:                             stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			Warning:                                 stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			ContainsEgg:                             boolPointer(fake.Bool()),
			ContainsDairy:                           boolPointer(fake.Bool()),
			ContainsPeanut:                          boolPointer(fake.Bool()),
			ContainsTreeNut:                         boolPointer(fake.Bool()),
			ContainsSoy:                             boolPointer(fake.Bool()),
			ContainsWheat:                           boolPointer(fake.Bool()),
			ContainsShellfish:                       boolPointer(fake.Bool()),
			ContainsSesame:                          boolPointer(fake.Bool()),
			ContainsFish:                            boolPointer(fake.Bool()),
			ContainsGluten:                          boolPointer(fake.Bool()),
			AnimalFlesh:                             boolPointer(fake.Bool()),
			IsMeasuredVolumetrically:                boolPointer(fake.Bool()),
			IconPath:                                stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			PluralName:                              stringPointer(fake.LoremIpsumSentence(exampleQuantity)),
			AnimalDerived:                           boolPointer(fake.Bool()),
			RestrictToPreparations:                  boolPointer(fake.Bool()),
			MinimumIdealStorageTemperatureInCelsius: float32Pointer(fake.Float32()),
			MaximumIdealStorageTemperatureInCelsius: float32Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
