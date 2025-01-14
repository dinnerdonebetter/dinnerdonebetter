package types

type (
	ValidEnumerationDataManager interface {
		ValidIngredientGroupDataManager
		ValidIngredientMeasurementUnitDataManager
		ValidIngredientPreparationDataManager
		ValidIngredientDataManager
		ValidIngredientStateIngredientDataManager
		ValidIngredientStateDataManager
		ValidMeasurementUnitDataManager
		ValidInstrumentDataManager
		ValidMeasurementUnitConversionDataManager
		ValidPreparationInstrumentDataManager
		ValidPreparationDataManager
		ValidPreparationVesselDataManager
		ValidVesselDataManager
	}

	ValidEnumerationDataService interface {
		ValidIngredientGroupDataService
		ValidIngredientMeasurementUnitDataService
		ValidIngredientPreparationDataService
		ValidIngredientDataService
		ValidIngredientStateIngredientDataService
		ValidIngredientStateDataService
		ValidMeasurementUnitDataService
		ValidInstrumentDataService
		ValidMeasurementUnitConversionDataService
		ValidPreparationInstrumentDataService
		ValidPreparationDataService
		ValidPreparationVesselDataService
		ValidVesselDataService
	}
)
