package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeStep_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStep{}

		updated := &RecipeStepUpdateInput{
			Index:                     uint(fake.Uint32()),
			PreparationID:             uint64(fake.Uint32()),
			PrerequisiteStep:          uint64(fake.Uint32()),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
			Notes:                     fake.Word(),
			Why:                       fake.Word(),
			RecipeID:                  uint64(fake.Uint32()),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "Index",
				OldValue:  x.Index,
				NewValue:  updated.Index,
			},
			{
				FieldName: "PreparationID",
				OldValue:  x.PreparationID,
				NewValue:  updated.PreparationID,
			},
			{
				FieldName: "PrerequisiteStep",
				OldValue:  x.PrerequisiteStep,
				NewValue:  updated.PrerequisiteStep,
			},
			{
				FieldName: "MinEstimatedTimeInSeconds",
				OldValue:  x.MinEstimatedTimeInSeconds,
				NewValue:  updated.MinEstimatedTimeInSeconds,
			},
			{
				FieldName: "MaxEstimatedTimeInSeconds",
				OldValue:  x.MaxEstimatedTimeInSeconds,
				NewValue:  updated.MaxEstimatedTimeInSeconds,
			},
			{
				FieldName: "TemperatureInCelsius",
				OldValue:  x.TemperatureInCelsius,
				NewValue:  updated.TemperatureInCelsius,
			},
			{
				FieldName: "Notes",
				OldValue:  x.Notes,
				NewValue:  updated.Notes,
			},
			{
				FieldName: "Why",
				OldValue:  x.Why,
				NewValue:  updated.Why,
			},
			{
				FieldName: "RecipeID",
				OldValue:  x.RecipeID,
				NewValue:  updated.RecipeID,
			},
		}
		actual := x.Update(updated)

		expectedJSONBytes, err := json.Marshal(expected)
		require.NoError(t, err)

		actualJSONBytes, err := json.Marshal(actual)
		require.NoError(t, err)

		expectedJSON, actualJSON := string(expectedJSONBytes), string(actualJSONBytes)

		assert.Equal(t, expectedJSON, actualJSON)

		assert.Equal(t, updated.Index, x.Index)
		assert.Equal(t, updated.PreparationID, x.PreparationID)
		assert.Equal(t, updated.PrerequisiteStep, x.PrerequisiteStep)
		assert.Equal(t, updated.MinEstimatedTimeInSeconds, x.MinEstimatedTimeInSeconds)
		assert.Equal(t, updated.MaxEstimatedTimeInSeconds, x.MaxEstimatedTimeInSeconds)
		assert.Equal(t, updated.TemperatureInCelsius, x.TemperatureInCelsius)
		assert.Equal(t, updated.Notes, x.Notes)
		assert.Equal(t, updated.Why, x.Why)
		assert.Equal(t, updated.RecipeID, x.RecipeID)
	})
}

func TestRecipeStepCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationInput{
			Index:                     uint(fake.Uint32()),
			PreparationID:             uint64(fake.Uint32()),
			PrerequisiteStep:          uint64(fake.Uint32()),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
			Notes:                     fake.Word(),
			Why:                       fake.Word(),
			RecipeID:                  uint64(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateInput{
			Index:                     uint(fake.Uint32()),
			PreparationID:             uint64(fake.Uint32()),
			PrerequisiteStep:          uint64(fake.Uint32()),
			MinEstimatedTimeInSeconds: fake.Uint32(),
			MaxEstimatedTimeInSeconds: fake.Uint32(),
			TemperatureInCelsius:      func(x uint16) *uint16 { return &x }(fake.Uint16()),
			Notes:                     fake.Word(),
			Why:                       fake.Word(),
			RecipeID:                  uint64(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
