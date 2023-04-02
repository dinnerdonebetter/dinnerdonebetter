package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pointers"

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
			MinimumIdealStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
			MaximumIdealStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
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
			Name:                                    pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Description:                             pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Warning:                                 pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			ContainsEgg:                             pointers.Bool(fake.Bool()),
			ContainsDairy:                           pointers.Bool(fake.Bool()),
			ContainsPeanut:                          pointers.Bool(fake.Bool()),
			ContainsTreeNut:                         pointers.Bool(fake.Bool()),
			ContainsSoy:                             pointers.Bool(fake.Bool()),
			ContainsWheat:                           pointers.Bool(fake.Bool()),
			ContainsShellfish:                       pointers.Bool(fake.Bool()),
			ContainsSesame:                          pointers.Bool(fake.Bool()),
			ContainsFish:                            pointers.Bool(fake.Bool()),
			ContainsGluten:                          pointers.Bool(fake.Bool()),
			AnimalFlesh:                             pointers.Bool(fake.Bool()),
			IsMeasuredVolumetrically:                pointers.Bool(fake.Bool()),
			IconPath:                                pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			PluralName:                              pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			AnimalDerived:                           pointers.Bool(fake.Bool()),
			RestrictToPreparations:                  pointers.Bool(fake.Bool()),
			MinimumIdealStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
			MaximumIdealStorageTemperatureInCelsius: pointers.Float32(fake.Float32()),
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
