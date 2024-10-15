package postgres

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
)

// DeleteExpiredOAuth2ClientTokens deletes expired oauth2 client tokens.
func (q *Querier) DeleteExpiredOAuth2ClientTokens(ctx context.Context) error {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	if _, err := q.generatedQuerier.DeleteExpiredOAuth2ClientTokens(ctx, q.db); err != nil {
		return observability.PrepareError(err, span, "deleting expired oauth2 client tokens")
	}

	return nil
}
