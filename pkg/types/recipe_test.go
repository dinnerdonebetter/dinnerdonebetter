package types

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"

	fake "github.com/brianvoe/gofakeit/v5"
	"github.com/stretchr/testify/assert"
)

func TestRecipe_FindStepForIndex(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: t.Name(),
				},
			},
		}

		assert.NotNil(t, x.FindStepForIndex(0))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Nil(t, x.FindStepForIndex(0))
	})
}

func TestRecipe_FindStepByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: t.Name(),
				},
			},
		}

		assert.NotNil(t, x.FindStepByID(t.Name()))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Nil(t, x.FindStepByID("whatever"))
	})
}

func TestRecipe_FindStepForRecipeStepProductID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: "whatever",
					Products: []*RecipeStepProduct{
						{
							ID: t.Name(),
						},
					},
				},
			},
		}

		assert.NotNil(t, x.FindStepForRecipeStepProductID(t.Name()))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Nil(t, x.FindStepForRecipeStepProductID("whatever"))
	})
}

func TestRecipe_FindStepIndexByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{
			Steps: []*RecipeStep{
				{
					ID: t.Name(),
				},
			},
		}

		x.FindStepIndexByID(t.Name())
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}

		assert.Equal(t, -1, x.FindStepIndexByID("whatever"))
	})
}

func TestRecipe_Update(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &Recipe{}
		input := &RecipeUpdateRequestInput{}

		fake.Struct(&input)
		input.SealOfApproval = pointers.Pointer(true)
		input.EligibleForMeals = pointers.Pointer(true)

		x.Update(input)
	})
}

func TestRecipeCreationRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{
			Name:                t.Name(),
			Source:              t.Name(),
			Slug:                t.Name(),
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			Description:         t.Name(),
			YieldsComponentType: MealComponentTypesMain,
			InspiredByRecipeID:  pointers.Pointer(t.Name()),
			Steps: []*RecipeStepCreationRequestInput{
				buildValidRecipeStepCreationRequestInput(),
				buildValidRecipeStepCreationRequestInput(),
			},
			PrepTasks: []*RecipePrepTaskWithinRecipeCreationRequestInput{
				{
					TaskSteps: []*RecipePrepTaskStepWithinRecipeCreationRequestInput{
						{
							BelongsToRecipeStepIndex: 0,
						},
					},
				},
			},
			SealOfApproval:           fake.Bool(),
			MinimumEstimatedPortions: fake.Float32(),
		}

		assert.NoError(t, x.ValidateWithContext(context.Background()))
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})

	T.Run("with overreferenced task steps", func(t *testing.T) {
		t.Parallel()

		x := &RecipeCreationRequestInput{
			Name:                t.Name(),
			Source:              t.Name(),
			Slug:                t.Name(),
			PortionName:         t.Name(),
			PluralPortionName:   t.Name(),
			Description:         t.Name(),
			YieldsComponentType: MealComponentTypesMain,
			InspiredByRecipeID:  pointers.Pointer(t.Name()),
			Steps: []*RecipeStepCreationRequestInput{
				buildValidRecipeStepCreationRequestInput(),
				buildValidRecipeStepCreationRequestInput(),
			},
			PrepTasks: []*RecipePrepTaskWithinRecipeCreationRequestInput{
				{
					TaskSteps: []*RecipePrepTaskStepWithinRecipeCreationRequestInput{
						{
							BelongsToRecipeStepIndex: 0,
						},
						{
							BelongsToRecipeStepIndex: 0,
						},
					},
				},
			},
			SealOfApproval:           fake.Bool(),
			MinimumEstimatedPortions: fake.Float32(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeDatabaseCreationInput_FindStepByIndex(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Index: 0,
				},
			},
		}

		assert.NotNil(t, x.FindStepByIndex(0))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					Index: 0,
				},
			},
		}

		assert.Nil(t, x.FindStepByIndex(1))
	})
}

func TestRecipeDatabaseCreationInput_FindStepByID(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: t.Name(),
				},
			},
		}

		assert.NotNil(t, x.FindStepByID(t.Name()))
	})

	T.Run("with invalid value", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			Steps: []*RecipeStepDatabaseCreationInput{
				{
					ID: t.Name(),
				},
			},
		}

		assert.Nil(t, x.FindStepByID("whatever"))
	})
}

func TestRecipeDatabaseCreationInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{
			ID:            t.Name(),
			Name:          t.Name(),
			CreatedByUser: t.Name(),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with invalid structure", func(t *testing.T) {
		t.Parallel()

		x := &RecipeDatabaseCreationInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}

func TestRecipeUpdateRequestInput_Validate(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateRequestInput{
			Name:                     pointers.Pointer(t.Name()),
			Source:                   pointers.Pointer(t.Name()),
			Description:              pointers.Pointer(t.Name()),
			InspiredByRecipeID:       pointers.Pointer(t.Name()),
			SealOfApproval:           pointers.Pointer(fake.Bool()),
			MinimumEstimatedPortions: pointers.Pointer(fake.Float32()),
		}

		actual := x.ValidateWithContext(context.Background())
		assert.Nil(t, actual)
	})

	T.Run("with empty strings", func(t *testing.T) {
		t.Parallel()

		x := &RecipeUpdateRequestInput{}

		actual := x.ValidateWithContext(context.Background())
		assert.Error(t, actual)
	})
}
