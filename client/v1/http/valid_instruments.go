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
	validInstrumentsBasePath = "valid_instruments"
)

// BuildValidInstrumentExistsRequest builds an HTTP request for checking the existence of a valid instrument.
func (c *V1Client) BuildValidInstrumentExistsRequest(ctx context.Context, validInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildValidInstrumentExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validInstrumentsBasePath,
		strconv.FormatUint(validInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ValidInstrumentExists retrieves whether or not a valid instrument exists.
func (c *V1Client) ValidInstrumentExists(ctx context.Context, validInstrumentID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidInstrumentExists")
	defer span.End()

	req, err := c.BuildValidInstrumentExistsRequest(ctx, validInstrumentID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetValidInstrumentRequest builds an HTTP request for fetching a valid instrument.
func (c *V1Client) BuildGetValidInstrumentRequest(ctx context.Context, validInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validInstrumentsBasePath,
		strconv.FormatUint(validInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidInstrument retrieves a valid instrument.
func (c *V1Client) GetValidInstrument(ctx context.Context, validInstrumentID uint64) (validInstrument *models.ValidInstrument, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidInstrument")
	defer span.End()

	req, err := c.BuildGetValidInstrumentRequest(ctx, validInstrumentID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validInstrument); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validInstrument, nil
}

// BuildGetValidInstrumentsRequest builds an HTTP request for fetching valid instruments.
func (c *V1Client) BuildGetValidInstrumentsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidInstrumentsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidInstruments retrieves a list of valid instruments.
func (c *V1Client) GetValidInstruments(ctx context.Context, filter *models.QueryFilter) (validInstruments *models.ValidInstrumentList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidInstruments")
	defer span.End()

	req, err := c.BuildGetValidInstrumentsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validInstruments); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validInstruments, nil
}

// BuildCreateValidInstrumentRequest builds an HTTP request for creating a valid instrument.
func (c *V1Client) BuildCreateValidInstrumentRequest(ctx context.Context, input *models.ValidInstrumentCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateValidInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateValidInstrument creates a valid instrument.
func (c *V1Client) CreateValidInstrument(ctx context.Context, input *models.ValidInstrumentCreationInput) (validInstrument *models.ValidInstrument, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidInstrument")
	defer span.End()

	req, err := c.BuildCreateValidInstrumentRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &validInstrument)
	return validInstrument, err
}

// BuildUpdateValidInstrumentRequest builds an HTTP request for updating a valid instrument.
func (c *V1Client) BuildUpdateValidInstrumentRequest(ctx context.Context, validInstrument *models.ValidInstrument) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateValidInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validInstrumentsBasePath,
		strconv.FormatUint(validInstrument.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, validInstrument)
}

// UpdateValidInstrument updates a valid instrument.
func (c *V1Client) UpdateValidInstrument(ctx context.Context, validInstrument *models.ValidInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidInstrument")
	defer span.End()

	req, err := c.BuildUpdateValidInstrumentRequest(ctx, validInstrument)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &validInstrument)
}

// BuildArchiveValidInstrumentRequest builds an HTTP request for updating a valid instrument.
func (c *V1Client) BuildArchiveValidInstrumentRequest(ctx context.Context, validInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveValidInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validInstrumentsBasePath,
		strconv.FormatUint(validInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveValidInstrument archives a valid instrument.
func (c *V1Client) ArchiveValidInstrument(ctx context.Context, validInstrumentID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidInstrument")
	defer span.End()

	req, err := c.BuildArchiveValidInstrumentRequest(ctx, validInstrumentID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
