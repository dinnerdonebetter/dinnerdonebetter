package postgres

import (
	"context"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.AdminUserDataManager = (*SQLQuerier)(nil)

// UpdateUserAccountStatus updates a user's household status.
func (q *SQLQuerier) UpdateUserAccountStatus(ctx context.Context, userID string, input *types.UserAccountStatusUpdateInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := &generated.SetUserAccountStatusParams{
		UserAccountStatus:            string(input.NewStatus),
		UserAccountStatusExplanation: input.Reason,
		ID:                           input.TargetUserID,
	}

	if err := q.generatedQuerier.SetUserAccountStatus(ctx, args); err != nil {
		return observability.PrepareError(err, logger, span, "user status update")
	}

	logger.Info("user account status updated")

	return nil
}
