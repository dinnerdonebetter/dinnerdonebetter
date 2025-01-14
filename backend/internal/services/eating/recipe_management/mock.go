package recipemanagement

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	mocktypes "github.com/dinnerdonebetter/backend/pkg/types/mock"

	"github.com/stretchr/testify/mock"
)

type RecipeManagementDataManagerMock struct {
	*mocktypes.RecipeDataManagerMock
	*mocktypes.RecipeMediaDataManagerMock
	*mocktypes.RecipeStepCompletionConditionDataManagerMock
	*mocktypes.RecipeStepIngredientDataManagerMock
	*mocktypes.RecipeStepInstrumentDataManagerMock
	*mocktypes.RecipeStepProductDataManagerMock
	*mocktypes.RecipeStepDataManagerMock
	*mocktypes.RecipeStepVesselDataManagerMock
	*mocktypes.RecipeRatingDataManagerMock
	*mocktypes.RecipePrepTaskDataManagerMock
}

var _ types.RecipeManagementDataManager = (*RecipeManagementDataManagerMock)(nil)

func NewRecipeManagementDataManagerMock() *RecipeManagementDataManagerMock {
	return &RecipeManagementDataManagerMock{
		RecipeDataManagerMock:                        &mocktypes.RecipeDataManagerMock{},
		RecipeMediaDataManagerMock:                   &mocktypes.RecipeMediaDataManagerMock{},
		RecipeStepCompletionConditionDataManagerMock: &mocktypes.RecipeStepCompletionConditionDataManagerMock{},
		RecipeStepIngredientDataManagerMock:          &mocktypes.RecipeStepIngredientDataManagerMock{},
		RecipeStepInstrumentDataManagerMock:          &mocktypes.RecipeStepInstrumentDataManagerMock{},
		RecipeStepProductDataManagerMock:             &mocktypes.RecipeStepProductDataManagerMock{},
		RecipeStepDataManagerMock:                    &mocktypes.RecipeStepDataManagerMock{},
		RecipeStepVesselDataManagerMock:              &mocktypes.RecipeStepVesselDataManagerMock{},
		RecipeRatingDataManagerMock:                  &mocktypes.RecipeRatingDataManagerMock{},
		RecipePrepTaskDataManagerMock:                &mocktypes.RecipePrepTaskDataManagerMock{},
	}
}

func (m *RecipeManagementDataManagerMock) AssertExpectations(t mock.TestingT) bool {
	return mock.AssertExpectationsForObjects(t,
		m.RecipeDataManagerMock,
		m.RecipeMediaDataManagerMock,
		m.RecipeStepCompletionConditionDataManagerMock,
		m.RecipeStepIngredientDataManagerMock,
		m.RecipeStepInstrumentDataManagerMock,
		m.RecipeStepProductDataManagerMock,
		m.RecipeStepDataManagerMock,
		m.RecipeStepVesselDataManagerMock,
		m.RecipeRatingDataManagerMock,
		m.RecipePrepTaskDataManagerMock,
	)
}
