package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidIngredientPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparation{}

		updated := &ValidIngredientPreparationUpdateInput{
			Notes:              fake.Word(),
			ValidIngredientID:  uint64(fake.Uint32()),
			ValidPreparationID: uint64(fake.Uint32()),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "Notes",
				OldValue:  x.Notes,
				NewValue:  updated.Notes,
			},
			{
				FieldName: "ValidIngredientID",
				OldValue:  x.ValidIngredientID,
				NewValue:  updated.ValidIngredientID,
			},
			{
				FieldName: "ValidPreparationID",
				OldValue:  x.ValidPreparationID,
				NewValue:  updated.ValidPreparationID,
			},
		}
		actual := x.Update(updated)

		expectedJSONBytes, err := json.Marshal(expected)
		require.NoError(t, err)

		actualJSONBytes, err := json.Marshal(actual)
		require.NoError(t, err)

		expectedJSON, actualJSON := string(expectedJSONBytes), string(actualJSONBytes)

		assert.Equal(t, expectedJSON, actualJSON)

		assert.Equal(t, updated.Notes, x.Notes)
		assert.Equal(t, updated.ValidIngredientID, x.ValidIngredientID)
		assert.Equal(t, updated.ValidPreparationID, x.ValidPreparationID)
	})
}

func TestValidIngredientPreparationCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationCreationInput{
			Notes:              fake.Word(),
			ValidIngredientID:  uint64(fake.Uint32()),
			ValidPreparationID: uint64(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientPreparationUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationUpdateInput{
			Notes:              fake.Word(),
			ValidIngredientID:  uint64(fake.Uint32()),
			ValidPreparationID: uint64(fake.Uint32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientPreparationUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
