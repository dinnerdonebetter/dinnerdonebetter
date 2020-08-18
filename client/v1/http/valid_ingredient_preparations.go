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
	validIngredientPreparationsBasePath = "valid_ingredient_preparations"
)

// BuildValidIngredientPreparationExistsRequest builds an HTTP request for checking the existence of a valid ingredient preparation.
func (c *V1Client) BuildValidIngredientPreparationExistsRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildValidIngredientPreparationExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientPreparationsBasePath,
		strconv.FormatUint(validIngredientPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ValidIngredientPreparationExists retrieves whether or not a valid ingredient preparation exists.
func (c *V1Client) ValidIngredientPreparationExists(ctx context.Context, validIngredientPreparationID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientPreparationExists")
	defer span.End()

	req, err := c.BuildValidIngredientPreparationExistsRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetValidIngredientPreparationRequest builds an HTTP request for fetching a valid ingredient preparation.
func (c *V1Client) BuildGetValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidIngredientPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientPreparationsBasePath,
		strconv.FormatUint(validIngredientPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidIngredientPreparation retrieves a valid ingredient preparation.
func (c *V1Client) GetValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) (validIngredientPreparation *models.ValidIngredientPreparation, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparation")
	defer span.End()

	req, err := c.BuildGetValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredientPreparation); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredientPreparation, nil
}

// BuildGetValidIngredientPreparationsRequest builds an HTTP request for fetching valid ingredient preparations.
func (c *V1Client) BuildGetValidIngredientPreparationsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidIngredientPreparationsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validIngredientPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidIngredientPreparations retrieves a list of valid ingredient preparations.
func (c *V1Client) GetValidIngredientPreparations(ctx context.Context, filter *models.QueryFilter) (validIngredientPreparations *models.ValidIngredientPreparationList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientPreparations")
	defer span.End()

	req, err := c.BuildGetValidIngredientPreparationsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredientPreparations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredientPreparations, nil
}

// BuildCreateValidIngredientPreparationRequest builds an HTTP request for creating a valid ingredient preparation.
func (c *V1Client) BuildCreateValidIngredientPreparationRequest(ctx context.Context, input *models.ValidIngredientPreparationCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateValidIngredientPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateValidIngredientPreparation creates a valid ingredient preparation.
func (c *V1Client) CreateValidIngredientPreparation(ctx context.Context, input *models.ValidIngredientPreparationCreationInput) (validIngredientPreparation *models.ValidIngredientPreparation, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidIngredientPreparation")
	defer span.End()

	req, err := c.BuildCreateValidIngredientPreparationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &validIngredientPreparation)
	return validIngredientPreparation, err
}

// BuildUpdateValidIngredientPreparationRequest builds an HTTP request for updating a valid ingredient preparation.
func (c *V1Client) BuildUpdateValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparation *models.ValidIngredientPreparation) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateValidIngredientPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientPreparationsBasePath,
		strconv.FormatUint(validIngredientPreparation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, validIngredientPreparation)
}

// UpdateValidIngredientPreparation updates a valid ingredient preparation.
func (c *V1Client) UpdateValidIngredientPreparation(ctx context.Context, validIngredientPreparation *models.ValidIngredientPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidIngredientPreparation")
	defer span.End()

	req, err := c.BuildUpdateValidIngredientPreparationRequest(ctx, validIngredientPreparation)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &validIngredientPreparation)
}

// BuildArchiveValidIngredientPreparationRequest builds an HTTP request for updating a valid ingredient preparation.
func (c *V1Client) BuildArchiveValidIngredientPreparationRequest(ctx context.Context, validIngredientPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveValidIngredientPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientPreparationsBasePath,
		strconv.FormatUint(validIngredientPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveValidIngredientPreparation archives a valid ingredient preparation.
func (c *V1Client) ArchiveValidIngredientPreparation(ctx context.Context, validIngredientPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredientPreparation")
	defer span.End()

	req, err := c.BuildArchiveValidIngredientPreparationRequest(ctx, validIngredientPreparationID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
