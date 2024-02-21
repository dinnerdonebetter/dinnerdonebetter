package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

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

		input := &ValidIngredientUpdateRequestInput{}
		assert.NoError(t, fake.Struct(&input))
		input.ContainsEgg = pointers.Pointer(true)
		input.ContainsDairy = pointers.Pointer(true)
		input.ContainsPeanut = pointers.Pointer(true)
		input.ContainsTreeNut = pointers.Pointer(true)
		input.ContainsSoy = pointers.Pointer(true)
		input.ContainsWheat = pointers.Pointer(true)
		input.ContainsShellfish = pointers.Pointer(true)
		input.ContainsSesame = pointers.Pointer(true)
		input.ContainsFish = pointers.Pointer(true)
		input.ContainsGluten = pointers.Pointer(true)
		input.AnimalFlesh = pointers.Pointer(true)
		input.IsMeasuredVolumetrically = pointers.Pointer(true)
		input.AnimalDerived = pointers.Pointer(true)
		input.RestrictToPreparations = pointers.Pointer(true)
		input.ContainsAlcohol = pointers.Pointer(true)
		input.IsLiquid = pointers.Pointer(true)
		input.IsProtein = pointers.Pointer(true)
		input.IsStarch = pointers.Pointer(true)
		input.IsGrain = pointers.Pointer(true)
		input.IsFruit = pointers.Pointer(true)
		input.IsSalt = pointers.Pointer(true)
		input.IsFat = pointers.Pointer(true)
		input.IsAcid = pointers.Pointer(true)
		input.IsHeat = pointers.Pointer(true)

		actual.Update(input)
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
