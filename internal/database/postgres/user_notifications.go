package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var (
	_ types.UserNotificationDataManager = (*Querier)(nil)
)

// UserNotificationExists fetches whether a user notification exists from the database.
func (q *Querier) UserNotificationExists(ctx context.Context, userID, userNotificationID string) (exists bool, err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if userNotificationID == "" {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserNotificationIDKey, userNotificationID)
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, userNotificationID)

	result, err := q.generatedQuerier.CheckUserNotificationExistence(ctx, q.db, &generated.CheckUserNotificationExistenceParams{
		ID:            userNotificationID,
		BelongsToUser: userID,
	})
	if err != nil {
		return false, observability.PrepareAndLogError(err, logger, span, "performing user notification existence check")
	}

	return result, nil
}

// GetUserNotification fetches a user notification from the database.
func (q *Querier) GetUserNotification(ctx context.Context, userID, userNotificationID string) (*types.UserNotification, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if userNotificationID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserNotificationIDKey, userNotificationID)
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, userNotificationID)

	result, err := q.generatedQuerier.GetUserNotification(ctx, q.db, &generated.GetUserNotificationParams{
		BelongsToUser: userID,
		ID:            userNotificationID,
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "fetching user notification")
	}

	userNotification := &types.UserNotification{
		CreatedAt:     result.CreatedAt,
		LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
		ID:            result.ID,
		Content:       result.Content,
		Status:        string(result.Status),
		BelongsToUser: result.BelongsToUser,
	}

	return userNotification, nil
}

// GetUserNotifications fetches a list of user notifications from the database that meet a particular filter.
func (q *Querier) GetUserNotifications(ctx context.Context, userID string, filter *types.QueryFilter) (x *types.QueryFilteredResult[types.UserNotification], err error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.Clone()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachToSpan(span, keys.UserIDKey, userID)

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	x = &types.QueryFilteredResult[types.UserNotification]{
		Pagination: filter.ToPagination(),
	}

	results, err := q.generatedQuerier.GetUserNotificationsForUser(ctx, q.db, &generated.GetUserNotificationsForUserParams{
		UserID:        userID,
		CreatedBefore: nullTimeFromTimePointer(filter.CreatedBefore),
		CreatedAfter:  nullTimeFromTimePointer(filter.CreatedAfter),
		UpdatedBefore: nullTimeFromTimePointer(filter.UpdatedBefore),
		UpdatedAfter:  nullTimeFromTimePointer(filter.UpdatedAfter),
		QueryOffset:   nullInt32FromUint16(filter.QueryOffset()),
		QueryLimit:    nullInt32FromUint8Pointer(filter.Limit),
	})
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "executing user notifications list retrieval query")
	}

	for _, result := range results {
		userNotification := &types.UserNotification{
			CreatedAt:     result.CreatedAt,
			LastUpdatedAt: timePointerFromNullTime(result.LastUpdatedAt),
			ID:            result.ID,
		}

		x.Data = append(x.Data, userNotification)
		x.FilteredCount = uint64(result.FilteredCount)
		x.TotalCount = uint64(result.TotalCount)
	}

	return x, nil
}

// CreateUserNotification creates a user notification in the database.
func (q *Querier) CreateUserNotification(ctx context.Context, input *types.UserNotificationDatabaseCreationInput) (*types.UserNotification, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, input.ID)
	logger := q.logger.WithValue(keys.UserNotificationIDKey, input.ID)

	// create the user notification.
	if err := q.generatedQuerier.CreateUserNotification(ctx, q.db, &generated.CreateUserNotificationParams{
		ID:            input.ID,
		Content:       input.Content,
		BelongsToUser: input.BelongsToUser,
	}); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "performing user notification creation query")
	}

	x := &types.UserNotification{
		ID:            input.ID,
		CreatedAt:     q.currentTime(),
		Content:       input.Content,
		Status:        types.UserNotificationStatusTypeUnread,
		BelongsToUser: input.BelongsToUser,
	}
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, x.ID)
	logger.Info("user notification created")

	return x, nil
}

// UpdateUserNotification updates a particular user notification.
func (q *Querier) UpdateUserNotification(ctx context.Context, updated *types.UserNotification) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if updated == nil {
		return ErrNilInputProvided
	}
	logger := q.logger.WithValue(keys.UserNotificationIDKey, updated.ID)
	tracing.AttachToSpan(span, keys.UserNotificationIDKey, updated.ID)

	if _, err := q.generatedQuerier.UpdateUserNotification(ctx, q.db, &generated.UpdateUserNotificationParams{
		Status: generated.UserNotificationStatus(updated.Status),
		ID:     updated.ID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating user notification")
	}

	logger.Info("user notification updated")

	return nil
}
