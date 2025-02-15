package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

// DeleteExpiredOAuth2ClientTokens deletes expired oauth2 client tokens.
func (q *Querier) DeleteExpiredOAuth2ClientTokens(ctx context.Context) (int64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	deleted, err := q.generatedQuerier.DeleteExpiredOAuth2ClientTokens(ctx, q.db)
	if err != nil {
		return 0, observability.PrepareError(err, span, "deleting expired oauth2 client tokens")
	}

	return deleted, nil
}
