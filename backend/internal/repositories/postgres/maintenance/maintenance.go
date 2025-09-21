package maintenance

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/platform/observability"
)

// DeleteExpiredOAuth2ClientTokens deletes expired oauth2 client tokens.
func (q *repository) DeleteExpiredOAuth2ClientTokens(ctx context.Context) (int64, error) {
	ctx, span := q.tracer.StartSpan(ctx)
	defer span.End()

	logger := q.logger.WithSpan(span)
	logger.Info("Deleting expired oauth2 client tokens")

	deleted, err := q.generatedQuerier.DeleteExpiredOAuth2ClientTokens(ctx, q.db)
	if err != nil {
		return 0, observability.PrepareError(err, span, "deleting expired oauth2 client tokens")
	}

	return deleted, nil
}
