package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidPreparationInstrument_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrument{}

		updated := &ValidPreparationInstrumentUpdateInput{
			InstrumentID:  uint64(fake.Uint32()),
			PreparationID: uint64(fake.Uint32()),
			Notes:         fake.Word(),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "InstrumentID",
				OldValue:  x.InstrumentID,
				NewValue:  updated.InstrumentID,
			},
			{
				FieldName: "PreparationID",
				OldValue:  x.PreparationID,
				NewValue:  updated.PreparationID,
			},
			{
				FieldName: "Notes",
				OldValue:  x.Notes,
				NewValue:  updated.Notes,
			},
		}
		actual := x.Update(updated)

		expectedJSONBytes, err := json.Marshal(expected)
		require.NoError(t, err)

		actualJSONBytes, err := json.Marshal(actual)
		require.NoError(t, err)

		expectedJSON, actualJSON := string(expectedJSONBytes), string(actualJSONBytes)

		assert.Equal(t, expectedJSON, actualJSON)

		assert.Equal(t, updated.InstrumentID, x.InstrumentID)
		assert.Equal(t, updated.PreparationID, x.PreparationID)
		assert.Equal(t, updated.Notes, x.Notes)
	})
}

func TestValidPreparationInstrumentCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentCreationInput{
			InstrumentID:  uint64(fake.Uint32()),
			PreparationID: uint64(fake.Uint32()),
			Notes:         fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationInstrumentUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentUpdateInput{
			InstrumentID:  uint64(fake.Uint32()),
			PreparationID: uint64(fake.Uint32()),
			Notes:         fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationInstrumentUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
