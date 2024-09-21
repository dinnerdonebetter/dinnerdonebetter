package apiclient

import (
	"context"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"
)

// GetRecipeRating gets a recipe rating.
func (c *Client) GetRecipeRating(ctx context.Context, mealID, recipeRatingID string) (*types.RecipeRating, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	if recipeRatingID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	res, err := c.authedGeneratedClient.GetRecipeRating(ctx, mealID, recipeRatingID)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "get recipe rating")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeRating]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe rating")
	}

	return apiResponse.Data, nil
}

// GetRecipeRatings retrieves a list of recipe ratings.
func (c *Client) GetRecipeRatings(ctx context.Context, mealID string, filter *types.QueryFilter) (*types.QueryFilteredResult[types.RecipeRating], error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if filter == nil {
		filter = types.DefaultQueryFilter()
	}
	logger = filter.AttachToLogger(logger)
	tracing.AttachQueryFilterToSpan(span, filter)

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	params := &generated.GetRecipeRatingsParams{}
	c.copyType(params, filter)

	res, err := c.authedGeneratedClient.GetRecipeRatings(ctx, mealID, params)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "recipe ratings list")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[[]*types.RecipeRating]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "retrieving recipe ratings")
	}

	response := &types.QueryFilteredResult[types.RecipeRating]{
		Data:       apiResponse.Data,
		Pagination: *apiResponse.Pagination,
	}

	return response, nil
}

// CreateRecipeRating creates a recipe rating.
func (c *Client) CreateRecipeRating(ctx context.Context, mealID string, input *types.RecipeRatingCreationRequestInput) (*types.RecipeRating, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "validating input")
	}

	body := generated.CreateRecipeRatingJSONRequestBody{}
	c.copyType(&body, input)

	res, err := c.authedGeneratedClient.CreateRecipeRating(ctx, mealID, body)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "create recipe rating")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeRating]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "creating recipe rating")
	}

	return apiResponse.Data, nil
}

// UpdateRecipeRating updates a recipe rating.
func (c *Client) UpdateRecipeRating(ctx context.Context, recipeRating *types.RecipeRating) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if recipeRating == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRating.ID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRating.ID)

	body := generated.UpdateRecipeRatingJSONRequestBody{}
	c.copyType(&body, recipeRating)

	res, err := c.authedGeneratedClient.UpdateRecipeRating(ctx, recipeRating.RecipeID, recipeRating.ID, body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "update recipe rating")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeRating]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "updating recipe rating")
	}

	return nil
}

// ArchiveRecipeRating archives a recipe rating.
func (c *Client) ArchiveRecipeRating(ctx context.Context, mealID, recipeRatingID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.Clone()

	if mealID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.MealIDKey, mealID)
	tracing.AttachToSpan(span, keys.MealIDKey, mealID)

	if recipeRatingID == "" {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.RecipeRatingIDKey, recipeRatingID)
	tracing.AttachToSpan(span, keys.RecipeRatingIDKey, recipeRatingID)

	res, err := c.authedGeneratedClient.ArchiveRecipeRating(ctx, mealID, recipeRatingID)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archive recipe rating")
	}
	defer c.closeResponseBody(ctx, res)

	var apiResponse *types.APIResponse[*types.RecipeRating]
	if err = c.unmarshalBody(ctx, res, &apiResponse); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "archiving recipe rating")
	}

	return nil
}
