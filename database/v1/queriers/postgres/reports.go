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
	reportsTableName             = "reports"
	reportsTableReportTypeColumn = "report_type"
	reportsTableConcernColumn    = "concern"
	reportsUserOwnershipColumn   = "belongs_to_user"
)

var (
	reportsTableColumns = []string{
		fmt.Sprintf("%s.%s", reportsTableName, idColumn),
		fmt.Sprintf("%s.%s", reportsTableName, reportsTableReportTypeColumn),
		fmt.Sprintf("%s.%s", reportsTableName, reportsTableConcernColumn),
		fmt.Sprintf("%s.%s", reportsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", reportsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", reportsTableName, archivedOnColumn),
		fmt.Sprintf("%s.%s", reportsTableName, reportsUserOwnershipColumn),
	}
)

// scanReport takes a database Scanner (i.e. *sql.Row) and scans the result into a Report struct
func (p *Postgres) scanReport(scan database.Scanner) (*models.Report, error) {
	x := &models.Report{}

	targetVars := []interface{}{
		&x.ID,
		&x.ReportType,
		&x.Concern,
		&x.CreatedOn,
		&x.LastUpdatedOn,
		&x.ArchivedOn,
		&x.BelongsToUser,
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, err
	}

	return x, nil
}

// scanReports takes a logger and some database rows and turns them into a slice of reports.
func (p *Postgres) scanReports(rows database.ResultIterator) ([]models.Report, error) {
	var (
		list []models.Report
	)

	for rows.Next() {
		x, err := p.scanReport(rows)
		if err != nil {
			return nil, err
		}

		list = append(list, *x)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if closeErr := rows.Close(); closeErr != nil {
		p.logger.Error(closeErr, "closing database rows")
	}

	return list, nil
}

//
func (p *Postgres) buildReportExistsQuery(reportID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.%s", reportsTableName, idColumn)).
		Prefix(existencePrefix).
		From(reportsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", reportsTableName, idColumn): reportID,
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
			fmt.Sprintf("%s.%s", reportsTableName, idColumn): reportID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetReport fetches a report from the database.
func (p *Postgres) GetReport(ctx context.Context, reportID uint64) (*models.Report, error) {
	query, args := p.buildGetReportQuery(reportID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return p.scanReport(row)
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
				fmt.Sprintf("%s.%s", reportsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfReportsQuery returns a query that fetches every report in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfReportsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(reportsTableColumns...).
		From(reportsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", reportsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", reportsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllReports fetches every report from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllReports(ctx context.Context, resultChannel chan []models.Report) error {
	count, err := p.GetAllReportsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfReportsQuery(begin, end)
			logger := p.logger.WithValues(map[string]interface{}{
				"query": query,
				"begin": begin,
				"end":   end,
			})

			rows, err := p.db.Query(query, args...)
			if err == sql.ErrNoRows {
				return
			} else if err != nil {
				logger.Error(err, "querying for database rows")
				return
			}

			reports, err := p.scanReports(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- reports
		}(beginID, endID)
	}

	return nil
}

// buildGetReportsQuery builds a SQL query selecting reports that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetReportsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(reportsTableColumns...).
		From(reportsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", reportsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", reportsTableName, idColumn))

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

	reports, err := p.scanReports(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ReportList{
		Pagination: models.Pagination{
			Page:  filter.Page,
			Limit: filter.Limit,
		},
		Reports: reports,
	}

	return list, nil
}

// buildGetReportsWithIDsQuery builds a SQL query selecting reports
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetReportsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(reportsTableColumns...).
		From(reportsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(reportsTableColumns...).
		FromSelect(subqueryBuilder, reportsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", reportsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetReportsWithIDs fetches a list of reports from the database that exist within a given set of IDs.
func (p *Postgres) GetReportsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.Report, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetReportsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for reports")
	}

	reports, err := p.scanReports(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return reports, nil
}

// buildCreateReportQuery takes a report and returns a creation query for that report and the relevant arguments.
func (p *Postgres) buildCreateReportQuery(input *models.Report) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(reportsTableName).
		Columns(
			reportsTableReportTypeColumn,
			reportsTableConcernColumn,
			reportsUserOwnershipColumn,
		).
		Values(
			input.ReportType,
			input.Concern,
			input.BelongsToUser,
		).
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
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
		Set(reportsTableReportTypeColumn, input.ReportType).
		Set(reportsTableConcernColumn, input.Concern).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                   input.ID,
			reportsUserOwnershipColumn: input.BelongsToUser,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateReport updates a particular report. Note that UpdateReport expects the provided input to have a valid ID.
func (p *Postgres) UpdateReport(ctx context.Context, input *models.Report) error {
	query, args := p.buildUpdateReportQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveReportQuery returns a SQL query which marks a given report belonging to a given user as archived.
func (p *Postgres) buildArchiveReportQuery(reportID, userID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(reportsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:                   reportID,
			archivedOnColumn:           nil,
			reportsUserOwnershipColumn: userID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
