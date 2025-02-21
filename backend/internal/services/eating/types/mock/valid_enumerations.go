package mocktypes

import (
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"

	"github.com/stretchr/testify/mock"
)

type ValidEnumerationDataManagerMock struct {
	*ValidIngredientGroupDataManagerMock
	*ValidIngredientMeasurementUnitDataManagerMock
	*ValidIngredientPreparationDataManagerMock
	*ValidIngredientDataManagerMock
	*ValidIngredientStateIngredientDataManagerMock
	*ValidIngredientStateDataManagerMock
	*ValidMeasurementUnitDataManagerMock
	*ValidInstrumentDataManagerMock
	*ValidMeasurementUnitConversionDataManagerMock
	*ValidPreparationInstrumentDataManagerMock
	*ValidPreparationDataManagerMock
	*ValidPreparationVesselDataManagerMock
	*ValidVesselDataManagerMock
}

var _ types.ValidEnumerationDataManager = (*ValidEnumerationDataManagerMock)(nil)

func NewValidEnumerationDataManagerMock() *ValidEnumerationDataManagerMock {
	return &ValidEnumerationDataManagerMock{
		ValidIngredientGroupDataManagerMock:           &ValidIngredientGroupDataManagerMock{},
		ValidIngredientMeasurementUnitDataManagerMock: &ValidIngredientMeasurementUnitDataManagerMock{},
		ValidIngredientPreparationDataManagerMock:     &ValidIngredientPreparationDataManagerMock{},
		ValidIngredientDataManagerMock:                &ValidIngredientDataManagerMock{},
		ValidIngredientStateIngredientDataManagerMock: &ValidIngredientStateIngredientDataManagerMock{},
		ValidIngredientStateDataManagerMock:           &ValidIngredientStateDataManagerMock{},
		ValidMeasurementUnitDataManagerMock:           &ValidMeasurementUnitDataManagerMock{},
		ValidInstrumentDataManagerMock:                &ValidInstrumentDataManagerMock{},
		ValidMeasurementUnitConversionDataManagerMock: &ValidMeasurementUnitConversionDataManagerMock{},
		ValidPreparationInstrumentDataManagerMock:     &ValidPreparationInstrumentDataManagerMock{},
		ValidPreparationDataManagerMock:               &ValidPreparationDataManagerMock{},
		ValidPreparationVesselDataManagerMock:         &ValidPreparationVesselDataManagerMock{},
		ValidVesselDataManagerMock:                    &ValidVesselDataManagerMock{},
	}
}

func (m *ValidEnumerationDataManagerMock) AssertExpectations(t mock.TestingT) bool {
	return mock.AssertExpectationsForObjects(t,
		m.ValidIngredientGroupDataManagerMock,
		m.ValidIngredientMeasurementUnitDataManagerMock,
		m.ValidIngredientPreparationDataManagerMock,
		m.ValidIngredientDataManagerMock,
		m.ValidIngredientStateIngredientDataManagerMock,
		m.ValidIngredientStateDataManagerMock,
		m.ValidMeasurementUnitDataManagerMock,
		m.ValidInstrumentDataManagerMock,
		m.ValidMeasurementUnitConversionDataManagerMock,
		m.ValidPreparationInstrumentDataManagerMock,
		m.ValidPreparationDataManagerMock,
		m.ValidPreparationVesselDataManagerMock,
		m.ValidVesselDataManagerMock,
	)
}
