package requests

import (
	"context"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	invitationsBasePath = "invitations"
)

// BuildGetHouseholdInvitationRequest builds an HTTP request for fetching a webhook.
func (b *Builder) BuildGetHouseholdInvitationRequest(ctx context.Context, householdID, invitationID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if invitationID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := b.logger.WithValue(keys.HouseholdInvitationIDKey, invitationID)
	tracing.AttachHouseholdInvitationIDToSpan(span, invitationID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
		invitationsBasePath,
		invitationID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetPendingHouseholdInvitationsFromUserRequest builds an HTTP request that retrieves pending household invitations sent by a user.
func (b *Builder) BuildGetPendingHouseholdInvitationsFromUserRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	uri := b.BuildURL(ctx, filter.ToValues(), "household_invitations", "sent")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetPendingHouseholdInvitationsForUserRequest builds an HTTP request that retrieves pending household invitations received by a user.
func (b *Builder) BuildGetPendingHouseholdInvitationsForUserRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	uri := b.BuildURL(ctx, filter.ToValues(), "household_invitations", "received")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCancelHouseholdInvitationRequest builds an HTTP request that cancels a given household invitation.
func (b *Builder) BuildCancelHouseholdInvitationRequest(ctx context.Context, householdID, invitationID, note string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
		invitationsBasePath,
		invitationID,
	)
	logger = logger.WithValue(keys.URLKey, uri)

	input := &types.HouseholdInvitationUpdateRequestInput{Note: note}
	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildAcceptHouseholdInvitationRequest builds an HTTP request that accepts a given household invitation.
func (b *Builder) BuildAcceptHouseholdInvitationRequest(ctx context.Context, householdID, invitationID, note string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
		invitationsBasePath,
		invitationID,
	)
	logger = logger.WithValue(keys.URLKey, uri)

	input := &types.HouseholdInvitationUpdateRequestInput{Note: note}
	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildRejectHouseholdInvitationRequest builds an HTTP request that rejects a given household invitation.
func (b *Builder) BuildRejectHouseholdInvitationRequest(ctx context.Context, householdID, invitationID, note string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
		invitationsBasePath,
		invitationID,
	)
	logger = logger.WithValue(keys.URLKey, uri)

	input := &types.HouseholdInvitationUpdateRequestInput{Note: note}
	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}
