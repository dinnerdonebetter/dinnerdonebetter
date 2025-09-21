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
	_ types.ValidInstrumentDataManager = (*repository)(nil)
)

// ValidInstrumentExists fetches whether a valid instrument exists from the database.
func (r *repository) ValidInstrumentExists(ctx context.Context, validInstrumentID string) (exists bool, err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validInstrumentID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	result, err := r.generatedQuerier.CheckValidInstrumentExistence(ctx, r.db, validInstrumentID)
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing valid instrument existence check")
	}

	return result, nil
}

// GetValidInstrument fetches a valid instrument from the database.
func (r *repository) GetValidInstrument(ctx context.Context, validInstrumentID string) (*types.ValidInstrument, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validInstrumentID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	result, err := r.generatedQuerier.GetValidInstrument(ctx, r.db, validInstrumentID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "getting valid instrument")
	}

	validInstrument := &types.ValidInstrument{
		CreatedAt:                      result.CreatedAt,
		LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
		IconPath:                       result.IconPath,
		ID:                             result.ID,
		Name:                           result.Name,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Slug:                           result.Slug,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		UsableForStorage:               result.UsableForStorage,
	}

	return validInstrument, nil
}

// GetRandomValidInstrument fetches a valid instrument from the database.
func (r *repository) GetRandomValidInstrument(ctx context.Context) (*types.ValidInstrument, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	result, err := r.generatedQuerier.GetRandomValidInstrument(ctx, r.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "scanning validInstrument")
	}

	validInstrument := &types.ValidInstrument{
		CreatedAt:                      result.CreatedAt,
		LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
		ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
		IconPath:                       result.IconPath,
		ID:                             result.ID,
		Name:                           result.Name,
		PluralName:                     result.PluralName,
		Description:                    result.Description,
		Slug:                           result.Slug,
		DisplayInSummaryLists:          result.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
		UsableForStorage:               result.UsableForStorage,
	}

	return validInstrument, nil
}

// SearchForValidInstruments fetches a valid instrument from the database.
func (r *repository) SearchForValidInstruments(ctx context.Context, query string) ([]*types.ValidInstrument, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if query == "" {
		return nil, database.ErrEmptyInputProvided
	}
	logger = logger.WithValue(keys.SearchQueryKey, query)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, query)

	results, err := r.generatedQuerier.SearchForValidInstruments(ctx, r.db, query)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments list retrieval query")
	}

	validInstruments := []*types.ValidInstrument{}
	for _, result := range results {
		validInstrument := &types.ValidInstrument{
			CreatedAt:                      result.CreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			IconPath:                       result.IconPath,
			ID:                             result.ID,
			Name:                           result.Name,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Slug:                           result.Slug,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			UsableForStorage:               result.UsableForStorage,
		}
		validInstruments = append(validInstruments, validInstrument)
	}

	return validInstruments, nil
}

// GetValidInstruments fetches a list of valid instruments from the database that meet a particular filter.
func (r *repository) GetValidInstruments(ctx context.Context, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.ValidInstrument], err error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &filtering.QueryFilteredResult[types.ValidInstrument]{
		Pagination: filter.ToPagination(),
	}

	results, err := r.generatedQuerier.GetValidInstruments(ctx, r.db, &generated.GetValidInstrumentsParams{
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:     database.NullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:      database.NullInt32FromUint8Pointer(filter.PageSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments list retrieval query")
	}

	for _, result := range results {
		x.Data = append(x.Data, &types.ValidInstrument{
			CreatedAt:                      result.CreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			IconPath:                       result.IconPath,
			ID:                             result.ID,
			Name:                           result.Name,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Slug:                           result.Slug,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			UsableForStorage:               result.UsableForStorage,
		})
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// GetValidInstrumentsWithIDs fetches a list of valid instruments from the database that meet a particular filter.
func (r *repository) GetValidInstrumentsWithIDs(ctx context.Context, ids []string) ([]*types.ValidInstrument, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	results, err := r.generatedQuerier.GetValidInstrumentsWithIDs(ctx, r.db, ids)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing valid instruments id list retrieval query")
	}

	instruments := []*types.ValidInstrument{}
	for _, result := range results {
		instruments = append(instruments, &types.ValidInstrument{
			CreatedAt:                      result.CreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.LastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ArchivedAt),
			IconPath:                       result.IconPath,
			ID:                             result.ID,
			Name:                           result.Name,
			PluralName:                     result.PluralName,
			Description:                    result.Description,
			Slug:                           result.Slug,
			DisplayInSummaryLists:          result.DisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.IncludeInGeneratedInstructions,
			UsableForStorage:               result.UsableForStorage,
		})
	}

	return instruments, nil
}

// GetValidInstrumentIDsThatNeedSearchIndexing fetches a list of valid instruments from the database that meet a particular filter.
func (r *repository) GetValidInstrumentIDsThatNeedSearchIndexing(ctx context.Context) ([]string, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	results, err := r.generatedQuerier.GetValidInstrumentsNeedingIndexing(ctx, r.db)
	if err != nil {
		return nil, observability.PrepareError(err, span, "executing valid instruments list retrieval query")
	}

	return results, nil
}

// CreateValidInstrument creates a valid instrument in the database.
func (r *repository) CreateValidInstrument(ctx context.Context, input *types.ValidInstrumentDatabaseCreationInput) (*types.ValidInstrument, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, input.ID)
	logger := r.logger.WithValue(keys.ValidInstrumentIDKey, input.ID)

	// create the valid instrument.
	if err := r.generatedQuerier.CreateValidInstrument(ctx, r.db, &generated.CreateValidInstrumentParams{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		Slug:                           input.Slug,
		UsableForStorage:               input.UsableForStorage,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing valid instrument creation query")
	}

	x := &types.ValidInstrument{
		ID:                             input.ID,
		Name:                           input.Name,
		PluralName:                     input.PluralName,
		Description:                    input.Description,
		IconPath:                       input.IconPath,
		UsableForStorage:               input.UsableForStorage,
		Slug:                           input.Slug,
		DisplayInSummaryLists:          input.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: input.IncludeInGeneratedInstructions,
		CreatedAt:                      r.CurrentTime(),
	}

	logger.Info("valid instrument created")

	return x, nil
}

// UpdateValidInstrument updates a particular valid instrument.
func (r *repository) UpdateValidInstrument(ctx context.Context, updated *types.ValidInstrument) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := r.logger.WithValue(keys.ValidInstrumentIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, updated.ID)

	if _, err := r.generatedQuerier.UpdateValidInstrument(ctx, r.db, &generated.UpdateValidInstrumentParams{
		Name:                           updated.Name,
		PluralName:                     updated.PluralName,
		Description:                    updated.Description,
		IconPath:                       updated.IconPath,
		Slug:                           updated.Slug,
		ID:                             updated.ID,
		UsableForStorage:               updated.UsableForStorage,
		DisplayInSummaryLists:          updated.DisplayInSummaryLists,
		IncludeInGeneratedInstructions: updated.IncludeInGeneratedInstructions,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating valid instrument")
	}

	logger.Info("valid instrument updated")

	return nil
}

// MarkValidInstrumentAsIndexed updates a particular valid instrument's last_indexed_at value.
func (r *repository) MarkValidInstrumentAsIndexed(ctx context.Context, validInstrumentID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validInstrumentID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	if _, err := r.generatedQuerier.UpdateValidInstrumentLastIndexedAt(ctx, r.db, validInstrumentID); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "marking valid instrument as indexed")
	}

	logger.Info("valid instrument marked as indexed")

	return nil
}

// ArchiveValidInstrument archives a valid instrument from the database by its ID.
func (r *repository) ArchiveValidInstrument(ctx context.Context, validInstrumentID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()

	if validInstrumentID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ValidInstrumentIDKey, validInstrumentID)
	tracing.AttachToSpan(span, keys.ValidInstrumentIDKey, validInstrumentID)

	rowsAffected, err := r.generatedQuerier.ArchiveValidInstrument(ctx, r.db, validInstrumentID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving valid instrument")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
