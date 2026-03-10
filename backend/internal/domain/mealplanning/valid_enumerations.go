package mealplanning

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/uploadedmedia"
)

type (
	// UploadedMediaFetcher fetches uploaded media by IDs (used for enriching preparations/ingredients with media).
	UploadedMediaFetcher interface {
		GetUploadedMediaWithIDs(ctx context.Context, ids []string) ([]*uploadedmedia.UploadedMedia, error)
	}

	ValidEnumerationDataManager interface {
		ValidIngredientGroupDataManager
		ValidIngredientMeasurementUnitDataManager
		ValidIngredientPreparationDataManager
		ValidPrepTaskConfigDataManager
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
		PreparationMediaDataManager
		IngredientMediaDataManager
		UploadedMediaFetcher
	}

	ValidEnumerationDataService interface {
		ValidIngredientGroupDataService
		ValidIngredientMeasurementUnitDataService
		ValidIngredientPreparationDataService
		ValidPrepTaskConfigDataService
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
