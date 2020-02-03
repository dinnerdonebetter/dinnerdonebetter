package dbclient

import (
	"context"
	"strconv"

	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"go.opencensus.io/trace"
)

var _ models.ReportDataManager = (*Client)(nil)

// attachReportIDToSpan provides a consistent way to attach a report's ID to a span
func attachReportIDToSpan(span *trace.Span, reportID uint64) {
	if span != nil {
		span.AddAttributes(trace.StringAttribute("report_id", strconv.FormatUint(reportID, 10)))
	}
}

// GetReport fetches a report from the database
func (c *Client) GetReport(ctx context.Context, reportID, userID uint64) (*models.Report, error) {
	ctx, span := trace.StartSpan(ctx, "GetReport")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachReportIDToSpan(span, reportID)

	c.logger.WithValues(map[string]interface{}{
		"report_id": reportID,
		"user_id":   userID,
	}).Debug("GetReport called")

	return c.querier.GetReport(ctx, reportID, userID)
}

// GetReportCount fetches the count of reports from the database that meet a particular filter
func (c *Client) GetReportCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetReportCount")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetReportCount called")

	return c.querier.GetReportCount(ctx, filter, userID)
}

// GetAllReportsCount fetches the count of reports from the database that meet a particular filter
func (c *Client) GetAllReportsCount(ctx context.Context) (count uint64, err error) {
	ctx, span := trace.StartSpan(ctx, "GetAllReportsCount")
	defer span.End()

	c.logger.Debug("GetAllReportsCount called")

	return c.querier.GetAllReportsCount(ctx)
}

// GetReports fetches a list of reports from the database that meet a particular filter
func (c *Client) GetReports(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.ReportList, error) {
	ctx, span := trace.StartSpan(ctx, "GetReports")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachFilterToSpan(span, filter)

	c.logger.WithValue("user_id", userID).Debug("GetReports called")

	reportList, err := c.querier.GetReports(ctx, filter, userID)

	return reportList, err
}

// GetAllReportsForUser fetches a list of reports from the database that meet a particular filter
func (c *Client) GetAllReportsForUser(ctx context.Context, userID uint64) ([]models.Report, error) {
	ctx, span := trace.StartSpan(ctx, "GetAllReportsForUser")
	defer span.End()

	attachUserIDToSpan(span, userID)
	c.logger.WithValue("user_id", userID).Debug("GetAllReportsForUser called")

	reportList, err := c.querier.GetAllReportsForUser(ctx, userID)

	return reportList, err
}

// CreateReport creates a report in the database
func (c *Client) CreateReport(ctx context.Context, input *models.ReportCreationInput) (*models.Report, error) {
	ctx, span := trace.StartSpan(ctx, "CreateReport")
	defer span.End()

	c.logger.WithValue("input", input).Debug("CreateReport called")

	return c.querier.CreateReport(ctx, input)
}

// UpdateReport updates a particular report. Note that UpdateReport expects the
// provided input to have a valid ID.
func (c *Client) UpdateReport(ctx context.Context, input *models.Report) error {
	ctx, span := trace.StartSpan(ctx, "UpdateReport")
	defer span.End()

	attachReportIDToSpan(span, input.ID)
	c.logger.WithValue("report_id", input.ID).Debug("UpdateReport called")

	return c.querier.UpdateReport(ctx, input)
}

// ArchiveReport archives a report from the database by its ID
func (c *Client) ArchiveReport(ctx context.Context, reportID, userID uint64) error {
	ctx, span := trace.StartSpan(ctx, "ArchiveReport")
	defer span.End()

	attachUserIDToSpan(span, userID)
	attachReportIDToSpan(span, reportID)

	c.logger.WithValues(map[string]interface{}{
		"report_id": reportID,
		"user_id":   userID,
	}).Debug("ArchiveReport called")

	return c.querier.ArchiveReport(ctx, reportID, userID)
}
