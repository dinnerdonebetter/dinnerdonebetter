package types

import (
	"context"
	"testing"

	"github.com/prixfixeco/backend/internal/pointers"

	"github.com/stretchr/testify/assert"
)

func TestMealPlanGroceryListItemCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanGroceryListItemCreationRequestInput{
			BelongsToMealPlan:      t.Name(),
			ValidIngredientID:      t.Name(),
			ValidMeasurementUnitID: t.Name(),
			MinimumQuantityNeeded:  1.23,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItemDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanGroceryListItemDatabaseCreationInput{
			ID:                     t.Name(),
			BelongsToMealPlan:      t.Name(),
			ValidIngredientID:      t.Name(),
			ValidMeasurementUnitID: t.Name(),
			MinimumQuantityNeeded:  1.23,
			MaximumQuantityNeeded:  pointers.Float32(1.23),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItemUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		x := &MealPlanGroceryListItemUpdateRequestInput{
			BelongsToMealPlan:      pointers.String(t.Name()),
			ValidIngredientID:      pointers.String(t.Name()),
			ValidMeasurementUnitID: pointers.String(t.Name()),
			MinimumQuantityNeeded:  pointers.Float32(1.23),
			MaximumQuantityNeeded:  pointers.Float32(1.23),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItem_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		actual := &MealPlanGroceryListItem{}
		expected := &MealPlanGroceryListItem{}
		input := &MealPlanGroceryListItemUpdateRequestInput{}

		actual.Update(input)

		assert.Equal(t, expected, actual)
	})
}
