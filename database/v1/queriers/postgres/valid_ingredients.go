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
	validIngredientsTableName                     = "valid_ingredients"
	validIngredientsTableNameColumn               = "name"
	validIngredientsTableVariantColumn            = "variant"
	validIngredientsTableDescriptionColumn        = "description"
	validIngredientsTableWarningColumn            = "warning"
	validIngredientsTableContainsEggColumn        = "contains_egg"
	validIngredientsTableContainsDairyColumn      = "contains_dairy"
	validIngredientsTableContainsPeanutColumn     = "contains_peanut"
	validIngredientsTableContainsTreeNutColumn    = "contains_tree_nut"
	validIngredientsTableContainsSoyColumn        = "contains_soy"
	validIngredientsTableContainsWheatColumn      = "contains_wheat"
	validIngredientsTableContainsShellfishColumn  = "contains_shellfish"
	validIngredientsTableContainsSesameColumn     = "contains_sesame"
	validIngredientsTableContainsFishColumn       = "contains_fish"
	validIngredientsTableContainsGlutenColumn     = "contains_gluten"
	validIngredientsTableAnimalFleshColumn        = "animal_flesh"
	validIngredientsTableAnimalDerivedColumn      = "animal_derived"
	validIngredientsTableMeasurableByVolumeColumn = "measurable_by_volume"
	validIngredientsTableIconColumn               = "icon"
)

var (
	validIngredientsTableColumns = []string{
		fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableNameColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableVariantColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableDescriptionColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableWarningColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsEggColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsDairyColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsPeanutColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsTreeNutColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsSoyColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsWheatColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsShellfishColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsSesameColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsFishColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableContainsGlutenColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableAnimalFleshColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableAnimalDerivedColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableMeasurableByVolumeColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, validIngredientsTableIconColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, createdOnColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, lastUpdatedOnColumn),
		fmt.Sprintf("%s.%s", validIngredientsTableName, archivedOnColumn),
	}
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
		&x.LastUpdatedOn,
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
		Select(fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn)).
		Prefix(existencePrefix).
		From(validIngredientsTableName).
		Suffix(existenceSuffix).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn): validIngredientID,
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
			fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn): validIngredientID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredient fetches a valid ingredient from the database.
func (p *Postgres) GetValidIngredient(ctx context.Context, validIngredientID uint64) (*models.ValidIngredient, error) {
	query, args := p.buildGetValidIngredientQuery(validIngredientID)
	row := p.db.QueryRowContext(ctx, query, args...)
	vi, _, err := p.scanValidIngredient(row, false)
	return vi, err
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
				fmt.Sprintf("%s.%s", validIngredientsTableName, archivedOnColumn): nil,
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

// buildGetBatchOfValidIngredientsQuery returns a query that fetches every valid ingredient in the database within a bucketed range.
func (p *Postgres) buildGetBatchOfValidIngredientsQuery(beginID, endID uint64) (query string, args []interface{}) {
	query, args, err := p.sqlBuilder.
		Select(validIngredientsTableColumns...).
		From(validIngredientsTableName).
		Where(squirrel.Gt{
			fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn): beginID,
		}).
		Where(squirrel.Lt{
			fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn): endID,
		}).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// GetAllValidIngredients fetches every valid ingredient from the database and writes them to a channel. This method primarily exists
// to aid in administrative data tasks.
func (p *Postgres) GetAllValidIngredients(ctx context.Context, resultChannel chan []models.ValidIngredient) error {
	count, err := p.GetAllValidIngredientsCount(ctx)
	if err != nil {
		return err
	}

	for beginID := uint64(1); beginID <= count; beginID += defaultBucketSize {
		endID := beginID + defaultBucketSize
		go func(begin, end uint64) {
			query, args := p.buildGetBatchOfValidIngredientsQuery(begin, end)
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

			validIngredients, _, err := p.scanValidIngredients(rows)
			if err != nil {
				logger.Error(err, "scanning database rows")
				return
			}

			resultChannel <- validIngredients
		}(beginID, endID)
	}

	return nil
}

// buildGetValidIngredientsQuery builds a SQL query selecting valid ingredients that adhere to a given QueryFilter,
// and returns both the query and the relevant args to pass to the query executor.
func (p *Postgres) buildGetValidIngredientsQuery(filter *models.QueryFilter) (query string, args []interface{}) {
	var err error

	builder := p.sqlBuilder.
		Select(append(validIngredientsTableColumns, fmt.Sprintf("(%s)", p.buildGetAllValidIngredientsCountQuery()))...).
		From(validIngredientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientsTableName, archivedOnColumn): nil,
		}).
		OrderBy(fmt.Sprintf("%s.%s", validIngredientsTableName, idColumn))

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

// buildGetValidIngredientsWithIDsQuery builds a SQL query selecting validIngredients
// and have IDs that exist within a given set of IDs. Returns both the query and the relevant
// args to pass to the query executor. This function is primarily intended for use with a search
// index, which would provide a slice of string IDs to query against. This function accepts a
// slice of uint64s instead of a slice of strings in order to ensure all the provided strings
// are valid database IDs, because there's no way in squirrel to escape them in the unnest join,
// and if we accept strings we could leave ourselves vulnerable to SQL injection attacks.
func (p *Postgres) buildGetValidIngredientsWithIDsQuery(limit uint8, ids []uint64) (query string, args []interface{}) {
	var err error

	subqueryBuilder := p.sqlBuilder.Select(validIngredientsTableColumns...).
		From(validIngredientsTableName).
		Join(fmt.Sprintf("unnest('{%s}'::int[])", joinUint64s(ids))).
		Suffix(fmt.Sprintf("WITH ORDINALITY t(id, ord) USING (id) ORDER BY t.ord LIMIT %d", limit))
	builder := p.sqlBuilder.
		Select(validIngredientsTableColumns...).
		FromSelect(subqueryBuilder, validIngredientsTableName).
		Where(squirrel.Eq{
			fmt.Sprintf("%s.%s", validIngredientsTableName, archivedOnColumn): nil,
		})

	query, args, err = builder.ToSql()
	p.logQueryBuildingError(err)

	return query, args
}

// GetValidIngredientsWithIDs fetches a list of valid ingredients from the database that exist within a given set of IDs.
func (p *Postgres) GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]models.ValidIngredient, error) {
	if limit == 0 {
		limit = uint8(models.DefaultLimit)
	}

	query, args := p.buildGetValidIngredientsWithIDsQuery(limit, ids)

	rows, err := p.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, buildError(err, "querying database for valid ingredients")
	}

	validIngredients, _, err := p.scanValidIngredients(rows)
	if err != nil {
		return nil, fmt.Errorf("scanning response from database: %w", err)
	}

	return validIngredients, nil
}

// buildCreateValidIngredientQuery takes a valid ingredient and returns a creation query for that valid ingredient and the relevant arguments.
func (p *Postgres) buildCreateValidIngredientQuery(input *models.ValidIngredient) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Insert(validIngredientsTableName).
		Columns(
			validIngredientsTableNameColumn,
			validIngredientsTableVariantColumn,
			validIngredientsTableDescriptionColumn,
			validIngredientsTableWarningColumn,
			validIngredientsTableContainsEggColumn,
			validIngredientsTableContainsDairyColumn,
			validIngredientsTableContainsPeanutColumn,
			validIngredientsTableContainsTreeNutColumn,
			validIngredientsTableContainsSoyColumn,
			validIngredientsTableContainsWheatColumn,
			validIngredientsTableContainsShellfishColumn,
			validIngredientsTableContainsSesameColumn,
			validIngredientsTableContainsFishColumn,
			validIngredientsTableContainsGlutenColumn,
			validIngredientsTableAnimalFleshColumn,
			validIngredientsTableAnimalDerivedColumn,
			validIngredientsTableMeasurableByVolumeColumn,
			validIngredientsTableIconColumn,
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
		Suffix(fmt.Sprintf("RETURNING %s, %s", idColumn, createdOnColumn)).
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
		Set(validIngredientsTableNameColumn, input.Name).
		Set(validIngredientsTableVariantColumn, input.Variant).
		Set(validIngredientsTableDescriptionColumn, input.Description).
		Set(validIngredientsTableWarningColumn, input.Warning).
		Set(validIngredientsTableContainsEggColumn, input.ContainsEgg).
		Set(validIngredientsTableContainsDairyColumn, input.ContainsDairy).
		Set(validIngredientsTableContainsPeanutColumn, input.ContainsPeanut).
		Set(validIngredientsTableContainsTreeNutColumn, input.ContainsTreeNut).
		Set(validIngredientsTableContainsSoyColumn, input.ContainsSoy).
		Set(validIngredientsTableContainsWheatColumn, input.ContainsWheat).
		Set(validIngredientsTableContainsShellfishColumn, input.ContainsShellfish).
		Set(validIngredientsTableContainsSesameColumn, input.ContainsSesame).
		Set(validIngredientsTableContainsFishColumn, input.ContainsFish).
		Set(validIngredientsTableContainsGlutenColumn, input.ContainsGluten).
		Set(validIngredientsTableAnimalFleshColumn, input.AnimalFlesh).
		Set(validIngredientsTableAnimalDerivedColumn, input.AnimalDerived).
		Set(validIngredientsTableMeasurableByVolumeColumn, input.MeasurableByVolume).
		Set(validIngredientsTableIconColumn, input.Icon).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn: input.ID,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", lastUpdatedOnColumn)).
		ToSql()

	p.logQueryBuildingError(err)

	return query, args
}

// UpdateValidIngredient updates a particular valid ingredient. Note that UpdateValidIngredient expects the provided input to have a valid ID.
func (p *Postgres) UpdateValidIngredient(ctx context.Context, input *models.ValidIngredient) error {
	query, args := p.buildUpdateValidIngredientQuery(input)
	return p.db.QueryRowContext(ctx, query, args...).Scan(&input.LastUpdatedOn)
}

// buildArchiveValidIngredientQuery returns a SQL query which marks a given valid ingredient as archived.
func (p *Postgres) buildArchiveValidIngredientQuery(validIngredientID uint64) (query string, args []interface{}) {
	var err error

	query, args, err = p.sqlBuilder.
		Update(validIngredientsTableName).
		Set(lastUpdatedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Set(archivedOnColumn, squirrel.Expr(currentUnixTimeQuery)).
		Where(squirrel.Eq{
			idColumn:         validIngredientID,
			archivedOnColumn: nil,
		}).
		Suffix(fmt.Sprintf("RETURNING %s", archivedOnColumn)).
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
