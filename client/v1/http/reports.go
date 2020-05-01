package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	reportsBasePath = "reports"
)

// BuildReportExistsRequest builds an HTTP request for checking the existence of a report.
func (c *V1Client) BuildReportExistsRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildReportExistsRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		reportsBasePath,
		strconv.FormatUint(reportID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodHead, uri, nil)
}

// ReportExists retrieves whether or not a report exists.
func (c *V1Client) ReportExists(ctx context.Context, reportID uint64) (exists bool, err error) {
	ctx, span := tracing.StartSpan(ctx, "ReportExists")
	defer span.End()

	req, err := c.BuildReportExistsRequest(ctx, reportID)
	if err != nil {
		return false, fmt.Errorf("building request: %w", err)
	}

	return c.checkExistence(ctx, req)
}

// BuildGetReportRequest builds an HTTP request for fetching a report.
func (c *V1Client) BuildGetReportRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetReportRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		reportsBasePath,
		strconv.FormatUint(reportID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetReport retrieves a report.
func (c *V1Client) GetReport(ctx context.Context, reportID uint64) (report *models.Report, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetReport")
	defer span.End()

	req, err := c.BuildGetReportRequest(ctx, reportID)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &report); retrieveErr != nil {
		return nil, retrieveErr
	}

	return report, nil
}

// BuildGetReportsRequest builds an HTTP request for fetching reports.
func (c *V1Client) BuildGetReportsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildGetReportsRequest")
	defer span.End()

	uri := c.BuildURL(
		filter.ToValues(),
		reportsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodGet, uri, nil)
}

// GetReports retrieves a list of reports.
func (c *V1Client) GetReports(ctx context.Context, filter *models.QueryFilter) (reports *models.ReportList, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetReports")
	defer span.End()

	req, err := c.BuildGetReportsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &reports); retrieveErr != nil {
		return nil, retrieveErr
	}

	return reports, nil
}

// BuildCreateReportRequest builds an HTTP request for creating a report.
func (c *V1Client) BuildCreateReportRequest(ctx context.Context, input *models.ReportCreationInput) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildCreateReportRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		reportsBasePath,
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPost, uri, input)
}

// CreateReport creates a report.
func (c *V1Client) CreateReport(ctx context.Context, input *models.ReportCreationInput) (report *models.Report, err error) {
	ctx, span := tracing.StartSpan(ctx, "CreateReport")
	defer span.End()

	req, err := c.BuildCreateReportRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &report)
	return report, err
}

// BuildUpdateReportRequest builds an HTTP request for updating a report.
func (c *V1Client) BuildUpdateReportRequest(ctx context.Context, report *models.Report) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildUpdateReportRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		reportsBasePath,
		strconv.FormatUint(report.ID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return c.buildDataRequest(ctx, http.MethodPut, uri, report)
}

// UpdateReport updates a report.
func (c *V1Client) UpdateReport(ctx context.Context, report *models.Report) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateReport")
	defer span.End()

	req, err := c.BuildUpdateReportRequest(ctx, report)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &report)
}

// BuildArchiveReportRequest builds an HTTP request for updating a report.
func (c *V1Client) BuildArchiveReportRequest(ctx context.Context, reportID uint64) (*http.Request, error) {
	ctx, span := tracing.StartSpan(ctx, "BuildArchiveReportRequest")
	defer span.End()

	uri := c.BuildURL(
		nil,
		reportsBasePath,
		strconv.FormatUint(reportID, 10),
	)
	tracing.AttachRequestURIToSpan(span, uri)

	return http.NewRequestWithContext(ctx, http.MethodDelete, uri, nil)
}

// ArchiveReport archives a report.
func (c *V1Client) ArchiveReport(ctx context.Context, reportID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveReport")
	defer span.End()

	req, err := c.BuildArchiveReportRequest(ctx, reportID)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
