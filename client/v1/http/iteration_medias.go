package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	iterationMediasBasePath = "iteration_medias"
)

// BuildGetIterationMediaRequest builds an HTTP request for fetching an iteration media
func (c *V1Client) BuildGetIterationMediaRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, iterationMediasBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetIterationMedia retrieves an iteration media
func (c *V1Client) GetIterationMedia(ctx context.Context, id uint64) (iterationMedia *models.IterationMedia, err error) {
	req, err := c.BuildGetIterationMediaRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &iterationMedia); retrieveErr != nil {
		return nil, retrieveErr
	}

	return iterationMedia, nil
}

// BuildGetIterationMediasRequest builds an HTTP request for fetching iteration medias
func (c *V1Client) BuildGetIterationMediasRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), iterationMediasBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetIterationMedias retrieves a list of iteration medias
func (c *V1Client) GetIterationMedias(ctx context.Context, filter *models.QueryFilter) (iterationMedias *models.IterationMediaList, err error) {
	req, err := c.BuildGetIterationMediasRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &iterationMedias); retrieveErr != nil {
		return nil, retrieveErr
	}

	return iterationMedias, nil
}

// BuildCreateIterationMediaRequest builds an HTTP request for creating an iteration media
func (c *V1Client) BuildCreateIterationMediaRequest(ctx context.Context, body *models.IterationMediaCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, iterationMediasBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateIterationMedia creates an iteration media
func (c *V1Client) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (iterationMedia *models.IterationMedia, err error) {
	req, err := c.BuildCreateIterationMediaRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &iterationMedia)
	return iterationMedia, err
}

// BuildUpdateIterationMediaRequest builds an HTTP request for updating an iteration media
func (c *V1Client) BuildUpdateIterationMediaRequest(ctx context.Context, updated *models.IterationMedia) (*http.Request, error) {
	uri := c.BuildURL(nil, iterationMediasBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateIterationMedia updates an iteration media
func (c *V1Client) UpdateIterationMedia(ctx context.Context, updated *models.IterationMedia) error {
	req, err := c.BuildUpdateIterationMediaRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveIterationMediaRequest builds an HTTP request for updating an iteration media
func (c *V1Client) BuildArchiveIterationMediaRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, iterationMediasBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveIterationMedia archives an iteration media
func (c *V1Client) ArchiveIterationMedia(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveIterationMediaRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
