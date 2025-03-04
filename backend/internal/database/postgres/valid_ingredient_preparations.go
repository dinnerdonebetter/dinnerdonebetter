package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.ValidIngredientPreparationDataManager = (*Querier)(nil)
)

// ValidIngredientPreparationExists fetches whether a valid ingredient preparation exists from the database.
func (q *Querier) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	result, err := q.generatedQuerier.CheckValidIngredientPreparationExistence(ctx, q.db, validIngredientPreparationID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient preparation existence check")
	}

	return result, nil
}

// GetValidIngredientPreparation fetches a valid ingredient preparation from the database.
func (q *Querier) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	result, err := q.generatedQuerier.GetValidIngredientPreparation(ctx, q.db, validIngredientPreparationID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient preparation retrieval")
	}

	validIngredientPreparation := &types.ValidIngredientPreparation{
		CreatedAt:     result.ValidIngredientPreparationCreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientPreparationLastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientPreparationArchivedAt),
		Notes:         result.ValidIngredientPreparationNotes,
		ID:            result.ValidIngredientPreparationID,
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
		Preparation: types.ValidPreparation{
			CreatedAt: result.ValidPreparationCreatedAt,
			InstrumentCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				Min: uint16(result.ValidPreparationMinimumInstrumentCount),
			},
			IngredientCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				Min: uint16(result.ValidPreparationMinimumInstrumentCount),
			},
			VesselCount: types.Uint16RangeWithOptionalMax{
				Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				Min: uint16(result.ValidPreparationMinimumVesselCount),
			},
			ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
			LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
			IconPath:                    result.ValidPreparationIconPath,
			PastTense:                   result.ValidPreparationPastTense,
			ID:                          result.ValidPreparationID,
			Name:                        result.ValidPreparationName,
			Description:                 result.ValidPreparationDescription,
			Slug:                        result.ValidPreparationSlug,
			RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
			TemperatureRequired:         result.ValidPreparationTemperatureRequired,
			TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
			ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
			ConsumesVessel:              result.ValidPreparationConsumesVessel,
			OnlyForVessels:              result.ValidPreparationOnlyForVessels,
			YieldsNothing:               result.ValidPreparationYieldsNothing,
		},
	}

	return validIngredientPreparation, nil
}

// GetValidIngredientPreparations fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientPreparations(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidIngredientPreparation], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidIngredientPreparation]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientPreparations(ctx, q.db, &generated.GetValidIngredientPreparationsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientPreparation{
			CreatedAt:     result.ValidIngredientPreparationCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientPreparationLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientPreparationArchivedAt),
			Notes:         result.ValidIngredientPreparationNotes,
			ID:            result.ValidIngredientPreparationID,
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
			Preparation: types.ValidPreparation{
				CreatedAt: result.ValidPreparationCreatedAt,
				InstrumentCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				IngredientCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				VesselCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
					Min: uint16(result.ValidPreparationMinimumVesselCount),
				},
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientPreparationsForPreparation fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientPreparationsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidIngredientPreparation], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, preparationID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidIngredientPreparation]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientPreparationsForPreparation(ctx, q.db, &generated.GetValidIngredientPreparationsForPreparationParams{
		ID:              preparationID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientPreparation{
			CreatedAt:     result.ValidIngredientPreparationCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientPreparationLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientPreparationArchivedAt),
			Notes:         result.ValidIngredientPreparationNotes,
			ID:            result.ValidIngredientPreparationID,
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
			Preparation: types.ValidPreparation{
				CreatedAt: result.ValidPreparationCreatedAt,
				InstrumentCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				IngredientCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				VesselCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
					Min: uint16(result.ValidPreparationMinimumVesselCount),
				},
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidIngredientPreparationsForIngredient fetches a list of valid ingredient preparations from the database that meet a particular filter.
func (q *Querier) GetValidIngredientPreparationsForIngredient(ctx context.Context, ingredientID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidIngredientPreparation], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if ingredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, ingredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, ingredientID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidIngredientPreparation]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidIngredientPreparationsForIngredient(ctx, q.db, &generated.GetValidIngredientPreparationsForIngredientParams{
		ID:              ingredientID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid ingredient preparations list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidIngredientPreparation{
			CreatedAt:     result.ValidIngredientPreparationCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidIngredientPreparationLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidIngredientPreparationArchivedAt),
			Notes:         result.ValidIngredientPreparationNotes,
			ID:            result.ValidIngredientPreparationID,
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
			Preparation: types.ValidPreparation{
				CreatedAt: result.ValidPreparationCreatedAt,
				InstrumentCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				IngredientCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
					Min: uint16(result.ValidPreparationMinimumInstrumentCount),
				},
				VesselCount: types.Uint16RangeWithOptionalMax{
					Max: database.Uint16PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
					Min: uint16(result.ValidPreparationMinimumVesselCount),
				},
				ArchivedAt:                  database.TimePointerFromNullTime(result.ValidPreparationArchivedAt),
				LastUpdatedAt:               database.TimePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				RestrictToIngredients:       result.ValidPreparationRestrictToIngredients,
				TemperatureRequired:         result.ValidPreparationTemperatureRequired,
				TimeEstimateRequired:        result.ValidPreparationTimeEstimateRequired,
				ConditionExpressionRequired: result.ValidPreparationConditionExpressionRequired,
				ConsumesVessel:              result.ValidPreparationConsumesVessel,
				OnlyForVessels:              result.ValidPreparationOnlyForVessels,
				YieldsNothing:               result.ValidPreparationYieldsNothing,
			},
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateValidIngredientPreparation creates a valid ingredient preparation in the database.
func (q *Querier) CreateValidIngredientPreparation(ctx context.Context, input *types.ValidIngredientPreparationDatabaseCreationInput) (*types.ValidIngredientPreparation, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, input.ID)
	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, input.ID)

	// create the valid ingredient preparation.
	if err := q.generatedQuerier.CreateValidIngredientPreparation(ctx, q.db, &generated.CreateValidIngredientPreparationParams{
		ID:                 input.ID,
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidIngredientID:  input.ValidIngredientID,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid ingredient preparation creation query")
	}

	x := &types.ValidIngredientPreparation{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Ingredient:  types.ValidIngredient{ID: input.ValidIngredientID},
		CreatedAt:   q.currentTime(),
	}

	preparation, err := q.GetValidPreparation(ctx, input.ValidPreparationID)
	if err != nil {
		// basically impossible for this to happen and not error out earlier
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid preparation for valid ingredient preparation")
	}
	if preparation != nil {
		x.Preparation = *preparation
	}

	ingredient, err := q.GetValidIngredient(ctx, input.ValidIngredientID)
	if err != nil {
		// basically impossible for this to happen and not error out earlier
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid ingredient for valid ingredient preparation")
	}
	if ingredient != nil {
		x.Ingredient = *ingredient
	}

	return x, nil
}

// UpdateValidIngredientPreparation updates a particular valid ingredient preparation.
func (q *Querier) UpdateValidIngredientPreparation(ctx context.Context, updated *types.ValidIngredientPreparation) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidIngredientPreparationIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidIngredientPreparation(ctx, q.db, &generated.UpdateValidIngredientPreparationParams{
		Notes:              updated.Notes,
		ValidPreparationID: updated.Preparation.ID,
		ValidIngredientID:  updated.Ingredient.ID,
		ID:                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation updated")

	return nil
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation from the database by its ID.
func (q *Querier) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validIngredientPreparationID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)
	tracing.AttachToSpan(span, keys.ValidIngredientPreparationIDKey, validIngredientPreparationID)

	if _, err := q.generatedQuerier.ArchiveValidIngredientPreparation(ctx, q.db, validIngredientPreparationID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid ingredient preparation")
	}

	logger.Info("valid ingredient preparation archived")

	return nil
}
