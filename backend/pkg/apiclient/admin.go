package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/jinzhu/copier"
)

// UpdateUserAccountStatus updates a user's account status.
func (c *Client) UpdateUserAccountStatus(ctx context.Context, input *types.UserAccountStatusUpdateInput) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return ErrNilInputProvided
	}

	tracing.AttachToSpan(span, keys.UserIDKey, input.TargetUserID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return observability.PrepareError(err, span, "validating input")
	}

	var body generated.AdminUpdateUserStatusJSONRequestBody
	if err := copier.Copy(&body, input); err != nil {
		return observability.PrepareError(err, span, "copying input")
	}
	res, err := c.authedGeneratedClient.AdminUpdateUserStatus(ctx, body)
	if err != nil {
		return observability.PrepareError(err, span, "updating user account status")
	}

	c.closeResponseBody(ctx, res)

	if err = errorFromResponse(res); err != nil {
		return observability.PrepareError(err, span, "invalid response status")
	}

	return nil
}
