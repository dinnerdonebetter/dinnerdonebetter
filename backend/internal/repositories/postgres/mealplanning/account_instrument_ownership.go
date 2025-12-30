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
	_ types.AccountInstrumentOwnershipDataManager = (*repository)(nil)
)

// AccountInstrumentOwnershipExists fetches whether an account instrument ownership exists from the database.
func (q *repository) AccountInstrumentOwnershipExists(ctx context.Context, accountInstrumentOwnershipID, accountID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountInstrumentOwnershipID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)

	if accountID == "" {
		return false, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	result, err := q.generatedQuerier.CheckAccountInstrumentOwnershipExistence(ctx, q.db, &generated.CheckAccountInstrumentOwnershipExistenceParams{
		ID:               accountInstrumentOwnershipID,
		BelongsToAccount: accountID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing account instrument ownership existence check")
	}

	return result, nil
}

// GetAccountInstrumentOwnership fetches an account instrument ownership from the database.
func (q *repository) GetAccountInstrumentOwnership(ctx context.Context, accountInstrumentOwnershipID, accountID string) (*types.AccountInstrumentOwnership, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountInstrumentOwnershipID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)

	if accountID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	result, err := q.generatedQuerier.GetAccountInstrumentOwnership(ctx, q.db, &generated.GetAccountInstrumentOwnershipParams{
		ID:               accountInstrumentOwnershipID,
		BelongsToAccount: accountID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching account instrument ownership")
	}

	accountInstrumentOwnership := &types.AccountInstrumentOwnership{
		CreatedAt:        result.CreatedAt,
		ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
		LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
		ID:               result.ID,
		Notes:            result.Notes,
		BelongsToAccount: result.BelongsToAccount,
		Quantity:         uint16(result.Quantity),
		Instrument: types.ValidInstrument{
			CreatedAt:                      result.ValidInstrumentCreatedAt,
			LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
			ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
			IconPath:                       result.ValidInstrumentIconPath,
			ID:                             result.ValidInstrumentID,
			Name:                           result.ValidInstrumentName,
			PluralName:                     result.ValidInstrumentPluralName,
			Description:                    result.ValidInstrumentDescription,
			Slug:                           result.ValidInstrumentSlug,
			DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists,
			IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions,
			UsableForStorage:               result.ValidInstrumentUsableForStorage,
		},
	}

	return accountInstrumentOwnership, nil
}

// GetAccountInstrumentOwnerships fetches a list of account instrument ownerships from the database that meet a particular filter.
func (q *repository) GetAccountInstrumentOwnerships(ctx context.Context, accountID string, filter *filtering.QueryFilter) (x *filtering.QueryFilteredResult[types.AccountInstrumentOwnership], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	var (
		data          []*types.AccountInstrumentOwnership
		filteredCount uint64
		totalCount    uint64
	)

	results, err := q.generatedQuerier.GetAccountInstrumentOwnerships(ctx, q.db, &generated.GetAccountInstrumentOwnershipsParams{
		AccountID:       accountID,
		CreatedBefore:   database.NullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:    database.NullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore:   database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:    database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:          database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:     database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived: database.NullBoolFromBoolPointer(filter.IncludeArchived),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing account instrument ownerships list retrieval query")
	}

	for _, result := range results {
		if totalCount == 0 {
			filteredCount = uint64(result.FilteredCount)
			totalCount = uint64(result.TotalCount)
		}
		data = append(data, &types.AccountInstrumentOwnership{
			CreatedAt:        result.CreatedAt,
			ArchivedAt:       database.TimePointerFromNullTime(result.ArchivedAt),
			LastUpdatedAt:    database.TimePointerFromNullTime(result.LastUpdatedAt),
			ID:               result.ID,
			Notes:            result.Notes,
			BelongsToAccount: result.BelongsToAccount,
			Quantity:         uint16(result.Quantity),
			Instrument: types.ValidInstrument{
				CreatedAt:                      result.ValidInstrumentCreatedAt,
				LastUpdatedAt:                  database.TimePointerFromNullTime(result.ValidInstrumentLastUpdatedAt),
				ArchivedAt:                     database.TimePointerFromNullTime(result.ValidInstrumentArchivedAt),
				IconPath:                       result.ValidInstrumentIconPath,
				ID:                             result.ValidInstrumentID,
				Name:                           result.ValidInstrumentName,
				PluralName:                     result.ValidInstrumentPluralName,
				Description:                    result.ValidInstrumentDescription,
				Slug:                           result.ValidInstrumentSlug,
				DisplayInSummaryLists:          result.ValidInstrumentDisplayInSummaryLists,
				IncludeInGeneratedInstructions: result.ValidInstrumentIncludeInGeneratedInstructions,
				UsableForStorage:               result.ValidInstrumentUsableForStorage,
			},
		})
	}

	x = filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(aio *types.AccountInstrumentOwnership) string { return aio.ID },
		filter,
	)

	return x, nil
}

// CreateAccountInstrumentOwnership creates an account instrument ownership in the database.
func (q *repository) CreateAccountInstrumentOwnership(ctx context.Context, input *types.AccountInstrumentOwnershipDatabaseCreationInput) (*types.AccountInstrumentOwnership, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, database.ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, input.ID)
	logger := q.logger.WithValue(keys.AccountInstrumentOwnershipIDKey, input.ID)

	// create the account instrument ownership.
	if err := q.generatedQuerier.CreateAccountInstrumentOwnership(ctx, q.db, &generated.CreateAccountInstrumentOwnershipParams{
		ID:                input.ID,
		Notes:             input.Notes,
		ValidInstrumentID: input.ValidInstrumentID,
		BelongsToAccount:  input.BelongsToAccount,
		Quantity:          int32(input.Quantity),
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing account instrument ownership creation query")
	}

	x := &types.AccountInstrumentOwnership{
		ID:               input.ID,
		Notes:            input.Notes,
		Quantity:         input.Quantity,
		Instrument:       types.ValidInstrument{ID: input.ValidInstrumentID},
		BelongsToAccount: input.BelongsToAccount,
		CreatedAt:        q.CurrentTime(),
	}

	logger.Info("account instrument ownership created")

	return x, nil
}

// UpdateAccountInstrumentOwnership updates a particular account instrument ownership.
func (q *repository) UpdateAccountInstrumentOwnership(ctx context.Context, updated *types.AccountInstrumentOwnership) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return database.ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.AccountInstrumentOwnershipIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateAccountInstrumentOwnership(ctx, q.db, &generated.UpdateAccountInstrumentOwnershipParams{
		Notes:             updated.Notes,
		ValidInstrumentID: updated.Instrument.ID,
		ID:                updated.ID,
		BelongsToAccount:  updated.BelongsToAccount,
		Quantity:          int32(updated.Quantity),
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating account instrument ownership")
	}

	logger.Info("account instrument ownership updated")

	return nil
}

// ArchiveAccountInstrumentOwnership archives an account instrument ownership from the database by its MealPlanTaskID.
func (q *repository) ArchiveAccountInstrumentOwnership(ctx context.Context, accountInstrumentOwnershipID, accountID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if accountInstrumentOwnershipID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)
	tracing.AttachToSpan(span, keys.AccountInstrumentOwnershipIDKey, accountInstrumentOwnershipID)

	if accountID == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AccountIDKey, accountID)
	tracing.AttachToSpan(span, keys.AccountIDKey, accountID)

	rowsAffected, err := q.generatedQuerier.ArchiveAccountInstrumentOwnership(ctx, q.db, &generated.ArchiveAccountInstrumentOwnershipParams{
		ID:               accountInstrumentOwnershipID,
		BelongsToAccount: accountID,
	})
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving account instrument ownership")
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
