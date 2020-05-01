package dbclient

import (
	"context"

	"gitlab.com/prixfixe/prixfixe/internal/v1/tracing"
	models "gitlab.com/prixfixe/prixfixe/models/v1"
)

var _ models.ReportDataManager = (*Client)(nil)

// ReportExists fetches whether or not a report exists from the database.
func (c *Client) ReportExists(ctx context.Context, reportID uint64) (bool, error) {
	ctx, span := tracing.StartSpan(ctx, "ReportExists")
	defer span.End()

	tracing.AttachReportIDToSpan(span, reportID)

	c.logger.WithValues(map[string]interface{}{
		"report_id": reportID,
	}).Debug("ReportExists called")

	return c.querier.ReportExists(ctx, reportID)
}

// GetReport fetches a report from the database.
func (c *Client) GetReport(ctx context.Context, reportID uint64) (*models.Report, error) {
	ctx, span := tracing.StartSpan(ctx, "GetReport")
	defer span.End()

	tracing.AttachReportIDToSpan(span, reportID)

	c.logger.WithValues(map[string]interface{}{
		"report_id": reportID,
	}).Debug("GetReport called")

	return c.querier.GetReport(ctx, reportID)
}

// GetAllReportsCount fetches the count of reports from the database that meet a particular filter.
func (c *Client) GetAllReportsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := tracing.StartSpan(ctx, "GetAllReportsCount")
	defer span.End()

	c.logger.Debug("GetAllReportsCount called")

	return c.querier.GetAllReportsCount(ctx)
}

// GetReports fetches a list of reports from the database that meet a particular filter.
func (c *Client) GetReports(ctx context.Context, filter *models.QueryFilter) (*models.ReportList, error) {
	ctx, span := tracing.StartSpan(ctx, "GetReports")
	defer span.End()

	tracing.AttachFilterToSpan(span, filter)

	c.logger.Debug("GetReports called")

	reportList, err := c.querier.GetReports(ctx, filter)

	return reportList, err
}

// CreateReport creates a report in the database.
func (c *Client) CreateReport(ctx context.Context, input *models.ReportCreationInput) (*models.Report, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateReport")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateReport called")

	return c.querier.CreateReport(ctx, input)
}

// UpdateReport updates a particular report. Note that UpdateReport expects the
// provided input to have a valid ID.
func (c *Client) UpdateReport(ctx context.Context, updated *models.Report) error {
	ctx, span := tracing.StartSpan(ctx, "UpdateReport")
	defer span.End()

	tracing.AttachReportIDToSpan(span, updated.ID)
	c.logger.WithValue("report_id", updated.ID).Debug("UpdateReport called")

	return c.querier.UpdateReport(ctx, updated)
}

// ArchiveReport archives a report from the database by its ID.
func (c *Client) ArchiveReport(ctx context.Context, reportID, userID uint64) error {
	ctx, span := tracing.StartSpan(ctx, "ArchiveReport")
	defer span.End()

	tracing.AttachUserIDToSpan(span, userID)
	tracing.AttachReportIDToSpan(span, reportID)

	c.logger.WithValues(map[string]interface{}{
		"report_id": reportID,
		"user_id":   userID,
	}).Debug("ArchiveReport called")

	return c.querier.ArchiveReport(ctx, reportID, userID)
}
