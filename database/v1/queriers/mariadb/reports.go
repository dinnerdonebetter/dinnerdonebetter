package mariadb

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
func (m *MariaDB) buildGetReportQuery(reportID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(reportsTableColumns...).
		From(reportsTableName).
		Where(squirrel.Eq{
			"id":         reportID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetReport fetches a report from the mariadb database
func (m *MariaDB) GetReport(ctx context.Context, reportID, userID uint64) (*models.Report, error) {
	query, args := m.buildGetReportQuery(reportID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanReport(row)
}

// buildGetReportCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of reports belonging to a given user that meet a given query
func (m *MariaDB) buildGetReportCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
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
	m.logQueryBuildingError(err)

	return query, args
}

// GetReportCount will fetch the count of reports from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetReportCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetReportCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allReportsCountQueryBuilder sync.Once
	allReportsCountQuery        string
)

// buildGetAllReportsCountQuery returns a query that fetches the total number of reports in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllReportsCountQuery() string {
	allReportsCountQueryBuilder.Do(func() {
		var err error
		allReportsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(reportsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allReportsCountQuery
}

// GetAllReportsCount will fetch the count of reports from the database
func (m *MariaDB) GetAllReportsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllReportsCountQuery()).Scan(&count)
	return count, err
}

// buildGetReportsQuery builds a SQL query selecting reports that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetReportsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
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
	m.logQueryBuildingError(err)

	return query, args
}

// GetReports fetches a list of reports from the database that meet a particular filter
func (m *MariaDB) GetReports(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.ReportList, error) {
	query, args := m.buildGetReportsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for reports")
	}

	list, err := scanReports(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetReportCount(ctx, filter, userID)
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
func (m *MariaDB) GetAllReportsForUser(ctx context.Context, userID uint64) ([]models.Report, error) {
	query, args := m.buildGetReportsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching reports for user")
	}

	list, err := scanReports(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateReportQuery takes a report and returns a creation query for that report and the relevant arguments.
func (m *MariaDB) buildCreateReportQuery(input *models.Report) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(reportsTableName).
		Columns(
			"report_type",
			"concern",
			"belongs_to",
			"created_on",
		).
		Values(
			input.ReportType,
			input.Concern,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildReportCreationTimeQuery takes a report and returns a creation query for that report and the relevant arguments
func (m *MariaDB) buildReportCreationTimeQuery(reportID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(reportsTableName).
		Where(squirrel.Eq{"id": reportID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateReport creates a report in the database
func (m *MariaDB) CreateReport(ctx context.Context, input *models.ReportCreationInput) (*models.Report, error) {
	x := &models.Report{
		ReportType: input.ReportType,
		Concern:    input.Concern,
		BelongsTo:  input.BelongsTo,
	}

	query, args := m.buildCreateReportQuery(x)

	// create the report
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing report creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildReportCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateReportQuery takes a report and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateReportQuery(input *models.Report) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(reportsTableName).
		Set("report_type", input.ReportType).
		Set("concern", input.Concern).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateReport updates a particular report. Note that UpdateReport expects the provided input to have a valid ID.
func (m *MariaDB) UpdateReport(ctx context.Context, input *models.Report) error {
	query, args := m.buildUpdateReportQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveReportQuery returns a SQL query which marks a given report belonging to a given user as archived.
func (m *MariaDB) buildArchiveReportQuery(reportID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(reportsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          reportID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveReport marks a report as archived in the database
func (m *MariaDB) ArchiveReport(ctx context.Context, reportID, userID uint64) error {
	query, args := m.buildArchiveReportQuery(reportID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
