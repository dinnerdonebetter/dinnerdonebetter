package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	preparationsBasePath = "preparations"
)

// BuildGetPreparationRequest builds an HTTP request for fetching a preparation
func (c *V1Client) BuildGetPreparationRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, preparationsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetPreparation retrieves a preparation
func (c *V1Client) GetPreparation(ctx context.Context, id uint64) (preparation *models.Preparation, err error) {
	req, err := c.BuildGetPreparationRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &preparation); retrieveErr != nil {
		return nil, retrieveErr
	}

	return preparation, nil
}

// BuildGetPreparationsRequest builds an HTTP request for fetching preparations
func (c *V1Client) BuildGetPreparationsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), preparationsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetPreparations retrieves a list of preparations
func (c *V1Client) GetPreparations(ctx context.Context, filter *models.QueryFilter) (preparations *models.PreparationList, err error) {
	req, err := c.BuildGetPreparationsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &preparations); retrieveErr != nil {
		return nil, retrieveErr
	}

	return preparations, nil
}

// BuildCreatePreparationRequest builds an HTTP request for creating a preparation
func (c *V1Client) BuildCreatePreparationRequest(ctx context.Context, body *models.PreparationCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, preparationsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreatePreparation creates a preparation
func (c *V1Client) CreatePreparation(ctx context.Context, input *models.PreparationCreationInput) (preparation *models.Preparation, err error) {
	req, err := c.BuildCreatePreparationRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &preparation)
	return preparation, err
}

// BuildUpdatePreparationRequest builds an HTTP request for updating a preparation
func (c *V1Client) BuildUpdatePreparationRequest(ctx context.Context, updated *models.Preparation) (*http.Request, error) {
	uri := c.BuildURL(nil, preparationsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdatePreparation updates a preparation
func (c *V1Client) UpdatePreparation(ctx context.Context, updated *models.Preparation) error {
	req, err := c.BuildUpdatePreparationRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchivePreparationRequest builds an HTTP request for updating a preparation
func (c *V1Client) BuildArchivePreparationRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, preparationsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchivePreparation archives a preparation
func (c *V1Client) ArchivePreparation(ctx context.Context, id uint64) error {
	req, err := c.BuildArchivePreparationRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
