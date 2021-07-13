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
	reportsBasePath = "reports"
)

// BuildReportExistsRequest builds an HTTP request for checking the existence of a report.
func (b *Builder) BuildReportExistsRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	uri := b.BuildURL(
		ctx,
		nil,
		reportsBasePath,
		id(reportID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetReportRequest builds an HTTP request for fetching a report.
func (b *Builder) BuildGetReportRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	uri := b.BuildURL(
		ctx,
		nil,
		reportsBasePath,
		id(reportID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetReportsRequest builds an HTTP request for fetching a list of reports.
func (b *Builder) BuildGetReportsRequest(ctx context.Context, filter *types.QueryFilter) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := filter.AttachToLogger(b.logger)

	uri := b.BuildURL(
		ctx,
		filter.ToValues(),
		reportsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildCreateReportRequest builds an HTTP request for creating a report.
func (b *Builder) BuildCreateReportRequest(ctx context.Context, input *types.ReportCreationInput) (*http.Request, error) {
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
		reportsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPost, uri, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildUpdateReportRequest builds an HTTP request for updating a report.
func (b *Builder) BuildUpdateReportRequest(ctx context.Context, report *types.Report) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if report == nil {
		return nil, ErrNilInputProvided
	}

	logger = logger.WithValue(keys.ReportIDKey, report.ID)
	tracing.AttachReportIDToSpan(span, report.ID)

	uri := b.BuildURL(
		ctx,
		nil,
		reportsBasePath,
		strconv.FormatUint(report.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := b.buildDataRequest(ctx, http.MethodPut, uri, report)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building request")
	}

	return req, nil
}

// BuildArchiveReportRequest builds an HTTP request for archiving a report.
func (b *Builder) BuildArchiveReportRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	uri := b.BuildURL(
		ctx,
		nil,
		reportsBasePath,
		id(reportID),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}

// BuildGetAuditLogForReportRequest builds an HTTP request for fetching a list of audit log entries pertaining to a report.
func (b *Builder) BuildGetAuditLogForReportRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := b.tracer.StartSpan(ctx)
	defer span.End()

	logger := b.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	uri := b.BuildURL(
		ctx,
		nil,
		reportsBasePath,
		id(reportID),
		"audit",
	)
	tracing.AttachRequestURIToSpan(span, uri)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building user status request")
	}

	return req, nil
}
