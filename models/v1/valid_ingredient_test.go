package models

import (
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		i := &ValidIngredient{}

		expected := &ValidIngredientUpdateInput{
			Name:              fake.Word(),
			Variant:           fake.Word(),
			Description:       fake.Word(),
			Warning:           fake.Word(),
			ContainsEgg:       fake.Bool(),
			ContainsDairy:     fake.Bool(),
			ContainsPeanut:    fake.Bool(),
			ContainsTreeNut:   fake.Bool(),
			ContainsSoy:       fake.Bool(),
			ContainsWheat:     fake.Bool(),
			ContainsShellfish: fake.Bool(),
			ContainsSesame:    fake.Bool(),
			ContainsFish:      fake.Bool(),
			ContainsGluten:    fake.Bool(),
			AnimalFlesh:       fake.Bool(),
			AnimalDerived:     fake.Bool(),
			ConsideredStaple:  fake.Bool(),
			Icon:              fake.Word(),
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

func TestValidIngredient_ToUpdateInput(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		validIngredient := &ValidIngredient{
			Name:              fake.Word(),
			Variant:           fake.Word(),
			Description:       fake.Word(),
			Warning:           fake.Word(),
			ContainsEgg:       fake.Bool(),
			ContainsDairy:     fake.Bool(),
			ContainsPeanut:    fake.Bool(),
			ContainsTreeNut:   fake.Bool(),
			ContainsSoy:       fake.Bool(),
			ContainsWheat:     fake.Bool(),
			ContainsShellfish: fake.Bool(),
			ContainsSesame:    fake.Bool(),
			ContainsFish:      fake.Bool(),
			ContainsGluten:    fake.Bool(),
			AnimalFlesh:       fake.Bool(),
			AnimalDerived:     fake.Bool(),
			ConsideredStaple:  fake.Bool(),
			Icon:              fake.Word(),
		}

		expected := &ValidIngredientUpdateInput{
			Name:              validIngredient.Name,
			Variant:           validIngredient.Variant,
			Description:       validIngredient.Description,
			Warning:           validIngredient.Warning,
			ContainsEgg:       validIngredient.ContainsEgg,
			ContainsDairy:     validIngredient.ContainsDairy,
			ContainsPeanut:    validIngredient.ContainsPeanut,
			ContainsTreeNut:   validIngredient.ContainsTreeNut,
			ContainsSoy:       validIngredient.ContainsSoy,
			ContainsWheat:     validIngredient.ContainsWheat,
			ContainsShellfish: validIngredient.ContainsShellfish,
			ContainsSesame:    validIngredient.ContainsSesame,
			ContainsFish:      validIngredient.ContainsFish,
			ContainsGluten:    validIngredient.ContainsGluten,
			AnimalFlesh:       validIngredient.AnimalFlesh,
			AnimalDerived:     validIngredient.AnimalDerived,
			ConsideredStaple:  validIngredient.ConsideredStaple,
			Icon:              validIngredient.Icon,
		}
		actual := validIngredient.ToUpdateInput()

		assert.Equal(t, expected, actual)
	})
}
