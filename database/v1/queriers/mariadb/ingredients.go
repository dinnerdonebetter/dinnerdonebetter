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
	ingredientsTableName = "ingredients"
)

var (
	ingredientsTableColumns = []string{
		"id",
		"name",
		"variant",
		"description",
		"warning",
		"contains_egg",
		"contains_dairy",
		"contains_peanut",
		"contains_tree_nut",
		"contains_soy",
		"contains_wheat",
		"contains_shellfish",
		"contains_sesame",
		"contains_fish",
		"contains_gluten",
		"animal_flesh",
		"animal_derived",
		"considered_staple",
		"icon",
		"created_on",
		"updated_on",
		"archived_on",
		"belongs_to",
	}
)

// scanIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into an Ingredient struct
func scanIngredient(scan database.Scanner) (*models.Ingredient, error) {
	x := &models.Ingredient{}

	if err := scan.Scan(
		&x.ID,
		&x.Name,
		&x.Variant,
		&x.Description,
		&x.Warning,
		&x.ContainsEgg,
		&x.ContainsDairy,
		&x.ContainsPeanut,
		&x.ContainsTreeNut,
		&x.ContainsSoy,
		&x.ContainsWheat,
		&x.ContainsShellfish,
		&x.ContainsSesame,
		&x.ContainsFish,
		&x.ContainsGluten,
		&x.AnimalFlesh,
		&x.AnimalDerived,
		&x.ConsideredStaple,
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

// scanIngredients takes a logger and some database rows and turns them into a slice of ingredients
func scanIngredients(logger logging.Logger, rows *sql.Rows) ([]models.Ingredient, error) {
	var list []models.Ingredient

	for rows.Next() {
		x, err := scanIngredient(rows)
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

// buildGetIngredientQuery constructs a SQL query for fetching an ingredient with a given ID belong to a user with a given ID.
func (m *MariaDB) buildGetIngredientQuery(ingredientID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Select(ingredientsTableColumns...).
		From(ingredientsTableName).
		Where(squirrel.Eq{
			"id":         ingredientID,
			"belongs_to": userID,
		}).ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// GetIngredient fetches an ingredient from the mariadb database
func (m *MariaDB) GetIngredient(ctx context.Context, ingredientID, userID uint64) (*models.Ingredient, error) {
	query, args := m.buildGetIngredientQuery(ingredientID, userID)
	row := m.db.QueryRowContext(ctx, query, args...)
	return scanIngredient(row)
}

// buildGetIngredientCountQuery takes a QueryFilter and a user ID and returns a SQL query (and the relevant arguments) for
// fetching the number of ingredients belonging to a given user that meet a given query
func (m *MariaDB) buildGetIngredientCountQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(CountQuery).
		From(ingredientsTableName).
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

// GetIngredientCount will fetch the count of ingredients from the database that meet a particular filter and belong to a particular user.
func (m *MariaDB) GetIngredientCount(ctx context.Context, filter *models.QueryFilter, userID uint64) (count uint64, err error) {
	query, args := m.buildGetIngredientCountQuery(filter, userID)
	err = m.db.QueryRowContext(ctx, query, args...).Scan(&count)
	return count, err
}

var (
	allIngredientsCountQueryBuilder sync.Once
	allIngredientsCountQuery        string
)

// buildGetAllIngredientsCountQuery returns a query that fetches the total number of ingredients in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (m *MariaDB) buildGetAllIngredientsCountQuery() string {
	allIngredientsCountQueryBuilder.Do(func() {
		var err error
		allIngredientsCountQuery, _, err = m.sqlBuilder.
			Select(CountQuery).
			From(ingredientsTableName).
			Where(squirrel.Eq{"archived_on": nil}).
			ToSql()
		m.logQueryBuildingError(err)
	})

	return allIngredientsCountQuery
}

// GetAllIngredientsCount will fetch the count of ingredients from the database
func (m *MariaDB) GetAllIngredientsCount(ctx context.Context) (count uint64, err error) {
	err = m.db.QueryRowContext(ctx, m.buildGetAllIngredientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetIngredientsQuery builds a SQL query selecting ingredients that adhere to a given QueryFilter and belong to a given user,
// and returns both the query and the relevant args to pass to the query executor.
func (m *MariaDB) buildGetIngredientsQuery(filter *models.QueryFilter, userID uint64) (query string, args []interface{}) {
	var err error
	builder := m.sqlBuilder.
		Select(ingredientsTableColumns...).
		From(ingredientsTableName).
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

// GetIngredients fetches a list of ingredients from the database that meet a particular filter
func (m *MariaDB) GetIngredients(ctx context.Context, filter *models.QueryFilter, userID uint64) (*models.IngredientList, error) {
	query, args := m.buildGetIngredientsQuery(filter, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for ingredients")
	}

	list, err := scanIngredients(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	count, err := m.GetIngredientCount(ctx, filter, userID)
	if err != nil {
		return nil, fmt.Errorf("fetching ingredient count: %w", err)
	}

	x := &models.IngredientList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		Ingredients: list,
	}

	return x, nil
}

// GetAllIngredientsForUser fetches every ingredient belonging to a user
func (m *MariaDB) GetAllIngredientsForUser(ctx context.Context, userID uint64) ([]models.Ingredient, error) {
	query, args := m.buildGetIngredientsQuery(nil, userID)

	rows, err := m.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "fetching ingredients for user")
	}

	list, err := scanIngredients(m.logger, rows)
	if err != nil {
		return nil, fmt.Errorf("parsing database results: %w", err)
	}

	return list, nil
}

// buildCreateIngredientQuery takes an ingredient and returns a creation query for that ingredient and the relevant arguments.
func (m *MariaDB) buildCreateIngredientQuery(input *models.Ingredient) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Insert(ingredientsTableName).
		Columns(
			"name",
			"variant",
			"description",
			"warning",
			"contains_egg",
			"contains_dairy",
			"contains_peanut",
			"contains_tree_nut",
			"contains_soy",
			"contains_wheat",
			"contains_shellfish",
			"contains_sesame",
			"contains_fish",
			"contains_gluten",
			"animal_flesh",
			"animal_derived",
			"considered_staple",
			"icon",
			"belongs_to",
			"created_on",
		).
		Values(
			input.Name,
			input.Variant,
			input.Description,
			input.Warning,
			input.ContainsEgg,
			input.ContainsDairy,
			input.ContainsPeanut,
			input.ContainsTreeNut,
			input.ContainsSoy,
			input.ContainsWheat,
			input.ContainsShellfish,
			input.ContainsSesame,
			input.ContainsFish,
			input.ContainsGluten,
			input.AnimalFlesh,
			input.AnimalDerived,
			input.ConsideredStaple,
			input.Icon,
			input.BelongsTo,
			squirrel.Expr(CurrentUnixTimeQuery),
		).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// buildIngredientCreationTimeQuery takes an ingredient and returns a creation query for that ingredient and the relevant arguments
func (m *MariaDB) buildIngredientCreationTimeQuery(ingredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = m.sqlBuilder.
		Select("created_on").
		From(ingredientsTableName).
		Where(squirrel.Eq{"id": ingredientID}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// CreateIngredient creates an ingredient in the database
func (m *MariaDB) CreateIngredient(ctx context.Context, input *models.IngredientCreationInput) (*models.Ingredient, error) {
	x := &models.Ingredient{
		Name:              input.Name,
		Variant:           input.Variant,
		Description:       input.Description,
		Warning:           input.Warning,
		ContainsEgg:       input.ContainsEgg,
		ContainsDairy:     input.ContainsDairy,
		ContainsPeanut:    input.ContainsPeanut,
		ContainsTreeNut:   input.ContainsTreeNut,
		ContainsSoy:       input.ContainsSoy,
		ContainsWheat:     input.ContainsWheat,
		ContainsShellfish: input.ContainsShellfish,
		ContainsSesame:    input.ContainsSesame,
		ContainsFish:      input.ContainsFish,
		ContainsGluten:    input.ContainsGluten,
		AnimalFlesh:       input.AnimalFlesh,
		AnimalDerived:     input.AnimalDerived,
		ConsideredStaple:  input.ConsideredStaple,
		Icon:              input.Icon,
		BelongsTo:         input.BelongsTo,
	}

	query, args := m.buildCreateIngredientQuery(x)

	// create the ingredient
	res, err := m.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("error executing ingredient creation query: %w", err)
	}

	// fetch the last inserted ID
	id, idErr := res.LastInsertId()
	if idErr == nil {
		x.ID = uint64(id)

		query, args := m.buildIngredientCreationTimeQuery(x.ID)
		m.logCreationTimeRetrievalError(m.db.QueryRowContext(ctx, query, args...).Scan(&x.CreatedOn))
	}

	return x, nil
}

// buildUpdateIngredientQuery takes an ingredient and returns an update SQL query, with the relevant query parameters
func (m *MariaDB) buildUpdateIngredientQuery(input *models.Ingredient) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(ingredientsTableName).
		Set("name", input.Name).
		Set("variant", input.Variant).
		Set("description", input.Description).
		Set("warning", input.Warning).
		Set("contains_egg", input.ContainsEgg).
		Set("contains_dairy", input.ContainsDairy).
		Set("contains_peanut", input.ContainsPeanut).
		Set("contains_tree_nut", input.ContainsTreeNut).
		Set("contains_soy", input.ContainsSoy).
		Set("contains_wheat", input.ContainsWheat).
		Set("contains_shellfish", input.ContainsShellfish).
		Set("contains_sesame", input.ContainsSesame).
		Set("contains_fish", input.ContainsFish).
		Set("contains_gluten", input.ContainsGluten).
		Set("animal_flesh", input.AnimalFlesh).
		Set("animal_derived", input.AnimalDerived).
		Set("considered_staple", input.ConsideredStaple).
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

// UpdateIngredient updates a particular ingredient. Note that UpdateIngredient expects the provided input to have a valid ID.
func (m *MariaDB) UpdateIngredient(ctx context.Context, input *models.Ingredient) error {
	query, args := m.buildUpdateIngredientQuery(input)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}

// buildArchiveIngredientQuery returns a SQL query which marks a given ingredient belonging to a given user as archived.
func (m *MariaDB) buildArchiveIngredientQuery(ingredientID, userID uint64) (query string, args []interface{}) {
	var err error
	query, args, err = m.sqlBuilder.
		Update(ingredientsTableName).
		Set("updated_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(CurrentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          ingredientID,
			"archived_on": nil,
			"belongs_to":  userID,
		}).
		ToSql()

	m.logQueryBuildingError(err)

	return query, args
}

// ArchiveIngredient marks an ingredient as archived in the database
func (m *MariaDB) ArchiveIngredient(ctx context.Context, ingredientID, userID uint64) error {
	query, args := m.buildArchiveIngredientQuery(ingredientID, userID)
	_, err := m.db.ExecContext(ctx, query, args...)
	return err
}
