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
	requiredPreparationInstrumentsBasePath = "required_preparation_instruments"
)

// BuildRequiredPreparationInstrumentExistsRequest builds an HTTP request for checking the existence of a required preparation instrument.
func (c *V1Client) BuildRequiredPreparationInstrumentExistsRequest(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRequiredPreparationInstrumentExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
		requiredPreparationInstrumentsBasePath,
		strconv.FormatUint(requiredPreparationInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RequiredPreparationInstrumentExists retrieves whether or not a required preparation instrument exists.
func (c *V1Client) RequiredPreparationInstrumentExists(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RequiredPreparationInstrumentExists")
	defer span.End()

	req, err := c.BuildRequiredPreparationInstrumentExistsRequest(ctx, validPreparationID, requiredPreparationInstrumentID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRequiredPreparationInstrumentRequest builds an HTTP request for fetching a required preparation instrument.
func (c *V1Client) BuildGetRequiredPreparationInstrumentRequest(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRequiredPreparationInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
		requiredPreparationInstrumentsBasePath,
		strconv.FormatUint(requiredPreparationInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRequiredPreparationInstrument retrieves a required preparation instrument.
func (c *V1Client) GetRequiredPreparationInstrument(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (requiredPreparationInstrument *models.RequiredPreparationInstrument, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRequiredPreparationInstrument")
	defer span.End()

	req, err := c.BuildGetRequiredPreparationInstrumentRequest(ctx, validPreparationID, requiredPreparationInstrumentID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &requiredPreparationInstrument); retrieveErr != nil {
		return nil, retrieveErr
	}

	return requiredPreparationInstrument, nil
}

// BuildGetRequiredPreparationInstrumentsRequest builds an HTTP request for fetching required preparation instruments.
func (c *V1Client) BuildGetRequiredPreparationInstrumentsRequest(ctx context.Context, validPreparationID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRequiredPreparationInstrumentsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
		requiredPreparationInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRequiredPreparationInstruments retrieves a list of required preparation instruments.
func (c *V1Client) GetRequiredPreparationInstruments(ctx context.Context, validPreparationID uint64, filter *models.QueryFilter) (requiredPreparationInstruments *models.RequiredPreparationInstrumentList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRequiredPreparationInstruments")
	defer span.End()

	req, err := c.BuildGetRequiredPreparationInstrumentsRequest(ctx, validPreparationID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &requiredPreparationInstruments); retrieveErr != nil {
		return nil, retrieveErr
	}

	return requiredPreparationInstruments, nil
}

// BuildCreateRequiredPreparationInstrumentRequest builds an HTTP request for creating a required preparation instrument.
func (c *V1Client) BuildCreateRequiredPreparationInstrumentRequest(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRequiredPreparationInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(input.BelongsToValidPreparation, 10),
		requiredPreparationInstrumentsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRequiredPreparationInstrument creates a required preparation instrument.
func (c *V1Client) CreateRequiredPreparationInstrument(ctx context.Context, input *models.RequiredPreparationInstrumentCreationInput) (requiredPreparationInstrument *models.RequiredPreparationInstrument, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRequiredPreparationInstrument")
	defer span.End()

	req, err := c.BuildCreateRequiredPreparationInstrumentRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &requiredPreparationInstrument)
	return requiredPreparationInstrument, err
}

// BuildUpdateRequiredPreparationInstrumentRequest builds an HTTP request for updating a required preparation instrument.
func (c *V1Client) BuildUpdateRequiredPreparationInstrumentRequest(ctx context.Context, requiredPreparationInstrument *models.RequiredPreparationInstrument) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRequiredPreparationInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(requiredPreparationInstrument.BelongsToValidPreparation, 10),
		requiredPreparationInstrumentsBasePath,
		strconv.FormatUint(requiredPreparationInstrument.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, requiredPreparationInstrument)
}

// UpdateRequiredPreparationInstrument updates a required preparation instrument.
func (c *V1Client) UpdateRequiredPreparationInstrument(ctx context.Context, requiredPreparationInstrument *models.RequiredPreparationInstrument) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRequiredPreparationInstrument")
	defer span.End()

	req, err := c.BuildUpdateRequiredPreparationInstrumentRequest(ctx, requiredPreparationInstrument)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &requiredPreparationInstrument)
}

// BuildArchiveRequiredPreparationInstrumentRequest builds an HTTP request for updating a required preparation instrument.
func (c *V1Client) BuildArchiveRequiredPreparationInstrumentRequest(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRequiredPreparationInstrumentRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validPreparationsBasePath,
		strconv.FormatUint(validPreparationID, 10),
		requiredPreparationInstrumentsBasePath,
		strconv.FormatUint(requiredPreparationInstrumentID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRequiredPreparationInstrument archives a required preparation instrument.
func (c *V1Client) ArchiveRequiredPreparationInstrument(ctx context.Context, validPreparationID, requiredPreparationInstrumentID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRequiredPreparationInstrument")
	defer span.End()

	req, err := c.BuildArchiveRequiredPreparationInstrumentRequest(ctx, validPreparationID, requiredPreparationInstrumentID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
