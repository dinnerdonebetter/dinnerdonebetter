package identity

import (
	"context"
	"database/sql"

	"github.com/dinnerdonebetter/backend/internal/domain/identity"
	"github.com/dinnerdonebetter/backend/internal/platform/database"
	"github.com/dinnerdonebetter/backend/internal/platform/observability"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/platform/observability/tracing"
)

var (
	_ identity.DataPrivacyDataManager = (*Querier)(nil)
)

// DeleteUser archives a user.
func (q *Querier) DeleteUser(ctx context.Context, userID string) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return database.ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.UserIDKey, userID)
	logger := q.logger.WithValue(keys.UserIDKey, userID)

	changed, err := q.generatedQuerier.DeleteUser(ctx, q.db, userID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving user")
	}

	if changed == 0 {
		return sql.ErrNoRows
	}

	logger.Info("user deleted")

	return nil
}
