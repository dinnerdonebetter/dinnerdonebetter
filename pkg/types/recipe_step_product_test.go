package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeStepProduct_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProduct{}

		updated := &RecipeStepProductUpdateInput{
			Name:          fake.Word(),
			QuantityType:  fake.Word(),
			QuantityValue: fake.Float32(),
			QuantityNotes: fake.Word(),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "Name",
				OldValue:  x.Name,
				NewValue:  updated.Name,
			},
			{
				FieldName: "QuantityType",
				OldValue:  x.QuantityType,
				NewValue:  updated.QuantityType,
			},
			{
				FieldName: "QuantityValue",
				OldValue:  x.QuantityValue,
				NewValue:  updated.QuantityValue,
			},
			{
				FieldName: "QuantityNotes",
				OldValue:  x.QuantityNotes,
				NewValue:  updated.QuantityNotes,
			},
		}
		actual := x.Update(updated)

		expectedJSONBytes, err := json.Marshal(expected)
		require.NoError(t, err)

		actualJSONBytes, err := json.Marshal(actual)
		require.NoError(t, err)

		expectedJSON, actualJSON := string(expectedJSONBytes), string(actualJSONBytes)

		assert.Equal(t, expectedJSON, actualJSON)

		assert.Equal(t, updated.Name, x.Name)
		assert.Equal(t, updated.QuantityType, x.QuantityType)
		assert.Equal(t, updated.QuantityValue, x.QuantityValue)
		assert.Equal(t, updated.QuantityNotes, x.QuantityNotes)
	})
}

func TestRecipeStepProductCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationInput{
			Name:          fake.Word(),
			QuantityType:  fake.Word(),
			QuantityValue: fake.Float32(),
			QuantityNotes: fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepProductUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateInput{
			Name:          fake.Word(),
			QuantityType:  fake.Word(),
			QuantityValue: fake.Float32(),
			QuantityNotes: fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepProductUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
