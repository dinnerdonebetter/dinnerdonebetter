package requests

import (
	"context"
	"net/http"
	"strconv"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

const (
	invitationsBasePath = "invitations"
)

// BuildInvitationExistsRequest builds an HTTP request for checking the existence of an invitation.
func (b *Builder) BuildInvitationExistsRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	uri := b.BuildURL(
		ctx,
		nil,
		invitationsBasePath,
		id(invitationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetInvitationRequest builds an HTTP request for fetching an invitation.
func (b *Builder) BuildGetInvitationRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	uri := b.BuildURL(
		ctx,
		nil,
		invitationsBasePath,
		id(invitationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetInvitationsRequest builds an HTTP request for fetching a list of invitations.
func (b *Builder) BuildGetInvitationsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		invitationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateInvitationRequest builds an HTTP request for creating an invitation.
func (b *Builder) BuildCreateInvitationRequest(ctx context.Context, input *types.InvitationCreationInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	uri := b.BuildURL(
		ctx,
		nil,
		invitationsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateInvitationRequest builds an HTTP request for updating an invitation.
func (b *Builder) BuildUpdateInvitationRequest(ctx context.Context, invitation *types.Invitation) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if invitation == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.InvitationIDKey, invitation.ID)
	tracing.AttachInvitationIDToSpan(span, invitation.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		invitationsBasePath,
		strconv.FormatUint(invitation.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, invitation)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveInvitationRequest builds an HTTP request for archiving an invitation.
func (b *Builder) BuildArchiveInvitationRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	uri := b.BuildURL(
		ctx,
		nil,
		invitationsBasePath,
		id(invitationID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetAuditLogForInvitationRequest builds an HTTP request for fetching a list of audit log entries pertaining to an invitation.
func (b *Builder) BuildGetAuditLogForInvitationRequest(ctx context.Context, invitationID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if invitationID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.InvitationIDKey, invitationID)
	tracing.AttachInvitationIDToSpan(span, invitationID)

	uri := b.BuildURL(
		ctx,
		nil,
		invitationsBasePath,
		id(invitationID),
		"audit",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}
