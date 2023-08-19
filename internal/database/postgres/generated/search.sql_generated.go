// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: search.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const searchForServiceSettings = `-- name: SearchForServiceSettings :many

SELECT
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.admins_only,
    service_settings.enumeration,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at
FROM service_settings
WHERE service_settings.archived_at IS NULL
	AND service_settings.name ILIKE $1
LIMIT 50
`

type SearchForServiceSettingsRow struct {
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
	ID            string
	Name          string
	Type          SettingType
	Description   string
	Enumeration   string
	DefaultValue  sql.NullString
	AdminsOnly    bool
}

func (q *Queries) SearchForServiceSettings(ctx context.Context, db DBTX, name string) ([]*SearchForServiceSettingsRow, error) {
	rows, err := db.QueryContext(ctx, searchForServiceSettings, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForServiceSettingsRow{}
	for rows.Next() {
		var i SearchForServiceSettingsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Type,
			&i.Description,
			&i.DefaultValue,
			&i.AdminsOnly,
			&i.Enumeration,
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

const searchForValidIngredientGroups = `-- name: SearchForValidIngredientGroups :many

SELECT
	valid_ingredient_groups.id,
	valid_ingredient_groups.name,
	valid_ingredient_groups.description,
	valid_ingredient_groups.slug,
	valid_ingredient_groups.created_at,
	valid_ingredient_groups.last_updated_at,
	valid_ingredient_groups.archived_at,
	valid_ingredient_group_members.id,
    valid_ingredient_group_members.belongs_to_group,
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
    valid_ingredients.volumetric,
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
    valid_ingredients.created_at,
    valid_ingredients.last_updated_at,
    valid_ingredients.archived_at,
    valid_ingredient_group_members.created_at,
    valid_ingredient_group_members.archived_at
FROM valid_ingredient_groups
 JOIN valid_ingredient_group_members ON valid_ingredient_group_members.belongs_to_group=valid_ingredient_groups.id
  JOIN valid_ingredients ON valid_ingredients.id = valid_ingredient_group_members.valid_ingredient
WHERE valid_ingredient_groups.name ILIKE $1
AND valid_ingredient_groups.archived_at IS NULL
AND valid_ingredient_group_members.archived_at IS NULL
LIMIT 50
`

type SearchForValidIngredientGroupsRow struct {
	CreatedAt_2                             time.Time
	CreatedAt                               time.Time
	CreatedAt_3                             time.Time
	ArchivedAt_3                            sql.NullTime
	LastUpdatedAt_2                         sql.NullTime
	ArchivedAt_2                            sql.NullTime
	LastUpdatedAt                           sql.NullTime
	ArchivedAt                              sql.NullTime
	Warning                                 string
	Name                                    string
	Name_2                                  string
	Description_2                           string
	IconPath                                string
	BelongsToGroup                          string
	ID_2                                    string
	Slug                                    string
	Description                             string
	ID_3                                    string
	ShoppingSuggestions                     string
	Slug_2                                  string
	StorageInstructions                     string
	ID                                      string
	PluralName                              string
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	Volumetric                              bool
	AnimalDerived                           bool
	AnimalFlesh                             bool
	RestrictToPreparations                  bool
	ContainsGluten                          bool
	ContainsFish                            bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsAlcohol                         bool
	ContainsWheat                           bool
	IsStarch                                bool
	IsProtein                               bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsSoy                             bool
	ContainsTreeNut                         bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) SearchForValidIngredientGroups(ctx context.Context, db DBTX, name string) ([]*SearchForValidIngredientGroupsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidIngredientGroups, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientGroupsRow{}
	for rows.Next() {
		var i SearchForValidIngredientGroupsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Slug,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
			&i.ID_2,
			&i.BelongsToGroup,
			&i.ID_3,
			&i.Name_2,
			&i.Description_2,
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
			&i.Volumetric,
			&i.IsLiquid,
			&i.IconPath,
			&i.AnimalDerived,
			&i.PluralName,
			&i.RestrictToPreparations,
			&i.MinimumIdealStorageTemperatureInCelsius,
			&i.MaximumIdealStorageTemperatureInCelsius,
			&i.StorageInstructions,
			&i.Slug_2,
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
			&i.CreatedAt_2,
			&i.LastUpdatedAt_2,
			&i.ArchivedAt_2,
			&i.CreatedAt_3,
			&i.ArchivedAt_3,
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

const searchForValidIngredientStates = `-- name: SearchForValidIngredientStates :many

SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.slug,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.name ILIKE $1
LIMIT 50
`

type SearchForValidIngredientStatesRow struct {
	ID            string
	Name          string
	Description   string
	IconPath      string
	Slug          string
	PastTense     string
	AttributeType IngredientAttributeType
	CreatedAt     time.Time
	LastUpdatedAt sql.NullTime
	ArchivedAt    sql.NullTime
}

func (q *Queries) SearchForValidIngredientStates(ctx context.Context, db DBTX, name string) ([]*SearchForValidIngredientStatesRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidIngredientStates, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidIngredientStatesRow{}
	for rows.Next() {
		var i SearchForValidIngredientStatesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.Slug,
			&i.PastTense,
			&i.AttributeType,
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
	valid_ingredients.volumetric,
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
	valid_ingredients.created_at,
	valid_ingredients.last_updated_at,
	valid_ingredients.archived_at
FROM valid_ingredients
WHERE valid_ingredients.name ILIKE $1
	AND valid_ingredients.archived_at IS NULL
LIMIT 50
`

type SearchForValidIngredientsRow struct {
	CreatedAt                               time.Time
	ArchivedAt                              sql.NullTime
	LastUpdatedAt                           sql.NullTime
	Warning                                 string
	Description                             string
	Name                                    string
	ShoppingSuggestions                     string
	Slug                                    string
	StorageInstructions                     string
	PluralName                              string
	ID                                      string
	IconPath                                string
	MaximumIdealStorageTemperatureInCelsius sql.NullString
	MinimumIdealStorageTemperatureInCelsius sql.NullString
	IsLiquid                                sql.NullBool
	AnimalDerived                           bool
	ContainsTreeNut                         bool
	AnimalFlesh                             bool
	ContainsGluten                          bool
	ContainsFish                            bool
	RestrictToPreparations                  bool
	ContainsSesame                          bool
	ContainsShellfish                       bool
	ContainsWheat                           bool
	ContainsSoy                             bool
	ContainsAlcohol                         bool
	Volumetric                              bool
	IsStarch                                bool
	IsProtein                               bool
	IsGrain                                 bool
	IsFruit                                 bool
	IsSalt                                  bool
	IsFat                                   bool
	IsAcid                                  bool
	IsHeat                                  bool
	ContainsPeanut                          bool
	ContainsDairy                           bool
	ContainsEgg                             bool
}

func (q *Queries) SearchForValidIngredients(ctx context.Context, db DBTX, name string) ([]*SearchForValidIngredientsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidIngredients, name)
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
			&i.Volumetric,
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

const searchForValidInstruments = `-- name: SearchForValidInstruments :many

SELECT
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at
FROM valid_instruments
WHERE valid_instruments.archived_at IS NULL
	AND valid_instruments.name ILIKE '%' || $1::text || '%'
    LIMIT 50
`

type SearchForValidInstrumentsRow struct {
	CreatedAt                      time.Time
	LastUpdatedAt                  sql.NullTime
	ArchivedAt                     sql.NullTime
	ID                             string
	Name                           string
	PluralName                     string
	Description                    string
	IconPath                       string
	Slug                           string
	UsableForStorage               bool
	DisplayInSummaryLists          bool
	IncludeInGeneratedInstructions bool
}

func (q *Queries) SearchForValidInstruments(ctx context.Context, db DBTX, query string) ([]*SearchForValidInstrumentsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidInstruments, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidInstrumentsRow{}
	for rows.Next() {
		var i SearchForValidInstrumentsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PluralName,
			&i.Description,
			&i.IconPath,
			&i.UsableForStorage,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.Slug,
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

const searchForValidMeasurementUnits = `-- name: SearchForValidMeasurementUnits :many

SELECT
	valid_measurement_units.id,
	valid_measurement_units.name,
	valid_measurement_units.description,
	valid_measurement_units.volumetric,
	valid_measurement_units.icon_path,
	valid_measurement_units.universal,
	valid_measurement_units.metric,
	valid_measurement_units.imperial,
	valid_measurement_units.slug,
	valid_measurement_units.plural_name,
	valid_measurement_units.created_at,
	valid_measurement_units.last_updated_at,
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE (valid_measurement_units.name ILIKE '%' || $1::text || '%' OR valid_measurement_units.universal is TRUE)
AND valid_measurement_units.archived_at IS NULL
LIMIT 50
`

type SearchForValidMeasurementUnitsRow struct {
	CreatedAt     time.Time
	ArchivedAt    sql.NullTime
	LastUpdatedAt sql.NullTime
	PluralName    string
	Name          string
	Description   string
	ID            string
	IconPath      string
	Slug          string
	Volumetric    sql.NullBool
	Imperial      bool
	Metric        bool
	Universal     bool
}

func (q *Queries) SearchForValidMeasurementUnits(ctx context.Context, db DBTX, query string) ([]*SearchForValidMeasurementUnitsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidMeasurementUnits, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidMeasurementUnitsRow{}
	for rows.Next() {
		var i SearchForValidMeasurementUnitsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.Volumetric,
			&i.IconPath,
			&i.Universal,
			&i.Metric,
			&i.Imperial,
			&i.Slug,
			&i.PluralName,
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

const searchForValidPreparations = `-- name: SearchForValidPreparations :many

SELECT
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.minimum_ingredient_count,
	valid_preparations.maximum_ingredient_count,
	valid_preparations.minimum_instrument_count,
	valid_preparations.maximum_instrument_count,
	valid_preparations.temperature_required,
	valid_preparations.time_estimate_required,
	valid_preparations.condition_expression_required,
    valid_preparations.consumes_vessel,
    valid_preparations.only_for_vessels,
    valid_preparations.minimum_vessel_count,
    valid_preparations.maximum_vessel_count,
	valid_preparations.slug,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at
FROM valid_preparations
WHERE valid_preparations.archived_at IS NULL
	AND valid_preparations.name ILIKE '%' || $1::text || '%'
LIMIT 50
`

type SearchForValidPreparationsRow struct {
	CreatedAt                   time.Time
	LastUpdatedAt               sql.NullTime
	ArchivedAt                  sql.NullTime
	Name                        string
	Description                 string
	IconPath                    string
	ID                          string
	Slug                        string
	PastTense                   string
	MaximumInstrumentCount      sql.NullInt32
	MaximumIngredientCount      sql.NullInt32
	MaximumVesselCount          sql.NullInt32
	MinimumVesselCount          int32
	MinimumInstrumentCount      int32
	MinimumIngredientCount      int32
	RestrictToIngredients       bool
	OnlyForVessels              bool
	ConsumesVessel              bool
	ConditionExpressionRequired bool
	TimeEstimateRequired        bool
	TemperatureRequired         bool
	YieldsNothing               bool
}

func (q *Queries) SearchForValidPreparations(ctx context.Context, db DBTX, query string) ([]*SearchForValidPreparationsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidPreparations, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidPreparationsRow{}
	for rows.Next() {
		var i SearchForValidPreparationsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.IconPath,
			&i.YieldsNothing,
			&i.RestrictToIngredients,
			&i.MinimumIngredientCount,
			&i.MaximumIngredientCount,
			&i.MinimumInstrumentCount,
			&i.MaximumInstrumentCount,
			&i.TemperatureRequired,
			&i.TimeEstimateRequired,
			&i.ConditionExpressionRequired,
			&i.ConsumesVessel,
			&i.OnlyForVessels,
			&i.MinimumVesselCount,
			&i.MaximumVesselCount,
			&i.Slug,
			&i.PastTense,
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

const searchForValidVessels = `-- name: SearchForValidVessels :many

SELECT
	valid_vessels.id,
    valid_vessels.name,
    valid_vessels.plural_name,
    valid_vessels.description,
    valid_vessels.icon_path,
    valid_vessels.usable_for_storage,
    valid_vessels.slug,
    valid_vessels.display_in_summary_lists,
    valid_vessels.include_in_generated_instructions,
    valid_vessels.capacity::float,
    valid_vessels.capacity_unit,
    valid_vessels.width_in_millimeters::float,
    valid_vessels.length_in_millimeters::float,
    valid_vessels.height_in_millimeters::float,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at
FROM valid_vessels
WHERE valid_vessels.archived_at IS NULL
    AND valid_vessels.name ILIKE '%' || $1::text || '%'
	LIMIT 50
`

type SearchForValidVesselsRow struct {
	CreatedAt                       time.Time
	ArchivedAt                      sql.NullTime
	LastUpdatedAt                   sql.NullTime
	ID                              string
	Name                            string
	PluralName                      string
	Description                     string
	IconPath                        string
	Slug                            string
	Shape                           VesselShape
	CapacityUnit                    sql.NullString
	ValidVesselsLengthInMillimeters float64
	ValidVesselsWidthInMillimeters  float64
	ValidVesselsHeightInMillimeters float64
	ValidVesselsCapacity            float64
	IncludeInGeneratedInstructions  bool
	DisplayInSummaryLists           bool
	UsableForStorage                bool
}

func (q *Queries) SearchForValidVessels(ctx context.Context, db DBTX, query string) ([]*SearchForValidVesselsRow, error) {
	rows, err := db.QueryContext(ctx, searchForValidVessels, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*SearchForValidVesselsRow{}
	for rows.Next() {
		var i SearchForValidVesselsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.PluralName,
			&i.Description,
			&i.IconPath,
			&i.UsableForStorage,
			&i.Slug,
			&i.DisplayInSummaryLists,
			&i.IncludeInGeneratedInstructions,
			&i.ValidVesselsCapacity,
			&i.CapacityUnit,
			&i.ValidVesselsWidthInMillimeters,
			&i.ValidVesselsLengthInMillimeters,
			&i.ValidVesselsHeightInMillimeters,
			&i.Shape,
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
