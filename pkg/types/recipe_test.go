package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecipe_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		updated := &RecipeUpdateInput{
			Name:               fake.Word(),
			Source:             fake.Word(),
			Description:        fake.Word(),
			InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "Name",
				OldValue:  x.Name,
				NewValue:  updated.Name,
			},
			{
				FieldName: "Source",
				OldValue:  x.Source,
				NewValue:  updated.Source,
			},
			{
				FieldName: "Description",
				OldValue:  x.Description,
				NewValue:  updated.Description,
			},
			{
				FieldName: "InspiredByRecipeID",
				OldValue:  x.InspiredByRecipeID,
				NewValue:  updated.InspiredByRecipeID,
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
		assert.Equal(t, updated.Source, x.Source)
		assert.Equal(t, updated.Description, x.Description)
		assert.Equal(t, updated.InspiredByRecipeID, x.InspiredByRecipeID)
	})
}

func TestRecipeCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationInput{
			Name:               fake.Word(),
			Source:             fake.Word(),
			Description:        fake.Word(),
			InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateInput{
			Name:               fake.Word(),
			Source:             fake.Word(),
			Description:        fake.Word(),
			InspiredByRecipeID: func(x uint64) *uint64 { return &x }(uint64(fake.Uint32())),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
