package payments

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/payments"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/database/filtering"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
	generated "github.com/dinnerdonebetter/backend/internal/repositories/postgres/payments/generated"
)

func (r *repository) CreateSubscription(ctx context.Context, input *payments.SubscriptionDatabaseCreationInput) (*payments.Subscription, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	logger = logger.WithValue(keys.SubscriptionIDKey, input.ID)
	tracing.AttachToSpan(span, keys.SubscriptionIDKey, input.ID)

	arg := &generated.CreateSubscriptionParams{
		ID:                     input.ID,
		BelongsToAccount:       input.BelongsToAccount,
		ProductID:              input.ProductID,
		ExternalSubscriptionID: database.NullStringFromString(input.ExternalSubscriptionID),
		Status:                 generated.SubscriptionStatus(input.Status),
		CurrentPeriodStart:     input.CurrentPeriodStart,
		CurrentPeriodEnd:       input.CurrentPeriodEnd,
	}

	if err := r.generatedQuerier.CreateSubscription(ctx, r.writeDB, arg); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating subscription")
	}

	return r.GetSubscription(ctx, input.ID)
}

func (r *repository) GetSubscription(ctx context.Context, id string) (*payments.Subscription, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return nil, database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SubscriptionIDKey, id)
	tracing.AttachToSpan(span, keys.SubscriptionIDKey, id)

	result, err := r.generatedQuerier.GetSubscription(ctx, r.readDB, id)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching subscription")
	}

	return convertSubscriptionFromGenerated(result), nil
}

func (r *repository) GetSubscriptionByExternalID(ctx context.Context, externalID string) (*payments.Subscription, error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	logger = logger.WithValue("external_subscription_id", externalID)
	tracing.AttachToSpan(span, "external_subscription_id", externalID)

	arg := sql.NullString{String: externalID, Valid: externalID != ""}
	result, err := r.generatedQuerier.GetSubscriptionByExternalID(ctx, r.readDB, arg)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching subscription by external ID")
	}

	return convertSubscriptionFromGenerated(result), nil
}

func (r *repository) GetSubscriptionsForAccount(ctx context.Context, accountID string, filter *filtering.QueryFilter) (*filtering.QueryFilteredResult[payments.Subscription], error) {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if filter == nil {
		filter = filtering.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	params := &generated.GetSubscriptionsForAccountParams{
		BelongsToAccount: accountID,
		CreatedAfter:     database.NullTimeFromTimePointer(filter.CreatedAfter),
		CreatedBefore:    database.NullTimeFromTimePointer(filter.CreatedBefore),
		UpdatedBefore:    database.NullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:     database.NullTimeFromTimePointer(filter.UpdatedAfter),
		Cursor:           database.NullStringFromStringPointer(filter.Cursor),
		ResultLimit:      database.NullInt32FromUint8Pointer(filter.MaxResponseSize),
		IncludeArchived:  database.NullBoolFromBoolPointer(filter.IncludeArchived),
	}

	results, err := r.generatedQuerier.GetSubscriptionsForAccount(ctx, r.readDB, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching subscriptions for account")
	}

	return convertSubscriptionsResult(results, filter), nil
}

func (r *repository) UpdateSubscription(ctx context.Context, sub *payments.Subscription) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	logger = logger.WithValue(keys.SubscriptionIDKey, sub.ID)
	tracing.AttachToSpan(span, keys.SubscriptionIDKey, sub.ID)

	arg := &generated.UpdateSubscriptionParams{
		ID:                     sub.ID,
		ExternalSubscriptionID: database.NullStringFromString(sub.ExternalSubscriptionID),
		Status:                 generated.SubscriptionStatus(sub.Status),
		CurrentPeriodStart:     sub.CurrentPeriodStart,
		CurrentPeriodEnd:       sub.CurrentPeriodEnd,
	}

	_, err := r.generatedQuerier.UpdateSubscription(ctx, r.writeDB, arg)
	return observability.PrepareAndLogError(err, logger, span, "updating subscription")
}

func (r *repository) UpdateSubscriptionStatus(ctx context.Context, id, status string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SubscriptionIDKey, id)
	tracing.AttachToSpan(span, keys.SubscriptionIDKey, id)

	arg := &generated.UpdateSubscriptionStatusParams{
		ID:     id,
		Status: generated.SubscriptionStatus(status),
	}

	_, err := r.generatedQuerier.UpdateSubscriptionStatus(ctx, r.writeDB, arg)
	return observability.PrepareAndLogError(err, logger, span, "updating subscription status")
}

func (r *repository) ArchiveSubscription(ctx context.Context, id string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	logger := r.logger.Clone()
	if id == "" {
		return database.ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.SubscriptionIDKey, id)
	tracing.AttachToSpan(span, keys.SubscriptionIDKey, id)

	_, err := r.generatedQuerier.ArchiveSubscription(ctx, r.writeDB, id)
	return observability.PrepareAndLogError(err, logger, span, "archiving subscription")
}

func convertSubscriptionFromGenerated(m *generated.Subscriptions) *payments.Subscription {
	if m == nil {
		return nil
	}
	return &payments.Subscription{
		ID:                     m.ID,
		BelongsToAccount:       m.BelongsToAccount,
		ProductID:              m.ProductID,
		ExternalSubscriptionID: database.StringFromNullString(m.ExternalSubscriptionID),
		Status:                 string(m.Status),
		CurrentPeriodStart:     m.CurrentPeriodStart,
		CurrentPeriodEnd:       m.CurrentPeriodEnd,
		CreatedAt:              m.CreatedAt,
		LastUpdatedAt:          database.TimePointerFromNullTime(m.LastUpdatedAt),
		ArchivedAt:             database.TimePointerFromNullTime(m.ArchivedAt),
	}
}

func convertSubscriptionsResult(rows []*generated.GetSubscriptionsForAccountRow, filter *filtering.QueryFilter) *filtering.QueryFilteredResult[payments.Subscription] {
	data := make([]*payments.Subscription, 0, len(rows))
	var filteredCount, totalCount uint64
	for _, row := range rows {
		data = append(data, &payments.Subscription{
			ID:                     row.ID,
			BelongsToAccount:       row.BelongsToAccount,
			ProductID:              row.ProductID,
			ExternalSubscriptionID: database.StringFromNullString(row.ExternalSubscriptionID),
			Status:                 string(row.Status),
			CurrentPeriodStart:     row.CurrentPeriodStart,
			CurrentPeriodEnd:       row.CurrentPeriodEnd,
			CreatedAt:              row.CreatedAt,
			LastUpdatedAt:          database.TimePointerFromNullTime(row.LastUpdatedAt),
			ArchivedAt:             database.TimePointerFromNullTime(row.ArchivedAt),
		})
		filteredCount = uint64(row.FilteredCount)
		totalCount = uint64(row.TotalCount)
	}
	return filtering.NewQueryFilteredResult(
		data,
		filteredCount,
		totalCount,
		func(p *payments.Subscription) string { return p.ID },
		filter,
	)
}
