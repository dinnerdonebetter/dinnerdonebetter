package requests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/internal/observability/keys"
	"github.com/dinnerdonebetter/backend/internal/observability/tracing"
)

const (
	recipesBasePath = "recipes"

	imagePNG  = "image/png"
	imageJPEG = "image/jpeg"
	imageGIF  = "image/gif"
)

func (b *Builder) prepareUploads(ctx context.Context, files map[string][]byte) (io.Reader, string, error) {
	_, span := b.tracer.StartSpan(ctx)
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

// BuildMultipleRecipeMediaUploadRequest builds an HTTP request that sets a user's avatar to the provided content.
func (b *Builder) BuildMultipleRecipeMediaUploadRequest(ctx context.Context, files map[string][]byte, recipeID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if recipeID == "" {
		return nil, ErrInvalidIDProvided
	}
	tracing.AttachToSpan(span, keys.RecipeIDKey, recipeID)

	body, formDataContentType, err := b.prepareUploads(ctx, files)
	if err != nil {
		return nil, observability.PrepareError(err, span, "preparing upload request")
	}

	uri := b.BuildURL(ctx, nil, recipesBasePath, recipeID, "images")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "media upload")
	}

	req.Header.Set("Content-Type", formDataContentType)

	return req, nil
}
