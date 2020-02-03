package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1"
)

const (
	reportsTableName = "reports"
)

var (
	reportsTableColumns = []string{
		"id",
		"report_type",
		"concern",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanReport takes a database Scanner (i.e. *sql.Row) and scans the result into a Report struct
func scanReport(scan database.Scanner) (*models.Report, error) {
	x := &models.Report{}

	if err := scan.Scan(
		&x.ID,
		&x.ReportType,
		&x.Concern,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanReports takes a logger and some database rows and turns them into a slice of reports
func scanReports(logger logging.Logger, rows *sql.Rows) ([]models.Report, error) {
	var list []models.Report

	for rows.Next() {
		x, err := scanReport(rows)
		if err != nil {
			return nil, err
		}
		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

// buildGetReportQuery constructs a SQL query for fetching a report with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetReportQuery(reportID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(reportsTableColumns...).
		From(reportsTableName).
		Where(squirrel.Eq{
			"id":         reportID,
			"belongs_to": userID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetReport fetches a report from the postgres database
func (p *Postgres) GetReport(ctx context.Context, reportID, userID uint64) (*models.Report, error) {
	query, args := p.buildGetReportQuery(reportID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanReport(row)
}

// buildGetReportCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of reports belonging to a given user that meet a given query
func (p *Postgres) buildGetReportCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(reportsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetReportCount will fetch the count of reports from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetReportCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetReportCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allReportsCountQueryBuilder sync.Once
	allReportsCountQuery        string
)

// buildGetAllReportsCountQuery returns a query that fetches the total number of reports in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllReportsCountQuery() string {
	allReportsCountQueryBuilder.Do(func() {
		var err error
		allReportsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(reportsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allReportsCountQuery
}

// GetAllReportsCount will fetch the count of reports from the database
func (p *Postgres) GetAllReportsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllReportsCountQuery()).Scan(&count)
	return count, err
}

// buildGetReportsQuery builds a SQL query selecting reports that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetReportsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(reportsTableColumns...).
		From(reportsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
			"belongs_to":  userID,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetReports fetches a list of reports from the database that meet a particular filter
func (p *Postgres) GetReports(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.ReportList, error) {
	query, args := p.buildGetReportsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for reports")
	}

	list, err := scanReports(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetReportCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching report count: %w", err)
	}

	x := &models.ReportList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Reports: list,
	}

	return x, nil
}

// GetAllReportsForUser fetches every report belonging to a user
func (p *Postgres) GetAllReportsForUser(ctx context.Context, userID uint64) ([]models.Report, error) {
	query, args := p.buildGetReportsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching reports for user")
	}

	list, err := scanReports(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateReportQuery takes a report and returns a creation query for that report and the relevant arguments.
func (p *Postgres) buildCreateReportQuery(input *models.Report) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(reportsTableName).
		Columns(
			"report_type",
			"concern",
			"belongs_to",
		).
		Values(
			input.ReportType,
			input.Concern,
			input.BelongsTo,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateReport creates a report in the database
func (p *Postgres) CreateReport(ctx context.Context, input *models.ReportCreationInput) (*models.Report, error) {
	x := &models.Report{
		ReportType: input.ReportType,
		Concern:    input.Concern,
		BelongsTo:  input.BelongsTo,
	}

	query, args := p.buildCreateReportQuery(x)

	// create the report
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing report creation query: %w", err)
	}

	return x, nil
}

// buildUpdateReportQuery takes a report and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdateReportQuery(input *models.Report) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(reportsTableName).
		Set("report_type", input.ReportType).
		Set("concern", input.Concern).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateReport updates a particular report. Note that UpdateReport expects the provided input to have a valid ID.
func (p *Postgres) UpdateReport(ctx context.Context, input *models.Report) error {
	query, args := p.buildUpdateReportQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveReportQuery returns a SQL query which marks a given report belonging to a given user as archived.
func (p *Postgres) buildArchiveReportQuery(reportID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(reportsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          reportID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveReport marks a report as archived in the database
func (p *Postgres) ArchiveReport(ctx context.Context, reportID, userID uint64) error {
	query, args := p.buildArchiveReportQuery(reportID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
