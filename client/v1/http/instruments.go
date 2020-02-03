package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	instrumentsBasePath = "instruments"
)

// BuildGetInstrumentRequest builds an HTTP request for fetching an instrument
func (c *V1Client) BuildGetInstrumentRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, instrumentsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetInstrument retrieves an instrument
func (c *V1Client) GetInstrument(ctx context.Context, id uint64) (instrument *models.Instrument, err error) {
	req, err := c.BuildGetInstrumentRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &instrument); retrieveErr != nil {
		return nil, retrieveErr
	}

	return instrument, nil
}

// BuildGetInstrumentsRequest builds an HTTP request for fetching instruments
func (c *V1Client) BuildGetInstrumentsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), instrumentsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetInstruments retrieves a list of instruments
func (c *V1Client) GetInstruments(ctx context.Context, filter *models.QueryFilter) (instruments *models.InstrumentList, err error) {
	req, err := c.BuildGetInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &instruments); retrieveErr != nil {
		return nil, retrieveErr
	}

	return instruments, nil
}

// BuildCreateInstrumentRequest builds an HTTP request for creating an instrument
func (c *V1Client) BuildCreateInstrumentRequest(ctx context.Context, body *models.InstrumentCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, instrumentsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateInstrument creates an instrument
func (c *V1Client) CreateInstrument(ctx context.Context, input *models.InstrumentCreationInput) (instrument *models.Instrument, err error) {
	req, err := c.BuildCreateInstrumentRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &instrument)
	return instrument, err
}

// BuildUpdateInstrumentRequest builds an HTTP request for updating an instrument
func (c *V1Client) BuildUpdateInstrumentRequest(ctx context.Context, updated *models.Instrument) (*http.Request, error) {
	uri := c.BuildURL(nil, instrumentsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateInstrument updates an instrument
func (c *V1Client) UpdateInstrument(ctx context.Context, updated *models.Instrument) error {
	req, err := c.BuildUpdateInstrumentRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveInstrumentRequest builds an HTTP request for updating an instrument
func (c *V1Client) BuildArchiveInstrumentRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, instrumentsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveInstrument archives an instrument
func (c *V1Client) ArchiveInstrument(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveInstrumentRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
