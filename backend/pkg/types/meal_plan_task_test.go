package types

import (
	"sort"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/lib/pointer"

	fake "github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMealPlanTaskListSort(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		expected := MealPlanTaskList{
			{
				ID: "a",
			},
			{
				ID: "b",
			},
			{
				ID: "c",
			},
		}

		actual := MealPlanTaskList{
			expected[2],
			expected[0],
			expected[1],
		}

		sort.Sort(actual)

		assert.Equal(t, expected, actual)
	})
}

func TestMealPlanTask_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &MealPlanTask{}
		input := &MealPlanTaskStatusChangeRequestInput{}

		assert.NoError(t, fake.Struct(&input))

		x.Update(input)
	})
}

func TestMealPlanTaskCreationRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := MealPlanTaskCreationRequestInput{
			MealPlanOptionID: t.Name(),
			RecipePrepTaskID: t.Name(),
			Status:           MealPlanTaskStatusUnfinished,
		}

		require.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanTaskDatabaseCreationInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := MealPlanTaskDatabaseCreationInput{
			ID:               t.Name(),
			MealPlanOptionID: t.Name(),
			RecipePrepTaskID: t.Name(),
		}

		require.NoError(t, x.ValidateWithContext(ctx))
	})
}

func TestMealPlanTaskStatusChangeRequestInput_ValidateWithContext(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		x := MealPlanTaskStatusChangeRequestInput{
			ID:     t.Name(),
			Status: pointer.To(MealPlanTaskStatusUnfinished),
		}

		require.NoError(t, x.ValidateWithContext(ctx))
	})
}
