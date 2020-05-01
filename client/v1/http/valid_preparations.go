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
	validPreparationsBasePath = "valid_preparations"
)

// BuildValidPreparationExistsRequest builds an HTTP request for checking the existence of a valid preparation.
func (c *V1Client) BuildValidPreparationExistsRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildValidPreparationExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ValidPreparationExists retrieves whether or not a valid preparation exists.
func (c *V1Client) ValidPreparationExists(ctx context.Context, validPreparationID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidPreparationExists")
	defer span.End()

	req, err := c.BuildValidPreparationExistsRequest(ctx, validPreparationID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetValidPreparationRequest builds an HTTP request for fetching a valid preparation.
func (c *V1Client) BuildGetValidPreparationRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidPreparation retrieves a valid preparation.
func (c *V1Client) GetValidPreparation(ctx context.Context, validPreparationID uint64) (validPreparation *models.ValidPreparation, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidPreparation")
	defer span.End()

	req, err := c.BuildGetValidPreparationRequest(ctx, validPreparationID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validPreparation); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validPreparation, nil
}

// BuildGetValidPreparationsRequest builds an HTTP request for fetching valid preparations.
func (c *V1Client) BuildGetValidPreparationsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidPreparationsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidPreparations retrieves a list of valid preparations.
func (c *V1Client) GetValidPreparations(ctx context.Context, filter *models.QueryFilter) (validPreparations *models.ValidPreparationList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidPreparations")
	defer span.End()

	req, err := c.BuildGetValidPreparationsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validPreparations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validPreparations, nil
}

// BuildCreateValidPreparationRequest builds an HTTP request for creating a valid preparation.
func (c *V1Client) BuildCreateValidPreparationRequest(ctx context.Context, input *models.ValidPreparationCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateValidPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateValidPreparation creates a valid preparation.
func (c *V1Client) CreateValidPreparation(ctx context.Context, input *models.ValidPreparationCreationInput) (validPreparation *models.ValidPreparation, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidPreparation")
	defer span.End()

	req, err := c.BuildCreateValidPreparationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &validPreparation)
	return validPreparation, err
}

// BuildUpdateValidPreparationRequest builds an HTTP request for updating a valid preparation.
func (c *V1Client) BuildUpdateValidPreparationRequest(ctx context.Context, validPreparation *models.ValidPreparation) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateValidPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, validPreparation)
}

// UpdateValidPreparation updates a valid preparation.
func (c *V1Client) UpdateValidPreparation(ctx context.Context, validPreparation *models.ValidPreparation) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidPreparation")
	defer span.End()

	req, err := c.BuildUpdateValidPreparationRequest(ctx, validPreparation)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &validPreparation)
}

// BuildArchiveValidPreparationRequest builds an HTTP request for updating a valid preparation.
func (c *V1Client) BuildArchiveValidPreparationRequest(ctx context.Context, validPreparationID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveValidPreparationRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveValidPreparation archives a valid preparation.
func (c *V1Client) ArchiveValidPreparation(ctx context.Context, validPreparationID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidPreparation")
	defer span.End()

	req, err := c.BuildArchiveValidPreparationRequest(ctx, validPreparationID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
