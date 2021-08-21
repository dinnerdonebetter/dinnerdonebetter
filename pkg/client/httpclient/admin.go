package httpclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/observability"
	"gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// UpdateUserReputation updates a user's reputation.
func (c *Client) UpdateUserReputation(ctx context.Context, input *types.UserReputationUpdateInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrNilInputProvided
	}

	logger := c.logger.WithValue(keys.HouseholdIDKey, input.TargetUserID)
	tracing.AttachHouseholdIDToSpan(span, input.TargetUserID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildUserReputationUpdateInputRequest(ctx, input)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building user reputation update request")
	}

	res, err := c.fetchResponseToRequest(ctx, c.authedClient, req)
	if err != nil {
		return observability.PrepareError(err, logger, span, "updating user reputation")
	}

	c.closeResponseBody(ctx, res)

	if err = errorFromResponse(res); err != nil {
		return observability.PrepareError(err, logger, span, "invalid response status")
	}

	return nil
}
