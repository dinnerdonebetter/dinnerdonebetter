package postgres

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
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
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

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
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	result, err := q.generatedQuerier.GetValidPreparationVessel(ctx, q.db, validPreparationVesselID)
	if err != nil {
		return nil, observability.PrepareError(err, span, "reading valid preparation vessel from database")
	}

	validPreparationVessel := &types.ValidPreparationVessel{
		CreatedAt:     result.ValidPreparationVesselCreatedAt,
		LastUpdatedAt: timePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
		ArchivedAt:    timePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
		ID:            result.ValidPreparationVesselID,
		Notes:         result.ValidPreparationVesselNotes,
		Vessel: types.ValidVessel{
			CreatedAt:     result.ValidVesselCreatedAt,
			ArchivedAt:    timePointerFromNullTime(result.ValidVesselArchivedAt),
			LastUpdatedAt: timePointerFromNullTime(result.ValidVesselLastUpdatedAt),
			CapacityUnit: &types.ValidMeasurementUnit{
				CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
				LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
				ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
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
			WidthInMillimeters:             float32(result.ValidVesselWidthInMillimeters),
			LengthInMillimeters:            float32(result.ValidVesselLengthInMillimeters),
			HeightInMillimeters:            float32(result.ValidVesselHeightInMillimeters),
			Capacity:                       float32(result.ValidVesselCapacity),
			IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
			DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
			UsableForStorage:               result.ValidVesselUsableForStorage,
		},
		Preparation: types.ValidPreparation{
			CreatedAt:                   result.ValidPreparationCreatedAt,
			MaximumInstrumentCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
			ArchivedAt:                  timePointerFromNullTime(result.ValidPreparationArchivedAt),
			MaximumIngredientCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
			LastUpdatedAt:               timePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
			MaximumVesselCount:          int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
			IconPath:                    result.ValidPreparationIconPath,
			PastTense:                   result.ValidPreparationPastTense,
			ID:                          result.ValidPreparationID,
			Name:                        result.ValidPreparationName,
			Description:                 result.ValidPreparationDescription,
			Slug:                        result.ValidPreparationSlug,
			MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
			MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
			MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
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
func (q *Querier) GetValidPreparationVessels(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetValidPreparationVessels(ctx, q.db, &generated.GetValidPreparationVesselsParams{
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
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
			LastUpdatedAt: timePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
			ID:            result.ValidPreparationVesselID,
			Notes:         result.ValidPreparationVesselNotes,
			Vessel: types.ValidVessel{
				CreatedAt:     result.ValidVesselCreatedAt,
				ArchivedAt:    timePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt: timePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit: &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
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
				WidthInMillimeters:             float32(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            float32(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            float32(result.ValidVesselHeightInMillimeters),
				Capacity:                       float32(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
				UsableForStorage:               result.ValidVesselUsableForStorage,
			},
			Preparation: types.ValidPreparation{
				CreatedAt:                   result.ValidPreparationCreatedAt,
				MaximumInstrumentCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				ArchivedAt:                  timePointerFromNullTime(result.ValidPreparationArchivedAt),
				MaximumIngredientCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				LastUpdatedAt:               timePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				MaximumVesselCount:          int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
				MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
				MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
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
func (q *Querier) GetValidPreparationVesselsForPreparation(ctx context.Context, preparationID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if preparationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidPreparationIDKey, preparationID)
	tracing.AttachValidPreparationVesselIDToSpan(span, preparationID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetValidPreparationVesselsForPreparation(ctx, q.db, &generated.GetValidPreparationVesselsForPreparationParams{
		Ids:           []string{preparationID},
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
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
			LastUpdatedAt: timePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
			ID:            result.ValidPreparationVesselID,
			Notes:         result.ValidPreparationVesselNotes,
			Vessel: types.ValidVessel{
				CreatedAt:     result.ValidVesselCreatedAt,
				ArchivedAt:    timePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt: timePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit: &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
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
				WidthInMillimeters:             float32(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            float32(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            float32(result.ValidVesselHeightInMillimeters),
				Capacity:                       float32(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
				UsableForStorage:               result.ValidVesselUsableForStorage,
			},
			Preparation: types.ValidPreparation{
				CreatedAt:                   result.ValidPreparationCreatedAt,
				MaximumInstrumentCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				ArchivedAt:                  timePointerFromNullTime(result.ValidPreparationArchivedAt),
				MaximumIngredientCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				LastUpdatedAt:               timePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				MaximumVesselCount:          int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
				MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
				MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
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
func (q *Querier) GetValidPreparationVesselsForVessel(ctx context.Context, vesselID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidPreparationVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if vesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachValidVesselIDToSpan(span, vesselID)
	logger = logger.WithValue(keys.ValidVesselIDKey, vesselID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.ValidPreparationVessel]{
		Pagination: filter.ToPagination(),
	}
	tracing.AttachQueryFilterToSpan(span, filter)

	results, err := q.generatedQuerier.GetValidPreparationVesselsForVessel(ctx, q.db, &generated.GetValidPreparationVesselsForVesselParams{
		Ids:           []string{vesselID},
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
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
			LastUpdatedAt: timePointerFromNullTime(result.ValidPreparationVesselLastUpdatedAt),
			ArchivedAt:    timePointerFromNullTime(result.ValidPreparationVesselArchivedAt),
			ID:            result.ValidPreparationVesselID,
			Notes:         result.ValidPreparationVesselNotes,
			Vessel: types.ValidVessel{
				CreatedAt:     result.ValidVesselCreatedAt,
				ArchivedAt:    timePointerFromNullTime(result.ValidVesselArchivedAt),
				LastUpdatedAt: timePointerFromNullTime(result.ValidVesselLastUpdatedAt),
				CapacityUnit: &types.ValidMeasurementUnit{
					CreatedAt:     result.ValidMeasurementUnitCreatedAt.Time,
					LastUpdatedAt: timePointerFromNullTime(result.ValidMeasurementUnitLastUpdatedAt),
					ArchivedAt:    timePointerFromNullTime(result.ValidMeasurementUnitArchivedAt),
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
				WidthInMillimeters:             float32(result.ValidVesselWidthInMillimeters),
				LengthInMillimeters:            float32(result.ValidVesselLengthInMillimeters),
				HeightInMillimeters:            float32(result.ValidVesselHeightInMillimeters),
				Capacity:                       float32(result.ValidVesselCapacity),
				IncludeInGeneratedInstructions: result.ValidVesselIncludeInGeneratedInstructions,
				DisplayInSummaryLists:          result.ValidVesselDisplayInSummaryLists,
				UsableForStorage:               result.ValidVesselUsableForStorage,
			},
			Preparation: types.ValidPreparation{
				CreatedAt:                   result.ValidPreparationCreatedAt,
				MaximumInstrumentCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumInstrumentCount),
				ArchivedAt:                  timePointerFromNullTime(result.ValidPreparationArchivedAt),
				MaximumIngredientCount:      int32PointerFromNullInt32(result.ValidPreparationMaximumIngredientCount),
				LastUpdatedAt:               timePointerFromNullTime(result.ValidPreparationLastUpdatedAt),
				MaximumVesselCount:          int32PointerFromNullInt32(result.ValidPreparationMaximumVesselCount),
				IconPath:                    result.ValidPreparationIconPath,
				PastTense:                   result.ValidPreparationPastTense,
				ID:                          result.ValidPreparationID,
				Name:                        result.ValidPreparationName,
				Description:                 result.ValidPreparationDescription,
				Slug:                        result.ValidPreparationSlug,
				MinimumIngredientCount:      result.ValidPreparationMinimumIngredientCount,
				MinimumInstrumentCount:      result.ValidPreparationMinimumInstrumentCount,
				MinimumVesselCount:          result.ValidPreparationMinimumVesselCount,
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
	tracing.AttachValidPreparationVesselIDToSpan(span, input.ID)

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
	tracing.AttachValidPreparationVesselIDToSpan(span, updated.ID)

	if err := q.generatedQuerier.UpdateValidPreparationVessel(ctx, q.db, &generated.UpdateValidPreparationVesselParams{
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
	tracing.AttachValidPreparationVesselIDToSpan(span, validPreparationVesselID)

	if err := q.generatedQuerier.ArchiveValidPreparationVessel(ctx, q.db, validPreparationVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid preparation vessel")
	}

	logger.Info("valid preparation vessel archived")

	return nil
}
