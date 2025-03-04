package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

var (
	_ types.ValidIngredientMeasurementUnitDataManager = (*Querier)(nil)
)

// ValidIngredientMeasurementUnitExists fetches whether a valid ingredient measurement unit exists from the database.
func (q *Querier) ValidIngredientMeasurementUnitExists(ctx context.Context, validIngredientMeasurementUnitID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	result, err := q.generatedQuerier.CheckValidIngredientMeasurementUnitExistence(ctx, q.db, validIngredientMeasurementUnitID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient measurement unit existence check")
	}

	return result, nil
}

// GetValidIngredientMeasurementUnit fetches a valid ingredient measurement unit from the database.
func (q *Querier) GetValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	result, err := q.generatedQuerier.GetValidIngredientMeasurementUnit(ctx, q.db, validIngredientMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning validIngredientMeasurementUnit")
	}

	validIngredientMeasurementUnit := &types.ValidIngredientMeasurementUnit{
		CreatedAt:     result.ValidIngredientMeasurementUnitCreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitLastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitArchivedAt),
		AllowableQuantity: types.Float32RangeWithOptionalMax{
			Max: database.Float32PointerFromNullString(result.ValidIngredientMeasurementUnitMaximumAllowableQuantity),
			Min: database.Float32FromString(result.ValidIngredientMeasurementUnitMinimumAllowableQuantity),
		},
		Notes: result.ValidIngredientMeasurementUnitNotes,
		ID:    result.ValidIngredientMeasurementUnitID,
		MeasurementUnit: types.ValidMeasurementUnit{
			CreatedAt:     result.ValidMeasurementUnitCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
			Name:          result.ValidMeasurementUnitName,
			IconPath:      result.ValidMeasurementUnitIconPath,
			ID:            result.ValidMeasurementUnitID,
			Description:   result.ValidMeasurementUnitDescription,
			PluralName:    result.ValidMeasurementUnitPluralName,
			Slug:          result.ValidMeasurementUnitSlug,
			Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
			Universal:     result.ValidMeasurementUnitUniversal,
			Metric:        result.ValidMeasurementUnitMetric,
			Imperial:      result.ValidMeasurementUnitImperial,
		},
		Ingredient: types.ValidIngredient{
			CreatedAt:     result.ValidIngredientCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
			StorageTemperatureInCelsius: types.OptionalFloat32Range{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
				Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
			},
			IconPath:               result.ValidIngredientIconPath,
			Warning:                result.ValidIngredientWarning,
			PluralName:             result.ValidIngredientPluralName,
			StorageInstructions:    result.ValidIngredientStorageInstructions,
			Name:                   result.ValidIngredientName,
			ID:                     result.ValidIngredientID,
			Description:            result.ValidIngredientDescription,
			Slug:                   result.ValidIngredientSlug,
			ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
			ContainsShellfish:      result.ValidIngredientContainsShellfish,
			IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
			ContainsPeanut:         result.ValidIngredientContainsPeanut,
			ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
			ContainsEgg:            result.ValidIngredientContainsEgg,
			ContainsWheat:          result.ValidIngredientContainsWheat,
			ContainsSoy:            result.ValidIngredientContainsSoy,
			AnimalDerived:          result.ValidIngredientAnimalDerived,
			RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
			ContainsSesame:         result.ValidIngredientContainsSesame,
			ContainsFish:           result.ValidIngredientContainsFish,
			ContainsGluten:         result.ValidIngredientContainsGluten,
			ContainsDairy:          result.ValidIngredientContainsDairy,
			ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
			AnimalFlesh:            result.ValidIngredientAnimalFlesh,
			IsStarch:               result.ValidIngredientIsStarch,
			IsProtein:              result.ValidIngredientIsProtein,
			IsGrain:                result.ValidIngredientIsGrain,
			IsFruit:                result.ValidIngredientIsFruit,
			IsSalt:                 result.ValidIngredientIsSalt,
			IsFat:                  result.ValidIngredientIsFat,
			IsAcid:                 result.ValidIngredientIsAcid,
			IsHeat:                 result.ValidIngredientIsHeat,
		},
	}

	return validIngredientMeasurementUnit, nil
}

// GetValidIngredientMeasurementUnitsForIngredient fetches a list of valid measurement units from the database that belong to a given ingredient ID.
func (q *Querier) GetValidIngredientMeasurementUnitsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, ingredientID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Pagination: filter.ToPagination(),
		Data:       []*types.ValidIngredientMeasurementUnit{},
	}

	results, err := q.generatedQuerier.GetValidIngredientMeasurementUnitsForIngredient(ctx, q.db, &generated.GetValidIngredientMeasurementUnitsForIngredientParams{
		ValidIngredientID: ingredientID,
		CreatedBefore:     database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:      database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:     database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:      database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:       database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:        database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived:   database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient measurement units list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientMeasurementUnit{
			CreatedAt:     result.ValidIngredientMeasurementUnitCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitArchivedAt),
			Notes:         result.ValidIngredientMeasurementUnitNotes,
			ID:            result.ValidIngredientMeasurementUnitID,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			Ingredient: types.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath,
				Warning:                result.ValidIngredientWarning,
				PluralName:             result.ValidIngredientPluralName,
				StorageInstructions:    result.ValidIngredientStorageInstructions,
				Name:                   result.ValidIngredientName,
				ID:                     result.ValidIngredientID,
				Description:            result.ValidIngredientDescription,
				Slug:                   result.ValidIngredientSlug,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:      result.ValidIngredientContainsShellfish,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
				ContainsEgg:            result.ValidIngredientContainsEgg,
				ContainsWheat:          result.ValidIngredientContainsWheat,
				ContainsSoy:            result.ValidIngredientContainsSoy,
				AnimalDerived:          result.ValidIngredientAnimalDerived,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
				ContainsSesame:         result.ValidIngredientContainsSesame,
				ContainsFish:           result.ValidIngredientContainsFish,
				ContainsGluten:         result.ValidIngredientContainsGluten,
				ContainsDairy:          result.ValidIngredientContainsDairy,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh,
				IsStarch:               result.ValidIngredientIsStarch,
				IsProtein:              result.ValidIngredientIsProtein,
				IsGrain:                result.ValidIngredientIsGrain,
				IsFruit:                result.ValidIngredientIsFruit,
				IsSalt:                 result.ValidIngredientIsSalt,
				IsFat:                  result.ValidIngredientIsFat,
				IsAcid:                 result.ValidIngredientIsAcid,
				IsHeat:                 result.ValidIngredientIsHeat,
			},
			AllowableQuantity: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMeasurementUnitMaximumAllowableQuantity),
				Min: database.Float32FromString(result.ValidIngredientMeasurementUnitMinimumAllowableQuantity),
			},
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientMeasurementUnitsForMeasurementUnit fetches a list of valid measurement units from the database that belong to a given ingredient ID.
func (q *Querier) GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx context.Context, validMeasurementUnitID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidPreparationInstrumentIDKey, validMeasurementUnitID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Pagination: filter.ToPagination(),
		Data:       []*types.ValidIngredientMeasurementUnit{},
	}

	results, err := q.generatedQuerier.GetValidIngredientMeasurementUnitsForMeasurementUnit(ctx, q.db, &generated.GetValidIngredientMeasurementUnitsForMeasurementUnitParams{
		ValidMeasurementUnitID: validMeasurementUnitID,
		CreatedBefore:          database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:           database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:          database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:           database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:            database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:             database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived:        database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient measurement units list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientMeasurementUnit{
			CreatedAt:     result.ValidIngredientMeasurementUnitCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitArchivedAt),
			Notes:         result.ValidIngredientMeasurementUnitNotes,
			ID:            result.ValidIngredientMeasurementUnitID,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			Ingredient: types.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath,
				Warning:                result.ValidIngredientWarning,
				PluralName:             result.ValidIngredientPluralName,
				StorageInstructions:    result.ValidIngredientStorageInstructions,
				Name:                   result.ValidIngredientName,
				ID:                     result.ValidIngredientID,
				Description:            result.ValidIngredientDescription,
				Slug:                   result.ValidIngredientSlug,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:      result.ValidIngredientContainsShellfish,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
				ContainsEgg:            result.ValidIngredientContainsEgg,
				ContainsWheat:          result.ValidIngredientContainsWheat,
				ContainsSoy:            result.ValidIngredientContainsSoy,
				AnimalDerived:          result.ValidIngredientAnimalDerived,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
				ContainsSesame:         result.ValidIngredientContainsSesame,
				ContainsFish:           result.ValidIngredientContainsFish,
				ContainsGluten:         result.ValidIngredientContainsGluten,
				ContainsDairy:          result.ValidIngredientContainsDairy,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh,
				IsStarch:               result.ValidIngredientIsStarch,
				IsProtein:              result.ValidIngredientIsProtein,
				IsGrain:                result.ValidIngredientIsGrain,
				IsFruit:                result.ValidIngredientIsFruit,
				IsSalt:                 result.ValidIngredientIsSalt,
				IsFat:                  result.ValidIngredientIsFat,
				IsAcid:                 result.ValidIngredientIsAcid,
				IsHeat:                 result.ValidIngredientIsHeat,
			},
			AllowableQuantity: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMeasurementUnitMaximumAllowableQuantity),
				Min: database.Float32FromString(result.ValidIngredientMeasurementUnitMinimumAllowableQuantity),
			},
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientMeasurementUnits fetches a list of valid ingredient measurement units from the database that meet a particular filter.
func (q *Querier) GetValidIngredientMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit], error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x := &filtering.QueryFilteredResult[types.ValidIngredientMeasurementUnit]{
		Pagination: filter.ToPagination(),
		Data:       []*types.ValidIngredientMeasurementUnit{},
	}

	results, err := q.generatedQuerier.GetValidIngredientMeasurementUnits(ctx, q.db, &generated.GetValidIngredientMeasurementUnitsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient measurement units list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientMeasurementUnit{
			CreatedAt:     result.ValidIngredientMeasurementUnitCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientMeasurementUnitArchivedAt),
			AllowableQuantity: types.Float32RangeWithOptionalMax{
				Max: database.Float32PointerFromNullString(result.ValidIngredientMeasurementUnitMaximumAllowableQuantity),
				Min: database.Float32FromString(result.ValidIngredientMeasurementUnitMinimumAllowableQuantity),
			},
			Notes: result.ValidIngredientMeasurementUnitNotes,
			ID:    result.ValidIngredientMeasurementUnitID,
			MeasurementUnit: types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName,
				IconPath:      result.ValidMeasurementUnitIconPath,
				ID:            result.ValidMeasurementUnitID,
				Description:   result.ValidMeasurementUnitDescription,
				PluralName:    result.ValidMeasurementUnitPluralName,
				Slug:          result.ValidMeasurementUnitSlug,
				Volumetric:    database.BoolFromNullBool(result.ValidMeasurementUnitVolumetric),
				Universal:     result.ValidMeasurementUnitUniversal,
				Metric:        result.ValidMeasurementUnitMetric,
				Imperial:      result.ValidMeasurementUnitImperial,
			},
			Ingredient: types.ValidIngredient{
				CreatedAt:     result.ValidIngredientCreatedAt,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientArchivedAt),
				StorageTemperatureInCelsius: types.OptionalFloat32Range{
					Max: database.Float32PointerFromNullString(result.ValidIngredientMaximumIdealStorageTemperatureInCelsius),
					Min: database.Float32PointerFromNullString(result.ValidIngredientMinimumIdealStorageTemperatureInCelsius),
				},
				IconPath:               result.ValidIngredientIconPath,
				Warning:                result.ValidIngredientWarning,
				PluralName:             result.ValidIngredientPluralName,
				StorageInstructions:    result.ValidIngredientStorageInstructions,
				Name:                   result.ValidIngredientName,
				ID:                     result.ValidIngredientID,
				Description:            result.ValidIngredientDescription,
				Slug:                   result.ValidIngredientSlug,
				ShoppingSuggestions:    result.ValidIngredientShoppingSuggestions,
				ContainsShellfish:      result.ValidIngredientContainsShellfish,
				IsLiquid:               database.BoolFromNullBool(result.ValidIngredientIsLiquid),
				ContainsPeanut:         result.ValidIngredientContainsPeanut,
				ContainsTreeNut:        result.ValidIngredientContainsTreeNut,
				ContainsEgg:            result.ValidIngredientContainsEgg,
				ContainsWheat:          result.ValidIngredientContainsWheat,
				ContainsSoy:            result.ValidIngredientContainsSoy,
				AnimalDerived:          result.ValidIngredientAnimalDerived,
				RestrictToPreparations: result.ValidIngredientRestrictToPreparations,
				ContainsSesame:         result.ValidIngredientContainsSesame,
				ContainsFish:           result.ValidIngredientContainsFish,
				ContainsGluten:         result.ValidIngredientContainsGluten,
				ContainsDairy:          result.ValidIngredientContainsDairy,
				ContainsAlcohol:        result.ValidIngredientContainsAlcohol,
				AnimalFlesh:            result.ValidIngredientAnimalFlesh,
				IsStarch:               result.ValidIngredientIsStarch,
				IsProtein:              result.ValidIngredientIsProtein,
				IsGrain:                result.ValidIngredientIsGrain,
				IsFruit:                result.ValidIngredientIsFruit,
				IsSalt:                 result.ValidIngredientIsSalt,
				IsFat:                  result.ValidIngredientIsFat,
				IsAcid:                 result.ValidIngredientIsAcid,
				IsHeat:                 result.ValidIngredientIsHeat,
			},
		})

		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateValidIngredientMeasurementUnit creates a valid ingredient measurement unit in the database.
func (q *Querier) CreateValidIngredientMeasurementUnit(ctx context.Context, input *types.ValidIngredientMeasurementUnitDatabaseCreationInput) (*types.ValidIngredientMeasurementUnit, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, input.ID)

	// create the valid ingredient measurement unit.
	if err := q.generatedQuerier.CreateValidIngredientMeasurementUnit(ctx, q.db, &generated.CreateValidIngredientMeasurementUnitParams{
		ID:                       input.ID,
		Notes:                    input.Notes,
		ValidMeasurementUnitID:   input.ValidMeasurementUnitID,
		ValidIngredientID:        input.ValidIngredientID,
		MinimumAllowableQuantity: database.StringFromFloat32(input.AllowableQuantity.Min),
		MaximumAllowableQuantity: database.NullStringFromFloat32Pointer(input.AllowableQuantity.Max),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient measurement unit creation query")
	}

	x := &types.ValidIngredientMeasurementUnit{
		ID:              input.ID,
		Notes:           input.Notes,
		MeasurementUnit: types.ValidMeasurementUnit{ID: input.ValidMeasurementUnitID},
		Ingredient:      types.ValidIngredient{ID: input.ValidIngredientID},
		AllowableQuantity: types.Float32RangeWithOptionalMax{
			Max: input.AllowableQuantity.Max,
			Min: input.AllowableQuantity.Min,
		},
		CreatedAt: q.currentTime(),
	}

	ingredient, err := q.GetValidIngredient(ctx, input.ValidIngredientID)
	if err != nil {
		// basically impossible for this to happen and not error out earlier
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient for valid ingredient measurement unit")
	}
	if ingredient != nil {
		x.Ingredient = *ingredient
	}

	measurementUnit, err := q.GetValidMeasurementUnit(ctx, input.ValidMeasurementUnitID)
	if err != nil {
		// basically impossible for this to happen and not error out earlier
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid measurement unit for valid ingredient measurement unit")
	}
	if measurementUnit != nil {
		x.MeasurementUnit = *measurementUnit
	}

	logger.Info("valid ingredient measurement unit created")

	return x, nil
}

// UpdateValidIngredientMeasurementUnit updates a particular valid ingredient measurement unit.
func (q *Querier) UpdateValidIngredientMeasurementUnit(ctx context.Context, updated *types.ValidIngredientMeasurementUnit) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidIngredientMeasurementUnit(ctx, q.db, &generated.UpdateValidIngredientMeasurementUnitParams{
		Notes:                    updated.Notes,
		ValidMeasurementUnitID:   updated.MeasurementUnit.ID,
		ValidIngredientID:        updated.Ingredient.ID,
		MinimumAllowableQuantity: database.StringFromFloat32(updated.AllowableQuantity.Min),
		MaximumAllowableQuantity: database.NullStringFromFloat32Pointer(updated.AllowableQuantity.Max),
		ID:                       updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient measurement unit")
	}

	logger.Info("valid ingredient measurement unit updated")

	return nil
}

// ArchiveValidIngredientMeasurementUnit archives a valid ingredient measurement unit from the database by its ID.
func (q *Querier) ArchiveValidIngredientMeasurementUnit(ctx context.Context, validIngredientMeasurementUnitID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientMeasurementUnitID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidIngredientMeasurementUnitIDKey, validIngredientMeasurementUnitID)

	if _, err := q.generatedQuerier.ArchiveValidIngredientMeasurementUnit(ctx, q.db, validIngredientMeasurementUnitID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient measurement unit")
	}

	logger.Info("valid ingredient measurement unit archived")

	return nil
}
