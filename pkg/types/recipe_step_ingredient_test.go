package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipeStepIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredient{}

		updated := &RecipeStepIngredientUpdateInput{
			IngredientID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			Name:                fake.Word(),
			QuantityType:        fake.Word(),
			QuantityValue:       fake.Float32(),
			QuantityNotes:       fake.Word(),
			ProductOfRecipeStep: true,
			IngredientNotes:     fake.Word(),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "IngredientID",
				OldValue:  x.IngredientID,
				NewValue:  updated.IngredientID,
			},
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
			{
				FieldName: "ProductOfRecipeStep",
				OldValue:  x.ProductOfRecipeStep,
				NewValue:  updated.ProductOfRecipeStep,
			},
			{
				FieldName: "IngredientNotes",
				OldValue:  x.IngredientNotes,
				NewValue:  updated.IngredientNotes,
			},
		}
		actual := x.Update(updated)

		expectedJSONBytes, err := json.Marshal(expected)
		require.NoError(t, err)

		actualJSONBytes, err := json.Marshal(actual)
		require.NoError(t, err)

		expectedJSON, actualJSON := string(expectedJSONBytes), string(actualJSONBytes)

		assert.Equal(t, expectedJSON, actualJSON)

		assert.Equal(t, updated.IngredientID, x.IngredientID)
		assert.Equal(t, updated.Name, x.Name)
		assert.Equal(t, updated.QuantityType, x.QuantityType)
		assert.Equal(t, updated.QuantityValue, x.QuantityValue)
		assert.Equal(t, updated.QuantityNotes, x.QuantityNotes)
		assert.Equal(t, updated.ProductOfRecipeStep, x.ProductOfRecipeStep)
		assert.Equal(t, updated.IngredientNotes, x.IngredientNotes)
	})
}

func TestRecipeStepIngredientCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationInput{
			IngredientID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			Name:                fake.Word(),
			QuantityType:        fake.Word(),
			QuantityValue:       fake.Float32(),
			QuantityNotes:       fake.Word(),
			ProductOfRecipeStep: fake.Bool(),
			IngredientNotes:     fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeStepIngredientUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateInput{
			IngredientID:        func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
			Name:                fake.Word(),
			QuantityType:        fake.Word(),
			QuantityValue:       fake.Float32(),
			QuantityNotes:       fake.Word(),
			ProductOfRecipeStep: fake.Bool(),
			IngredientNotes:     fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeStepIngredientUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
