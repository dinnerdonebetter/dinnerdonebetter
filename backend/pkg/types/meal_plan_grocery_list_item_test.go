package types

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

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
			QuantityNeeded:         Float32RangeWithOptionalMax{Min: 1.23},
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
			QuantityNeeded: Float32RangeWithOptionalMax{
				Min: 1.23,
				Max: pointer.To(float32(1.23)),
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
			BelongsToMealPlan:      pointer.To(t.Name()),
			ValidIngredientID:      pointer.To(t.Name()),
			ValidMeasurementUnitID: pointer.To(t.Name()),
			QuantityNeeded: Float32RangeWithOptionalMaxUpdateRequestInput{
				Min: pointer.To(float32(1.23)),
				Max: pointer.To(float32(1.23)),
			},
			Status: pointer.To(MealPlanGroceryListItemStatusUnknown),
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
			QuantityNeeded:           Float32RangeWithOptionalMax{Max: pointer.To(float32(1.23))},
		}
		input := &MealPlanGroceryListItemUpdateRequestInput{}

		assert.NoError(t, fake.Struct(&input))
		input.PurchasedMeasurementUnitID = pointer.To(t.Name())
		input.QuantityNeeded.Max = pointer.To(float32(3.21))

		x.Update(input)
	})
}
