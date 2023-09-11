package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v6"
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
			Status:                 MealPlanGroceryListItemStatusUnknown,
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
			Status:                 MealPlanGroceryListItemStatusUnknown,
			MaximumQuantityNeeded:  pointers.Pointer(float32(1.23)),
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
			BelongsToMealPlan:      pointers.Pointer(t.Name()),
			ValidIngredientID:      pointers.Pointer(t.Name()),
			ValidMeasurementUnitID: pointers.Pointer(t.Name()),
			MinimumQuantityNeeded:  pointers.Pointer(float32(1.23)),
			MaximumQuantityNeeded:  pointers.Pointer(float32(1.23)),
			Status:                 pointers.Pointer(MealPlanGroceryListItemStatusUnknown),
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
			MaximumQuantityNeeded:    pointers.Pointer(float32(1.23)),
		}
		input := &MealPlanGroceryListItemUpdateRequestInput{}

		fake.Struct(&input)
		input.PurchasedMeasurementUnitID = pointers.Pointer(t.Name())
		input.MaximumQuantityNeeded = pointers.Pointer(float32(3.21))

		x.Update(input)
	})
}
