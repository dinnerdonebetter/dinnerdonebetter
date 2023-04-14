package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"

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

func TestValidIngredientCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidIngredientCreationRequestInput{
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidIngredientDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidIngredientDatabaseCreationInput{
			ID:   t.Name(),
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidIngredientUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidIngredientUpdateRequestInput{
			Name: pointers.String(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := &ValidIngredient{
			MaximumIdealStorageTemperatureInCelsius: nil,
			MinimumIdealStorageTemperatureInCelsius: nil,
			IconPath:                                "",
			Warning:                                 "",
			PluralName:                              "",
			StorageInstructions:                     "",
			Name:                                    "",
			Description:                             "",
			Slug:                                    "",
			ShoppingSuggestions:                     "",
			ContainsShellfish:                       false,
			IsMeasuredVolumetrically:                false,
			IsLiquid:                                false,
			ContainsPeanut:                          false,
			ContainsTreeNut:                         false,
			ContainsEgg:                             false,
			ContainsWheat:                           false,
			ContainsSoy:                             false,
			AnimalDerived:                           false,
			RestrictToPreparations:                  false,
			ContainsSesame:                          false,
			ContainsFish:                            false,
			ContainsGluten:                          false,
			ContainsDairy:                           false,
			ContainsAlcohol:                         false,
			AnimalFlesh:                             false,
		}

		input := &ValidIngredientUpdateRequestInput{
			MinimumIdealStorageTemperatureInCelsius: pointers.Float32(1.23),
			MaximumIdealStorageTemperatureInCelsius: pointers.Float32(1.23),
			IconPath:                                pointers.String(t.Name()),
			Warning:                                 pointers.String(t.Name()),
			PluralName:                              pointers.String(t.Name()),
			StorageInstructions:                     pointers.String(t.Name()),
			Name:                                    pointers.String(t.Name()),
			Description:                             pointers.String(t.Name()),
			Slug:                                    pointers.String(t.Name()),
			ShoppingSuggestions:                     pointers.String(t.Name()),
			ContainsShellfish:                       pointers.Bool(true),
			IsMeasuredVolumetrically:                pointers.Bool(true),
			IsLiquid:                                pointers.Bool(true),
			ContainsPeanut:                          pointers.Bool(true),
			ContainsTreeNut:                         pointers.Bool(true),
			ContainsEgg:                             pointers.Bool(true),
			ContainsWheat:                           pointers.Bool(true),
			ContainsSoy:                             pointers.Bool(true),
			AnimalDerived:                           pointers.Bool(true),
			RestrictToPreparations:                  pointers.Bool(true),
			ContainsSesame:                          pointers.Bool(true),
			ContainsFish:                            pointers.Bool(true),
			ContainsGluten:                          pointers.Bool(true),
			ContainsDairy:                           pointers.Bool(true),
			ContainsAlcohol:                         pointers.Bool(true),
			AnimalFlesh:                             pointers.Bool(true),
		}

		expected := &ValidIngredient{
			MaximumIdealStorageTemperatureInCelsius: input.MaximumIdealStorageTemperatureInCelsius,
			MinimumIdealStorageTemperatureInCelsius: input.MinimumIdealStorageTemperatureInCelsius,
			IconPath:                                *input.IconPath,
			Warning:                                 *input.Warning,
			PluralName:                              *input.PluralName,
			StorageInstructions:                     *input.StorageInstructions,
			Name:                                    *input.Name,
			Description:                             *input.Description,
			Slug:                                    *input.Slug,
			ShoppingSuggestions:                     *input.ShoppingSuggestions,
			ContainsShellfish:                       *input.ContainsShellfish,
			IsMeasuredVolumetrically:                *input.IsMeasuredVolumetrically,
			IsLiquid:                                *input.IsLiquid,
			ContainsPeanut:                          *input.ContainsPeanut,
			ContainsTreeNut:                         *input.ContainsTreeNut,
			ContainsEgg:                             *input.ContainsEgg,
			ContainsWheat:                           *input.ContainsWheat,
			ContainsSoy:                             *input.ContainsSoy,
			AnimalDerived:                           *input.AnimalDerived,
			RestrictToPreparations:                  *input.RestrictToPreparations,
			ContainsSesame:                          *input.ContainsSesame,
			ContainsFish:                            *input.ContainsFish,
			ContainsGluten:                          *input.ContainsGluten,
			ContainsDairy:                           *input.ContainsDairy,
			ContainsAlcohol:                         *input.ContainsAlcohol,
			AnimalFlesh:                             *input.AnimalFlesh,
		}

		actual.Update(input)

		assert.Equal(t, actual, expected)
	})
}
