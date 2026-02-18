package payments

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/domain/payments"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	generated "github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments/generated"
)

func (r *repository) CreatePurchase(ctx context.Context, input *payments.PurchaseDatabaseCreationInput) (*payments.Purchase, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	logger = logger.WithValue(keys.PurchaseIDKey, input.ID)
	tracing.AttachToSpan(span, keys.PurchaseIDKey, input.ID)

	arg := &generated.CreatePurchaseParams{
		ID:                    input.ID,
		BelongsToAccount:      input.BelongsToAccount,
		ProductID:             input.ProductID,
		AmountCents:           input.AmountCents,
		Currency:              input.Currency,
		CompletedAt:           database.NullTimeFromTimePointer(input.CompletedAt),
		ExternalTransactionID: database.NullStringFromString(input.ExternalTransactionID),
	}

	if err := r.generatedQuerier.CreatePurchase(ctx, r.writeDB, arg); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating purchase")
	}

	return r.GetPurchase(ctx, input.ID)
}

func (r *repository) GetPurchase(ctx context.Context, id string) (*payments.Purchase, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.PurchaseIDKey, id)
	tracing.AttachToSpan(span, keys.PurchaseIDKey, id)

	result, err := r.generatedQuerier.GetPurchase(ctx, r.readDB, id)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching purchase")
	}

	return convertPurchaseFromGenerated(result), nil
}

func (r *repository) GetPurchasesForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Purchase], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetPurchasesForAccountParams{
		BelongsToAccount: accountID,
		CreatedAfter:     database.NullTimeFromTimePointer(filter.CreatedAfter),
		CreatedBefore:    database.NullTimeFromTimePointer(filter.CreatedBefore),
		UpdatedBefore:    database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:     database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:           database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:      database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived:  database.NullBoolFromBoolPointer(filter.IncludeArchived),
	}

	results, err := r.generatedQuerier.GetPurchasesForAccount(ctx, r.readDB, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching purchases for account")
	}

	return convertPurchasesResult(results, filter), nil
}

func convertPurchaseFromGenerated(m *generated.Purchases) *payments.Purchase {
	if m == nil {
		return nil
	}
	return &payments.Purchase{
		ID:                    m.ID,
		BelongsToAccount:      m.BelongsToAccount,
		ProductID:             m.ProductID,
		AmountCents:           m.AmountCents,
		Currency:              m.Currency,
		CompletedAt:           database.TimePointerFromNullTime(m.CompletedAt),
		ExternalTransactionID: database.StringFromNullString(m.ExternalTransactionID),
		CreatedAt:             m.CreatedAt,
		LastUpdatedAt:         database.TimePointerFromNullTime(m.LastUpdatedAt),
		ArchivedAt:            database.TimePointerFromNullTime(m.ArchivedAt),
	}
}

func convertPurchasesResult(rows []*generated.GetPurchasesForAccountRow, filter *filtering.QueryFilter) *filtering.QueryFilteredResult[payments.Purchase] {
	data := make([]*payments.Purchase, 0, len(rows))
	var filteredCount, totalCount uint64
	for _, row := range rows {
		data = append(data, &payments.Purchase{
			ID:                    row.ID,
			BelongsToAccount:      row.BelongsToAccount,
			ProductID:             row.ProductID,
			AmountCents:           row.AmountCents,
			Currency:              row.Currency,
			CompletedAt:           database.TimePointerFromNullTime(row.CompletedAt),
			ExternalTransactionID: database.StringFromNullString(row.ExternalTransactionID),
			CreatedAt:             row.CreatedAt,
			LastUpdatedAt:         database.TimePointerFromNullTime(row.LastUpdatedAt),
			ArchivedAt:            database.TimePointerFromNullTime(row.ArchivedAt),
		})
		filteredCount = uint64(row.FilteredCount)
		totalCount = uint64(row.TotalCount)
	}
	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(p *payments.Purchase) string { return p.ID },
		filter,
	)
}
