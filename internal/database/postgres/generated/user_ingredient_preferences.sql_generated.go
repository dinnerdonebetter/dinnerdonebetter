// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: user_ingredient_preferences.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveUserIngredientPreference = `-- name: ArchiveUserIngredientPreference :execrows

UPDATE user_ingredient_preferences SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_user = $1 AND id = $2
`

type ArchiveUserIngredientPreferenceParams struct {
	BelongsToUser string
	ID            string
}

func (q *Queries) ArchiveUserIngredientPreference(ctx context.Context, db DBTX, arg *ArchiveUserIngredientPreferenceParams) (int64, error) {
	result, err := db.ExecContext(ctx, archiveUserIngredientPreference, arg.BelongsToUser, arg.ID)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkUserIngredientPreferenceExistence = `-- name: CheckUserIngredientPreferenceExistence :one

SELECT EXISTS (
	SELECT user_ingredient_preferences.id
	FROM user_ingredient_preferences
	WHERE user_ingredient_preferences.archived_at IS NULL
		AND user_ingredient_preferences.id = $1
		AND user_ingredient_preferences.belongs_to_user = $2
)
`

type CheckUserIngredientPreferenceExistenceParams struct {
	ID            string
	BelongsToUser string
}

func (q *Queries) CheckUserIngredientPreferenceExistence(ctx context.Context, db DBTX, arg *CheckUserIngredientPreferenceExistenceParams) (bool, error) {
	row := db.QueryRowContext(ctx, checkUserIngredientPreferenceExistence, arg.ID, arg.BelongsToUser)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createUserIngredientPreference = `-- name: CreateUserIngredientPreference :exec

INSERT INTO user_ingredient_preferences (
	id,
	ingredient,
	rating,
	notes,
	allergy,
	belongs_to_user
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
`

type CreateUserIngredientPreferenceParams struct {
	ID            string
	Ingredient    string
	Notes         string
	BelongsToUser string
	Rating        int16
	Allergy       bool
}

func (q *Queries) CreateUserIngredientPreference(ctx context.Context, db DBTX, arg *CreateUserIngredientPreferenceParams) error {
	_, err := db.ExecContext(ctx, createUserIngredientPreference,
		arg.ID,
		arg.Ingredient,
		arg.Rating,
		arg.Notes,
		arg.Allergy,
		arg.BelongsToUser,
	)
	return err
}

const getUserIngredientPreference = `-- name: GetUserIngredientPreference :one

SELECT
	user_ingredient_preferences.id,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
	valid_ingredients.is_starch as valid_ingredient_is_starch,
	valid_ingredients.is_protein as valid_ingredient_is_protein,
	valid_ingredients.is_grain as valid_ingredient_is_grain,
	valid_ingredients.is_fruit as valid_ingredient_is_fruit,
	valid_ingredients.is_salt as valid_ingredient_is_salt,
	valid_ingredients.is_fat as valid_ingredient_is_fat,
	valid_ingredients.is_acid as valid_ingredient_is_acid,
	valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
	user_ingredient_preferences.rating,
	user_ingredient_preferences.notes,
	user_ingredient_preferences.allergy,
	user_ingredient_preferences.created_at,
	user_ingredient_preferences.last_updated_at,
	user_ingredient_preferences.archived_at,
	user_ingredient_preferences.belongs_to_user
FROM user_ingredient_preferences
	JOIN valid_ingredients ON valid_ingredients.id = user_ingredient_preferences.ingredient
WHERE user_ingredient_preferences.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND user_ingredient_preferences.id = $1
	AND user_ingredient_preferences.belongs_to_user = $2
`

type GetUserIngredientPreferenceParams struct {
	UserIngredientPreferenceID string
	UserID                     string
}

type GetUserIngredientPreferenceRow struct {
	ValidIngredientCreatedAt                               time.Time
	CreatedAt                                              time.Time
	ArchivedAt                                             sql.NullTime
	ValidIngredientLastUpdatedAt                           sql.NullTime
	ValidIngredientArchivedAt                              sql.NullTime
	LastUpdatedAt                                          sql.NullTime
	ValidIngredientShoppingSuggestions                     string
	ValidIngredientSlug                                    string
	ValidIngredientWarning                                 string
	Notes                                                  string
	ValidIngredientPluralName                              string
	ID                                                     string
	ValidIngredientDescription                             string
	ValidIngredientName                                    string
	ValidIngredientID                                      string
	ValidIngredientStorageInstructions                     string
	BelongsToUser                                          string
	ValidIngredientIconPath                                string
	ValidIngredientMinimumIdealStorageTemperatureInCelsius sql.NullString
	ValidIngredientMaximumIdealStorageTemperatureInCelsius sql.NullString
	Rating                                                 int16
	ValidIngredientIsLiquid                                sql.NullBool
	ValidIngredientContainsShellfish                       bool
	ValidIngredientIsFat                                   bool
	ValidIngredientAnimalDerived                           bool
	ValidIngredientVolumetric                              bool
	ValidIngredientContainsAlcohol                         bool
	ValidIngredientAnimalFlesh                             bool
	ValidIngredientIsStarch                                bool
	ValidIngredientIsProtein                               bool
	ValidIngredientIsGrain                                 bool
	ValidIngredientIsFruit                                 bool
	ValidIngredientIsSalt                                  bool
	ValidIngredientRestrictToPreparations                  bool
	ValidIngredientIsAcid                                  bool
	ValidIngredientIsHeat                                  bool
	ValidIngredientContainsGluten                          bool
	ValidIngredientContainsFish                            bool
	ValidIngredientContainsSesame                          bool
	ValidIngredientContainsWheat                           bool
	ValidIngredientContainsSoy                             bool
	Allergy                                                bool
	ValidIngredientContainsTreeNut                         bool
	ValidIngredientContainsPeanut                          bool
	ValidIngredientContainsDairy                           bool
	ValidIngredientContainsEgg                             bool
}

func (q *Queries) GetUserIngredientPreference(ctx context.Context, db DBTX, arg *GetUserIngredientPreferenceParams) (*GetUserIngredientPreferenceRow, error) {
	row := db.QueryRowContext(ctx, getUserIngredientPreference, arg.UserIngredientPreferenceID, arg.UserID)
	var i GetUserIngredientPreferenceRow
	err := row.Scan(
		&i.ID,
		&i.ValidIngredientID,
		&i.ValidIngredientName,
		&i.ValidIngredientDescription,
		&i.ValidIngredientWarning,
		&i.ValidIngredientContainsEgg,
		&i.ValidIngredientContainsDairy,
		&i.ValidIngredientContainsPeanut,
		&i.ValidIngredientContainsTreeNut,
		&i.ValidIngredientContainsSoy,
		&i.ValidIngredientContainsWheat,
		&i.ValidIngredientContainsShellfish,
		&i.ValidIngredientContainsSesame,
		&i.ValidIngredientContainsFish,
		&i.ValidIngredientContainsGluten,
		&i.ValidIngredientAnimalFlesh,
		&i.ValidIngredientVolumetric,
		&i.ValidIngredientIsLiquid,
		&i.ValidIngredientIconPath,
		&i.ValidIngredientAnimalDerived,
		&i.ValidIngredientPluralName,
		&i.ValidIngredientRestrictToPreparations,
		&i.ValidIngredientMinimumIdealStorageTemperatureInCelsius,
		&i.ValidIngredientMaximumIdealStorageTemperatureInCelsius,
		&i.ValidIngredientStorageInstructions,
		&i.ValidIngredientSlug,
		&i.ValidIngredientContainsAlcohol,
		&i.ValidIngredientShoppingSuggestions,
		&i.ValidIngredientIsStarch,
		&i.ValidIngredientIsProtein,
		&i.ValidIngredientIsGrain,
		&i.ValidIngredientIsFruit,
		&i.ValidIngredientIsSalt,
		&i.ValidIngredientIsFat,
		&i.ValidIngredientIsAcid,
		&i.ValidIngredientIsHeat,
		&i.ValidIngredientCreatedAt,
		&i.ValidIngredientLastUpdatedAt,
		&i.ValidIngredientArchivedAt,
		&i.Rating,
		&i.Notes,
		&i.Allergy,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.BelongsToUser,
	)
	return &i, err
}

const getUserIngredientPreferencesForUser = `-- name: GetUserIngredientPreferencesForUser :many

SELECT
	user_ingredient_preferences.id,
	valid_ingredients.id as valid_ingredient_id,
	valid_ingredients.name as valid_ingredient_name,
	valid_ingredients.description as valid_ingredient_description,
	valid_ingredients.warning as valid_ingredient_warning,
	valid_ingredients.contains_egg as valid_ingredient_contains_egg,
	valid_ingredients.contains_dairy as valid_ingredient_contains_dairy,
	valid_ingredients.contains_peanut as valid_ingredient_contains_peanut,
	valid_ingredients.contains_tree_nut as valid_ingredient_contains_tree_nut,
	valid_ingredients.contains_soy as valid_ingredient_contains_soy,
	valid_ingredients.contains_wheat as valid_ingredient_contains_wheat,
	valid_ingredients.contains_shellfish as valid_ingredient_contains_shellfish,
	valid_ingredients.contains_sesame as valid_ingredient_contains_sesame,
	valid_ingredients.contains_fish as valid_ingredient_contains_fish,
	valid_ingredients.contains_gluten as valid_ingredient_contains_gluten,
	valid_ingredients.animal_flesh as valid_ingredient_animal_flesh,
	valid_ingredients.volumetric as valid_ingredient_volumetric,
	valid_ingredients.is_liquid as valid_ingredient_is_liquid,
	valid_ingredients.icon_path as valid_ingredient_icon_path,
	valid_ingredients.animal_derived as valid_ingredient_animal_derived,
	valid_ingredients.plural_name as valid_ingredient_plural_name,
	valid_ingredients.restrict_to_preparations as valid_ingredient_restrict_to_preparations,
	valid_ingredients.minimum_ideal_storage_temperature_in_celsius as valid_ingredient_minimum_ideal_storage_temperature_in_celsius,
	valid_ingredients.maximum_ideal_storage_temperature_in_celsius as valid_ingredient_maximum_ideal_storage_temperature_in_celsius,
	valid_ingredients.storage_instructions as valid_ingredient_storage_instructions,
	valid_ingredients.slug as valid_ingredient_slug,
	valid_ingredients.contains_alcohol as valid_ingredient_contains_alcohol,
	valid_ingredients.shopping_suggestions as valid_ingredient_shopping_suggestions,
	valid_ingredients.is_starch as valid_ingredient_is_starch,
	valid_ingredients.is_protein as valid_ingredient_is_protein,
	valid_ingredients.is_grain as valid_ingredient_is_grain,
	valid_ingredients.is_fruit as valid_ingredient_is_fruit,
	valid_ingredients.is_salt as valid_ingredient_is_salt,
	valid_ingredients.is_fat as valid_ingredient_is_fat,
	valid_ingredients.is_acid as valid_ingredient_is_acid,
	valid_ingredients.is_heat as valid_ingredient_is_heat,
	valid_ingredients.created_at as valid_ingredient_created_at,
	valid_ingredients.last_updated_at as valid_ingredient_last_updated_at,
	valid_ingredients.archived_at as valid_ingredient_archived_at,
	user_ingredient_preferences.rating,
	user_ingredient_preferences.notes,
	user_ingredient_preferences.allergy,
	user_ingredient_preferences.created_at,
	user_ingredient_preferences.last_updated_at,
	user_ingredient_preferences.archived_at,
	user_ingredient_preferences.belongs_to_user,
	(
		SELECT
			COUNT(user_ingredient_preferences.id)
		FROM
			user_ingredient_preferences
		WHERE
			user_ingredient_preferences.archived_at IS NULL
			AND user_ingredient_preferences.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
			AND user_ingredient_preferences.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
			AND (
				user_ingredient_preferences.last_updated_at IS NULL
				OR user_ingredient_preferences.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
			)
			AND (
				user_ingredient_preferences.last_updated_at IS NULL
				OR user_ingredient_preferences.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
			)
	) AS filtered_count,
	(
		SELECT
			COUNT(user_ingredient_preferences.id)
		FROM
			user_ingredient_preferences
		WHERE
			user_ingredient_preferences.archived_at IS NULL
	) AS total_count
FROM user_ingredient_preferences
	JOIN valid_ingredients ON valid_ingredients.id = user_ingredient_preferences.ingredient
WHERE user_ingredient_preferences.archived_at IS NULL
	AND valid_ingredients.archived_at IS NULL
	AND user_ingredient_preferences.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
	AND user_ingredient_preferences.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
	AND (
		user_ingredient_preferences.last_updated_at IS NULL
		OR user_ingredient_preferences.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
	)
	AND (
		user_ingredient_preferences.last_updated_at IS NULL
		OR user_ingredient_preferences.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
	)
	AND user_ingredient_preferences.belongs_to_user = $5
OFFSET $6
LIMIT $7
`

type GetUserIngredientPreferencesForUserParams struct {
	CreatedAfter  sql.NullTime
	CreatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	UpdatedBefore sql.NullTime
	UserID        string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetUserIngredientPreferencesForUserRow struct {
	ValidIngredientCreatedAt                               time.Time
	CreatedAt                                              time.Time
	ArchivedAt                                             sql.NullTime
	LastUpdatedAt                                          sql.NullTime
	ValidIngredientArchivedAt                              sql.NullTime
	ValidIngredientLastUpdatedAt                           sql.NullTime
	ValidIngredientStorageInstructions                     string
	ValidIngredientSlug                                    string
	ID                                                     string
	ValidIngredientPluralName                              string
	BelongsToUser                                          string
	ValidIngredientID                                      string
	ValidIngredientName                                    string
	ValidIngredientDescription                             string
	Notes                                                  string
	ValidIngredientIconPath                                string
	ValidIngredientWarning                                 string
	ValidIngredientShoppingSuggestions                     string
	ValidIngredientMaximumIdealStorageTemperatureInCelsius sql.NullString
	ValidIngredientMinimumIdealStorageTemperatureInCelsius sql.NullString
	FilteredCount                                          int64
	TotalCount                                             int64
	Rating                                                 int16
	ValidIngredientIsLiquid                                sql.NullBool
	ValidIngredientContainsDairy                           bool
	ValidIngredientContainsWheat                           bool
	ValidIngredientContainsAlcohol                         bool
	ValidIngredientContainsPeanut                          bool
	ValidIngredientIsStarch                                bool
	ValidIngredientIsProtein                               bool
	ValidIngredientIsGrain                                 bool
	ValidIngredientIsFruit                                 bool
	ValidIngredientIsSalt                                  bool
	ValidIngredientIsFat                                   bool
	ValidIngredientContainsTreeNut                         bool
	ValidIngredientAnimalDerived                           bool
	ValidIngredientVolumetric                              bool
	ValidIngredientContainsEgg                             bool
	ValidIngredientRestrictToPreparations                  bool
	ValidIngredientAnimalFlesh                             bool
	ValidIngredientContainsGluten                          bool
	Allergy                                                bool
	ValidIngredientContainsFish                            bool
	ValidIngredientContainsSesame                          bool
	ValidIngredientContainsShellfish                       bool
	ValidIngredientIsHeat                                  bool
	ValidIngredientContainsSoy                             bool
	ValidIngredientIsAcid                                  bool
}

func (q *Queries) GetUserIngredientPreferencesForUser(ctx context.Context, db DBTX, arg *GetUserIngredientPreferencesForUserParams) ([]*GetUserIngredientPreferencesForUserRow, error) {
	rows, err := db.QueryContext(ctx, getUserIngredientPreferencesForUser,
		arg.CreatedAfter,
		arg.CreatedBefore,
		arg.UpdatedAfter,
		arg.UpdatedBefore,
		arg.UserID,
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetUserIngredientPreferencesForUserRow{}
	for rows.Next() {
		var i GetUserIngredientPreferencesForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.ValidIngredientID,
			&i.ValidIngredientName,
			&i.ValidIngredientDescription,
			&i.ValidIngredientWarning,
			&i.ValidIngredientContainsEgg,
			&i.ValidIngredientContainsDairy,
			&i.ValidIngredientContainsPeanut,
			&i.ValidIngredientContainsTreeNut,
			&i.ValidIngredientContainsSoy,
			&i.ValidIngredientContainsWheat,
			&i.ValidIngredientContainsShellfish,
			&i.ValidIngredientContainsSesame,
			&i.ValidIngredientContainsFish,
			&i.ValidIngredientContainsGluten,
			&i.ValidIngredientAnimalFlesh,
			&i.ValidIngredientVolumetric,
			&i.ValidIngredientIsLiquid,
			&i.ValidIngredientIconPath,
			&i.ValidIngredientAnimalDerived,
			&i.ValidIngredientPluralName,
			&i.ValidIngredientRestrictToPreparations,
			&i.ValidIngredientMinimumIdealStorageTemperatureInCelsius,
			&i.ValidIngredientMaximumIdealStorageTemperatureInCelsius,
			&i.ValidIngredientStorageInstructions,
			&i.ValidIngredientSlug,
			&i.ValidIngredientContainsAlcohol,
			&i.ValidIngredientShoppingSuggestions,
			&i.ValidIngredientIsStarch,
			&i.ValidIngredientIsProtein,
			&i.ValidIngredientIsGrain,
			&i.ValidIngredientIsFruit,
			&i.ValidIngredientIsSalt,
			&i.ValidIngredientIsFat,
			&i.ValidIngredientIsAcid,
			&i.ValidIngredientIsHeat,
			&i.ValidIngredientCreatedAt,
			&i.ValidIngredientLastUpdatedAt,
			&i.ValidIngredientArchivedAt,
			&i.Rating,
			&i.Notes,
			&i.Allergy,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.BelongsToUser,
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

const updateUserIngredientPreference = `-- name: UpdateUserIngredientPreference :execrows

UPDATE user_ingredient_preferences SET
	ingredient = $1,
	rating = $2,
	notes = $3,
	allergy = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $5
	AND id = $6
`

type UpdateUserIngredientPreferenceParams struct {
	Ingredient    string
	Notes         string
	BelongsToUser string
	ID            string
	Rating        int16
	Allergy       bool
}

func (q *Queries) UpdateUserIngredientPreference(ctx context.Context, db DBTX, arg *UpdateUserIngredientPreferenceParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateUserIngredientPreference,
		arg.Ingredient,
		arg.Rating,
		arg.Notes,
		arg.Allergy,
		arg.BelongsToUser,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
