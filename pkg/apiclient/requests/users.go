package requests

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	usersBasePath = "users"
)

// BuildGetUserRequest builds an HTTP request for fetching a user.
func (b *Builder) BuildGetUserRequest(ctx context.Context, userID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	uri := b.BuildURL(ctx, nil, usersBasePath, userID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildGetUsersRequest builds an HTTP request for fetching a list of users.
func (b *Builder) BuildGetUsersRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	tracing.AttachQueryFilterToSpan(span, filter)
	uri := b.BuildURL(ctx, filter.ToValues(), usersBasePath)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildSearchForUsersByUsernameRequest builds an HTTP request that searches for a user by their username.
func (b *Builder) BuildSearchForUsersByUsernameRequest(ctx context.Context, username string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if username == "" {
		return nil, ErrEmptyUsernameProvided
	}

	tracing.AttachUsernameToSpan(span, username)

	u := b.buildAPIV1URL(ctx, nil, usersBasePath, "search")
	q := u.Query()

	q.Set(types.SearchQueryKey, username)
	u.RawQuery = q.Encode()
	uri := u.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildCreateUserRequest builds an HTTP request for creating a user.
func (b *Builder) BuildCreateUserRequest(ctx context.Context, input *types.UserRegistrationInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	tracing.AttachUsernameToSpan(span, input.Username)

	qp := url.Values{}
	if input.InvitationID != "" && input.InvitationToken != "" {
		qp.Set("i", input.InvitationID)
		qp.Set("t", input.InvitationToken)
	}

	// deliberately not validating here
	uri := b.buildUnversionedURL(ctx, qp, usersBasePath)

	tracing.AttachRequestURIToSpan(span, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildArchiveUserRequest builds an HTTP request for archiving a user.
func (b *Builder) BuildArchiveUserRequest(ctx context.Context, userID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if userID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachUserIDToSpan(span, userID)

	// deliberately not validating here, maybe there should make a client-side validate method vs a server-side?

	uri := b.buildAPIV1URL(ctx, nil, usersBasePath, userID).String()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, http.NoBody)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}

// BuildAvatarUploadRequest builds an HTTP request that sets a user's avatar to the provided content.
func (b *Builder) BuildAvatarUploadRequest(ctx context.Context, avatar []byte, extension string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if len(avatar) == 0 {
		return nil, ErrNilInputProvided
	}

	var ct string

	switch strings.ToLower(strings.TrimSpace(extension)) {
	case "jpeg":
		ct = "image/jpeg"
	case "png":
		ct = "image/png"
	case "gif":
		ct = "image/gif"
	default:
		return nil, fmt.Errorf("%s: %w", extension, ErrInvalidPhotoEncodingForUpload)
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("avatar", fmt.Sprintf("avatar.%s", extension))
	if err != nil {
		return nil, observability.PrepareError(err, span, "creating form file")
	}

	if _, err = io.Copy(part, bytes.NewReader(avatar)); err != nil {
		return nil, observability.PrepareError(err, span, "copying file contents to request")
	}

	if err = writer.Close(); err != nil {
		return nil, observability.PrepareError(err, span, "closing avatar writer")
	}

	uri := b.BuildURL(ctx, nil, usersBasePath, "avatar", "upload")

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building avatar upload request")
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-Upload-Content-Type", ct)

	return req, nil
}

// BuildCheckUserPermissionsRequests builds an HTTP request for checking a user's permissions.
func (b *Builder) BuildCheckUserPermissionsRequests(ctx context.Context, permissions ...string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if len(permissions) == 0 {
		return nil, ErrNilInputProvided
	}

	uri := b.BuildURL(ctx, nil, usersBasePath, "permissions", "check")
	body := &types.UserPermissionsRequestInput{Permissions: permissions}

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, body)
	if err != nil {
		return nil, observability.PrepareError(err, span, "building request")
	}

	return req, nil
}
