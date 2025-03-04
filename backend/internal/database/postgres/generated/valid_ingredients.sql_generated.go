// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: valid_ingredients.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const archiveValidIngredient = `-- name: ArchiveValidIngredient :execrows
UPDATE valid_ingredients SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1
`

func (q *Queries) ArchiveValidIngredient(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveValidIngredient, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkValidIngredientExistence = `-- name: CheckValidIngredientExistence :one
SELECT EXISTS (
	SELECT valid_ingredients.id
	FROM valid_ingredients
	WHERE valid_ingredients.archived_at IS NULL
		AND valid_ingredients.id = $1
)
`

func (q *Queries) CheckValidIngredientExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkValidIngredientExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createValidIngredient = `-- name: CreateValidIngredient :exec
INSERT INTO valid_ingredients (
	id,
	name,
	description,
	warning,
	contains_egg,
	contains_dairy,
	contains_peanut,
	contains_tree_nut,
	contains_soy,
	contains_wheat,
	contains_shellfish,
	contains_sesame,
	contains_fish,
	contains_gluten,
	animal_flesh,
	is_liquid,
	icon_path,
	animal_derived,
	plural_name,
	restrict_to_preparations,
	minimum_ideal_storage_temperature_in_celsius,
	maximum_ideal_storage_temperature_in_celsius,
	storage_instructions,
	slug,
	contains_alcohol,
	shopping_suggestions,
	is_starch,
	is_protein,
	is_grain,
	is_fruit,
	is_salt,
	is_fat,
	is_acid,
	is_heat
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11,
	$12,
	$13,
	$14,
	$15,
	$16,
	$17,
	$18,
	$19,
	$20,
	$21,
	$22,
	$23,
	$24,
	$25,
	$26,
	$27,
	$28,
	$29,
	$30,
	$31,
	$32,
	$33,
	$34
)
`

type CreateValidIngredientParams struct {
	Name                                    string
	Description                             string
	Warning                                 string
	ShoppingSuggestions                     string
	Slug                                    string
	StorageInstructions                     string
	ID                                      string
	PluralName                              string
	IconPath                                string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsSoy                             bool
	ContainsDairy                           bool
	AnimalFlesh                             bool
	ContainsFish                            bool
	ContainsSesame                          bool
	AnimalDerived                           bool
	ContainsShellfish                       bool
	RestrictToPreparations                  bool
	ContainsWheat                           bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsGluten                          bool
	ContainsAlcohol                         bool
	ContainsEgg                             bool
	IsStarch                                bool
	IsProtein                               bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
}

func (q *Queries) CreateValidIngredient(ctx context.Context, db DBTX, arg *CreateValidIngredientParams) error {
	_, err := db.ExecContext(ctx, createValidIngredient,
		arg.ID,
		arg.Name,
		arg.Description,
		arg.Warning,
		arg.ContainsEgg,
		arg.ContainsDairy,
		arg.ContainsPeanut,
		arg.ContainsTreeNut,
		arg.ContainsSoy,
		arg.ContainsWheat,
		arg.ContainsShellfish,
		arg.ContainsSesame,
		arg.ContainsFish,
		arg.ContainsGluten,
		arg.AnimalFlesh,
		arg.IsLiquid,
		arg.IconPath,
		arg.AnimalDerived,
		arg.PluralName,
		arg.RestrictToPreparations,
		arg.MinimumIdealStorageTemperatureInCelsius,
		arg.MaximumIdealStorageTemperatureInCelsius,
		arg.StorageInstructions,
		arg.Slug,
		arg.ContainsAlcohol,
		arg.ShoppingSuggestions,
		arg.IsStarch,
		arg.IsProtein,
		arg.IsGrain,
		arg.IsFruit,
		arg.IsSalt,
		arg.IsFat,
		arg.IsAcid,
		arg.IsHeat,
	)
	return err
}

const getRandomValidIngredient = `-- name: GetRandomValidIngredient :one
SELECT
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
	valid_ingredients.is_starch,
	valid_ingredients.is_protein,
	valid_ingredients.is_grain,
	valid_ingredients.is_fruit,
	valid_ingredients.is_salt,
	valid_ingredients.is_fat,
	valid_ingredients.is_acid,
	valid_ingredients.is_heat,
	valid_ingredients.last_indexed_at,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at
FROM valid_ingredients
WHERE valid_ingredients.archived_at IS NULL
ORDER BY RANDOM() LIMIT 1
`

type GetRandomValidIngredientRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastIndexedAt                           sql.NullTime
	LastUpdatedAt                           sql.NullTime
	ShoppingSuggestions                     string
	Warning                                 string
	Description                             string
	Name                                    string
	IconPath                                string
	ID                                      string
	Slug                                    string
	StorageInstructions                     string
	PluralName                              string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsWheat                           bool
	IsProtein                               bool
	AnimalFlesh                             bool
	RestrictToPreparations                  bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsAlcohol                         bool
	ContainsSoy                             bool
	IsStarch                                bool
	AnimalDerived                           bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) GetRandomValidIngredient(ctx context.Context, db DBTX) (*GetRandomValidIngredientRow, error) {
	row := db.QueryRowContext(ctx, getRandomValidIngredient)
	var i GetRandomValidIngredientRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Warning,
		&i.ContainsEgg,
		&i.ContainsDairy,
		&i.ContainsPeanut,
		&i.ContainsTreeNut,
		&i.ContainsSoy,
		&i.ContainsWheat,
		&i.ContainsShellfish,
		&i.ContainsSesame,
		&i.ContainsFish,
		&i.ContainsGluten,
		&i.AnimalFlesh,
		&i.IsLiquid,
		&i.IconPath,
		&i.AnimalDerived,
		&i.PluralName,
		&i.RestrictToPreparations,
		&i.MinimumIdealStorageTemperatureInCelsius,
		&i.MaximumIdealStorageTemperatureInCelsius,
		&i.StorageInstructions,
		&i.Slug,
		&i.ContainsAlcohol,
		&i.ShoppingSuggestions,
		&i.IsStarch,
		&i.IsProtein,
		&i.IsGrain,
		&i.IsFruit,
		&i.IsSalt,
		&i.IsFat,
		&i.IsAcid,
		&i.IsHeat,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidIngredient = `-- name: GetValidIngredient :one
SELECT
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
	valid_ingredients.is_starch,
	valid_ingredients.is_protein,
	valid_ingredients.is_grain,
	valid_ingredients.is_fruit,
	valid_ingredients.is_salt,
	valid_ingredients.is_fat,
	valid_ingredients.is_acid,
	valid_ingredients.is_heat,
	valid_ingredients.last_indexed_at,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at
FROM valid_ingredients
WHERE valid_ingredients.archived_at IS NULL
AND valid_ingredients.id = $1
`

type GetValidIngredientRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastIndexedAt                           sql.NullTime
	LastUpdatedAt                           sql.NullTime
	ShoppingSuggestions                     string
	Warning                                 string
	Description                             string
	Name                                    string
	IconPath                                string
	ID                                      string
	Slug                                    string
	StorageInstructions                     string
	PluralName                              string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsWheat                           bool
	IsProtein                               bool
	AnimalFlesh                             bool
	RestrictToPreparations                  bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsAlcohol                         bool
	ContainsSoy                             bool
	IsStarch                                bool
	AnimalDerived                           bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) GetValidIngredient(ctx context.Context, db DBTX, id string) (*GetValidIngredientRow, error) {
	row := db.QueryRowContext(ctx, getValidIngredient, id)
	var i GetValidIngredientRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.Warning,
		&i.ContainsEgg,
		&i.ContainsDairy,
		&i.ContainsPeanut,
		&i.ContainsTreeNut,
		&i.ContainsSoy,
		&i.ContainsWheat,
		&i.ContainsShellfish,
		&i.ContainsSesame,
		&i.ContainsFish,
		&i.ContainsGluten,
		&i.AnimalFlesh,
		&i.IsLiquid,
		&i.IconPath,
		&i.AnimalDerived,
		&i.PluralName,
		&i.RestrictToPreparations,
		&i.MinimumIdealStorageTemperatureInCelsius,
		&i.MaximumIdealStorageTemperatureInCelsius,
		&i.StorageInstructions,
		&i.Slug,
		&i.ContainsAlcohol,
		&i.ShoppingSuggestions,
		&i.IsStarch,
		&i.IsProtein,
		&i.IsGrain,
		&i.IsFruit,
		&i.IsSalt,
		&i.IsFat,
		&i.IsAcid,
		&i.IsHeat,
		&i.LastIndexedAt,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getValidIngredients = `-- name: GetValidIngredients :many
SELECT
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
	valid_ingredients.is_starch,
	valid_ingredients.is_protein,
	valid_ingredients.is_grain,
	valid_ingredients.is_fruit,
	valid_ingredients.is_salt,
	valid_ingredients.is_fat,
	valid_ingredients.is_acid,
	valid_ingredients.is_heat,
	valid_ingredients.last_indexed_at,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at,
	(
		SELECT COUNT(valid_ingredients.id)
		FROM valid_ingredients
		WHERE valid_ingredients.archived_at IS NULL
			AND
			valid_ingredients.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
			AND valid_ingredients.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				valid_ingredients.last_updated_at IS NULL
				OR valid_ingredients.last_updated_at > COALESCE($3, (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				valid_ingredients.last_updated_at IS NULL
				OR valid_ingredients.last_updated_at < COALESCE($4, (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (NOT COALESCE($5, false)::boolean OR valid_ingredients.archived_at = NULL)
	) AS filtered_count,
	(
		SELECT COUNT(valid_ingredients.id)
		FROM valid_ingredients
		WHERE valid_ingredients.archived_at IS NULL
	) AS total_count
FROM valid_ingredients
WHERE
	valid_ingredients.archived_at IS NULL
	AND valid_ingredients.created_at > COALESCE($1, (SELECT NOW() - '999 years'::INTERVAL))
	AND valid_ingredients.created_at < COALESCE($2, (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		valid_ingredients.last_updated_at IS NULL
		OR valid_ingredients.last_updated_at > COALESCE($4, (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		valid_ingredients.last_updated_at IS NULL
		OR valid_ingredients.last_updated_at < COALESCE($3, (SELECT NOW() + '999 years'::INTERVAL))
	)
			AND (NOT COALESCE($5, false)::boolean OR valid_ingredients.archived_at = NULL)
GROUP BY valid_ingredients.id
ORDER BY valid_ingredients.id
LIMIT $7
OFFSET $6
`

type GetValidIngredientsParams struct {
	CreatedAfter    sql.NullTime
	CreatedBefore   sql.NullTime
	UpdatedBefore   sql.NullTime
	UpdatedAfter    sql.NullTime
	IncludeArchived sql.NullBool
	QueryOffset     sql.NullInt32
	QueryLimit      sql.NullInt32
}

type GetValidIngredientsRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastIndexedAt                           sql.NullTime
	LastUpdatedAt                           sql.NullTime
	Slug                                    string
	ShoppingSuggestions                     string
	IconPath                                string
	Warning                                 string
	Description                             string
	Name                                    string
	ID                                      string
	StorageInstructions                     string
	PluralName                              string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	FilteredCount                           int64
	TotalCount                              int64
	IsLiquid                                sql.NullBool
	ContainsShellfish                       bool
	IsFruit                                 bool
	AnimalDerived                           bool
	AnimalFlesh                             bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsAlcohol                         bool
	ContainsSesame                          bool
	IsStarch                                bool
	IsProtein                               bool
	IsGrain                                 bool
	RestrictToPreparations                  bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsWheat                           bool
	ContainsSoy                             bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) GetValidIngredients(ctx context.Context, db DBTX, arg *GetValidIngredientsParams) ([]*GetValidIngredientsRow, error) {
	rows, err := db.QueryContext(ctx, getValidIngredients,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		arg.IncludeArchived,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidIngredientsRow{}
	for rows.Next() {
		var i GetValidIngredientsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Warning,
			&i.ContainsEgg,
			&i.ContainsDairy,
			&i.ContainsPeanut,
			&i.ContainsTreeNut,
			&i.ContainsSoy,
			&i.ContainsWheat,
			&i.ContainsShellfish,
			&i.ContainsSesame,
			&i.ContainsFish,
			&i.ContainsGluten,
			&i.AnimalFlesh,
			&i.IsLiquid,
			&i.IconPath,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug,
			&i.ContainsAlcohol,
			&i.ShoppingSuggestions,
			&i.IsStarch,
			&i.IsProtein,
			&i.IsGrain,
			&i.IsFruit,
			&i.IsSalt,
			&i.IsFat,
			&i.IsAcid,
			&i.IsHeat,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.FilteredCount,
			&i.TotalCount,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getValidIngredientsNeedingIndexing = `-- name: GetValidIngredientsNeedingIndexing :many
SELECT valid_ingredients.id
FROM valid_ingredients
WHERE valid_ingredients.archived_at IS NULL
	AND (
	valid_ingredients.last_indexed_at IS NULL
	OR valid_ingredients.last_indexed_at < NOW() - '24 hours'::INTERVAL
)
`

func (q *Queries) GetValidIngredientsNeedingIndexing(ctx context.Context, db DBTX) ([]string, error) {
	rows, err := db.QueryContext(ctx, getValidIngredientsNeedingIndexing)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getValidIngredientsWithIDs = `-- name: GetValidIngredientsWithIDs :many
SELECT
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
	valid_ingredients.is_starch,
	valid_ingredients.is_protein,
	valid_ingredients.is_grain,
	valid_ingredients.is_fruit,
	valid_ingredients.is_salt,
	valid_ingredients.is_fat,
	valid_ingredients.is_acid,
	valid_ingredients.is_heat,
	valid_ingredients.last_indexed_at,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at
FROM valid_ingredients
WHERE valid_ingredients.archived_at IS NULL
	AND valid_ingredients.id = ANY($1::text[])
`

type GetValidIngredientsWithIDsRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastIndexedAt                           sql.NullTime
	LastUpdatedAt                           sql.NullTime
	ShoppingSuggestions                     string
	Warning                                 string
	Description                             string
	Name                                    string
	IconPath                                string
	ID                                      string
	Slug                                    string
	StorageInstructions                     string
	PluralName                              string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsWheat                           bool
	IsProtein                               bool
	AnimalFlesh                             bool
	RestrictToPreparations                  bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsAlcohol                         bool
	ContainsSoy                             bool
	IsStarch                                bool
	AnimalDerived                           bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) GetValidIngredientsWithIDs(ctx context.Context, db DBTX, ids []string) ([]*GetValidIngredientsWithIDsRow, error) {
	rows, err := db.QueryContext(ctx, getValidIngredientsWithIDs, pq.Array(ids))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidIngredientsWithIDsRow{}
	for rows.Next() {
		var i GetValidIngredientsWithIDsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Warning,
			&i.ContainsEgg,
			&i.ContainsDairy,
			&i.ContainsPeanut,
			&i.ContainsTreeNut,
			&i.ContainsSoy,
			&i.ContainsWheat,
			&i.ContainsShellfish,
			&i.ContainsSesame,
			&i.ContainsFish,
			&i.ContainsGluten,
			&i.AnimalFlesh,
			&i.IsLiquid,
			&i.IconPath,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug,
			&i.ContainsAlcohol,
			&i.ShoppingSuggestions,
			&i.IsStarch,
			&i.IsProtein,
			&i.IsGrain,
			&i.IsFruit,
			&i.IsSalt,
			&i.IsFat,
			&i.IsAcid,
			&i.IsHeat,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchForValidIngredients = `-- name: SearchForValidIngredients :many
SELECT
	valid_ingredients.id,
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
	valid_ingredients.is_starch,
	valid_ingredients.is_protein,
	valid_ingredients.is_grain,
	valid_ingredients.is_fruit,
	valid_ingredients.is_salt,
	valid_ingredients.is_fat,
	valid_ingredients.is_acid,
	valid_ingredients.is_heat,
	valid_ingredients.last_indexed_at,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at
FROM valid_ingredients
WHERE valid_ingredients.name ILIKE '%' || $1::text || '%'
	AND valid_ingredients.archived_at IS NULL
LIMIT 50
`

type SearchForValidIngredientsRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastIndexedAt                           sql.NullTime
	LastUpdatedAt                           sql.NullTime
	ShoppingSuggestions                     string
	Warning                                 string
	Description                             string
	Name                                    string
	IconPath                                string
	ID                                      string
	Slug                                    string
	StorageInstructions                     string
	PluralName                              string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsWheat                           bool
	IsProtein                               bool
	AnimalFlesh                             bool
	RestrictToPreparations                  bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsAlcohol                         bool
	ContainsSoy                             bool
	IsStarch                                bool
	AnimalDerived                           bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) SearchForValidIngredients(ctx context.Context, db DBTX, nameQuery string) ([]*SearchForValidIngredientsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidIngredients, nameQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientsRow{}
	for rows.Next() {
		var i SearchForValidIngredientsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Warning,
			&i.ContainsEgg,
			&i.ContainsDairy,
			&i.ContainsPeanut,
			&i.ContainsTreeNut,
			&i.ContainsSoy,
			&i.ContainsWheat,
			&i.ContainsShellfish,
			&i.ContainsSesame,
			&i.ContainsFish,
			&i.ContainsGluten,
			&i.AnimalFlesh,
			&i.IsLiquid,
			&i.IconPath,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug,
			&i.ContainsAlcohol,
			&i.ShoppingSuggestions,
			&i.IsStarch,
			&i.IsProtein,
			&i.IsGrain,
			&i.IsFruit,
			&i.IsSalt,
			&i.IsFat,
			&i.IsAcid,
			&i.IsHeat,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const searchValidIngredientsByPreparationAndIngredientName = `-- name: SearchValidIngredientsByPreparationAndIngredientName :many
SELECT
	DISTINCT(valid_ingredients.id),
	valid_ingredients.name,
	valid_ingredients.description,
	valid_ingredients.warning,
	valid_ingredients.contains_egg,
	valid_ingredients.contains_dairy,
	valid_ingredients.contains_peanut,
	valid_ingredients.contains_tree_nut,
	valid_ingredients.contains_soy,
	valid_ingredients.contains_wheat,
	valid_ingredients.contains_shellfish,
	valid_ingredients.contains_sesame,
	valid_ingredients.contains_fish,
	valid_ingredients.contains_gluten,
	valid_ingredients.animal_flesh,
	valid_ingredients.is_liquid,
	valid_ingredients.icon_path,
	valid_ingredients.animal_derived,
	valid_ingredients.plural_name,
	valid_ingredients.restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions,
	valid_ingredients.slug,
	valid_ingredients.contains_alcohol,
	valid_ingredients.shopping_suggestions,
	valid_ingredients.is_starch,
	valid_ingredients.is_protein,
	valid_ingredients.is_grain,
	valid_ingredients.is_fruit,
	valid_ingredients.is_salt,
	valid_ingredients.is_fat,
	valid_ingredients.is_acid,
	valid_ingredients.is_heat,
	valid_ingredients.last_indexed_at,
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at
FROM valid_ingredient_preparations
	JOIN valid_ingredients ON valid_ingredient_preparations.valid_ingredient_id = valid_ingredients.id
	JOIN valid_preparations ON valid_ingredient_preparations.valid_preparation_id = valid_preparations.id
WHERE valid_ingredient_preparations.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND (
		valid_ingredient_preparations.valid_preparation_id = $1
		OR valid_preparations.restrict_to_ingredients IS FALSE
	)
	AND valid_ingredients.name ILIKE '%' || $2::text || '%'
`

type SearchValidIngredientsByPreparationAndIngredientNameParams struct {
	ValidPreparationID string
	NameQuery          string
}

type SearchValidIngredientsByPreparationAndIngredientNameRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastIndexedAt                           sql.NullTime
	LastUpdatedAt                           sql.NullTime
	ShoppingSuggestions                     string
	Warning                                 string
	Description                             string
	Name                                    string
	IconPath                                string
	ID                                      string
	Slug                                    string
	StorageInstructions                     string
	PluralName                              string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsWheat                           bool
	IsProtein                               bool
	AnimalFlesh                             bool
	RestrictToPreparations                  bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsAlcohol                         bool
	ContainsSoy                             bool
	IsStarch                                bool
	AnimalDerived                           bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) SearchValidIngredientsByPreparationAndIngredientName(ctx context.Context, db DBTX, arg *SearchValidIngredientsByPreparationAndIngredientNameParams) ([]*SearchValidIngredientsByPreparationAndIngredientNameRow, error) {
	rows, err := db.QueryContext(ctx, searchValidIngredientsByPreparationAndIngredientName, arg.ValidPreparationID, arg.NameQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchValidIngredientsByPreparationAndIngredientNameRow{}
	for rows.Next() {
		var i SearchValidIngredientsByPreparationAndIngredientNameRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Warning,
			&i.ContainsEgg,
			&i.ContainsDairy,
			&i.ContainsPeanut,
			&i.ContainsTreeNut,
			&i.ContainsSoy,
			&i.ContainsWheat,
			&i.ContainsShellfish,
			&i.ContainsSesame,
			&i.ContainsFish,
			&i.ContainsGluten,
			&i.AnimalFlesh,
			&i.IsLiquid,
			&i.IconPath,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug,
			&i.ContainsAlcohol,
			&i.ShoppingSuggestions,
			&i.IsStarch,
			&i.IsProtein,
			&i.IsGrain,
			&i.IsFruit,
			&i.IsSalt,
			&i.IsFat,
			&i.IsAcid,
			&i.IsHeat,
			&i.LastIndexedAt,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateValidIngredient = `-- name: UpdateValidIngredient :execrows
UPDATE valid_ingredients SET
	name = $1,
	description = $2,
	warning = $3,
	contains_egg = $4,
	contains_dairy = $5,
	contains_peanut = $6,
	contains_tree_nut = $7,
	contains_soy = $8,
	contains_wheat = $9,
	contains_shellfish = $10,
	contains_sesame = $11,
	contains_fish = $12,
	contains_gluten = $13,
	animal_flesh = $14,
	is_liquid = $15,
	icon_path = $16,
	animal_derived = $17,
	plural_name = $18,
	restrict_to_preparations = $19,
	minimum_ideal_storage_temperature_in_celsius = $20,
	maximum_ideal_storage_temperature_in_celsius = $21,
	storage_instructions = $22,
	slug = $23,
	contains_alcohol = $24,
	shopping_suggestions = $25,
	is_starch = $26,
	is_protein = $27,
	is_grain = $28,
	is_fruit = $29,
	is_salt = $30,
	is_fat = $31,
	is_acid = $32,
	is_heat = $33,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $34
`

type UpdateValidIngredientParams struct {
	Description                             string
	Warning                                 string
	ID                                      string
	ShoppingSuggestions                     string
	Slug                                    string
	StorageInstructions                     string
	Name                                    string
	PluralName                              string
	IconPath                                string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	ContainsWheat                           bool
	ContainsAlcohol                         bool
	ContainsGluten                          bool
	ContainsFish                            bool
	AnimalDerived                           bool
	ContainsSesame                          bool
	RestrictToPreparations                  bool
	ContainsShellfish                       bool
	ContainsSoy                             bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	AnimalFlesh                             bool
	ContainsDairy                           bool
	IsStarch                                bool
	IsProtein                               bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsEgg                             bool
}

func (q *Queries) UpdateValidIngredient(ctx context.Context, db DBTX, arg *UpdateValidIngredientParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidIngredient,
		arg.Name,
		arg.Description,
		arg.Warning,
		arg.ContainsEgg,
		arg.ContainsDairy,
		arg.ContainsPeanut,
		arg.ContainsTreeNut,
		arg.ContainsSoy,
		arg.ContainsWheat,
		arg.ContainsShellfish,
		arg.ContainsSesame,
		arg.ContainsFish,
		arg.ContainsGluten,
		arg.AnimalFlesh,
		arg.IsLiquid,
		arg.IconPath,
		arg.AnimalDerived,
		arg.PluralName,
		arg.RestrictToPreparations,
		arg.MinimumIdealStorageTemperatureInCelsius,
		arg.MaximumIdealStorageTemperatureInCelsius,
		arg.StorageInstructions,
		arg.Slug,
		arg.ContainsAlcohol,
		arg.ShoppingSuggestions,
		arg.IsStarch,
		arg.IsProtein,
		arg.IsGrain,
		arg.IsFruit,
		arg.IsSalt,
		arg.IsFat,
		arg.IsAcid,
		arg.IsHeat,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const updateValidIngredientLastIndexedAt = `-- name: UpdateValidIngredientLastIndexedAt :execrows
UPDATE valid_ingredients SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL
`

func (q *Queries) UpdateValidIngredientLastIndexedAt(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, updateValidIngredientLastIndexedAt, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
