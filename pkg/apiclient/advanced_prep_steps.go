package apiclient

import (
	"context"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

// GetAdvancedPrepStep gets an advanced prep step.
func (c *Client) GetAdvancedPrepStep(ctx context.Context, validIngredientID string) (*types.AdvancedPrepStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if validIngredientID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.AdvancedPrepStepIDKey, validIngredientID)
	tracing.AttachAdvancedPrepStepIDToSpan(span, validIngredientID)

	req, err := c.requestBuilder.BuildGetAdvancedPrepStepRequest(ctx, validIngredientID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building get advanced prep step request")
	}

	var validIngredient *types.AdvancedPrepStep
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving advanced prep step")
	}

	return validIngredient, nil
}

// UpdateAdvancedPrepStepStatus updates an advanced prep step.
func (c *Client) UpdateAdvancedPrepStepStatus(ctx context.Context, input *types.AdvancedPrepStepStatusChangeRequestInput) (*types.AdvancedPrepStep, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildChangeAdvancedPrepStepStatusRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "building create advanced prep step request")
	}

	var validIngredient *types.AdvancedPrepStep
	if err = c.fetchAndUnmarshal(ctx, req, &validIngredient); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating advanced prep step")
	}

	return validIngredient, nil
}
