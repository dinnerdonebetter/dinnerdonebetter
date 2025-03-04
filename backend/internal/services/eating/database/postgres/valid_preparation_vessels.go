package postgres

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/lib/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database"
	"github.com/dinnerdonebetter/backend/internal/services/eating/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/services/eating/types"
)

var (
	_ types.ValidPreparationVesselDataManager = (*Querier)(nil)
)

// ValidPreparationVesselExists fetches whether a valid preparation vessel exists from the database.
func (q *Querier) ValidPreparationVesselExists(ctx context.Context, validPreparationVesselID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return false, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	exists, err = q.generatedQuerier.CheckValidPreparationVesselExistence(ctx, q.db, validPreparationVesselID)
	if err != nil {
		return false, observability.PrepareError(err, span, "performing valid preparation vessel existence check")
	}

	return exists, nil
}

// GetValidPreparationVessel fetches a valid preparation vessel from the database.
func (q *Querier) GetValidPreparationVessel(ctx context.Context, validPreparationVesselID string) (*types.ValidPreparationVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	result, err := q.generatedQuerier.GetValidPreparationVessel(ctx, q.db, validPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "reading valid preparation vessel from database")
	}

	validPreparationVessel := &types.ValidPreparationVessel{
		CreatedAt:     result.ValidPreparationVesselCreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
		ID:            result.ValidPreparationVesselID,
		Notes:         result.ValidPreparationVesselNotes,
		Vessel: types.ValidVessel{
			CreatedAt:     result.ValidVesselCreatedAt,
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
			CapacityUnit: &types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
				Name:          result.ValidMeasurementUnitName.String,
				IconPath:      result.ValidMeasurementUnitIconPath.String,
				ID:            result.ValidMeasurementUnitID.String,
				Description:   result.ValidMeasurementUnitDescription.String,
				PluralName:    result.ValidMeasurementUnitPluralName.String,
				Slug:          result.ValidMeasurementUnitSlug.String,
				Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
				Universal:     result.ValidMeasurementUnitUniversal.Bool,
				Metric:        result.ValidMeasurementUnitMetric.Bool,
				Imperial:      result.ValidMeasurementUnitImperial.Bool,
			},
			IconPath:                       result.ValidVesselIconPath,
			PluralName:                     result.ValidVesselPluralName,
			Description:                    result.ValidVesselDescription,
			Name:                           result.ValidVesselName,
			Slug:                           result.ValidVesselSlug,
			Shape:                          string(result.ValidVesselShape),
			ID:                             result.ValidVesselID,
			WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
			LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
			HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
			Capacity:                       database.Float32FromString(result.ValidVesselCapacity),
			IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
			DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
			UsableForStorage:               result.ValidVesselUsableForStorage,
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

	return validPreparationVessel, nil
}

// GetValidPreparationVessels fetches a list of valid preparation vessels from the database that meet a particular filter.
func (q *Querier) GetValidPreparationVessels(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetValidPreparationVessels(ctx, q.db, &generated.GetValidPreparationVesselsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation vessels list retrieval query")
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, result := range results {
		validPreparationVessel := &types.ValidPreparationVessel{
			CreatedAt:     result.ValidPreparationVesselCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
			ID:            result.ValidPreparationVesselID,
			Notes:         result.ValidPreparationVesselNotes,
			Vessel: types.ValidVessel{
				CreatedAt:     result.ValidVesselCreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit: &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
					Name:          result.ValidMeasurementUnitName.String,
					IconPath:      result.ValidMeasurementUnitIconPath.String,
					ID:            result.ValidMeasurementUnitID.String,
					Description:   result.ValidMeasurementUnitDescription.String,
					PluralName:    result.ValidMeasurementUnitPluralName.String,
					Slug:          result.ValidMeasurementUnitSlug.String,
					Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
					Universal:     result.ValidMeasurementUnitUniversal.Bool,
					Metric:        result.ValidMeasurementUnitMetric.Bool,
					Imperial:      result.ValidMeasurementUnitImperial.Bool,
				},
				IconPath:                       result.ValidVesselIconPath,
				PluralName:                     result.ValidVesselPluralName,
				Description:                    result.ValidVesselDescription,
				Name:                           result.ValidVesselName,
				Slug:                           result.ValidVesselSlug,
				Shape:                          string(result.ValidVesselShape),
				ID:                             result.ValidVesselID,
				WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
				Capacity:                       database.Float32FromString(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
				UsableForStorage:               result.ValidVesselUsableForStorage,
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
		if validPreparationVessel.Vessel.CapacityUnit.ID == "" {
			validPreparationVessel.Vessel.CapacityUnit = nil
		}

		x.Data = append(x.Data, validPreparationVessel)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidPreparationVesselsForPreparation fetches a list of valid preparation vessels from the database that meet a particular filter.
func (q *Querier) GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, preparationID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetValidPreparationVesselsForPreparation(ctx, q.db, &generated.GetValidPreparationVesselsForPreparationParams{
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
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation vessels list retrieval query")
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, result := range results {
		validPreparationVessel := &types.ValidPreparationVessel{
			CreatedAt:     result.ValidPreparationVesselCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
			ID:            result.ValidPreparationVesselID,
			Notes:         result.ValidPreparationVesselNotes,
			Vessel: types.ValidVessel{
				CreatedAt:     result.ValidVesselCreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit: &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
					Name:          result.ValidMeasurementUnitName.String,
					IconPath:      result.ValidMeasurementUnitIconPath.String,
					ID:            result.ValidMeasurementUnitID.String,
					Description:   result.ValidMeasurementUnitDescription.String,
					PluralName:    result.ValidMeasurementUnitPluralName.String,
					Slug:          result.ValidMeasurementUnitSlug.String,
					Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
					Universal:     result.ValidMeasurementUnitUniversal.Bool,
					Metric:        result.ValidMeasurementUnitMetric.Bool,
					Imperial:      result.ValidMeasurementUnitImperial.Bool,
				},
				IconPath:                       result.ValidVesselIconPath,
				PluralName:                     result.ValidVesselPluralName,
				Description:                    result.ValidVesselDescription,
				Name:                           result.ValidVesselName,
				Slug:                           result.ValidVesselSlug,
				Shape:                          string(result.ValidVesselShape),
				ID:                             result.ValidVesselID,
				WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
				Capacity:                       database.Float32FromString(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
				UsableForStorage:               result.ValidVesselUsableForStorage,
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
		if validPreparationVessel.Vessel.CapacityUnit.ID == "" {
			validPreparationVessel.Vessel.CapacityUnit = nil
		}

		x.Data = append(x.Data, validPreparationVessel)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidPreparationVesselsForVessel fetches a list of valid preparation vessels from the database that meet a particular filter.
func (q *Querier) GetValidPreparationVesselsForVessel(ctx context.Context, vesselID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if vesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, vesselID)
	logger = logger.WithValue(keys.ValidVesselIDKey, vesselID)

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetValidPreparationVesselsForVessel(ctx, q.db, &generated.GetValidPreparationVesselsForVesselParams{
		ID:              vesselID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.Limit),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid preparation vessels list retrieval query")
	}

	if len(results) == 0 {
		return nil, sql.ErrNoRows
	}

	for _, result := range results {
		validPreparationVessel := &types.ValidPreparationVessel{
			CreatedAt:     result.ValidPreparationVesselCreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
			ID:            result.ValidPreparationVesselID,
			Notes:         result.ValidPreparationVesselNotes,
			Vessel: types.ValidVessel{
				CreatedAt:     result.ValidVesselCreatedAt,
				ArchivedAt:    database.TimePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt: database.TimePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit: &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: database.TimePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    database.TimePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
					Name:          result.ValidMeasurementUnitName.String,
					IconPath:      result.ValidMeasurementUnitIconPath.String,
					ID:            result.ValidMeasurementUnitID.String,
					Description:   result.ValidMeasurementUnitDescription.String,
					PluralName:    result.ValidMeasurementUnitPluralName.String,
					Slug:          result.ValidMeasurementUnitSlug.String,
					Volumetric:    result.ValidMeasurementUnitVolumetric.Bool,
					Universal:     result.ValidMeasurementUnitUniversal.Bool,
					Metric:        result.ValidMeasurementUnitMetric.Bool,
					Imperial:      result.ValidMeasurementUnitImperial.Bool,
				},
				IconPath:                       result.ValidVesselIconPath,
				PluralName:                     result.ValidVesselPluralName,
				Description:                    result.ValidVesselDescription,
				Name:                           result.ValidVesselName,
				Slug:                           result.ValidVesselSlug,
				Shape:                          string(result.ValidVesselShape),
				ID:                             result.ValidVesselID,
				WidthInMillimeters:             database.Float32FromNullString(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            database.Float32FromNullString(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            database.Float32FromNullString(result.ValidVesselHeightInMillimeters),
				Capacity:                       database.Float32FromString(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
				UsableForStorage:               result.ValidVesselUsableForStorage,
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
		if validPreparationVessel.Vessel.CapacityUnit.ID == "" {
			validPreparationVessel.Vessel.CapacityUnit = nil
		}

		x.Data = append(x.Data, validPreparationVessel)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateValidPreparationVessel creates a valid preparation vessel in the database.
func (q *Querier) CreateValidPreparationVessel(ctx context.Context, input *types.ValidPreparationVesselDatabaseCreationInput) (*types.ValidPreparationVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationVesselIDKey, input.ID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, input.ID)

	// create the valid preparation vessel.
	if err := q.generatedQuerier.CreateValidPreparationVessel(ctx, q.db, &generated.CreateValidPreparationVesselParams{
		ID:                 input.ID,
		Notes:              input.Notes,
		ValidPreparationID: input.ValidPreparationID,
		ValidVesselID:      input.ValidVesselID,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid preparation vessel creation query")
	}

	x := &types.ValidPreparationVessel{
		ID:          input.ID,
		Notes:       input.Notes,
		Preparation: types.ValidPreparation{ID: input.ValidPreparationID},
		Vessel:      types.ValidVessel{ID: input.ValidVesselID},
		CreatedAt:   q.currentTime(),
	}

	preparation, err := q.GetValidPreparation(ctx, input.ValidPreparationID)
	if err != nil {
		// basically impossible for this to happen and not error out earlier
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid preparation for valid preparation vessel")
	}
	if preparation != nil {
		x.Preparation = *preparation
	}

	vessel, err := q.GetValidVessel(ctx, input.ValidVesselID)
	if err != nil {
		// basically impossible for this to happen and not error out earlier
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching valid vessel for valid preparation vessel")
	}
	if vessel != nil {
		x.Vessel = *vessel
	}

	logger.Info("valid preparation vessel created")

	return x, nil
}

// UpdateValidPreparationVessel updates a particular valid preparation vessel.
func (q *Querier) UpdateValidPreparationVessel(ctx context.Context, updated *types.ValidPreparationVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationVesselIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateValidPreparationVessel(ctx, q.db, &generated.UpdateValidPreparationVesselParams{
		Notes:              updated.Notes,
		ValidPreparationID: updated.Preparation.ID,
		ValidVesselID:      updated.Vessel.ID,
		ID:                 updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}

	logger.Info("valid preparation vessel updated")

	return nil
}

// ArchiveValidPreparationVessel archives a valid preparation vessel from the database by its ID.
func (q *Querier) ArchiveValidPreparationVessel(ctx context.Context, validPreparationVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if validPreparationVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger := q.logger.WithValue(keys.ValidPreparationVesselIDKey, validPreparationVesselID)
	tracing.AttachToSpan(span, keys.ValidVesselIDKey, validPreparationVesselID)

	if _, err := q.generatedQuerier.ArchiveValidPreparationVessel(ctx, q.db, validPreparationVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}

	logger.Info("valid preparation vessel archived")

	return nil
}
