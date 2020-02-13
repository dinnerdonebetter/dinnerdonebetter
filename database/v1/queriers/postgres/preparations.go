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
	preparationsTableName = "preparations"
)

var (
	preparationsTableColumns = []string{
		"id",
		"name",
		"variant",
		"description",
		"allergy_warning",
		"icon",
		"created_on",
		"updated_on",
		"archived_on",
	}
)

// scanPreparation takes a database Scanner (i.e. *sql.Row) and scans the result into a Preparation struct
func scanPreparation(scan database.Scanner) (*models.Preparation, error) {
	x := &models.Preparation{}

	if err := scan.Scan(
		&x.ID,
		&x.Name,
		&x.Variant,
		&x.Description,
		&x.AllergyWarning,
		&x.Icon,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanPreparations takes a logger and some database rows and turns them into a slice of preparations
func scanPreparations(logger logging.Logger, rows *sql.Rows) ([]models.Preparation, error) {
	var list []models.Preparation

	for rows.Next() {
		x, err := scanPreparation(rows)
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

// buildGetPreparationQuery constructs a SQL query for fetching a preparation with a given ID belong to a user with a given ID.
func (p *Postgres) buildGetPreparationQuery(preparationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Select(preparationsTableColumns...).
		From(preparationsTableName).
		Where(squirrel.Eq{
			"id": preparationID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetPreparation fetches a preparation from the postgres database
func (p *Postgres) GetPreparation(ctx context.Context, preparationID, userID uint64) (*models.Preparation, error) {
	query, args := p.buildGetPreparationQuery(preparationID, userID)
	row := p.db.QueryRowContext(ctx, query, args...)
	return scanPreparation(row)
}

// buildGetPreparationCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of preparations belonging to a given user that meet a given query
func (p *Postgres) buildGetPreparationCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(CountQuery).
		From(preparationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetPreparationCount will fetch the count of preparations from the database that meet a particular filter and belong to a particular user.
func (p *Postgres) GetPreparationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := p.buildGetPreparationCountQuery(filter, userID)
	err = p.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allPreparationsCountQueryBuilder sync.Once
	allPreparationsCountQuery        string
)

// buildGetAllPreparationsCountQuery returns a query that fetches the total number of preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllPreparationsCountQuery() string {
	allPreparationsCountQueryBuilder.Do(func() {
		var err error
		allPreparationsCountQuery, _, err = p.sqlBuilder.
			Select(CountQuery).
			From(preparationsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allPreparationsCountQuery
}

// GetAllPreparationsCount will fetch the count of preparations from the database
func (p *Postgres) GetAllPreparationsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllPreparationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetPreparationsQuery builds a SQL query selecting preparations that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetPreparationsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := p.sqlBuilder.
		Select(preparationsTableColumns...).
		From(preparationsTableName).
		Where(squirrel.Eq{
			"archived_on": nil,
		})

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetPreparations fetches a list of preparations from the database that meet a particular filter
func (p *Postgres) GetPreparations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.PreparationList, error) {
	query, args := p.buildGetPreparationsQuery(filter, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for preparations")
	}

	list, err := scanPreparations(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := p.GetPreparationCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching preparation count: %w", err)
	}

	x := &models.PreparationList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Preparations: list,
	}

	return x, nil
}

// GetAllPreparationsForUser fetches every preparation belonging to a user
func (p *Postgres) GetAllPreparationsForUser(ctx context.Context, userID uint64) ([]models.Preparation, error) {
	query, args := p.buildGetPreparationsQuery(nil, userID)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching preparations for user")
	}

	list, err := scanPreparations(p.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreatePreparationQuery takes a preparation and returns a creation query for that preparation and the relevant arguments.
func (p *Postgres) buildCreatePreparationQuery(input *models.Preparation) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Insert(preparationsTableName).
		Columns(
			"name",
			"variant",
			"description",
			"allergy_warning",
			"icon",
		).
		Values(
			input.Name,
			input.Variant,
			input.Description,
			input.AllergyWarning,
			input.Icon,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreatePreparation creates a preparation in the database
func (p *Postgres) CreatePreparation(ctx context.Context, input *models.PreparationCreationInput) (*models.Preparation, error) {
	x := &models.Preparation{
		Name:           input.Name,
		Variant:        input.Variant,
		Description:    input.Description,
		AllergyWarning: input.AllergyWarning,
		Icon:           input.Icon,
	}

	query, args := p.buildCreatePreparationQuery(x)

	// create the preparation
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing preparation creation query: %w", err)
	}

	return x, nil
}

// buildUpdatePreparationQuery takes a preparation and returns an update SQL query, with the relevant query parameters
func (p *Postgres) buildUpdatePreparationQuery(input *models.Preparation) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(preparationsTableName).
		Set("name", input.Name).
		Set("variant", input.Variant).
		Set("description", input.Description).
		Set("allergy_warning", input.AllergyWarning).
		Set("icon", input.Icon).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdatePreparation updates a particular preparation. Note that UpdatePreparation expects the provided input to have a valid ID.
func (p *Postgres) UpdatePreparation(ctx context.Context, input *models.Preparation) error {
	query, args := p.buildUpdatePreparationQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchivePreparationQuery returns a SQL query which marks a given preparation belonging to a given user as archived.
func (p *Postgres) buildArchivePreparationQuery(preparationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = p.sqlBuilder.
		Update(preparationsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          preparationID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchivePreparation marks a preparation as archived in the database
func (p *Postgres) ArchivePreparation(ctx context.Context, preparationID, userID uint64) error {
	query, args := p.buildArchivePreparationQuery(preparationID, userID)
	_, err := p.db.ExecContext(ctx, query, args...)
	return err
}
