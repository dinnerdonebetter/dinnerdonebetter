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
	validIngredientTagsBasePath = "valid_ingredient_tags"
)

// BuildValidIngredientTagExistsRequest builds an HTTP request for checking the existence of a valid ingredient tag.
func (c *V1Client) BuildValidIngredientTagExistsRequest(ctx context.Context, validIngredientTagID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildValidIngredientTagExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientTagsBasePath,
		strconv.FormatUint(validIngredientTagID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ValidIngredientTagExists retrieves whether or not a valid ingredient tag exists.
func (c *V1Client) ValidIngredientTagExists(ctx context.Context, validIngredientTagID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ValidIngredientTagExists")
	defer span.End()

	req, err := c.BuildValidIngredientTagExistsRequest(ctx, validIngredientTagID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetValidIngredientTagRequest builds an HTTP request for fetching a valid ingredient tag.
func (c *V1Client) BuildGetValidIngredientTagRequest(ctx context.Context, validIngredientTagID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidIngredientTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientTagsBasePath,
		strconv.FormatUint(validIngredientTagID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidIngredientTag retrieves a valid ingredient tag.
func (c *V1Client) GetValidIngredientTag(ctx context.Context, validIngredientTagID uint64) (validIngredientTag *models.ValidIngredientTag, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientTag")
	defer span.End()

	req, err := c.BuildGetValidIngredientTagRequest(ctx, validIngredientTagID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredientTag); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredientTag, nil
}

// BuildGetValidIngredientTagsRequest builds an HTTP request for fetching valid ingredient tags.
func (c *V1Client) BuildGetValidIngredientTagsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetValidIngredientTagsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validIngredientTagsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetValidIngredientTags retrieves a list of valid ingredient tags.
func (c *V1Client) GetValidIngredientTags(ctx context.Context, filter *models.QueryFilter) (validIngredientTags *models.ValidIngredientTagList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetValidIngredientTags")
	defer span.End()

	req, err := c.BuildGetValidIngredientTagsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &validIngredientTags); retrieveErr != nil {
		return nil, retrieveErr
	}

	return validIngredientTags, nil
}

// BuildCreateValidIngredientTagRequest builds an HTTP request for creating a valid ingredient tag.
func (c *V1Client) BuildCreateValidIngredientTagRequest(ctx context.Context, input *models.ValidIngredientTagCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateValidIngredientTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientTagsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateValidIngredientTag creates a valid ingredient tag.
func (c *V1Client) CreateValidIngredientTag(ctx context.Context, input *models.ValidIngredientTagCreationInput) (validIngredientTag *models.ValidIngredientTag, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateValidIngredientTag")
	defer span.End()

	req, err := c.BuildCreateValidIngredientTagRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &validIngredientTag)
	return validIngredientTag, err
}

// BuildUpdateValidIngredientTagRequest builds an HTTP request for updating a valid ingredient tag.
func (c *V1Client) BuildUpdateValidIngredientTagRequest(ctx context.Context, validIngredientTag *models.ValidIngredientTag) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateValidIngredientTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientTagsBasePath,
		strconv.FormatUint(validIngredientTag.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, validIngredientTag)
}

// UpdateValidIngredientTag updates a valid ingredient tag.
func (c *V1Client) UpdateValidIngredientTag(ctx context.Context, validIngredientTag *models.ValidIngredientTag) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateValidIngredientTag")
	defer span.End()

	req, err := c.BuildUpdateValidIngredientTagRequest(ctx, validIngredientTag)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &validIngredientTag)
}

// BuildArchiveValidIngredientTagRequest builds an HTTP request for updating a valid ingredient tag.
func (c *V1Client) BuildArchiveValidIngredientTagRequest(ctx context.Context, validIngredientTagID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveValidIngredientTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientTagsBasePath,
		strconv.FormatUint(validIngredientTagID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveValidIngredientTag archives a valid ingredient tag.
func (c *V1Client) ArchiveValidIngredientTag(ctx context.Context, validIngredientTagID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveValidIngredientTag")
	defer span.End()

	req, err := c.BuildArchiveValidIngredientTagRequest(ctx, validIngredientTagID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
