package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredientCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientCreationRequestInput{
			Name:                                    t.Name(),
			Description:                             t.Name(),
			Warning:                                 t.Name(),
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
			IconPath:                                t.Name(),
			PluralName:                              t.Name(),
			AnimalDerived:                           fake.Bool(),
			RestrictToPreparations:                  fake.Bool(),
			MinimumIdealStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
			MaximumIdealStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
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
			Name:                                    pointers.Pointer(t.Name()),
			Description:                             pointers.Pointer(t.Name()),
			Warning:                                 pointers.Pointer(t.Name()),
			ContainsEgg:                             pointers.Pointer(fake.Bool()),
			ContainsDairy:                           pointers.Pointer(fake.Bool()),
			ContainsPeanut:                          pointers.Pointer(fake.Bool()),
			ContainsTreeNut:                         pointers.Pointer(fake.Bool()),
			ContainsSoy:                             pointers.Pointer(fake.Bool()),
			ContainsWheat:                           pointers.Pointer(fake.Bool()),
			ContainsShellfish:                       pointers.Pointer(fake.Bool()),
			ContainsSesame:                          pointers.Pointer(fake.Bool()),
			ContainsFish:                            pointers.Pointer(fake.Bool()),
			ContainsGluten:                          pointers.Pointer(fake.Bool()),
			AnimalFlesh:                             pointers.Pointer(fake.Bool()),
			IsMeasuredVolumetrically:                pointers.Pointer(fake.Bool()),
			IconPath:                                pointers.Pointer(t.Name()),
			PluralName:                              pointers.Pointer(t.Name()),
			AnimalDerived:                           pointers.Pointer(fake.Bool()),
			RestrictToPreparations:                  pointers.Pointer(fake.Bool()),
			MinimumIdealStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
			MaximumIdealStorageTemperatureInCelsius: pointers.Pointer(fake.Float32()),
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
			Name: pointers.Pointer(t.Name()),
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
			MinimumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(1.23)),
			MaximumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(1.23)),
			IconPath:                                pointers.Pointer(t.Name()),
			Warning:                                 pointers.Pointer(t.Name()),
			PluralName:                              pointers.Pointer(t.Name()),
			StorageInstructions:                     pointers.Pointer(t.Name()),
			Name:                                    pointers.Pointer(t.Name()),
			Description:                             pointers.Pointer(t.Name()),
			Slug:                                    pointers.Pointer(t.Name()),
			ShoppingSuggestions:                     pointers.Pointer(t.Name()),
			ContainsShellfish:                       pointers.Pointer(true),
			IsMeasuredVolumetrically:                pointers.Pointer(true),
			IsLiquid:                                pointers.Pointer(true),
			ContainsPeanut:                          pointers.Pointer(true),
			ContainsTreeNut:                         pointers.Pointer(true),
			ContainsEgg:                             pointers.Pointer(true),
			ContainsWheat:                           pointers.Pointer(true),
			ContainsSoy:                             pointers.Pointer(true),
			AnimalDerived:                           pointers.Pointer(true),
			RestrictToPreparations:                  pointers.Pointer(true),
			ContainsSesame:                          pointers.Pointer(true),
			ContainsFish:                            pointers.Pointer(true),
			ContainsGluten:                          pointers.Pointer(true),
			ContainsDairy:                           pointers.Pointer(true),
			ContainsAlcohol:                         pointers.Pointer(true),
			AnimalFlesh:                             pointers.Pointer(true),
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
