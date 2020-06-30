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
	validIngredientsTableName = "valid_ingredients"
)

var (
	validIngredientsTableColumns = []string{
		fmt.Sprintf("%s.%s", validIngredientsTableName, "id"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "name"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "variant"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "description"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "warning"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_egg"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_dairy"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_peanut"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_tree_nut"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_soy"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_wheat"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_shellfish"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_sesame"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_fish"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "contains_gluten"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "animal_flesh"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "animal_derived"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "measurable_by_volume"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "icon"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "created_on"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "updated_on"),
		fmt.Sprintf("%s.%s", validIngredientsTableName, "archived_on"),
	}

	validIngredientsOnIngredientTagMappingsJoinClause       = fmt.Sprintf("%s ON %s.%s=%s.id", validIngredientsTableName, ingredientTagMappingsTableName, ingredientTagMappingsTableOwnershipColumn, validIngredientsTableName)
	validIngredientsOnValidIngredientPreparationsJoinClause = fmt.Sprintf("%s ON %s.%s=%s.id", validIngredientsTableName, validIngredientPreparationsTableName, validIngredientPreparationsTableOwnershipColumn, validIngredientsTableName)
)

// scanValidIngredient takes a database Scanner (i.e. *sql.Row) and scans the result into a Valid Ingredient struct
func (p *Postgres) scanValidIngredient(scan database.Scanner, includeCount bool) (*models.ValidIngredient, uint64, error) {
	x := &models.ValidIngredient{}
	var count uint64

	targetVars := []interface{}{
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
		&x.MeasurableByVolume,
		&x.Icon,
		&x.CreatedOn,
		&x.UpdatedOn,
		&x.ArchivedOn,
	}

	if includeCount {
		targetVars = append(targetVars, &count)
	}

	if err := scan.Scan(targetVars...); err != nil {
		return nil, 0, err
	}

	return x, count, nil
}

// scanValidIngredients takes a logger and some database rows and turns them into a slice of valid ingredients.
func (p *Postgres) scanValidIngredients(rows database.ResultIterator) ([]models.ValidIngredient, uint64, error) {
	var (
		list  []models.ValidIngredient
		count uint64
	)

	for rows.Next() {
		x, c, err := p.scanValidIngredient(rows, true)
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

// buildValidIngredientExistsQuery constructs a SQL query for checking if a valid ingredient with a given ID exists
func (p *Postgres) buildValidIngredientExistsQuery(validIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(fmt.Sprintf("%s.id", validIngredientsTableName)).
		Prefix(existencePrefix).
		From(validIngredientsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validIngredientsTableName): validIngredientID,
		}).ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ValidIngredientExists queries the database to see if a given valid ingredient belonging to a given user exists.
func (p *Postgres) ValidIngredientExists(ctx context.Context, validIngredientID uint64) (exists bool, err error) {
	query, args := p.buildValidIngredientExistsQuery(validIngredientID)

	err = p.db.QueryRowContext(ctx, query, args...).Scan(&exists)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return exists, err
}

// buildGetValidIngredientQuery constructs a SQL query for fetching a valid ingredient with a given ID.
func (p *Postgres) buildGetValidIngredientQuery(validIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Select(validIngredientsTableColumns...).
		From(validIngredientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.id", validIngredientsTableName): validIngredientID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredient fetches a valid ingredient from the database.
func (p *Postgres) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*models.ValidIngredient, error) {
	query, args := p.buildGetValidIngredientQuery(validIngredientID)
	row := p.db.QueryRowContext(ctx, query, args...)

	validIngredient, _, err := p.scanValidIngredient(row, false)
	return validIngredient, err
}

var (
	allValidIngredientsCountQueryBuilder sync.Once
	allValidIngredientsCountQuery        string
)

// buildGetAllValidIngredientsCountQuery returns a query that fetches the total number of valid ingredients in the database.
// This query only gets generated once, and is otherwise returned from cache.
func (p *Postgres) buildGetAllValidIngredientsCountQuery() string {
	allValidIngredientsCountQueryBuilder.Do(func() {
		var err error

		allValidIngredientsCountQuery, _, err = p.sqlBuilder.
			Select(fmt.Sprintf(countQuery, validIngredientsTableName)).
			From(validIngredientsTableName).
			Where(squirrel.Eq{
				fmt.Sprintf("%s.archived_on", validIngredientsTableName): nil,
			}).
			ToSql()
		p.logQueryBuildingError(err)
	})

	return allValidIngredientsCountQuery
}

// GetAllValidIngredientsCount will fetch the count of valid ingredients from the database.
func (p *Postgres) GetAllValidIngredientsCount(ctx context.Context) (count uint64, err error) {
	err = p.db.QueryRowContext(ctx, p.buildGetAllValidIngredientsCountQuery()).Scan(&count)
	return count, err
}

// buildGetValidIngredientsQuery builds a SQL query selecting valid ingredients that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidIngredientsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	columnsToSelect := append(validIngredientsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidIngredientsCountQuery()))

	builder := p.sqlBuilder.
		Select(columnsToSelect...).
		From(validIngredientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.archived_on", validIngredientsTableName): nil,
		}).
		OrderBy(fmt.Sprintf("%s.id", validIngredientsTableName))

	if filter != nil {
		builder = filter.ApplyToQueryBuilder(builder, validIngredientsTableName)
	}

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredients fetches a list of valid ingredients from the database that meet a particular filter.
func (p *Postgres) GetValidIngredients(ctx context.Context, filter *models.QueryFilter) (*models.ValidIngredientList, error) {
	query, args := p.buildGetValidIngredientsQuery(filter)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid ingredients")
	}

	validIngredients, count, err := p.scanValidIngredients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	list := &models.ValidIngredientList{
		Pagination: models.Pagination{
			Page:       filter.Page,
			Limit:      filter.Limit,
			TotalCount: count,
		},
		ValidIngredients: validIngredients,
	}

	return list, nil
}

// buildCreateValidIngredientQuery takes a valid ingredient and returns a creation query for that valid ingredient and the relevant arguments.
func (p *Postgres) buildCreateValidIngredientQuery(input *models.ValidIngredient) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validIngredientsTableName).
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
			"measurable_by_volume",
			"icon",
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
			input.MeasurableByVolume,
			input.Icon,
		).
		Suffix("RETURNING id, created_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// CreateValidIngredient creates a valid ingredient in the database.
func (p *Postgres) CreateValidIngredient(ctx context.Context, input *models.ValidIngredientCreationInput) (*models.ValidIngredient, error) {
	x := &models.ValidIngredient{
		Name:               input.Name,
		Variant:            input.Variant,
		Description:        input.Description,
		Warning:            input.Warning,
		ContainsEgg:        input.ContainsEgg,
		ContainsDairy:      input.ContainsDairy,
		ContainsPeanut:     input.ContainsPeanut,
		ContainsTreeNut:    input.ContainsTreeNut,
		ContainsSoy:        input.ContainsSoy,
		ContainsWheat:      input.ContainsWheat,
		ContainsShellfish:  input.ContainsShellfish,
		ContainsSesame:     input.ContainsSesame,
		ContainsFish:       input.ContainsFish,
		ContainsGluten:     input.ContainsGluten,
		AnimalFlesh:        input.AnimalFlesh,
		AnimalDerived:      input.AnimalDerived,
		MeasurableByVolume: input.MeasurableByVolume,
		Icon:               input.Icon,
	}

	query, args := p.buildCreateValidIngredientQuery(x)

	// create the valid ingredient.
	err := p.db.QueryRowContext(ctx, query, args...).Scan(&x.ID, &x.CreatedOn)
	if err != nil {
		return nil, fmt.Errorf("error executing valid ingredient creation query: %w", err)
	}

	return x, nil
}

// buildUpdateValidIngredientQuery takes a valid ingredient and returns an update SQL query, with the relevant query parameters.
func (p *Postgres) buildUpdateValidIngredientQuery(input *models.ValidIngredient) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientsTableName).
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
		Set("measurable_by_volume", input.MeasurableByVolume).
		Set("icon", input.Icon).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id": input.ID,
		}).
		Suffix("RETURNING updated_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidIngredient updates a particular valid ingredient. Note that UpdateValidIngredient expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidIngredient(ctx context.Context, input *models.ValidIngredient) error {
	query, args := p.buildUpdateValidIngredientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.UpdatedOn)
}

// buildArchiveValidIngredientQuery returns a SQL query which marks a given valid ingredient as archived.
func (p *Postgres) buildArchiveValidIngredientQuery(validIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientsTableName).
		Set("updated_on", squirrel.Expr(currentUnixTimeQuery)).
		Set("archived_on", squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			"id":          validIngredientID,
			"archived_on": nil,
		}).
		Suffix("RETURNING archived_on").
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// ArchiveValidIngredient marks a valid ingredient as archived in the database.
func (p *Postgres) ArchiveValidIngredient(ctx context.Context, validIngredientID uint64) error {
	query, args := p.buildArchiveValidIngredientQuery(validIngredientID)

	res, err := p.db.ExecContext(ctx, query, args...)
	if res != nil {
		if rowCount, rowCountErr := res.RowsAffected(); rowCountErr == nil && rowCount == 0 {
			return sql.ErrNoRows
		}
	}

	return err
}
