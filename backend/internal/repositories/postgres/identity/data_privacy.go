package identity

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/dataprivacy"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var (
	_ dataprivacy.DataPrivacyDataManager = (*repository)(nil)
)

// DeleteUser archives a user.
func (r *repository) DeleteUser(ctx context.Context, userID string) error {
	ctx, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger := r.logger.WithValue(keys.UserIDKey, userID)

	changed, err := r.generatedQuerier.DeleteUser(ctx, r.writeDB, userID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving user")
	}

	if changed == 0 {
		return sql.ErrNoRows
	}

	logger.Info("user deleted")

	return nil
}

func (r *repository) FetchUserDataCollection(ctx context.Context, userID string) (*dataprivacy.UserDataCollectionResponse, error) {
	_, span := r.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger := observability.ObserveValues(map[string]any{
		keys.UserIDKey: userID,
	}, span, r.logger)

	logger.Info("TODO")

	return &dataprivacy.UserDataCollectionResponse{}, nil
}
