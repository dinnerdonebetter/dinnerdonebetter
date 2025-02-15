package apiclient

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
)

const (
	imagePNG  = "image/png"
	imageJPEG = "image/jpeg"
	imageGIF  = "image/gif"
)

func (c *Client) prepareUploads(ctx context.Context, files map[string][]byte) (io.Reader, string, error) {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	if files == nil {
		return nil, "", ErrNilInputProvided
	}

	var ct string
	for filename := range files {
		switch strings.ToLower(strings.TrimSpace(filepath.Ext(filename))) {
		case ".jpeg":
			if ct == "" {
				ct = imageJPEG
			} else if ct != imageJPEG {
				return nil, "", fmt.Errorf("all file uploads must be the same type: %w", ErrInvalidPhotoEncodingForUpload)
			}
		case ".png":
			if ct == "" {
				ct = imagePNG
			} else if ct != imagePNG {
				return nil, "", fmt.Errorf("all file uploads must be the same type: %w", ErrInvalidPhotoEncodingForUpload)
			}
		case ".gif":
			if ct == "" {
				ct = imageGIF
			} else if ct != imageGIF {
				return nil, "", fmt.Errorf("all file uploads must be the same type: %w", ErrInvalidPhotoEncodingForUpload)
			}
		}
	}

	if ct == "" {
		return nil, "", ErrInvalidPhotoEncodingForUpload
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	for filename, content := range files {
		part, err := writer.CreateFormFile("upload", filename)
		if err != nil {
			return nil, "", observability.PrepareError(err, span, "creating form file")
		}

		if _, err = io.Copy(part, bytes.NewReader(content)); err != nil {
			return nil, "", observability.PrepareError(err, span, "copying file contents to request")
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", observability.PrepareError(err, span, "closing media writer")
	}

	return body, writer.FormDataContentType(), nil
}

// buildMultipleRecipeMediaUploadRequest builds an HTTP request that sets a user's avatar to the provided content.
func (c *Client) buildMultipleRecipeMediaUploadRequest(ctx context.Context, uri string, files map[string][]byte) (*http.Request, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	body, formDataContentType, err := c.prepareUploads(ctx, files)
	if err != nil {
		return nil, observability.PrepareError(err, span, "preparing upload request")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "media upload")
	}

	req.Header.Set("Content-Type", formDataContentType)

	return req, nil
}
