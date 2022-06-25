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
			Name:              fake.LoremIpsumSentence(exampleQuantity),
			Description:       fake.LoremIpsumSentence(exampleQuantity),
			Warning:           fake.LoremIpsumSentence(exampleQuantity),
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
			Volumetric:        fake.Bool(),
			IconPath:          fake.LoremIpsumSentence(exampleQuantity),
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
			Name:              fake.LoremIpsumSentence(exampleQuantity),
			Description:       fake.LoremIpsumSentence(exampleQuantity),
			Warning:           fake.LoremIpsumSentence(exampleQuantity),
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
			Volumetric:        fake.Bool(),
			IconPath:          fake.LoremIpsumSentence(exampleQuantity),
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
