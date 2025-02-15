package images

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dinnerdonebetter/backend/internal/lib/observability"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/logging"
	"github.com/dinnerdonebetter/backend/internal/lib/observability/tracing"

	"github.com/hashicorp/go-multierror"
)

const (
	headerContentType = "Content-Type"

	imagePNG  = "image/png"
	imageJPEG = "image/jpeg"
	imageGIF  = "image/gif"
)

var (
	// ErrInvalidImageContentType is what we return to indicate the provided image was of the wrong type.
	ErrInvalidImageContentType = errors.New("invalid image content type")
)

type (
	// Upload is a helper struct for handling images.
	Upload struct {
		Filename    string
		ContentType string
		Data        []byte
		Size        int
	}

	// MediaUploadProcessor processes media uploads.
	MediaUploadProcessor interface {
		ProcessFile(ctx context.Context, req *http.Request, filename string) (*Upload, error)
		ProcessFiles(ctx context.Context, req *http.Request, filenamePrefix string) ([]*Upload, error)
	}

	uploadProcessor struct {
		tracer tracing.Tracer
		logger logging.Logger
	}
)

// DataURI converts image to base64 data URI.
func (i *Upload) DataURI() string {
	return fmt.Sprintf("data:%s;base64,%s", i.ContentType, base64.StdEncoding.EncodeToString(i.Data))
}

// Write image to HTTP response.
func (i *Upload) Write(w http.ResponseWriter) error {
	w.Header().Set(headerContentType, i.ContentType)
	w.Header().Set("RawHTML-Length", strconv.Itoa(i.Size))

	if _, err := w.Write(i.Data); err != nil {
		return fmt.Errorf("writing image to HTTP response: %w", err)
	}

	return nil
}

// Thumbnail creates a thumbnail from an image.
func (i *Upload) Thumbnail(width, height uint, filename string) (*Upload, error) {
	t, err := newThumbnailer(i.ContentType)
	if err != nil {
		return nil, err
	}

	return t.Thumbnail(i, width, height, filename)
}

// NewImageUploadProcessor provides a new MediaUploadProcessor.
func NewImageUploadProcessor(logger logging.Logger, tracerProvider tracing.TracerProvider) MediaUploadProcessor {
	return &uploadProcessor{
		logger: logging.EnsureLogger(logger).WithName("media_upload_processor"),
		tracer: tracing.NewTracer(tracing.EnsureTracerProvider(tracerProvider).Tracer("media_upload_processor")),
	}
}

// LimitFileSize limits the size of uploaded files, for use before ProcessFile.
func LimitFileSize(maxSize uint16, res http.ResponseWriter, req *http.Request) {
	if maxSize == 0 {
		maxSize = 4096
	}

	req.Body = http.MaxBytesReader(res, req.Body, int64(maxSize))
}

func contentTypeFromFilename(filename string) string {
	ext := filepath.Ext(filename)

	switch ext {
	case ".png":
		return imagePNG
	case ".jpeg":
		return imageJPEG
	case ".gif":
		return imageGIF
	default:
		return mime.TypeByExtension(ext)
	}
}

func isImage(filename string) bool {
	contentType := contentTypeFromFilename(filename)

	switch strings.TrimSpace(strings.ToLower(contentType)) {
	case imagePNG, imageJPEG, imageGIF:
		return true
	default:
		return false
	}
}

// ProcessFile extracts an image from an *http.Request.
func (p *uploadProcessor) ProcessFile(ctx context.Context, req *http.Request, filename string) (*Upload, error) {
	_, span := p.tracer.StartSpan(ctx)
	defer span.End()

	logger := p.logger.WithRequest(req).WithValue("filename", filename)

	file, info, err := req.FormFile(filename)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "parsing filename %q from request", filename)
	}

	return p.processFile(ctx, file, info, filename)
}

// processFile extracts a single file by name from an *http.Request.
func (p *uploadProcessor) processFile(ctx context.Context, file multipart.File, info *multipart.FileHeader, filename string) (*Upload, error) {
	_, span := p.tracer.StartSpan(ctx)
	defer span.End()

	logger := p.logger.WithValue("filename", filename)
	tracing.AttachToSpan(span, "filename", filename)

	bs, err := io.ReadAll(file)
	if err != nil {
		return nil, observability.PrepareAndLogError(err, logger, span, "reading filename %q from request", filename)
	}

	if isImage(info.Filename) {
		if _, _, err = image.Decode(bytes.NewReader(bs)); err != nil {
			return nil, observability.PrepareAndLogError(err, logger, span, "decoding the image data")
		}
	}

	i := &Upload{
		Filename:    info.Filename,
		ContentType: contentTypeFromFilename(filename),
		Data:        bs,
		Size:        len(bs),
	}

	return i, nil
}

const (
	defaultMaxMemory = 32 << 20 // 32 MB
)

// ProcessFiles extracts all the files from an *http.Request.
func (p *uploadProcessor) ProcessFiles(ctx context.Context, req *http.Request, filenamePrefix string) ([]*Upload, error) {
	_, span := p.tracer.StartSpan(ctx)
	defer span.End()

	logger := p.logger.WithRequest(req).WithValue("filename_prefix", filenamePrefix)

	if req.MultipartForm == nil {
		if err := req.ParseMultipartForm(defaultMaxMemory); err != nil {
			return nil, fmt.Errorf("parsing multipart form: %w", err)
		}
	}

	var (
		uploads = []*Upload{}
		errs    *multierror.Error
	)
	if req.MultipartForm != nil && req.MultipartForm.File != nil {
		for _, fileHeaders := range req.MultipartForm.File {
			for i, fileHeader := range fileHeaders {
				file, err := fileHeader.Open()
				if err != nil {
					errs = multierror.Append(errs, err)
				} else {
					upload, parseUploadErr := p.processFile(ctx, file, fileHeader, fmt.Sprintf("%s_%d", filenamePrefix, i))
					if parseUploadErr != nil {
						errs = multierror.Append(errs, parseUploadErr)
					} else {
						uploads = append(uploads, upload)
					}
				}
			}
		}
	}

	if errs != nil {
		logger.Error("processing image uploads", errs)
		return nil, errs
	}

	return uploads, nil
}
