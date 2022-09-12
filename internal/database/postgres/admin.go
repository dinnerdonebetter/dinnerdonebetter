package postgres

import (
	"context"
	_ "embed"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

var _ types.AdminUserDataManager = (*Querier)(nil)

//go:embed queries/admin_set_user_account_status.sql
var setUserAccountStatusQuery string

// UpdateUserAccountStatus updates a user's household status.
func (q *Querier) UpdateUserAccountStatus(ctx context.Context, userID string, input *types.UserAccountStatusUpdateInput) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithValue(keys.UserIDKey, userID)
	tracing.AttachUserIDToSpan(span, userID)

	args := []interface{}{
		input.NewStatus,
		input.Reason,
		input.TargetUserID,
	}

	if err := q.performWriteQuery(ctx, q.db, "user status update query", setUserAccountStatusQuery, args); err != nil {
		return observability.PrepareError(err, span, "user status update")
	}

	logger.Info("user account status updated")

	return nil
}
