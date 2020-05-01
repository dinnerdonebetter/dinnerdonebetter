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
	iterationMediasBasePath = "iteration_medias"
)

// BuildIterationMediaExistsRequest builds an HTTP request for checking the existence of an iteration media.
func (c *V1Client) BuildIterationMediaExistsRequest(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildIterationMediaExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
		iterationMediasBasePath,
		strconv.FormatUint(iterationMediaID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// IterationMediaExists retrieves whether or not an iteration media exists.
func (c *V1Client) IterationMediaExists(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "IterationMediaExists")
	defer span.End()

	req, err := c.BuildIterationMediaExistsRequest(ctx, recipeID, recipeIterationID, iterationMediaID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetIterationMediaRequest builds an HTTP request for fetching an iteration media.
func (c *V1Client) BuildGetIterationMediaRequest(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetIterationMediaRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
		iterationMediasBasePath,
		strconv.FormatUint(iterationMediaID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetIterationMedia retrieves an iteration media.
func (c *V1Client) GetIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (iterationMedia *models.IterationMedia, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetIterationMedia")
	defer span.End()

	req, err := c.BuildGetIterationMediaRequest(ctx, recipeID, recipeIterationID, iterationMediaID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &iterationMedia); retrieveErr != nil {
		return nil, retrieveErr
	}

	return iterationMedia, nil
}

// BuildGetIterationMediasRequest builds an HTTP request for fetching iteration medias.
func (c *V1Client) BuildGetIterationMediasRequest(ctx context.Context, recipeID, recipeIterationID uint64, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetIterationMediasRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
		iterationMediasBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetIterationMedias retrieves a list of iteration medias.
func (c *V1Client) GetIterationMedias(ctx context.Context, recipeID, recipeIterationID uint64, filter *models.QueryFilter) (iterationMedias *models.IterationMediaList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetIterationMedias")
	defer span.End()

	req, err := c.BuildGetIterationMediasRequest(ctx, recipeID, recipeIterationID, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &iterationMedias); retrieveErr != nil {
		return nil, retrieveErr
	}

	return iterationMedias, nil
}

// BuildCreateIterationMediaRequest builds an HTTP request for creating an iteration media.
func (c *V1Client) BuildCreateIterationMediaRequest(ctx context.Context, recipeID uint64, input *models.IterationMediaCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateIterationMediaRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(input.BelongsToRecipeIteration, 10),
		iterationMediasBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateIterationMedia creates an iteration media.
func (c *V1Client) CreateIterationMedia(ctx context.Context, recipeID uint64, input *models.IterationMediaCreationInput) (iterationMedia *models.IterationMedia, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateIterationMedia")
	defer span.End()

	req, err := c.BuildCreateIterationMediaRequest(ctx, recipeID, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &iterationMedia)
	return iterationMedia, err
}

// BuildUpdateIterationMediaRequest builds an HTTP request for updating an iteration media.
func (c *V1Client) BuildUpdateIterationMediaRequest(ctx context.Context, recipeID uint64, iterationMedia *models.IterationMedia) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateIterationMediaRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(iterationMedia.BelongsToRecipeIteration, 10),
		iterationMediasBasePath,
		strconv.FormatUint(iterationMedia.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, iterationMedia)
}

// UpdateIterationMedia updates an iteration media.
func (c *V1Client) UpdateIterationMedia(ctx context.Context, recipeID uint64, iterationMedia *models.IterationMedia) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateIterationMedia")
	defer span.End()

	req, err := c.BuildUpdateIterationMediaRequest(ctx, recipeID, iterationMedia)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &iterationMedia)
}

// BuildArchiveIterationMediaRequest builds an HTTP request for updating an iteration media.
func (c *V1Client) BuildArchiveIterationMediaRequest(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveIterationMediaRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		recipesBasePath,
		strconv.FormatUint(recipeID, 10),
		recipeIterationsBasePath,
		strconv.FormatUint(recipeIterationID, 10),
		iterationMediasBasePath,
		strconv.FormatUint(iterationMediaID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveIterationMedia archives an iteration media.
func (c *V1Client) ArchiveIterationMedia(ctx context.Context, recipeID, recipeIterationID, iterationMediaID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveIterationMedia")
	defer span.End()

	req, err := c.BuildArchiveIterationMediaRequest(ctx, recipeID, recipeIterationID, iterationMediaID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
