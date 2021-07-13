package types

import (
	"context"
	"encoding/json"
	"testing"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidIngredient_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredient{}

		updated := &ValidIngredientUpdateInput{
			Name:              fake.Word(),
			Variant:           fake.Word(),
			Description:       fake.Word(),
			Warning:           fake.Word(),
			ContainsEgg:       true,
			ContainsDairy:     true,
			ContainsPeanut:    true,
			ContainsTreeNut:   true,
			ContainsSoy:       true,
			ContainsWheat:     true,
			ContainsShellfish: true,
			ContainsSesame:    true,
			ContainsFish:      true,
			ContainsGluten:    true,
			AnimalFlesh:       true,
			AnimalDerived:     true,
			Volumetric:        true,
			IconPath:          fake.Word(),
		}

		expected := []*FieldChangeSummary{
			{
				FieldName: "Name",
				OldValue:  x.Name,
				NewValue:  updated.Name,
			},
			{
				FieldName: "Variant",
				OldValue:  x.Variant,
				NewValue:  updated.Variant,
			},
			{
				FieldName: "Description",
				OldValue:  x.Description,
				NewValue:  updated.Description,
			},
			{
				FieldName: "Warning",
				OldValue:  x.Warning,
				NewValue:  updated.Warning,
			},
			{
				FieldName: "ContainsEgg",
				OldValue:  x.ContainsEgg,
				NewValue:  updated.ContainsEgg,
			},
			{
				FieldName: "ContainsDairy",
				OldValue:  x.ContainsDairy,
				NewValue:  updated.ContainsDairy,
			},
			{
				FieldName: "ContainsPeanut",
				OldValue:  x.ContainsPeanut,
				NewValue:  updated.ContainsPeanut,
			},
			{
				FieldName: "ContainsTreeNut",
				OldValue:  x.ContainsTreeNut,
				NewValue:  updated.ContainsTreeNut,
			},
			{
				FieldName: "ContainsSoy",
				OldValue:  x.ContainsSoy,
				NewValue:  updated.ContainsSoy,
			},
			{
				FieldName: "ContainsWheat",
				OldValue:  x.ContainsWheat,
				NewValue:  updated.ContainsWheat,
			},
			{
				FieldName: "ContainsShellfish",
				OldValue:  x.ContainsShellfish,
				NewValue:  updated.ContainsShellfish,
			},
			{
				FieldName: "ContainsSesame",
				OldValue:  x.ContainsSesame,
				NewValue:  updated.ContainsSesame,
			},
			{
				FieldName: "ContainsFish",
				OldValue:  x.ContainsFish,
				NewValue:  updated.ContainsFish,
			},
			{
				FieldName: "ContainsGluten",
				OldValue:  x.ContainsGluten,
				NewValue:  updated.ContainsGluten,
			},
			{
				FieldName: "AnimalFlesh",
				OldValue:  x.AnimalFlesh,
				NewValue:  updated.AnimalFlesh,
			},
			{
				FieldName: "AnimalDerived",
				OldValue:  x.AnimalDerived,
				NewValue:  updated.AnimalDerived,
			},
			{
				FieldName: "Volumetric",
				OldValue:  x.Volumetric,
				NewValue:  updated.Volumetric,
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
		assert.Equal(t, updated.Variant, x.Variant)
		assert.Equal(t, updated.Description, x.Description)
		assert.Equal(t, updated.Warning, x.Warning)
		assert.Equal(t, updated.ContainsEgg, x.ContainsEgg)
		assert.Equal(t, updated.ContainsDairy, x.ContainsDairy)
		assert.Equal(t, updated.ContainsPeanut, x.ContainsPeanut)
		assert.Equal(t, updated.ContainsTreeNut, x.ContainsTreeNut)
		assert.Equal(t, updated.ContainsSoy, x.ContainsSoy)
		assert.Equal(t, updated.ContainsWheat, x.ContainsWheat)
		assert.Equal(t, updated.ContainsShellfish, x.ContainsShellfish)
		assert.Equal(t, updated.ContainsSesame, x.ContainsSesame)
		assert.Equal(t, updated.ContainsFish, x.ContainsFish)
		assert.Equal(t, updated.ContainsGluten, x.ContainsGluten)
		assert.Equal(t, updated.AnimalFlesh, x.AnimalFlesh)
		assert.Equal(t, updated.AnimalDerived, x.AnimalDerived)
		assert.Equal(t, updated.Volumetric, x.Volumetric)
		assert.Equal(t, updated.IconPath, x.IconPath)
	})
}

func TestValidIngredientCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientCreationInput{
			Name:              fake.Word(),
			Variant:           fake.Word(),
			Description:       fake.Word(),
			Warning:           fake.Word(),
			ContainsEgg:       fake.Bool(),
			ContainsDairy:     fake.Bool(),
			ContainsPeanut:    fake.Bool(),
			ContainsTreeNut:   fake.Bool(),
			ContainsSoy:       fake.Bool(),
			ContainsWheat:     fake.Bool(),
			ContainsShellfish: fake.Bool(),
			ContainsSesame:    fake.Bool(),
			ContainsFish:      fake.Bool(),
			ContainsGluten:    fake.Bool(),
			AnimalFlesh:       fake.Bool(),
			AnimalDerived:     fake.Bool(),
			Volumetric:        fake.Bool(),
			IconPath:          fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestValidIngredientUpdateInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientUpdateInput{
			Name:              fake.Word(),
			Variant:           fake.Word(),
			Description:       fake.Word(),
			Warning:           fake.Word(),
			ContainsEgg:       fake.Bool(),
			ContainsDairy:     fake.Bool(),
			ContainsPeanut:    fake.Bool(),
			ContainsTreeNut:   fake.Bool(),
			ContainsSoy:       fake.Bool(),
			ContainsWheat:     fake.Bool(),
			ContainsShellfish: fake.Bool(),
			ContainsSesame:    fake.Bool(),
			ContainsFish:      fake.Bool(),
			ContainsGluten:    fake.Bool(),
			AnimalFlesh:       fake.Bool(),
			AnimalDerived:     fake.Bool(),
			Volumetric:        fake.Bool(),
			IconPath:          fake.Word(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &ValidIngredientUpdateInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
