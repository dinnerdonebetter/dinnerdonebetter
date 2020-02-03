package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

const (
	reportsBasePath = "reports"
)

// BuildGetReportRequest builds an HTTP request for fetching a report
func (c *V1Client) BuildGetReportRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, reportsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetReport retrieves a report
func (c *V1Client) GetReport(ctx context.Context, id uint64) (report *models.Report, err error) {
	req, err := c.BuildGetReportRequest(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &report); retrieveErr != nil {
		return nil, retrieveErr
	}

	return report, nil
}

// BuildGetReportsRequest builds an HTTP request for fetching reports
func (c *V1Client) BuildGetReportsRequest(ctx context.Context, filter *models.QueryFilter) (*http.Request, error) {
	uri := c.BuildURL(filter.ToValues(), reportsBasePath)

	return http.NewRequest(http.MethodGet, uri, nil)
}

// GetReports retrieves a list of reports
func (c *V1Client) GetReports(ctx context.Context, filter *models.QueryFilter) (reports *models.ReportList, err error) {
	req, err := c.BuildGetReportsRequest(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	if retrieveErr := c.retrieve(ctx, req, &reports); retrieveErr != nil {
		return nil, retrieveErr
	}

	return reports, nil
}

// BuildCreateReportRequest builds an HTTP request for creating a report
func (c *V1Client) BuildCreateReportRequest(ctx context.Context, body *models.ReportCreationInput) (*http.Request, error) {
	uri := c.BuildURL(nil, reportsBasePath)

	return c.buildDataRequest(http.MethodPost, uri, body)
}

// CreateReport creates a report
func (c *V1Client) CreateReport(ctx context.Context, input *models.ReportCreationInput) (report *models.Report, err error) {
	req, err := c.BuildCreateReportRequest(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("building request: %w", err)
	}

	err = c.executeRequest(ctx, req, &report)
	return report, err
}

// BuildUpdateReportRequest builds an HTTP request for updating a report
func (c *V1Client) BuildUpdateReportRequest(ctx context.Context, updated *models.Report) (*http.Request, error) {
	uri := c.BuildURL(nil, reportsBasePath, strconv.FormatUint(updated.ID, 10))

	return c.buildDataRequest(http.MethodPut, uri, updated)
}

// UpdateReport updates a report
func (c *V1Client) UpdateReport(ctx context.Context, updated *models.Report) error {
	req, err := c.BuildUpdateReportRequest(ctx, updated)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, &updated)
}

// BuildArchiveReportRequest builds an HTTP request for updating a report
func (c *V1Client) BuildArchiveReportRequest(ctx context.Context, id uint64) (*http.Request, error) {
	uri := c.BuildURL(nil, reportsBasePath, strconv.FormatUint(id, 10))

	return http.NewRequest(http.MethodDelete, uri, nil)
}

// ArchiveReport archives a report
func (c *V1Client) ArchiveReport(ctx context.Context, id uint64) error {
	req, err := c.BuildArchiveReportRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("building request: %w", err)
	}

	return c.executeRequest(ctx, req, nil)
}
