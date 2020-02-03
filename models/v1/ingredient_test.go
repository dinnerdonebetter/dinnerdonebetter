package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &Ingredient{}

		expected := &IngredientUpdateInput{
			Name:              "example",
			Variant:           "example",
			Description:       "example",
			Warning:           "example",
			ContainsEgg:       false,
			ContainsDairy:     false,
			ContainsPeanut:    false,
			ContainsTreeNut:   false,
			ContainsSoy:       false,
			ContainsWheat:     false,
			ContainsShellfish: false,
			ContainsSesame:    false,
			ContainsFish:      false,
			ContainsGluten:    false,
			AnimalFlesh:       false,
			AnimalDerived:     false,
			ConsideredStaple:  false,
			Icon:              "example",
		}

		i.Update(expected)
		assert.Equal(t, expected.Name, i.Name)
		assert.Equal(t, expected.Variant, i.Variant)
		assert.Equal(t, expected.Description, i.Description)
		assert.Equal(t, expected.Warning, i.Warning)
		assert.Equal(t, expected.ContainsEgg, i.ContainsEgg)
		assert.Equal(t, expected.ContainsDairy, i.ContainsDairy)
		assert.Equal(t, expected.ContainsPeanut, i.ContainsPeanut)
		assert.Equal(t, expected.ContainsTreeNut, i.ContainsTreeNut)
		assert.Equal(t, expected.ContainsSoy, i.ContainsSoy)
		assert.Equal(t, expected.ContainsWheat, i.ContainsWheat)
		assert.Equal(t, expected.ContainsShellfish, i.ContainsShellfish)
		assert.Equal(t, expected.ContainsSesame, i.ContainsSesame)
		assert.Equal(t, expected.ContainsFish, i.ContainsFish)
		assert.Equal(t, expected.ContainsGluten, i.ContainsGluten)
		assert.Equal(t, expected.AnimalFlesh, i.AnimalFlesh)
		assert.Equal(t, expected.AnimalDerived, i.AnimalDerived)
		assert.Equal(t, expected.ConsideredStaple, i.ConsideredStaple)
		assert.Equal(t, expected.Icon, i.Icon)
	})
}
