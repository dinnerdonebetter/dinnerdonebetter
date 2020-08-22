package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	validIngredientsBasePath = "valid_ingredients"
)

// BuildValidIngredientExistsRequest builds an HTTP request for checking the existence of a valid ingredient.
func (c *V1Client) BuildValidIngredientExistsRequest(ctx context.Context, validIngredientID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildValidIngredientExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ValidIngredientExists retrieves whether or not a valid ingredient exists.
func (c *V1Client) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientExists")
	defer span.End()

	req, err := c.BuildValidIngredientExistsRequest(ctx, validIngredientID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetValidIngredientRequest builds an HTTP request for fetching a valid ingredient.
func (c *V1Client) BuildGetValidIngredientRequest(ctx context.Context, validIngredientID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidIngredient retrieves a valid ingredient.
func (c *V1Client) GetValidIngredient(ctx context.Context, validIngredientID uint64) (validIngredient *models.ValidIngredient, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredient")
	defer span.End()

	req, err := c.BuildGetValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredient); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredient, nil
}

// BuildSearchValidIngredientsRequest builds an HTTP request for querying valid ingredients.
func (c *V1Client) BuildSearchValidIngredientsRequest(ctx context.Context, query string, limit uint8) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildSearchValidIngredientsRequest")
	defer span.End()

	params := url.Values{}
	params.Set(models.SearchQueryKey, query)
	params.Set(models.LimitQueryKey, strconv.FormatUint(uint64(limit), 10))

	uri := c.BuildURL(
		params,
		validIngredientsBasePath,
		"search",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// SearchValidIngredients searches for a list of valid ingredients.
func (c *V1Client) SearchValidIngredients(ctx context.Context, query string, limit uint8) (validIngredients []models.ValidIngredient, err error) {
	ctx, span := tracing.StartSpan(ctx, "SearchValidIngredients")
	defer span.End()

	req, err := c.BuildSearchValidIngredientsRequest(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredients); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredients, nil
}

// BuildGetValidIngredientsRequest builds an HTTP request for fetching valid ingredients.
func (c *V1Client) BuildGetValidIngredientsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidIngredientsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidIngredients retrieves a list of valid ingredients.
func (c *V1Client) GetValidIngredients(ctx context.Context, filter *models.QueryFilter) (validIngredients *models.ValidIngredientList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredients")
	defer span.End()

	req, err := c.BuildGetValidIngredientsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredients); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredients, nil
}

// BuildCreateValidIngredientRequest builds an HTTP request for creating a valid ingredient.
func (c *V1Client) BuildCreateValidIngredientRequest(ctx context.Context, input *models.ValidIngredientCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateValidIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateValidIngredient creates a valid ingredient.
func (c *V1Client) CreateValidIngredient(ctx context.Context, input *models.ValidIngredientCreationInput) (validIngredient *models.ValidIngredient, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidIngredient")
	defer span.End()

	req, err := c.BuildCreateValidIngredientRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &validIngredient)
	return validIngredient, err
}

// BuildUpdateValidIngredientRequest builds an HTTP request for updating a valid ingredient.
func (c *V1Client) BuildUpdateValidIngredientRequest(ctx context.Context, validIngredient *models.ValidIngredient) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateValidIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredient.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, validIngredient)
}

// UpdateValidIngredient updates a valid ingredient.
func (c *V1Client) UpdateValidIngredient(ctx context.Context, validIngredient *models.ValidIngredient) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidIngredient")
	defer span.End()

	req, err := c.BuildUpdateValidIngredientRequest(ctx, validIngredient)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &validIngredient)
}

// BuildArchiveValidIngredientRequest builds an HTTP request for updating a valid ingredient.
func (c *V1Client) BuildArchiveValidIngredientRequest(ctx context.Context, validIngredientID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveValidIngredientRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveValidIngredient archives a valid ingredient.
func (c *V1Client) ArchiveValidIngredient(ctx context.Context, validIngredientID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredient")
	defer span.End()

	req, err := c.BuildArchiveValidIngredientRequest(ctx, validIngredientID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
