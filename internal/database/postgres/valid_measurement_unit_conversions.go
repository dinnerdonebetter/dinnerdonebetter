package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/pkg/pointers"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.ValidMeasurementUnitConversionDataManager = (*Querier)(nil)
)

// ValidMeasurementUnitConversionExists fetches whether a valid measurement conversion exists from the database.
func (q *Querier) ValidMeasurementUnitConversionExists(ctx context.Context, validMeasurementUnitConversionID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	result, err := q.generatedQuerier.CheckValidMeasurementUnitConversionExistence(ctx, q.db, validMeasurementUnitConversionID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid measurement conversion existence check")
	}

	return result, nil
}

// GetValidMeasurementUnitConversion fetches a valid measurement conversion from the database.
func (q *Querier) GetValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	result, err := q.generatedQuerier.GetValidMeasurementUnitConversion(ctx, q.db, validMeasurementUnitConversionID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement conversion")
	}

	validMeasurementUnitConversion := &types.ValidMeasurementUnitConversion{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
		OnlyForIngredient: &types.ValidIngredient{
			CreatedAt:                               result.ValidIngredientCreatedAt.Time,
			LastUpdatedAt:                           &result.ValidIngredientLastUpdatedAt.Time,
			ArchivedAt:                              &result.ValidIngredientArchivedAt.Time,
			MaximumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius)),
			MinimumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius)),
			IconPath:                                result.ValidIngredientIconPath.String,
			Warning:                                 result.ValidIngredientWarning.String,
			PluralName:                              result.ValidIngredientPluralName.String,
			StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
			Name:                                    result.ValidIngredientName.String,
			ID:                                      result.ValidIngredientID.String,
			Description:                             result.ValidIngredientDescription.String,
			Slug:                                    result.ValidIngredientSlug.String,
			ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
			ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
			IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
			IsLiquid:                                result.ValidIngredientIsLiquid.Bool,
			ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
			ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
			ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
			ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
			ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
			AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
			RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
			ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
			ContainsFish:                            result.ValidIngredientContainsFish.Bool,
			ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
			ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
			ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
			AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
			IsStarch:                                result.ValidIngredientIsStarch.Bool,
			IsProtein:                               result.ValidIngredientIsProtein.Bool,
			IsGrain:                                 result.ValidIngredientIsGrain.Bool,
			IsFruit:                                 result.ValidIngredientIsFruit.Bool,
			IsSalt:                                  result.ValidIngredientIsSalt.Bool,
			IsFat:                                   result.ValidIngredientIsFat.Bool,
			IsAcid:                                  result.ValidIngredientIsAcid.Bool,
			IsHeat:                                  result.ValidIngredientIsHeat.Bool,
		},
		Notes: result.Notes,
		ID:    result.ID,
		From: types.ValidMeasurementUnit{
			CreatedAt:     result.FromUnitCreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.FromUnitLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.FromUnitArchivedAt),
			Name:          result.FromUnitName,
			IconPath:      result.FromUnitIconPath,
			ID:            result.FromUnitID,
			Description:   result.FromUnitDescription,
			PluralName:    result.FromUnitPluralName,
			Slug:          result.FromUnitSlug,
			Volumetric:    boolFromNullBool(result.FromUnitVolumetric),
			Universal:     result.FromUnitUniversal,
			Metric:        result.FromUnitMetric,
			Imperial:      result.FromUnitImperial,
		},
		To: types.ValidMeasurementUnit{
			CreatedAt:     result.ToUnitCreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.ToUnitLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ToUnitArchivedAt),
			Name:          result.ToUnitName,
			IconPath:      result.ToUnitIconPath,
			ID:            result.ToUnitID,
			Description:   result.ToUnitDescription,
			PluralName:    result.ToUnitPluralName,
			Slug:          result.ToUnitSlug,
			Volumetric:    boolFromNullBool(result.ToUnitVolumetric),
			Universal:     result.ToUnitUniversal,
			Metric:        result.ToUnitMetric,
			Imperial:      result.ToUnitImperial,
		},
		Modifier: float32(result.ValidMeasurementConversionsModifier),
	}

	return validMeasurementUnitConversion, nil
}

// GetValidMeasurementUnitConversionsFromUnit fetches a valid measurement conversions from a given measurement unit.
func (q *Querier) GetValidMeasurementUnitConversionsFromUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	results, err := q.generatedQuerier.GetAllValidMeasurementUnitConversionsFromMeasurementUnit(ctx, q.db, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for valid measurement conversions")
	}

	validMeasurementUnitConversions := make([]*types.ValidMeasurementUnitConversion, len(results))
	for i, result := range results {
		validMeasurementUnitConversions[i] = &types.ValidMeasurementUnitConversion{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
			OnlyForIngredient: &types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt:                           &result.ValidIngredientLastUpdatedAt.Time,
				ArchivedAt:                              &result.ValidIngredientArchivedAt.Time,
				MaximumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius)),
				MinimumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius)),
				IconPath:                                result.ValidIngredientIconPath.String,
				Warning:                                 result.ValidIngredientWarning.String,
				PluralName:                              result.ValidIngredientPluralName.String,
				StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
				Name:                                    result.ValidIngredientName.String,
				ID:                                      result.ValidIngredientID.String,
				Description:                             result.ValidIngredientDescription.String,
				Slug:                                    result.ValidIngredientSlug.String,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
				IsLiquid:                                result.ValidIngredientIsLiquid.Bool,
				ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
				ContainsFish:                            result.ValidIngredientContainsFish.Bool,
				ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:                                result.ValidIngredientIsStarch.Bool,
				IsProtein:                               result.ValidIngredientIsProtein.Bool,
				IsGrain:                                 result.ValidIngredientIsGrain.Bool,
				IsFruit:                                 result.ValidIngredientIsFruit.Bool,
				IsSalt:                                  result.ValidIngredientIsSalt.Bool,
				IsFat:                                   result.ValidIngredientIsFat.Bool,
				IsAcid:                                  result.ValidIngredientIsAcid.Bool,
				IsHeat:                                  result.ValidIngredientIsHeat.Bool,
			},
			Notes: result.Notes,
			ID:    result.ID,
			From: types.ValidMeasurementUnit{
				CreatedAt:     result.FromUnitCreatedAt,
				LastUpdatedAt: timePointerFromNullTime(result.FromUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.FromUnitArchivedAt),
				Name:          result.FromUnitName,
				IconPath:      result.FromUnitIconPath,
				ID:            result.FromUnitID,
				Description:   result.FromUnitDescription,
				PluralName:    result.FromUnitPluralName,
				Slug:          result.FromUnitSlug,
				Volumetric:    boolFromNullBool(result.FromUnitVolumetric),
				Universal:     result.FromUnitUniversal,
				Metric:        result.FromUnitMetric,
				Imperial:      result.FromUnitImperial,
			},
			To: types.ValidMeasurementUnit{
				CreatedAt:     result.ToUnitCreatedAt,
				LastUpdatedAt: timePointerFromNullTime(result.ToUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ToUnitArchivedAt),
				Name:          result.ToUnitName,
				IconPath:      result.ToUnitIconPath,
				ID:            result.ToUnitID,
				Description:   result.ToUnitDescription,
				PluralName:    result.ToUnitPluralName,
				Slug:          result.ToUnitSlug,
				Volumetric:    boolFromNullBool(result.ToUnitVolumetric),
				Universal:     result.ToUnitUniversal,
				Metric:        result.ToUnitMetric,
				Imperial:      result.ToUnitImperial,
			},
			Modifier: float32(result.ValidMeasurementConversionsModifier),
		}
	}

	return validMeasurementUnitConversions, nil
}

// GetValidMeasurementUnitConversionsToUnit fetches a valid measurement conversions to a given measurement unit.
func (q *Querier) GetValidMeasurementUnitConversionsToUnit(ctx context.Context, validMeasurementUnitID string) ([]*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachValidMeasurementUnitIDToSpan(span, validMeasurementUnitID)

	results, err := q.generatedQuerier.GetAllValidMeasurementUnitConversionsToMeasurementUnit(ctx, q.db, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "querying for valid measurement conversions")
	}

	validMeasurementUnitConversions := make([]*types.ValidMeasurementUnitConversion, len(results))
	for i, result := range results {
		validMeasurementUnitConversions[i] = &types.ValidMeasurementUnitConversion{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
			OnlyForIngredient: &types.ValidIngredient{
				CreatedAt:                               result.ValidIngredientCreatedAt.Time,
				LastUpdatedAt:                           &result.ValidIngredientLastUpdatedAt.Time,
				ArchivedAt:                              &result.ValidIngredientArchivedAt.Time,
				MaximumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius)),
				MinimumIdealStorageTemperatureInCelsius: pointers.Pointer(float32(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius)),
				IconPath:                                result.ValidIngredientIconPath.String,
				Warning:                                 result.ValidIngredientWarning.String,
				PluralName:                              result.ValidIngredientPluralName.String,
				StorageInstructions:                     result.ValidIngredientStorageInstructions.String,
				Name:                                    result.ValidIngredientName.String,
				ID:                                      result.ValidIngredientID.String,
				Description:                             result.ValidIngredientDescription.String,
				Slug:                                    result.ValidIngredientSlug.String,
				ShoppingSuggestions:                     result.ValidIngredientShoppingSuggestions.String,
				ContainsShellfish:                       result.ValidIngredientContainsShellfish.Bool,
				IsMeasuredVolumetrically:                result.ValidIngredientVolumetric.Bool,
				IsLiquid:                                result.ValidIngredientIsLiquid.Bool,
				ContainsPeanut:                          result.ValidIngredientContainsPeanut.Bool,
				ContainsTreeNut:                         result.ValidIngredientContainsTreeNut.Bool,
				ContainsEgg:                             result.ValidIngredientContainsEgg.Bool,
				ContainsWheat:                           result.ValidIngredientContainsWheat.Bool,
				ContainsSoy:                             result.ValidIngredientContainsSoy.Bool,
				AnimalDerived:                           result.ValidIngredientAnimalDerived.Bool,
				RestrictToPreparations:                  result.ValidIngredientRestrictToPreparations.Bool,
				ContainsSesame:                          result.ValidIngredientContainsSesame.Bool,
				ContainsFish:                            result.ValidIngredientContainsFish.Bool,
				ContainsGluten:                          result.ValidIngredientContainsGluten.Bool,
				ContainsDairy:                           result.ValidIngredientContainsDairy.Bool,
				ContainsAlcohol:                         result.ValidIngredientContainsAlcohol.Bool,
				AnimalFlesh:                             result.ValidIngredientAnimalFlesh.Bool,
				IsStarch:                                result.ValidIngredientIsStarch.Bool,
				IsProtein:                               result.ValidIngredientIsProtein.Bool,
				IsGrain:                                 result.ValidIngredientIsGrain.Bool,
				IsFruit:                                 result.ValidIngredientIsFruit.Bool,
				IsSalt:                                  result.ValidIngredientIsSalt.Bool,
				IsFat:                                   result.ValidIngredientIsFat.Bool,
				IsAcid:                                  result.ValidIngredientIsAcid.Bool,
				IsHeat:                                  result.ValidIngredientIsHeat.Bool,
			},
			Notes: result.Notes,
			ID:    result.ID,
			From: types.ValidMeasurementUnit{
				CreatedAt:     result.FromUnitCreatedAt,
				LastUpdatedAt: timePointerFromNullTime(result.FromUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.FromUnitArchivedAt),
				Name:          result.FromUnitName,
				IconPath:      result.FromUnitIconPath,
				ID:            result.FromUnitID,
				Description:   result.FromUnitDescription,
				PluralName:    result.FromUnitPluralName,
				Slug:          result.FromUnitSlug,
				Volumetric:    boolFromNullBool(result.FromUnitVolumetric),
				Universal:     result.FromUnitUniversal,
				Metric:        result.FromUnitMetric,
				Imperial:      result.FromUnitImperial,
			},
			To: types.ValidMeasurementUnit{
				CreatedAt:     result.ToUnitCreatedAt,
				LastUpdatedAt: timePointerFromNullTime(result.ToUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ToUnitArchivedAt),
				Name:          result.ToUnitName,
				IconPath:      result.ToUnitIconPath,
				ID:            result.ToUnitID,
				Description:   result.ToUnitDescription,
				PluralName:    result.ToUnitPluralName,
				Slug:          result.ToUnitSlug,
				Volumetric:    boolFromNullBool(result.ToUnitVolumetric),
				Universal:     result.ToUnitUniversal,
				Metric:        result.ToUnitMetric,
				Imperial:      result.ToUnitImperial,
			},
			Modifier: float32(result.ValidMeasurementConversionsModifier),
		}
	}

	return validMeasurementUnitConversions, nil
}

// CreateValidMeasurementUnitConversion creates a valid measurement conversion in the database.
func (q *Querier) CreateValidMeasurementUnitConversion(ctx context.Context, input *types.ValidMeasurementUnitConversionDatabaseCreationInput) (*types.ValidMeasurementUnitConversion, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, input.ID)

	// create the valid measurement conversion.
	if err := q.generatedQuerier.CreateValidMeasurementUnitConversion(ctx, q.db, &generated.CreateValidMeasurementUnitConversionParams{
		ID:                input.ID,
		FromUnit:          input.From,
		ToUnit:            input.To,
		Modifier:          float64(input.Modifier),
		Notes:             input.Notes,
		OnlyForIngredient: nullStringFromStringPointer(input.OnlyForIngredient),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid measurement conversion creation query")
	}

	x := &types.ValidMeasurementUnitConversion{
		ID:        input.ID,
		From:      types.ValidMeasurementUnit{ID: input.From},
		To:        types.ValidMeasurementUnit{ID: input.To},
		Modifier:  input.Modifier,
		Notes:     input.Notes,
		CreatedAt: q.currentTime(),
	}

	if input.OnlyForIngredient != nil {
		x.OnlyForIngredient = &types.ValidIngredient{ID: *input.OnlyForIngredient}
	}

	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, x.ID)
	logger.Info("valid measurement conversion created")

	return x, nil
}

// UpdateValidMeasurementUnitConversion updates a particular valid measurement conversion.
func (q *Querier) UpdateValidMeasurementUnitConversion(ctx context.Context, updated *types.ValidMeasurementUnitConversion) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}

	logger := q.logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, updated.ID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, updated.ID)

	var ingredientID *string
	if updated.OnlyForIngredient != nil {
		ingredientID = &updated.OnlyForIngredient.ID
	}

	if err := q.generatedQuerier.UpdateValidMeasurementUnitConversion(ctx, q.db, &generated.UpdateValidMeasurementUnitConversionParams{
		FromUnit:          updated.From.ID,
		ToUnit:            updated.To.ID,
		OnlyForIngredient: nullStringFromStringPointer(ingredientID),
		Modifier:          float64(updated.Modifier),
		Notes:             updated.Notes,
		ID:                updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement conversion")
	}

	logger.Info("valid measurement conversion updated")

	return nil
}

// ArchiveValidMeasurementUnitConversion archives a valid measurement conversion from the database by its ID.
func (q *Querier) ArchiveValidMeasurementUnitConversion(ctx context.Context, validMeasurementUnitConversionID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitConversionID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitConversionIDKey, validMeasurementUnitConversionID)
	tracing.AttachValidMeasurementUnitConversionIDToSpan(span, validMeasurementUnitConversionID)

	if err := q.generatedQuerier.ArchiveValidMeasurementUnitConversion(ctx, q.db, validMeasurementUnitConversionID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement conversion")
	}

	logger.Info("valid measurement conversion archived")

	return nil
}
