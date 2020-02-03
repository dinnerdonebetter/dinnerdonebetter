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
	iterationMediasTableName = "iteration_medias"
)

var (
	iterationMediasTableColumns = []string{
		"id",
		"path",
		"mimetype",
		"recipe_iteration_id",
		"recipe_step_id",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanIterationMedia takes a database Scanner (i.e. *sql.Row) and scans the result into an Iteration Media struct
func scanIterationMedia(scan database.Scanner) (*models.IterationMedia, error) {
	x := &models.IterationMedia{}

	if err := scan.Scan(
		&x.ID,
		&x.Path,
		&x.Mimetype,
		&x.RecipeIterationID,
		&x.RecipeStepID,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
		&x.BelongsTo,
	); err != nil {
		return nil, err
	}

	return x, nil
}

// scanIterationMedias takes a logger and some database rows and turns them into a slice of iteration medias
func scanIterationMedias(logger logging.Logger, rows *sql.Rows) ([]models.IterationMedia, error) {
	var list []models.IterationMedia

	for rows.Next() {
		x, err := scanIterationMedia(rows)
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

// buildGetIterationMediaQuery constructs a SQL query for fetching an iteration media with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetIterationMediaQuery(iterationMediaID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(iterationMediasTableColumns...).
		From(iterationMediasTableName).
		Where(squirrel.Eq{
			"id":         iterationMediaID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetIterationMedia fetches an iteration media from the mariadb database
func (m *MariaDB) GetIterationMedia(ctx context.Context, iterationMediaID, userID uint64) (*models.IterationMedia, error) {
	query, args := m.buildGetIterationMediaQuery(iterationMediaID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanIterationMedia(row)
}

// buildGetIterationMediaCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of iteration medias belonging to a given user that meet a given query
func (m *MariaDB) buildGetIterationMediaCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(iterationMediasTableName).
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

// GetIterationMediaCount will fetch the count of iteration medias from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetIterationMediaCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetIterationMediaCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allIterationMediasCountQueryBuilder sync.Once
	allIterationMediasCountQuery        string
)

// buildGetAllIterationMediasCountQuery returns a query that fetches the total number of iteration medias in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllIterationMediasCountQuery() string {
	allIterationMediasCountQueryBuilder.Do(func() {
		var err error
		allIterationMediasCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(iterationMediasTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allIterationMediasCountQuery
}

// GetAllIterationMediasCount will fetch the count of iteration medias from the database
func (m *MariaDB) GetAllIterationMediasCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllIterationMediasCountQuery()).Scan(&count)
	return count, err
}

// buildGetIterationMediasQuery builds a SQL query selecting iteration medias that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetIterationMediasQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(iterationMediasTableColumns...).
		From(iterationMediasTableName).
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

// GetIterationMedias fetches a list of iteration medias from the database that meet a particular filter
func (m *MariaDB) GetIterationMedias(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.IterationMediaList, error) {
	query, args := m.buildGetIterationMediasQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for iteration medias")
	}

	list, err := scanIterationMedias(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetIterationMediaCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching iteration media count: %w", err)
	}

	x := &models.IterationMediaList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		IterationMedias: list,
	}

	return x, nil
}

// GetAllIterationMediasForUser fetches every iteration media belonging to a user
func (m *MariaDB) GetAllIterationMediasForUser(ctx context.Context, userID uint64) ([]models.IterationMedia, error) {
	query, args := m.buildGetIterationMediasQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching iteration medias for user")
	}

	list, err := scanIterationMedias(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateIterationMediaQuery takes an iteration media and returns a creation query for that iteration media and the relevant arguments.
func (m *MariaDB) buildCreateIterationMediaQuery(input *models.IterationMedia) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(iterationMediasTableName).
		Columns(
			"path",
			"mimetype",
			"recipe_iteration_id",
			"recipe_step_id",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Path,
			input.Mimetype,
			input.RecipeIterationID,
			input.RecipeStepID,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildIterationMediaCreationTimeQuery takes an iteration media and returns a creation query for that iteration media and the relevant arguments
func (m *MariaDB) buildIterationMediaCreationTimeQuery(iterationMediaID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(iterationMediasTableName).
		Where(squirrel.Eq{"id": iterationMediaID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateIterationMedia creates an iteration media in the database
func (m *MariaDB) CreateIterationMedia(ctx context.Context, input *models.IterationMediaCreationInput) (*models.IterationMedia, error) {
	x := &models.IterationMedia{
		Path:              input.Path,
		Mimetype:          input.Mimetype,
		RecipeIterationID: input.RecipeIterationID,
		RecipeStepID:      input.RecipeStepID,
		BelongsTo:         input.BelongsTo,
	}

	query, args := m.buildCreateIterationMediaQuery(x)

	// create the iteration media
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing iteration media creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildIterationMediaCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateIterationMediaQuery takes an iteration media and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateIterationMediaQuery(input *models.IterationMedia) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(iterationMediasTableName).
		Set("path", input.Path).
		Set("mimetype", input.Mimetype).
		Set("recipe_iteration_id", input.RecipeIterationID).
		Set("recipe_step_id", input.RecipeStepID).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":         input.ID,
			"belongs_to": input.BelongsTo,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// UpdateIterationMedia updates a particular iteration media. Note that UpdateIterationMedia expects the provided input to have a valid ID.
func (m *MariaDB) UpdateIterationMedia(ctx context.Context, input *models.IterationMedia) error {
	query, args := m.buildUpdateIterationMediaQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveIterationMediaQuery returns a SQL query which marks a given iteration media belonging to a given user as archived.
func (m *MariaDB) buildArchiveIterationMediaQuery(iterationMediaID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(iterationMediasTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          iterationMediaID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveIterationMedia marks an iteration media as archived in the database
func (m *MariaDB) ArchiveIterationMedia(ctx context.Context, iterationMediaID, userID uint64) error {
	query, args := m.buildArchiveIterationMediaQuery(iterationMediaID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
