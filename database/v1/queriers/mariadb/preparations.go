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
		"belongs_to",
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
		&x.BelongsTo,
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
func (m *MariaDB) buildGetPreparationQuery(preparationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(preparationsTableColumns...).
		From(preparationsTableName).
		Where(squirrel.Eq{
			"id":         preparationID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetPreparation fetches a preparation from the mariadb database
func (m *MariaDB) GetPreparation(ctx context.Context, preparationID, userID uint64) (*models.Preparation, error) {
	query, args := m.buildGetPreparationQuery(preparationID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanPreparation(row)
}

// buildGetPreparationCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of preparations belonging to a given user that meet a given query
func (m *MariaDB) buildGetPreparationCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(preparationsTableName).
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

// GetPreparationCount will fetch the count of preparations from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetPreparationCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetPreparationCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allPreparationsCountQueryBuilder sync.Once
	allPreparationsCountQuery        string
)

// buildGetAllPreparationsCountQuery returns a query that fetches the total number of preparations in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllPreparationsCountQuery() string {
	allPreparationsCountQueryBuilder.Do(func() {
		var err error
		allPreparationsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(preparationsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allPreparationsCountQuery
}

// GetAllPreparationsCount will fetch the count of preparations from the database
func (m *MariaDB) GetAllPreparationsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllPreparationsCountQuery()).Scan(&count)
	return count, err
}

// buildGetPreparationsQuery builds a SQL query selecting preparations that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetPreparationsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(preparationsTableColumns...).
		From(preparationsTableName).
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

// GetPreparations fetches a list of preparations from the database that meet a particular filter
func (m *MariaDB) GetPreparations(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.PreparationList, error) {
	query, args := m.buildGetPreparationsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for preparations")
	}

	list, err := scanPreparations(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetPreparationCount(ctx, filter, userID)
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
func (m *MariaDB) GetAllPreparationsForUser(ctx context.Context, userID uint64) ([]models.Preparation, error) {
	query, args := m.buildGetPreparationsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching preparations for user")
	}

	list, err := scanPreparations(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreatePreparationQuery takes a preparation and returns a creation query for that preparation and the relevant arguments.
func (m *MariaDB) buildCreatePreparationQuery(input *models.Preparation) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(preparationsTableName).
		Columns(
			"name",
			"variant",
			"description",
			"allergy_warning",
			"icon",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Name,
			input.Variant,
			input.Description,
			input.AllergyWarning,
			input.Icon,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildPreparationCreationTimeQuery takes a preparation and returns a creation query for that preparation and the relevant arguments
func (m *MariaDB) buildPreparationCreationTimeQuery(preparationID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(preparationsTableName).
		Where(squirrel.Eq{"id": preparationID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreatePreparation creates a preparation in the database
func (m *MariaDB) CreatePreparation(ctx context.Context, input *models.PreparationCreationInput) (*models.Preparation, error) {
	x := &models.Preparation{
		Name:           input.Name,
		Variant:        input.Variant,
		Description:    input.Description,
		AllergyWarning: input.AllergyWarning,
		Icon:           input.Icon,
		BelongsTo:      input.BelongsTo,
	}

	query, args := m.buildCreatePreparationQuery(x)

	// create the preparation
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing preparation creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildPreparationCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdatePreparationQuery takes a preparation and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdatePreparationQuery(input *models.Preparation) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(preparationsTableName).
		Set("name", input.Name).
		Set("variant", input.Variant).
		Set("description", input.Description).
		Set("allergy_warning", input.AllergyWarning).
		Set("icon", input.Icon).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdatePreparation updates a particular preparation. Note that UpdatePreparation expects the provided input to have a valid ID.
func (m *MariaDB) UpdatePreparation(ctx context.Context, input *models.Preparation) error {
	query, args := m.buildUpdatePreparationQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchivePreparationQuery returns a SQL query which marks a given preparation belonging to a given user as archived.
func (m *MariaDB) buildArchivePreparationQuery(preparationID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(preparationsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          preparationID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchivePreparation marks a preparation as archived in the database
func (m *MariaDB) ArchivePreparation(ctx context.Context, preparationID, userID uint64) error {
	query, args := m.buildArchivePreparationQuery(preparationID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
