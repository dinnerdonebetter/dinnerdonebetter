package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidPreparation_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparation{}

		updated := &ValidPreparationUpdateInput{
			Name:        fake.Word(),
			Description: fake.Word(),
			IconPath:    fake.Word(),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "Name",
				OldValue:  x.Name,
				NewValue:  updated.Name,
			},
			{
				FieldName: "Description",
				OldValue:  x.Description,
				NewValue:  updated.Description,
			},
			{
				FieldName: "IconPath",
				OldValue:  x.IconPath,
				NewValue:  updated.IconPath,
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
		assert.Equal(t, updated.Description, x.Description)
		assert.Equal(t, updated.IconPath, x.IconPath)
	})
}

func TestValidPreparationCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationInput{
			Name:        fake.Word(),
			Description: fake.Word(),
			IconPath:    fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidPreparationUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateInput{
			Name:        fake.Word(),
			Description: fake.Word(),
			IconPath:    fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidPreparationUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
