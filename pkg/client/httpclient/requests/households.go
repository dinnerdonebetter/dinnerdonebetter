package requests

import (
	"context"
	"fmt"
	"net/http"

	"github.com/prixfixeco/api_server/internal/observability"
	"github.com/prixfixeco/api_server/internal/observability/keys"
	"github.com/prixfixeco/api_server/internal/observability/tracing"
	"github.com/prixfixeco/api_server/pkg/types"
)

const (
	householdsBasePath = "households"
)

// BuildSwitchActiveHouseholdRequest builds an HTTP request for switching active households.
func (b *Builder) BuildSwitchActiveHouseholdRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	tracing.AttachHouseholdIDToSpan(span, householdID)

	uri := b.buildUnversionedURL(ctx, nil, usersBasePath, "household", "select")

	input := &types.ChangeActiveHouseholdInput{
		HouseholdID: householdID,
	}

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildGetCurrentHouseholdRequest builds an HTTP request for fetching a household.
func (b *Builder) BuildGetCurrentHouseholdRequest(ctx context.Context) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger
	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		"current",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetHouseholdRequest builds an HTTP request for fetching a household.
func (b *Builder) BuildGetHouseholdRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := b.logger.WithValue(keys.HouseholdIDKey, householdID)
	tracing.AttachHouseholdIDToSpan(span, householdID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetHouseholdsRequest builds an HTTP request for fetching a list of households.
func (b *Builder) BuildGetHouseholdsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)
	uri := b.BuildURL(ctx, filter.ToValues(), householdsBasePath)

	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateHouseholdRequest builds an HTTP request for creating a household.
func (b *Builder) BuildCreateHouseholdRequest(ctx context.Context, input *types.HouseholdCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := b.logger.WithValue(keys.NameKey, input.Name)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath)
	tracing.AttachRequestURIToSpan(span, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildUpdateHouseholdRequest builds an HTTP request for updating a household.
func (b *Builder) BuildUpdateHouseholdRequest(ctx context.Context, household *types.Household) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if household == nil {
		return nil, ErrNilInputProvided
	}

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		household.ID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return b.buildDataRequest(ctx, http.MethodPut, uri, household)
}

// BuildArchiveHouseholdRequest builds an HTTP request for archiving a household.
func (b *Builder) BuildArchiveHouseholdRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := b.logger.WithValue(keys.HouseholdIDKey, householdID)

	uri := b.BuildURL(
		ctx,
		nil,
		householdsBasePath,
		householdID,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildInviteUserToHouseholdRequest builds a request that adds a user to a household.
func (b *Builder) BuildInviteUserToHouseholdRequest(ctx context.Context, input *types.HouseholdInvitationCreationRequestInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if input == nil {
		return nil, ErrNilInputProvided
	}

	// we don't validate here because it needs to have the user ID

	uri := b.BuildURL(ctx, nil, householdsBasePath, input.DestinationHousehold, "invite")
	tracing.AttachRequestURIToSpan(span, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// BuildMarkAsDefaultRequest builds a request that marks a given household as the default for a given user.
func (b *Builder) BuildMarkAsDefaultRequest(ctx context.Context, householdID string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := b.logger.WithValue(keys.HouseholdIDKey, householdID)

	uri := b.BuildURL(ctx, nil, householdsBasePath, householdID, "default")
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildRemoveUserRequest builds a request that removes a user from a household.
func (b *Builder) BuildRemoveUserRequest(ctx context.Context, householdID, userID, reason string) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return nil, ErrInvalidIDProvided
	}

	logger := b.logger.WithValue(keys.HouseholdIDKey, householdID).
		WithValue(keys.UserIDKey, userID).
		WithValue(keys.ReasonKey, reason)

	u := b.buildAPIV1URL(ctx, nil, householdsBasePath, householdID, "members", userID)

	if reason != "" {
		q := u.Query()
		q.Set("reason", reason)
		u.RawQuery = q.Encode()
	}

	tracing.AttachURLToSpan(span, u)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u.String(), nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildModifyMemberPermissionsRequest builds a request that modifies a given user's permissions for a given household.
func (b *Builder) BuildModifyMemberPermissionsRequest(ctx context.Context, householdID, userID string, input *types.ModifyUserPermissionsInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" || userID == "" {
		return nil, ErrInvalidIDProvided
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := b.logger.WithValue(keys.UserIDKey, userID).WithValue(keys.HouseholdIDKey, householdID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath, householdID, "members", userID, "permissions")
	tracing.AttachRequestURIToSpan(span, uri)

	return b.buildDataRequest(ctx, http.MethodPatch, uri, input)
}

// BuildTransferHouseholdOwnershipRequest builds a request that transfers ownership of a household to a given user.
func (b *Builder) BuildTransferHouseholdOwnershipRequest(ctx context.Context, householdID string, input *types.HouseholdOwnershipTransferInput) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	if householdID == "" {
		return nil, fmt.Errorf("householdID: %w", ErrInvalidIDProvided)
	}

	if input == nil {
		return nil, ErrNilInputProvided
	}

	logger := b.logger.WithValue(keys.HouseholdIDKey, householdID)

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	uri := b.BuildURL(ctx, nil, householdsBasePath, householdID, "transfer")
	tracing.AttachRequestURIToSpan(span, uri)

	return b.buildDataRequest(ctx, http.MethodPost, uri, input)
}
