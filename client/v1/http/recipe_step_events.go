package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	recipeStepEventsBasePath = "recipe_step_events"
)

// BuildRecipeStepEventExistsRequest builds an HTTP request for checking the existence of a recipe step event.
func (c *V1Client) BuildRecipeStepEventExistsRequest(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeStepEventExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepEventsBasePath,
		strconv.FormatUint(recipeStepEventID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeStepEventExists retrieves whether or not a recipe step event exists.
func (c *V1Client) RecipeStepEventExists(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeStepEventExists")
	defer span.End()

	req, err := c.BuildRecipeStepEventExistsRequest(ctx, recipeID, recipeStepID, recipeStepEventID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeStepEventRequest builds an HTTP request for fetching a recipe step event.
func (c *V1Client) BuildGetRecipeStepEventRequest(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepEventRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepEventsBasePath,
		strconv.FormatUint(recipeStepEventID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepEvent retrieves a recipe step event.
func (c *V1Client) GetRecipeStepEvent(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (recipeStepEvent *models.RecipeStepEvent, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepEvent")
	defer span.End()

	req, err := c.BuildGetRecipeStepEventRequest(ctx, recipeID, recipeStepID, recipeStepEventID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepEvent); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepEvent, nil
}

// BuildGetRecipeStepEventsRequest builds an HTTP request for fetching recipe step events.
func (c *V1Client) BuildGetRecipeStepEventsRequest(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeStepEventsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepEventsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeStepEvents retrieves a list of recipe step events.
func (c *V1Client) GetRecipeStepEvents(ctx context.Context, recipeID, recipeStepID uint64, filter *models.QueryFilter) (recipeStepEvents *models.RecipeStepEventList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeStepEvents")
	defer span.End()

	req, err := c.BuildGetRecipeStepEventsRequest(ctx, recipeID, recipeStepID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeStepEvents); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeStepEvents, nil
}

// BuildCreateRecipeStepEventRequest builds an HTTP request for creating a recipe step event.
func (c *V1Client) BuildCreateRecipeStepEventRequest(ctx context.Context, recipeID uint64, input *models.RecipeStepEventCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeStepEventRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(input.BelongsToRecipeStep, 10),
		recipeStepEventsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeStepEvent creates a recipe step event.
func (c *V1Client) CreateRecipeStepEvent(ctx context.Context, recipeID uint64, input *models.RecipeStepEventCreationInput) (recipeStepEvent *models.RecipeStepEvent, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeStepEvent")
	defer span.End()

	req, err := c.BuildCreateRecipeStepEventRequest(ctx, recipeID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeStepEvent)
	return recipeStepEvent, err
}

// BuildUpdateRecipeStepEventRequest builds an HTTP request for updating a recipe step event.
func (c *V1Client) BuildUpdateRecipeStepEventRequest(ctx context.Context, recipeID uint64, recipeStepEvent *models.RecipeStepEvent) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeStepEventRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepEvent.BelongsToRecipeStep, 10),
		recipeStepEventsBasePath,
		strconv.FormatUint(recipeStepEvent.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeStepEvent)
}

// UpdateRecipeStepEvent updates a recipe step event.
func (c *V1Client) UpdateRecipeStepEvent(ctx context.Context, recipeID uint64, recipeStepEvent *models.RecipeStepEvent) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeStepEvent")
	defer span.End()

	req, err := c.BuildUpdateRecipeStepEventRequest(ctx, recipeID, recipeStepEvent)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeStepEvent)
}

// BuildArchiveRecipeStepEventRequest builds an HTTP request for updating a recipe step event.
func (c *V1Client) BuildArchiveRecipeStepEventRequest(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeStepEventRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeStepsBasePath,
		strconv.FormatUint(recipeStepID, 10),
		recipeStepEventsBasePath,
		strconv.FormatUint(recipeStepEventID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeStepEvent archives a recipe step event.
func (c *V1Client) ArchiveRecipeStepEvent(ctx context.Context, recipeID, recipeStepID, recipeStepEventID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeStepEvent")
	defer span.End()

	req, err := c.BuildArchiveRecipeStepEventRequest(ctx, recipeID, recipeStepID, recipeStepEventID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
