package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestValidPreparationCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{
			Name:                  fake.LoremIpsumSentence(exampleQuantity),
			Description:           fake.LoremIpsumSentence(exampleQuantity),
			IconPath:              fake.LoremIpsumSentence(exampleQuantity),
			PastTense:             fake.LoremIpsumSentence(exampleQuantity),
			YieldsNothing:         fake.Bool(),
			RestrictToIngredients: fake.Bool(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{
			Name:                  pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			Description:           pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			IconPath:              pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			PastTense:             pointers.Pointer(fake.LoremIpsumSentence(exampleQuantity)),
			YieldsNothing:         pointers.Pointer(fake.Bool()),
			RestrictToIngredients: pointers.Pointer(fake.Bool()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := &ValidPreparation{
			MaximumInstrumentCount:      pointers.Pointer(fake.Int32()),
			MaximumIngredientCount:      pointers.Pointer(fake.Int32()),
			MaximumVesselCount:          pointers.Pointer(fake.Int32()),
			IconPath:                    fake.LoremIpsumSentence(exampleQuantity),
			PastTense:                   fake.LoremIpsumSentence(exampleQuantity),
			Name:                        fake.LoremIpsumSentence(exampleQuantity),
			Description:                 fake.LoremIpsumSentence(exampleQuantity),
			Slug:                        fake.LoremIpsumSentence(exampleQuantity),
			MinimumIngredientCount:      fake.Int32(),
			MinimumInstrumentCount:      fake.Int32(),
			MinimumVesselCount:          fake.Int32(),
			RestrictToIngredients:       fake.Bool(),
			TemperatureRequired:         fake.Bool(),
			TimeEstimateRequired:        fake.Bool(),
			ConditionExpressionRequired: fake.Bool(),
			ConsumesVessel:              fake.Bool(),
			OnlyForVessels:              fake.Bool(),
			YieldsNothing:               fake.Bool(),
		}

		expected := &ValidPreparation{
			MaximumInstrumentCount:      pointers.Pointer(fake.Int32()),
			MaximumIngredientCount:      pointers.Pointer(fake.Int32()),
			MaximumVesselCount:          pointers.Pointer(fake.Int32()),
			IconPath:                    fake.LoremIpsumSentence(exampleQuantity),
			PastTense:                   fake.LoremIpsumSentence(exampleQuantity),
			Name:                        fake.LoremIpsumSentence(exampleQuantity),
			Description:                 fake.LoremIpsumSentence(exampleQuantity),
			Slug:                        fake.LoremIpsumSentence(exampleQuantity),
			MinimumIngredientCount:      fake.Int32(),
			MinimumInstrumentCount:      fake.Int32(),
			MinimumVesselCount:          fake.Int32(),
			RestrictToIngredients:       !actual.RestrictToIngredients,
			TemperatureRequired:         !actual.TemperatureRequired,
			TimeEstimateRequired:        !actual.TimeEstimateRequired,
			ConditionExpressionRequired: !actual.ConditionExpressionRequired,
			ConsumesVessel:              !actual.ConsumesVessel,
			OnlyForVessels:              !actual.OnlyForVessels,
			YieldsNothing:               !actual.YieldsNothing,
		}

		exampleUpdatedValue := &ValidPreparationUpdateRequestInput{
			MaximumInstrumentCount:      expected.MaximumInstrumentCount,
			MaximumIngredientCount:      expected.MaximumIngredientCount,
			MaximumVesselCount:          expected.MaximumVesselCount,
			IconPath:                    &expected.IconPath,
			PastTense:                   &expected.PastTense,
			Name:                        &expected.Name,
			Description:                 &expected.Description,
			Slug:                        &expected.Slug,
			MinimumIngredientCount:      &expected.MinimumIngredientCount,
			MinimumInstrumentCount:      &expected.MinimumInstrumentCount,
			MinimumVesselCount:          &expected.MinimumVesselCount,
			RestrictToIngredients:       &expected.RestrictToIngredients,
			TemperatureRequired:         &expected.TemperatureRequired,
			TimeEstimateRequired:        &expected.TimeEstimateRequired,
			ConditionExpressionRequired: &expected.ConditionExpressionRequired,
			ConsumesVessel:              &expected.ConsumesVessel,
			OnlyForVessels:              &expected.OnlyForVessels,
			YieldsNothing:               &expected.YieldsNothing,
		}

		actual.Update(exampleUpdatedValue)

		assert.Equal(t, expected, actual)
	})
}

func TestValidPreparationCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationCreationRequestInput{
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationDatabaseCreationInput{
			ID:   t.Name(),
			Name: t.Name(),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestValidPreparationUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &ValidPreparationUpdateRequestInput{
			Name: pointers.Pointer(t.Name()),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}
