package httpclient

import (
	"context"

	observability "gitlab.com/prixfixe/prixfixe/internal/observability"
	keys "gitlab.com/prixfixe/prixfixe/internal/observability/keys"
	"gitlab.com/prixfixe/prixfixe/internal/observability/tracing"
	"gitlab.com/prixfixe/prixfixe/pkg/types"
)

// ReportExists retrieves whether a report exists.
func (c *Client) ReportExists(ctx context.Context, reportID uint64) (bool, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if reportID == 0 {
		return false, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	req, err := c.requestBuilder.BuildReportExistsRequest(ctx, reportID)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "building report existence request")
	}

	exists, err := c.responseIsOK(ctx, req)
	if err != nil {
		return false, observability.PrepareError(err, logger, span, "checking existence for report #%d", reportID)
	}

	return exists, nil
}

// GetReport gets a report.
func (c *Client) GetReport(ctx context.Context, reportID uint64) (*types.Report, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	req, err := c.requestBuilder.BuildGetReportRequest(ctx, reportID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get report request")
	}

	var report *types.Report
	if err = c.fetchAndUnmarshal(ctx, req, &report); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving report")
	}

	return report, nil
}

// GetReports retrieves a list of reports.
func (c *Client) GetReports(ctx context.Context, filter *types.QueryFilter) (*types.ReportList, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.loggerWithFilter(filter)
	tracing.AttachQueryFilterToSpan(span, filter)

	req, err := c.requestBuilder.BuildGetReportsRequest(ctx, filter)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building reports list request")
	}

	var reports *types.ReportList
	if err = c.fetchAndUnmarshal(ctx, req, &reports); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving reports")
	}

	return reports, nil
}

// CreateReport creates a report.
func (c *Client) CreateReport(ctx context.Context, input *types.ReportCreationInput) (*types.Report, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if input == nil {
		return nil, ErrNilInputProvided
	}

	if err := input.ValidateWithContext(ctx); err != nil {
		return nil, observability.PrepareError(err, logger, span, "validating input")
	}

	req, err := c.requestBuilder.BuildCreateReportRequest(ctx, input)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building create report request")
	}

	var report *types.Report
	if err = c.fetchAndUnmarshal(ctx, req, &report); err != nil {
		return nil, observability.PrepareError(err, logger, span, "creating report")
	}

	return report, nil
}

// UpdateReport updates a report.
func (c *Client) UpdateReport(ctx context.Context, report *types.Report) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if report == nil {
		return ErrNilInputProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, report.ID)
	tracing.AttachReportIDToSpan(span, report.ID)

	req, err := c.requestBuilder.BuildUpdateReportRequest(ctx, report)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building update report request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, &report); err != nil {
		return observability.PrepareError(err, logger, span, "updating report #%d", report.ID)
	}

	return nil
}

// ArchiveReport archives a report.
func (c *Client) ArchiveReport(ctx context.Context, reportID uint64) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if reportID == 0 {
		return ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	req, err := c.requestBuilder.BuildArchiveReportRequest(ctx, reportID)
	if err != nil {
		return observability.PrepareError(err, logger, span, "building archive report request")
	}

	if err = c.fetchAndUnmarshal(ctx, req, nil); err != nil {
		return observability.PrepareError(err, logger, span, "archiving report #%d", reportID)
	}

	return nil
}

// GetAuditLogForReport retrieves a list of audit log entries pertaining to a report.
func (c *Client) GetAuditLogForReport(ctx context.Context, reportID uint64) ([]*types.AuditLogEntry, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger

	if reportID == 0 {
		return nil, ErrInvalidIDProvided
	}
	logger = logger.WithValue(keys.ReportIDKey, reportID)
	tracing.AttachReportIDToSpan(span, reportID)

	req, err := c.requestBuilder.BuildGetAuditLogForReportRequest(ctx, reportID)
	if err != nil {
		return nil, observability.PrepareError(err, logger, span, "building get audit log entries for report request")
	}

	var entries []*types.AuditLogEntry
	if err = c.fetchAndUnmarshal(ctx, req, &entries); err != nil {
		return nil, observability.PrepareError(err, logger, span, "retrieving plan")
	}

	return entries, nil
}
