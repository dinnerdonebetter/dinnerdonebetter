package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/jinzhu/copier"
)

// RunFinalizeMealPlansWorker runs a worker.
func (c *Client) RunFinalizeMealPlansWorker(ctx context.Context, input *types.FinalizeMealPlansRequest) (*types.FinalizeMealPlansResponse, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	var body generated.RunFinalizeMealPlanWorkerJSONRequestBody
	if err := copier.Copy(&body, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input")
	}

	res, err := c.authedGeneratedClient.RunFinalizeMealPlanWorker(ctx, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "finalizing meal plan")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.FinalizeMealPlansResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "loading finalize meal plan response")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return nil, err
	}

	return apiResponse.Data, nil
}

// RunMealPlanGroceryListInitializationWorker runs a worker.
func (c *Client) RunMealPlanGroceryListInitializationWorker(ctx context.Context) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	body := generated.RunMealPlanGroceryListInitializerWorkerJSONRequestBody{}
	res, err := c.authedGeneratedClient.RunMealPlanGroceryListInitializerWorker(ctx, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "running meal plan grocery list initializer worker")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.FinalizeMealPlansResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "initializing meal plan grocery lists")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}

// RunMealPlanTaskCreationWorker runs a worker.
func (c *Client) RunMealPlanTaskCreationWorker(ctx context.Context) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	body := generated.RunMealPlanTaskCreatorWorkerJSONRequestBody{}
	res, err := c.authedGeneratedClient.RunMealPlanTaskCreatorWorker(ctx, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "running meal plan task creation worker")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.FinalizeMealPlansResponse]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "creating meal plan tasks")
	}

	if err = apiResponse.Error.AsError(); err != nil {
		return err
	}

	return nil
}
