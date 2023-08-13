package postgres

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/dinnerdonebetter/backend/internal/database"
	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"
)

const (
	validVesselsTable = "valid_vessels"
)

var (
	_ types.ValidVesselDataManager = (*Querier)(nil)

	// validVesselsTableColumns are the columns for the valid_vessels table.
	validVesselsTableColumns = []string{
		"valid_vessels.id",
		"valid_vessels.name",
		"valid_vessels.plural_name",
		"valid_vessels.description",
		"valid_vessels.icon_path",
		"valid_vessels.usable_for_storage",
		"valid_vessels.slug",
		"valid_vessels.display_in_summary_lists",
		"valid_vessels.include_in_generated_instructions",
		"valid_vessels.capacity",
		"valid_measurement_units.id",
		"valid_measurement_units.name",
		"valid_measurement_units.description",
		"valid_measurement_units.volumetric",
		"valid_measurement_units.icon_path",
		"valid_measurement_units.universal",
		"valid_measurement_units.metric",
		"valid_measurement_units.imperial",
		"valid_measurement_units.slug",
		"valid_measurement_units.plural_name",
		"valid_measurement_units.created_at",
		"valid_measurement_units.last_updated_at",
		"valid_measurement_units.archived_at",
		"valid_vessels.width_in_millimeters",
		"valid_vessels.length_in_millimeters",
		"valid_vessels.height_in_millimeters",
		"valid_vessels.shape",
		"valid_vessels.created_at",
		"valid_vessels.last_updated_at",
		"valid_vessels.archived_at",
	}
)

// scanValidVessel takes a database Scanner (i.e. *sql.Row) and scans the result into a valid vessel struct.
func (q *Querier) scanValidVessel(ctx context.Context, scan database.Scanner, includeCounts bool) (x *types.ValidVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	z := &types.NullableValidVessel{
		CapacityUnit: &types.NullableValidMeasurementUnit{},
	}

	targetVars := []any{
		&z.ID,
		&z.Name,
		&z.PluralName,
		&z.Description,
		&z.IconPath,
		&z.UsableForStorage,
		&z.Slug,
		&z.DisplayInSummaryLists,
		&z.IncludeInGeneratedInstructions,
		&z.Capacity,
		&z.CapacityUnit.ID,
		&z.CapacityUnit.Name,
		&z.CapacityUnit.Description,
		&z.CapacityUnit.Volumetric,
		&z.CapacityUnit.IconPath,
		&z.CapacityUnit.Universal,
		&z.CapacityUnit.Metric,
		&z.CapacityUnit.Imperial,
		&z.CapacityUnit.Slug,
		&z.CapacityUnit.PluralName,
		&z.CapacityUnit.CreatedAt,
		&z.CapacityUnit.LastUpdatedAt,
		&z.CapacityUnit.ArchivedAt,
		&z.WidthInMillimeters,
		&z.LengthInMillimeters,
		&z.HeightInMillimeters,
		&z.Shape,
		&z.CreatedAt,
		&z.LastUpdatedAt,
		&z.ArchivedAt,
	}

	if includeCounts {
		targetVars = append(targetVars, &filteredCount, &totalCount)
	}

	if err = scan.Scan(targetVars...); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "")
	}

	x = converters.ConvertNullableValidVesselToValidVessel(z)

	if x.CapacityUnit != nil && x.CapacityUnit.ID == "" {
		x.CapacityUnit = nil
	}

	return x, filteredCount, totalCount, nil
}

// scanValidVessels takes some database rows and turns them into a slice of valid vessels.
func (q *Querier) scanValidVessels(ctx context.Context, rows database.ResultIterator, includeCounts bool) (validVessels []*types.ValidVessel, filteredCount, totalCount uint64, err error) {
	_, span := q.tracer.StartSpan(ctx)
	defer span.End()

	for rows.Next() {
		x, fc, tc, scanErr := q.scanValidVessel(ctx, rows, includeCounts)
		if scanErr != nil {
			return nil, 0, 0, scanErr
		}

		if includeCounts {
			if filteredCount == 0 {
				filteredCount = fc
			}

			if totalCount == 0 {
				totalCount = tc
			}
		}

		validVessels = append(validVessels, x)
	}

	if err = q.checkRowsForErrorAndClose(ctx, rows); err != nil {
		return nil, 0, 0, observability.PrepareError(err, span, "handling rows")
	}

	return validVessels, filteredCount, totalCount, nil
}

// ValidVesselExists fetches whether a valid vessel exists from the database.
func (q *Querier) ValidVesselExists(ctx context.Context, validVesselID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachValidVesselIDToSpan(span, validVesselID)

	result, err := q.generatedQuerier.CheckValidVesselExistence(ctx, q.db, validVesselID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid vessel existence check")
	}

	return result, nil
}

// GetValidVessel fetches a valid vessel from the database.
func (q *Querier) GetValidVessel(ctx context.Context, validVesselID string) (*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachValidVesselIDToSpan(span, validVesselID)

	result, err := q.generatedQuerier.GetValidVessel(ctx, q.db, validVesselID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid vessel existence check")
	}

	validVessel := &types.ValidVessel{
		CreatedAt:     result.CreatedAt,
		ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
		CapacityUnit: &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt_2,
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt_2),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt_2),
			Name:          result.Name_2,
			IconPath:      result.IconPath_2,
			ID:            result.ID_2,
			Description:   result.Description_2,
			PluralName:    result.PluralName_2,
			Slug:          result.Slug_2,
			Volumetric:    boolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		},
		IconPath:                       result.IconPath,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Name:                           result.Name,
		Slug:                           result.Slug,
		Shape:                          string(result.Shape),
		ID:                             result.ID,
		WidthInMillimeters:             float32(result.ValidVesselsWidthInMillimeters),
		LengthInMillimeters:            float32(result.ValidVesselsLengthInMillimeters),
		HeightInMillimeters:            float32(result.ValidVesselsHeightInMillimeters),
		Capacity:                       float32(result.ValidVesselsCapacity),
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		UsableForStorage:               result.UsableForStorage,
	}

	return validVessel, nil
}

// GetRandomValidVessel fetches a valid vessel from the database.
func (q *Querier) GetRandomValidVessel(ctx context.Context) (*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	result, err := q.generatedQuerier.GetRandomValidVessel(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "querying for random valid vessel")
	}

	validVessel := &types.ValidVessel{
		CreatedAt:     result.CreatedAt,
		ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
		CapacityUnit: &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt_2,
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt_2),
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt_2),
			Name:          result.Name_2,
			IconPath:      result.IconPath_2,
			ID:            result.ID_2,
			Description:   result.Description_2,
			PluralName:    result.PluralName_2,
			Slug:          result.Slug_2,
			Volumetric:    boolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		},
		IconPath:                       result.IconPath,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Name:                           result.Name,
		Slug:                           result.Slug,
		Shape:                          string(result.Shape),
		ID:                             result.ID,
		WidthInMillimeters:             float32(result.ValidVesselsWidthInMillimeters),
		LengthInMillimeters:            float32(result.ValidVesselsLengthInMillimeters),
		HeightInMillimeters:            float32(result.ValidVesselsHeightInMillimeters),
		Capacity:                       float32(result.ValidVesselsCapacity),
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		UsableForStorage:               result.UsableForStorage,
	}

	return validVessel, nil
}

//go:embed queries/valid_vessels/search.sql
var validVesselSearchQuery string

// SearchForValidVessels fetches a valid vessel from the database.
func (q *Querier) SearchForValidVessels(ctx context.Context, query string) ([]*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if query == "" {
		return nil, ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachValidVesselIDToSpan(span, query)

	args := []any{
		wrapQueryForILIKE(query),
	}

	rows, err := q.getRows(ctx, q.db, "valid vessels", validVesselSearchQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid vessels list retrieval query")
	}

	validVessels, _, _, err := q.scanValidVessels(ctx, rows, false)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid vessel")
	}

	return validVessels, nil
}

//go:embed queries/valid_vessels/get_many.sql
var getValidVesselsQuery string

// GetValidVessels fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetValidVessels(ctx context.Context, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.ValidVessel], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	x = &types.QueryFilteredResult[types.ValidVessel]{}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if filter != nil {
		if filter.Page != nil {
			x.Page = *filter.Page
		}

		if filter.Limit != nil {
			x.Limit = *filter.Limit
		}
	} else {
		filter = types.DefaultQueryFilter()
	}

	args := []any{
		filter.CreatedBefore,
		filter.CreatedAfter,
		filter.UpdatedBefore,
		filter.UpdatedAfter,
		filter.QueryOffset(),
		filter.Limit,
	}

	rows, err := q.getRows(ctx, q.db, "valid vessels", getValidVesselsQuery, args)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid vessels list retrieval query")
	}

	if x.Data, x.FilteredCount, x.TotalCount, err = q.scanValidVessels(ctx, rows, true); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid vessels")
	}

	return x, nil
}

// GetValidVesselsWithIDs fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetValidVesselsWithIDs(ctx context.Context, ids []string) ([]*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	results, err := q.generatedQuerier.GetValidVesselsWithIDs(ctx, q.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid vessels id list retrieval query")
	}

	validVessels := []*types.ValidVessel{}
	for _, result := range results {
		validVessels = append(validVessels, &types.ValidVessel{
			CreatedAt:     result.CreatedAt,
			ArchivedAt:    timePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			CapacityUnit: &types.ValidMeasurementUnit{
				CreatedAt:     result.CreatedAt_2,
				LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt_2),
				ArchivedAt:    timePointerFromNullTime(result.ArchivedAt_2),
				Name:          result.Name_2,
				IconPath:      result.IconPath_2,
				ID:            result.ID_2,
				Description:   result.Description_2,
				PluralName:    result.PluralName_2,
				Slug:          result.Slug_2,
				Volumetric:    boolFromNullBool(result.Volumetric),
				Universal:     result.Universal,
				Metric:        result.Metric,
				Imperial:      result.Imperial,
			},
			IconPath:                       result.IconPath,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Name:                           result.Name,
			Slug:                           result.Slug,
			Shape:                          string(result.Shape),
			ID:                             result.ID,
			WidthInMillimeters:             float32(result.ValidVesselsWidthInMillimeters),
			LengthInMillimeters:            float32(result.ValidVesselsLengthInMillimeters),
			HeightInMillimeters:            float32(result.ValidVesselsHeightInMillimeters),
			Capacity:                       float32(result.ValidVesselsCapacity),
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			UsableForStorage:               result.UsableForStorage,
		})
	}

	return validVessels, nil
}

// GetValidVesselIDsThatNeedSearchIndexing fetches a list of valid vessels from the database that meet a particular filter.
func (q *Querier) GetValidVesselIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	results, err := q.generatedQuerier.GetValidVesselIDsNeedingIndexing(ctx, q.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid vessels list retrieval query")
	}

	return results, nil
}

// CreateValidVessel creates a valid vessel in the database.
func (q *Querier) CreateValidVessel(ctx context.Context, input *types.ValidVesselDatabaseCreationInput) (*types.ValidVessel, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidVesselIDKey, input.ID)

	// create the valid vessel.
	if err := q.generatedQuerier.CreateValidVessel(ctx, q.db, &generated.CreateValidVesselParams{
		Slug:                           input.Slug,
		ID:                             input.ID,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		Shape:                          generated.VesselShape(input.Shape),
		Name:                           input.Name,
		Capacity:                       float64(input.Capacity),
		CapacityUnit:                   nullStringFromStringPointer(input.CapacityUnitID),
		WidthInMillimeters:             float64(input.WidthInMillimeters),
		LengthInMillimeters:            float64(input.LengthInMillimeters),
		HeightInMillimeters:            float64(input.HeightInMillimeters),
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		UsableForStorage:               input.UsableForStorage,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid vessel creation query")
	}

	x := &types.ValidVessel{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		Slug:                           input.Slug,
		Shape:                          input.Shape,
		WidthInMillimeters:             input.WidthInMillimeters,
		LengthInMillimeters:            input.LengthInMillimeters,
		HeightInMillimeters:            input.HeightInMillimeters,
		Capacity:                       input.Capacity,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		UsableForStorage:               input.UsableForStorage,
		CreatedAt:                      q.currentTime(),
	}

	if input.CapacityUnitID != nil {
		x.CapacityUnit = &types.ValidMeasurementUnit{ID: *input.CapacityUnitID}
	}

	tracing.AttachValidVesselIDToSpan(span, x.ID)
	logger.Info("valid vessel created")

	return x, nil
}

// UpdateValidVessel updates a particular valid vessel.
func (q *Querier) UpdateValidVessel(ctx context.Context, updated *types.ValidVessel) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.ValidVesselIDKey, updated.ID)
	tracing.AttachValidVesselIDToSpan(span, updated.ID)

	if updated.CapacityUnit == nil {
		return fmt.Errorf("capacity unit: %w", ErrNilInputProvided)
	}

	if err := q.generatedQuerier.UpdateValidVessel(ctx, q.db, &generated.UpdateValidVesselParams{
		Name:                           updated.Name,
		PluralName:                     updated.PluralName,
		Description:                    updated.Description,
		IconPath:                       updated.IconPath,
		UsableForStorage:               updated.UsableForStorage,
		Slug:                           updated.Slug,
		DisplayInSummaryLists:          updated.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: updated.IncludeInGeneratedInstructions,
		Capacity:                       float64(updated.Capacity),
		CapacityUnit:                   nullStringFromString(updated.CapacityUnit.ID),
		WidthInMillimeters:             float64(updated.WidthInMillimeters),
		LengthInMillimeters:            float64(updated.LengthInMillimeters),
		HeightInMillimeters:            float64(updated.HeightInMillimeters),
		Shape:                          generated.VesselShape(updated.Shape),
		ID:                             updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid vessel")
	}

	logger.Info("valid vessel updated")

	return nil
}

// MarkValidVesselAsIndexed updates a particular valid vessel's last_indexed_at value.
func (q *Querier) MarkValidVesselAsIndexed(ctx context.Context, validVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachValidVesselIDToSpan(span, validVesselID)

	if err := q.generatedQuerier.UpdateValidVesselLastIndexedAt(ctx, q.db, validVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid vessel as indexed")
	}

	logger.Info("valid vessel marked as indexed")

	return nil
}

// ArchiveValidVessel archives a valid vessel from the database by its ID.
func (q *Querier) ArchiveValidVessel(ctx context.Context, validVesselID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if validVesselID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidVesselIDKey, validVesselID)
	tracing.AttachValidVesselIDToSpan(span, validVesselID)

	if err := q.generatedQuerier.ArchiveValidVessel(ctx, q.db, validVesselID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid vessel")
	}

	logger.Info("valid vessel archived")

	return nil
}
