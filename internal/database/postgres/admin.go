package postgres

import (
	"context"
	_ "embed"

	"github.com/dinnerdonebetter/backend/internal/database/postgres/generated"
	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

var _ types.AdminUserDataManager = (*Querier)(nil)

// UpdateUserAccountStatus updates a user's household status.
func (q *Querier) UpdateUserAccountStatus(ctx context.Context, userID string, input *types.UserAccountStatusUpdateInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	if err := q.generatedQuerier.SetUserAccountStatus(ctx, q.db, &generated.SetUserAccountStatusParams{
		UserAccountStatus:            input.NewStatus,
		UserAccountStatusExplanation: input.Reason,
		ID:                           input.TargetUserID,
	}); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "user status update")
	}

	logger.Info("user account status updated")

	return nil
}
