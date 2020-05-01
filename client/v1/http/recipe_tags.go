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
	recipeTagsBasePath = "recipe_tags"
)

// BuildRecipeTagExistsRequest builds an HTTP request for checking the existence of a recipe tag.
func (c *V1Client) BuildRecipeTagExistsRequest(ctx context.Context, recipeID, recipeTagID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildRecipeTagExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeTagsBasePath,
		strconv.FormatUint(recipeTagID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// RecipeTagExists retrieves whether or not a recipe tag exists.
func (c *V1Client) RecipeTagExists(ctx context.Context, recipeID, recipeTagID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "RecipeTagExists")
	defer span.End()

	req, err := c.BuildRecipeTagExistsRequest(ctx, recipeID, recipeTagID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetRecipeTagRequest builds an HTTP request for fetching a recipe tag.
func (c *V1Client) BuildGetRecipeTagRequest(ctx context.Context, recipeID, recipeTagID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeTagsBasePath,
		strconv.FormatUint(recipeTagID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeTag retrieves a recipe tag.
func (c *V1Client) GetRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) (recipeTag *models.RecipeTag, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeTag")
	defer span.End()

	req, err := c.BuildGetRecipeTagRequest(ctx, recipeID, recipeTagID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeTag); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeTag, nil
}

// BuildGetRecipeTagsRequest builds an HTTP request for fetching recipe tags.
func (c *V1Client) BuildGetRecipeTagsRequest(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetRecipeTagsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeTagsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetRecipeTags retrieves a list of recipe tags.
func (c *V1Client) GetRecipeTags(ctx context.Context, recipeID uint64, filter *models.QueryFilter) (recipeTags *models.RecipeTagList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetRecipeTags")
	defer span.End()

	req, err := c.BuildGetRecipeTagsRequest(ctx, recipeID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &recipeTags); retrieveErr != nil {
		return nil, retrieveErr
	}

	return recipeTags, nil
}

// BuildCreateRecipeTagRequest builds an HTTP request for creating a recipe tag.
func (c *V1Client) BuildCreateRecipeTagRequest(ctx context.Context, input *models.RecipeTagCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateRecipeTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(input.BelongsToRecipe, 10),
		recipeTagsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateRecipeTag creates a recipe tag.
func (c *V1Client) CreateRecipeTag(ctx context.Context, input *models.RecipeTagCreationInput) (recipeTag *models.RecipeTag, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateRecipeTag")
	defer span.End()

	req, err := c.BuildCreateRecipeTagRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &recipeTag)
	return recipeTag, err
}

// BuildUpdateRecipeTagRequest builds an HTTP request for updating a recipe tag.
func (c *V1Client) BuildUpdateRecipeTagRequest(ctx context.Context, recipeTag *models.RecipeTag) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateRecipeTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeTag.BelongsToRecipe, 10),
		recipeTagsBasePath,
		strconv.FormatUint(recipeTag.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, recipeTag)
}

// UpdateRecipeTag updates a recipe tag.
func (c *V1Client) UpdateRecipeTag(ctx context.Context, recipeTag *models.RecipeTag) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateRecipeTag")
	defer span.End()

	req, err := c.BuildUpdateRecipeTagRequest(ctx, recipeTag)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &recipeTag)
}

// BuildArchiveRecipeTagRequest builds an HTTP request for updating a recipe tag.
func (c *V1Client) BuildArchiveRecipeTagRequest(ctx context.Context, recipeID, recipeTagID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveRecipeTagRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeTagsBasePath,
		strconv.FormatUint(recipeTagID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveRecipeTag archives a recipe tag.
func (c *V1Client) ArchiveRecipeTag(ctx context.Context, recipeID, recipeTagID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveRecipeTag")
	defer span.End()

	req, err := c.BuildArchiveRecipeTagRequest(ctx, recipeID, recipeTagID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
