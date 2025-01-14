package mocktypes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/stretchr/testify/mock"
)

type (
	MealPlanningDataManagerMock struct {
		*MealDataManagerMock
		*MealPlanDataManagerMock
		*MealPlanEventDataManagerMock
		*MealPlanOptionDataManagerMock
		*MealPlanOptionVoteDataManagerMock
		*MealPlanTaskDataManagerMock
		*MealPlanGroceryListItemDataManagerMock
		*UserIngredientPreferenceDataManagerMock
		*HouseholdInstrumentOwnershipDataManagerMock
	}
)

var _ types.MealPlanningDataManager = (*MealPlanningDataManagerMock)(nil)

func NewMealPlanningDataManagerMock() *MealPlanningDataManagerMock {
	return &MealPlanningDataManagerMock{
		MealDataManagerMock:                         &MealDataManagerMock{},
		MealPlanDataManagerMock:                     &MealPlanDataManagerMock{},
		MealPlanEventDataManagerMock:                &MealPlanEventDataManagerMock{},
		MealPlanOptionDataManagerMock:               &MealPlanOptionDataManagerMock{},
		MealPlanOptionVoteDataManagerMock:           &MealPlanOptionVoteDataManagerMock{},
		MealPlanTaskDataManagerMock:                 &MealPlanTaskDataManagerMock{},
		MealPlanGroceryListItemDataManagerMock:      &MealPlanGroceryListItemDataManagerMock{},
		UserIngredientPreferenceDataManagerMock:     &UserIngredientPreferenceDataManagerMock{},
		HouseholdInstrumentOwnershipDataManagerMock: &HouseholdInstrumentOwnershipDataManagerMock{},
	}
}

func (m *MealPlanningDataManagerMock) AssertExpectations(t mock.TestingT) bool {
	return mock.AssertExpectationsForObjects(t,
		m.MealDataManagerMock,
		m.MealPlanDataManagerMock,
		m.MealPlanEventDataManagerMock,
		m.MealPlanOptionDataManagerMock,
		m.MealPlanOptionVoteDataManagerMock,
		m.MealPlanTaskDataManagerMock,
		m.MealPlanGroceryListItemDataManagerMock,
		m.UserIngredientPreferenceDataManagerMock,
		m.HouseholdInstrumentOwnershipDataManagerMock,
	)
}
