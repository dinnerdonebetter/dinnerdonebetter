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
	ingredientTagMappingsBasePath = "ingredient_tag_mappings"
)

// BuildIngredientTagMappingExistsRequest builds an HTTP request for checking the existence of an ingredient tag mapping.
func (c *V1Client) BuildIngredientTagMappingExistsRequest(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildIngredientTagMappingExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
		ingredientTagMappingsBasePath,
		strconv.FormatUint(ingredientTagMappingID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// IngredientTagMappingExists retrieves whether or not an ingredient tag mapping exists.
func (c *V1Client) IngredientTagMappingExists(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "IngredientTagMappingExists")
	defer span.End()

	req, err := c.BuildIngredientTagMappingExistsRequest(ctx, validIngredientID, ingredientTagMappingID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetIngredientTagMappingRequest builds an HTTP request for fetching an ingredient tag mapping.
func (c *V1Client) BuildGetIngredientTagMappingRequest(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetIngredientTagMappingRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
		ingredientTagMappingsBasePath,
		strconv.FormatUint(ingredientTagMappingID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetIngredientTagMapping retrieves an ingredient tag mapping.
func (c *V1Client) GetIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (ingredientTagMapping *models.IngredientTagMapping, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetIngredientTagMapping")
	defer span.End()

	req, err := c.BuildGetIngredientTagMappingRequest(ctx, validIngredientID, ingredientTagMappingID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &ingredientTagMapping); retrieveErr != nil {
		return nil, retrieveErr
	}

	return ingredientTagMapping, nil
}

// BuildGetIngredientTagMappingsRequest builds an HTTP request for fetching ingredient tag mappings.
func (c *V1Client) BuildGetIngredientTagMappingsRequest(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetIngredientTagMappingsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
		ingredientTagMappingsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetIngredientTagMappings retrieves a list of ingredient tag mappings.
func (c *V1Client) GetIngredientTagMappings(ctx context.Context, validIngredientID uint64, filter *models.QueryFilter) (ingredientTagMappings *models.IngredientTagMappingList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetIngredientTagMappings")
	defer span.End()

	req, err := c.BuildGetIngredientTagMappingsRequest(ctx, validIngredientID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &ingredientTagMappings); retrieveErr != nil {
		return nil, retrieveErr
	}

	return ingredientTagMappings, nil
}

// BuildCreateIngredientTagMappingRequest builds an HTTP request for creating an ingredient tag mapping.
func (c *V1Client) BuildCreateIngredientTagMappingRequest(ctx context.Context, input *models.IngredientTagMappingCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateIngredientTagMappingRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(input.BelongsToValidIngredient, 10),
		ingredientTagMappingsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateIngredientTagMapping creates an ingredient tag mapping.
func (c *V1Client) CreateIngredientTagMapping(ctx context.Context, input *models.IngredientTagMappingCreationInput) (ingredientTagMapping *models.IngredientTagMapping, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateIngredientTagMapping")
	defer span.End()

	req, err := c.BuildCreateIngredientTagMappingRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &ingredientTagMapping)
	return ingredientTagMapping, err
}

// BuildUpdateIngredientTagMappingRequest builds an HTTP request for updating an ingredient tag mapping.
func (c *V1Client) BuildUpdateIngredientTagMappingRequest(ctx context.Context, ingredientTagMapping *models.IngredientTagMapping) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateIngredientTagMappingRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(ingredientTagMapping.BelongsToValidIngredient, 10),
		ingredientTagMappingsBasePath,
		strconv.FormatUint(ingredientTagMapping.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, ingredientTagMapping)
}

// UpdateIngredientTagMapping updates an ingredient tag mapping.
func (c *V1Client) UpdateIngredientTagMapping(ctx context.Context, ingredientTagMapping *models.IngredientTagMapping) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateIngredientTagMapping")
	defer span.End()

	req, err := c.BuildUpdateIngredientTagMappingRequest(ctx, ingredientTagMapping)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &ingredientTagMapping)
}

// BuildArchiveIngredientTagMappingRequest builds an HTTP request for updating an ingredient tag mapping.
func (c *V1Client) BuildArchiveIngredientTagMappingRequest(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveIngredientTagMappingRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		validIngredientsBasePath,
		strconv.FormatUint(validIngredientID, 10),
		ingredientTagMappingsBasePath,
		strconv.FormatUint(ingredientTagMappingID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveIngredientTagMapping archives an ingredient tag mapping.
func (c *V1Client) ArchiveIngredientTagMapping(ctx context.Context, validIngredientID, ingredientTagMappingID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveIngredientTagMapping")
	defer span.End()

	req, err := c.BuildArchiveIngredientTagMappingRequest(ctx, validIngredientID, ingredientTagMappingID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
