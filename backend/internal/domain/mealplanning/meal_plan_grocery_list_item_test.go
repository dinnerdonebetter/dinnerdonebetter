package mealplanning

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/platform/types"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestMealPlanGroceryListItemCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealPlanGroceryListItemCreationRequestInput{
			BelongsToMealPlan:      t.Name(),
			ValidIngredientID:      t.Name(),
			ValidMeasurementUnitID: t.Name(),
			Status:                 MealPlanGroceryListItemStatusUnknown,
			QuantityNeeded:         types.Float32RangeWithOptionalMax{Min: 1.23},
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItemDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealPlanGroceryListItemDatabaseCreationInput{
			ID:                     t.Name(),
			BelongsToMealPlan:      t.Name(),
			ValidIngredientID:      t.Name(),
			ValidMeasurementUnitID: t.Name(),
			QuantityNeeded: types.Float32RangeWithOptionalMax{
				Min: 1.23,
				Max: new(float32(1.23)),
			},
			Status: MealPlanGroceryListItemStatusUnknown,
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItemUpdateRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := &MealPlanGroceryListItemUpdateRequestInput{
			BelongsToMealPlan:      new(t.Name()),
			ValidIngredientID:      new(t.Name()),
			ValidMeasurementUnitID: new(t.Name()),
			QuantityNeeded: types.Float32RangeWithOptionalMaxUpdateRequestInput{
				Min: new(float32(1.23)),
				Max: new(float32(1.23)),
			},
			Status: new(MealPlanGroceryListItemStatusUnknown),
		}

		assert.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanGroceryListItem_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanGroceryListItem{
			PurchasedMeasurementUnit: &ValidMeasurementUnit{},
			QuantityNeeded:           types.Float32RangeWithOptionalMax{Max: new(float32(1.23))},
		}
		input := &MealPlanGroceryListItemUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.PurchasedMeasurementUnitID = new(t.Name())
		input.QuantityNeeded.Max = new(float32(3.21))

		x.Update(input)
	})
}
