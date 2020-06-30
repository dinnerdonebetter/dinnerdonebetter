package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	database "gitlab.com/prixfixe/prixfixe/database/v1"
	models "gitlab.com/prixfixe/prixfixe/models/v1"

	"github.com/Masterminds/squirrel"
)

const (
	reportsTableName           = "reports"
	reportsUserOwnershipColumn = "belongs_to_user"
)

var (
	reportsTableColumns = []string{
		fmt.Sprintf("%s.%s", reportsTableName, "id"),
		fmt.Sprintf("%s.%s", reportsTableName, "report_type"),
		fmt.Sprintf("%s.%s", reportsTableName, "concern"),
		fmt.Sprintf("%s.%s", reportsTableName, "created_on"),
		fmt.Sprintf("%s.%s", reportsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", reportsTableName, "archived_on"),
		fmt.Sprintf("%s.%s", reportsTableName, reportsUserOwnershipColumn),
	}
)

// scanReport takes a database Scanner (i.e. *sql.Row) and scans the result into a Report struct
func (p *Postgres) scanReport(scan database.Scanner, includeCount bool) (*models.Report, uint64, error) {
	x := &models.Report{}
	var count uint64

	targetVars := []interface{}{
		&x.ID,
		&x.ReportType,
		&x.Concern,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanReports takes a logger and some database rows and turns them into a slice of reports.
func (p *Postgres) scanReports(rows database.ResultIterator) ([]models.Report, uint64, error) {
	var (
		list  []models.Report
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanReport(rows, true)
		if err != nil {
			return nil, 0, err
		}

		if count == 0 {
			count = c
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, count, nil
}

//
func (p *Postgres) buildReportExistsQuery(reportID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", reportsTableName)).
		Prefix(existencePrefix).
		From(reportsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", reportsTableName): reportID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ReportExists queries the database to see if a given report belonging to a given user exists.
func (p *Postgres) ReportExists(ctx context.Context, reportID uint64) (exists bool, err error) {
	query, args := p.buildReportExistsQuery(reportID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

//
func (p *Postgres) buildGetReportQuery(reportID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(reportsTableColumns...).
		From(reportsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", reportsTableName): reportID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetReport fetches a report from the database.
func (p *Postgres) GetReport(ctx context.Context, reportID uint64) (*models.Report, error) {
	query, args := p.buildGetReportQuery(reportID)
	row := p.db.QueryRowContext(ctx, query, args...)

	report, _, err := p.scanReport(row, false)
	return report, err
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
			Select(fmt.Sprintf(countQuery, reportsTableName)).
			From(reportsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", reportsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allReportsCountQuery
}

// GetAllReportsCount will fetch the count of reports from the database.
func (p *Postgres) GetAllReportsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllReportsCountQuery()).Scan(&count)
	return count, err
}

// buildGetReportsQuery builds a SQL query selecting reports that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetReportsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(reportsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllReportsCountQuery()))...).
		From(reportsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", reportsTableName): nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", reportsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, reportsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetReports fetches a list of reports from the database that meet a particular filter.
func (p *Postgres) GetReports(ctx context.Context, filter *models.QueryFilter) (*models.ReportList, error) {
	query, args := p.buildGetReportsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for reports")
	}

	reports, count, err := p.scanReports(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ReportList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Reports: reports,
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
			reportsUserOwnershipColumn,
		).
		Values(
			input.ReportType,
			input.Concern,
			input.BelongsToUser,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateReport creates a report in the database.
func (p *Postgres) CreateReport(ctx context.Context, input *models.ReportCreationInput) (*models.Report, error) {
	x := &models.Report{
		ReportType:    input.ReportType,
		Concern:       input.Concern,
		BelongsToUser: input.BelongsToUser,
	}

	query, args := p.buildCreateReportQuery(x)

	// create the report.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing report creation query: %w", err)
	}

	return x, nil
}

// buildUpdateReportQuery takes a report and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateReportQuery(input *models.Report) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(reportsTableName).
		Set("report_type", input.ReportType).
		Set("concern", input.Concern).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                       input.ID,
			reportsUserOwnershipColumn: input.BelongsToUser,
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
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":                       reportID,
			"archived_on":              nil,
			reportsUserOwnershipColumn: userID,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveReport marks a report as archived in the database.
func (p *Postgres) ArchiveReport(ctx context.Context, reportID, userID uint64) error {
	query, args := p.buildArchiveReportQuery(reportID, userID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
