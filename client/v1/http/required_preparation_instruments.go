package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	requiredPreparationInstrumentsBasePath = "required_preparation_instruments"
)

// BuildGetRequiredPreparationInstrumentRequest builds an HTTP request for fetching a required preparation instrument
func (c *V1Client) BuildGetRequiredPreparationInstrumentRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, requiredPreparationInstrumentsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRequiredPreparationInstrument retrieves a required preparation instrument
func (c *V1Client) GetRequiredPreparationInstrument(ctx context.Context, id uint64) (requiredPreparationInstrument *models.RequiredPreparationInstrument, err error) {
	req, err := c.BuildGetRequiredPreparationInstrumentRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &requiredPreparationInstrument); retrieveErr != nil {
		return nil, retrieveErr
	}

	return requiredPreparationInstrument, nil
}

// BuildGetRequiredPreparationInstrumentsRequest builds an HTTP request for fetching required preparation instruments
func (c *V1Client) BuildGetRequiredPreparationInstrumentsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), requiredPreparationInstrumentsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetRequiredPreparationInstruments retrieves a list of required preparation instruments
func (c *V1Client) GetRequiredPreparationInstruments(ctx context.Context, filter *models.QueryFilter) (requiredPreparationInstruments *models.RequiredPreparationInstrumentList, err error) {
	req, err := c.BuildGetRequiredPreparationInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &requiredPreparationInstruments); retrieveErr != nil {
		return nil, retrieveErr
	}

	return requiredPreparationInstruments, nil
}

// BuildCreateRequiredPreparationInstrumentRequest builds an HTTP request for creating a required preparation instrument
func (c *V1Client) BuildCreateRequiredPreparationInstrumentRequest(ctx context.Context, body *models.RequiredPreparationInstrumentCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, requiredPreparationInstrumentsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateRequiredPreparationInstrument creates a required preparation instrument
func (c *V1Client) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (requiredPreparationInstrument *models.RequiredPreparationInstrument, err error) {
	req, err := c.BuildCreateRequiredPreparationInstrumentRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &requiredPreparationInstrument)
	return requiredPreparationInstrument, err
}

// BuildUpdateRequiredPreparationInstrumentRequest builds an HTTP request for updating a required preparation instrument
func (c *V1Client) BuildUpdateRequiredPreparationInstrumentRequest(ctx context.Context, updated *models.RequiredPreparationInstrument) (*http.Request, error) {
	uri := c.BuildURL(nil, requiredPreparationInstrumentsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateRequiredPreparationInstrument updates a required preparation instrument
func (c *V1Client) UpdateRequiredPreparationInstrument(ctx context.Context, updated *models.RequiredPreparationInstrument) error {
	req, err := c.BuildUpdateRequiredPreparationInstrumentRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveRequiredPreparationInstrumentRequest builds an HTTP request for updating a required preparation instrument
func (c *V1Client) BuildArchiveRequiredPreparationInstrumentRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, requiredPreparationInstrumentsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveRequiredPreparationInstrument archives a required preparation instrument
func (c *V1Client) ArchiveRequiredPreparationInstrument(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveRequiredPreparationInstrumentRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
