package types

import (
	"context"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"

	"github.com/prixfixeco/backend/internal/pointers"
)

func TestValidMeasurementUnitCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitCreationRequestInput{
			Name:        fake.LoremIpsumSentence(exampleQuantity),
			Description: fake.LoremIpsumSentence(exampleQuantity),
			Volumetric:  fake.Bool(),
			IconPath:    fake.LoremIpsumSentence(exampleQuantity),
			Universal:   fake.Bool(),
			Metric:      fake.Bool(),
			Imperial:    fake.Bool(),
			PluralName:  fake.LoremIpsumSentence(exampleQuantity),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidMeasurementUnitUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitUpdateRequestInput{
			Name:        pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Description: pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Volumetric:  pointers.Bool(fake.Bool()),
			IconPath:    pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
			Universal:   pointers.Bool(fake.Bool()),
			Metric:      pointers.Bool(fake.Bool()),
			Imperial:    pointers.Bool(fake.Bool()),
			PluralName:  pointers.String(fake.LoremIpsumSentence(exampleQuantity)),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidMeasurementUnitUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
