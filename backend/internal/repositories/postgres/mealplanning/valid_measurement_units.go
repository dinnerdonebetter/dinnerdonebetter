package mealplanning

import (
	"context"
	"database/sql"

	types "github.com/dinnerdonebetter/backend/internal/domain/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	"github.com/dinnerdonebetter/backend/internal/repositories/postgres/mealplanning/generated"
)

var (
	_ types.ValidMeasurementUnitDataManager = (*repository)(nil)
)

// ValidMeasurementUnitExists fetches whether a valid measurement unit exists from the database.
func (r *repository) ValidMeasurementUnitExists(ctx context.Context, validMeasurementUnitID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validMeasurementUnitID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	result, err := r.generatedQuerier.CheckValidMeasurementUnitExistence(ctx, r.db, validMeasurementUnitID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid measurement unit existence check")
	}

	return result, nil
}

// GetValidMeasurementUnit fetches a valid measurement unit from the database.
func (r *repository) GetValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) (*types.ValidMeasurementUnit, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validMeasurementUnitID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	result, err := r.generatedQuerier.GetValidMeasurementUnit(ctx, r.db, validMeasurementUnitID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "scanning valid measurement unit")
	}

	validMeasurementUnit := &types.ValidMeasurementUnit{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		Name:          result.Name,
		IconPath:      result.IconPath,
		ID:            result.ID,
		Description:   result.Description,
		PluralName:    result.PluralName,
		Slug:          result.Slug,
		Volumetric:    database.BoolFromNullBool(result.Volumetric),
		Universal:     result.Universal,
		Metric:        result.Metric,
		Imperial:      result.Imperial,
	}

	return validMeasurementUnit, nil
}

// GetRandomValidMeasurementUnit fetches a valid measurement unit from the database.
func (r *repository) GetRandomValidMeasurementUnit(ctx context.Context) (*types.ValidMeasurementUnit, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	result, err := r.generatedQuerier.GetRandomValidMeasurementUnit(ctx, r.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning valid measurement unit")
	}

	validMeasurementUnit := &types.ValidMeasurementUnit{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
		Name:          result.Name,
		IconPath:      result.IconPath,
		ID:            result.ID,
		Description:   result.Description,
		PluralName:    result.PluralName,
		Slug:          result.Slug,
		Volumetric:    database.BoolFromNullBool(result.Volumetric),
		Universal:     result.Universal,
		Metric:        result.Metric,
		Imperial:      result.Imperial,
	}

	return validMeasurementUnit, nil
}

// SearchForValidMeasurementUnits fetches a valid measurement unit from the database.
func (r *repository) SearchForValidMeasurementUnits(ctx context.Context, query string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if query == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, query)

	results, err := r.generatedQuerier.SearchForValidMeasurementUnits(ctx, r.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	x := []*types.ValidMeasurementUnit{}
	for _, result := range results {
		x = append(x, &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		})
	}

	return x, nil
}

// ValidMeasurementUnitsForIngredientID fetches a valid measurement unit from the database.
func (r *repository) ValidMeasurementUnitsForIngredientID(ctx context.Context, validIngredientID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[types.ValidMeasurementUnit], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x := &filtering.QueryFilteredResult[types.ValidMeasurementUnit]{
		Pagination: filter.ToPagination(),
	}

	if validIngredientID == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.ValidIngredientIDKey, validIngredientID)
	tracing.AttachToSpan(span, keys.ValidIngredientIDKey, validIngredientID)

	results, err := r.generatedQuerier.SearchValidMeasurementUnitsByIngredientID(ctx, r.db, &generated.SearchValidMeasurementUnitsByIngredientIDParams{
		CreatedBefore:     database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:      database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:     database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:      database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:       database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:        database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived:   database.NullBoolFromBoolPointer(filter.IncludeArchived),
		ValidIngredientID: validIngredientID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	for _, result := range results {
		validMeasurementUnit := &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		}

		x.Data = append(x.Data, validMeasurementUnit)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// GetValidMeasurementUnits fetches a list of valid measurement units from the database that meet a particular filter.
func (r *repository) GetValidMeasurementUnits(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidMeasurementUnit], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidMeasurementUnit]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetValidMeasurementUnits(ctx, r.db, &generated.GetValidMeasurementUnitsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement units list retrieval query")
	}

	for _, result := range results {
		validMeasurementUnit := &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		}

		x.Data = append(x.Data, validMeasurementUnit)
		x.TotalCount = uint64(result.TotalCount)
		x.FilteredCount = uint64(result.FilteredCount)
	}

	return x, nil
}

// GetValidMeasurementUnitsWithIDs fetches a list of valid measurement unit from the database that meet a particular filter.
func (r *repository) GetValidMeasurementUnitsWithIDs(ctx context.Context, ids []string) ([]*types.ValidMeasurementUnit, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	results, err := r.generatedQuerier.GetValidMeasurementUnitsWithIDs(ctx, r.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid measurement unit id list retrieval query")
	}

	x := []*types.ValidMeasurementUnit{}
	for _, result := range results {
		x = append(x, &types.ValidMeasurementUnit{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:    database.TimePointerFromNullTime(result.ArchivedAt),
			Name:          result.Name,
			IconPath:      result.IconPath,
			ID:            result.ID,
			Description:   result.Description,
			PluralName:    result.PluralName,
			Slug:          result.Slug,
			Volumetric:    database.BoolFromNullBool(result.Volumetric),
			Universal:     result.Universal,
			Metric:        result.Metric,
			Imperial:      result.Imperial,
		})
	}

	return x, nil
}

// GetValidMeasurementUnitIDsThatNeedSearchIndexing fetches a list of valid measurement units from the database that meet a particular filter.
func (r *repository) GetValidMeasurementUnitIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	results, err := r.generatedQuerier.GetValidMeasurementUnitsNeedingIndexing(ctx, r.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid measurement units list retrieval query")
	}

	return results, nil
}

// CreateValidMeasurementUnit creates a valid measurement unit in the database.
func (r *repository) CreateValidMeasurementUnit(ctx context.Context, input *types.ValidMeasurementUnitDatabaseCreationInput) (*types.ValidMeasurementUnit, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, input.ID)
	logger := r.logger.WithValue(keys.ValidMeasurementUnitIDKey, input.ID)

	// create the valid measurement unit.
	if err := r.generatedQuerier.CreateValidMeasurementUnit(ctx, r.db, &generated.CreateValidMeasurementUnitParams{
		Name:        input.Name,
		Description: input.Description,
		IconPath:    input.IconPath,
		Slug:        input.Slug,
		PluralName:  input.PluralName,
		ID:          input.ID,
		Volumetric:  database.NullBoolFromBool(input.Volumetric),
		Universal:   input.Universal,
		Metric:      input.Metric,
		Imperial:    input.Imperial,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid measurement unit creation query")
	}

	x := &types.ValidMeasurementUnit{
		ID:          input.ID,
		Name:        input.Name,
		Description: input.Description,
		Volumetric:  input.Volumetric,
		IconPath:    input.IconPath,
		Universal:   input.Universal,
		Metric:      input.Metric,
		Imperial:    input.Imperial,
		Slug:        input.Slug,
		PluralName:  input.PluralName,
		CreatedAt:   r.CurrentTime(),
	}

	logger.Info("valid measurement unit created")

	return x, nil
}

// UpdateValidMeasurementUnit updates a particular valid measurement unit.
func (r *repository) UpdateValidMeasurementUnit(ctx context.Context, updated *types.ValidMeasurementUnit) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.ValidMeasurementUnitIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateValidMeasurementUnit(ctx, r.db, &generated.UpdateValidMeasurementUnitParams{
		Name:        updated.Name,
		Description: updated.Description,
		IconPath:    updated.IconPath,
		Slug:        updated.Slug,
		PluralName:  updated.PluralName,
		ID:          updated.ID,
		Volumetric:  database.NullBoolFromBool(updated.Volumetric),
		Universal:   updated.Universal,
		Metric:      updated.Metric,
		Imperial:    updated.Imperial,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid measurement unit")
	}

	logger.Info("valid measurement unit updated")

	return nil
}

// MarkValidMeasurementUnitAsIndexed updates a particular valid measurement unit's last_indexed_at value.
func (r *repository) MarkValidMeasurementUnitAsIndexed(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validMeasurementUnitID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	if _, err := r.generatedQuerier.UpdateValidMeasurementUnitLastIndexedAt(ctx, r.db, validMeasurementUnitID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid measurement unit as indexed")
	}

	logger.Info("valid measurement unit marked as indexed")

	return nil
}

// ArchiveValidMeasurementUnit archives a valid measurement unit from the database by its ID.
func (r *repository) ArchiveValidMeasurementUnit(ctx context.Context, validMeasurementUnitID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validMeasurementUnitID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)
	tracing.AttachToSpan(span, keys.ValidMeasurementUnitIDKey, validMeasurementUnitID)

	rowsAffected, err := r.generatedQuerier.ArchiveValidMeasurementUnit(ctx, r.db, validMeasurementUnitID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid measurement unit")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
