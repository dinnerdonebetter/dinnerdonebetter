package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
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
			Name:        pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Description: pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Volumetric:  pointers.Pointer(fake.Bool()),
			IconPath:    pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Universal:   pointers.Pointer(fake.Bool()),
			Metric:      pointers.Pointer(fake.Bool()),
			Imperial:    pointers.Pointer(fake.Bool()),
			PluralName:  pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
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
