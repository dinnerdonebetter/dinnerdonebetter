// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_random.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const GetRandomValidIngredient = `-- name: GetRandomValidIngredient :one

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
WHERE valid_ingredients.archived_at IS NULL
ORDER BY random() LIMIT 1
`

type GetRandomValidIngredientRow struct {
	CreatedAt                               time.Time      `db:"created_at"`
	ArchivedAt                              sql.NullTime   `db:"archived_at"`
	LastUpdatedAt                           sql.NullTime   `db:"last_updated_at"`
	Warning                                 string         `db:"warning"`
	Description                             string         `db:"description"`
	Name                                    string         `db:"name"`
	ShoppingSuggestions                     string         `db:"shopping_suggestions"`
	Slug                                    string         `db:"slug"`
	StorageInstructions                     string         `db:"storage_instructions"`
	PluralName                              string         `db:"plural_name"`
	ID                                      string         `db:"id"`
	IconPath                                string         `db:"icon_path"`
	MaximumIdealStorageTemperatureInCelsius sql.NullString `db:"maximum_ideal_storage_temperature_in_celsius"`
	MinimumIdealStorageTemperatureInCelsius sql.NullString `db:"minimum_ideal_storage_temperature_in_celsius"`
	IsLiquid                                sql.NullBool   `db:"is_liquid"`
	AnimalDerived                           bool           `db:"animal_derived"`
	ContainsTreeNut                         bool           `db:"contains_tree_nut"`
	AnimalFlesh                             bool           `db:"animal_flesh"`
	ContainsGluten                          bool           `db:"contains_gluten"`
	ContainsFish                            bool           `db:"contains_fish"`
	RestrictToPreparations                  bool           `db:"restrict_to_preparations"`
	ContainsSesame                          bool           `db:"contains_sesame"`
	ContainsShellfish                       bool           `db:"contains_shellfish"`
	ContainsWheat                           bool           `db:"contains_wheat"`
	ContainsSoy                             bool           `db:"contains_soy"`
	ContainsAlcohol                         bool           `db:"contains_alcohol"`
	Volumetric                              bool           `db:"volumetric"`
	IsStarch                                bool           `db:"is_starch"`
	IsProtein                               bool           `db:"is_protein"`
	IsGrain                                 bool           `db:"is_grain"`
	IsFruit                                 bool           `db:"is_fruit"`
	IsSalt                                  bool           `db:"is_salt"`
	IsFat                                   bool           `db:"is_fat"`
	IsAcid                                  bool           `db:"is_acid"`
	IsHeat                                  bool           `db:"is_heat"`
	ContainsPeanut                          bool           `db:"contains_peanut"`
	ContainsDairy                           bool           `db:"contains_dairy"`
	ContainsEgg                             bool           `db:"contains_egg"`
}

func (q *Queries) GetRandomValidIngredient(ctx context.Context, db DBTX) (*GetRandomValidIngredientRow, error) {
	row := db.QueryRowContext(ctx, GetRandomValidIngredient)
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
	)
	return &i, err
}

const GetRandomValidInstrument = `-- name: GetRandomValidInstrument :one

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
	ORDER BY random() LIMIT 1
`

type GetRandomValidInstrumentRow struct {
	CreatedAt                      time.Time    `db:"created_at"`
	LastUpdatedAt                  sql.NullTime `db:"last_updated_at"`
	ArchivedAt                     sql.NullTime `db:"archived_at"`
	ID                             string       `db:"id"`
	Name                           string       `db:"name"`
	PluralName                     string       `db:"plural_name"`
	Description                    string       `db:"description"`
	IconPath                       string       `db:"icon_path"`
	Slug                           string       `db:"slug"`
	UsableForStorage               bool         `db:"usable_for_storage"`
	DisplayInSummaryLists          bool         `db:"display_in_summary_lists"`
	IncludeInGeneratedInstructions bool         `db:"include_in_generated_instructions"`
}

func (q *Queries) GetRandomValidInstrument(ctx context.Context, db DBTX) (*GetRandomValidInstrumentRow, error) {
	row := db.QueryRowContext(ctx, GetRandomValidInstrument)
	var i GetRandomValidInstrumentRow
	err := row.Scan(
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
	)
	return &i, err
}

const GetRandomValidMeasurementUnit = `-- name: GetRandomValidMeasurementUnit :one

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
WHERE valid_measurement_units.archived_at IS NULL
	ORDER BY random() LIMIT 1
`

type GetRandomValidMeasurementUnitRow struct {
	CreatedAt     time.Time    `db:"created_at"`
	ArchivedAt    sql.NullTime `db:"archived_at"`
	LastUpdatedAt sql.NullTime `db:"last_updated_at"`
	PluralName    string       `db:"plural_name"`
	Name          string       `db:"name"`
	Description   string       `db:"description"`
	ID            string       `db:"id"`
	IconPath      string       `db:"icon_path"`
	Slug          string       `db:"slug"`
	Volumetric    sql.NullBool `db:"volumetric"`
	Imperial      bool         `db:"imperial"`
	Metric        bool         `db:"metric"`
	Universal     bool         `db:"universal"`
}

func (q *Queries) GetRandomValidMeasurementUnit(ctx context.Context, db DBTX) (*GetRandomValidMeasurementUnitRow, error) {
	row := db.QueryRowContext(ctx, GetRandomValidMeasurementUnit)
	var i GetRandomValidMeasurementUnitRow
	err := row.Scan(
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
	)
	return &i, err
}

const GetRandomValidPreparation = `-- name: GetRandomValidPreparation :one

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
	ORDER BY random() LIMIT 1
`

type GetRandomValidPreparationRow struct {
	CreatedAt                   time.Time     `db:"created_at"`
	LastUpdatedAt               sql.NullTime  `db:"last_updated_at"`
	ArchivedAt                  sql.NullTime  `db:"archived_at"`
	Name                        string        `db:"name"`
	Description                 string        `db:"description"`
	IconPath                    string        `db:"icon_path"`
	ID                          string        `db:"id"`
	Slug                        string        `db:"slug"`
	PastTense                   string        `db:"past_tense"`
	MaximumInstrumentCount      sql.NullInt32 `db:"maximum_instrument_count"`
	MaximumIngredientCount      sql.NullInt32 `db:"maximum_ingredient_count"`
	MaximumVesselCount          sql.NullInt32 `db:"maximum_vessel_count"`
	MinimumVesselCount          int32         `db:"minimum_vessel_count"`
	MinimumInstrumentCount      int32         `db:"minimum_instrument_count"`
	MinimumIngredientCount      int32         `db:"minimum_ingredient_count"`
	RestrictToIngredients       bool          `db:"restrict_to_ingredients"`
	OnlyForVessels              bool          `db:"only_for_vessels"`
	ConsumesVessel              bool          `db:"consumes_vessel"`
	ConditionExpressionRequired bool          `db:"condition_expression_required"`
	TimeEstimateRequired        bool          `db:"time_estimate_required"`
	TemperatureRequired         bool          `db:"temperature_required"`
	YieldsNothing               bool          `db:"yields_nothing"`
}

func (q *Queries) GetRandomValidPreparation(ctx context.Context, db DBTX) (*GetRandomValidPreparationRow, error) {
	row := db.QueryRowContext(ctx, GetRandomValidPreparation)
	var i GetRandomValidPreparationRow
	err := row.Scan(
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
	)
	return &i, err
}

const GetRandomValidVessel = `-- name: GetRandomValidVessel :one

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
    valid_vessels.capacity,
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
	valid_measurement_units.archived_at,
    valid_vessels.width_in_millimeters,
    valid_vessels.length_in_millimeters,
    valid_vessels.height_in_millimeters,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at
FROM valid_vessels
	 JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE valid_vessels.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	ORDER BY random() LIMIT 1
`

type GetRandomValidVesselRow struct {
	CreatedAt                      time.Time      `db:"created_at"`
	CreatedAt_2                    time.Time      `db:"created_at_2"`
	ArchivedAt_2                   sql.NullTime   `db:"archived_at_2"`
	LastUpdatedAt_2                sql.NullTime   `db:"last_updated_at_2"`
	ArchivedAt                     sql.NullTime   `db:"archived_at"`
	LastUpdatedAt                  sql.NullTime   `db:"last_updated_at"`
	IconPath_2                     string         `db:"icon_path_2"`
	IconPath                       string         `db:"icon_path"`
	Name                           string         `db:"name"`
	Capacity                       string         `db:"capacity"`
	ID_2                           string         `db:"id_2"`
	Name_2                         string         `db:"name_2"`
	Description_2                  string         `db:"description_2"`
	PluralName                     string         `db:"plural_name"`
	ID                             string         `db:"id"`
	Description                    string         `db:"description"`
	Shape                          VesselShape    `db:"shape"`
	Slug                           string         `db:"slug"`
	Slug_2                         string         `db:"slug_2"`
	PluralName_2                   string         `db:"plural_name_2"`
	WidthInMillimeters             sql.NullString `db:"width_in_millimeters"`
	LengthInMillimeters            sql.NullString `db:"length_in_millimeters"`
	HeightInMillimeters            sql.NullString `db:"height_in_millimeters"`
	Volumetric                     sql.NullBool   `db:"volumetric"`
	Imperial                       bool           `db:"imperial"`
	UsableForStorage               bool           `db:"usable_for_storage"`
	DisplayInSummaryLists          bool           `db:"display_in_summary_lists"`
	Metric                         bool           `db:"metric"`
	Universal                      bool           `db:"universal"`
	IncludeInGeneratedInstructions bool           `db:"include_in_generated_instructions"`
}

func (q *Queries) GetRandomValidVessel(ctx context.Context, db DBTX) (*GetRandomValidVesselRow, error) {
	row := db.QueryRowContext(ctx, GetRandomValidVessel)
	var i GetRandomValidVesselRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.PluralName,
		&i.Description,
		&i.IconPath,
		&i.UsableForStorage,
		&i.Slug,
		&i.DisplayInSummaryLists,
		&i.IncludeInGeneratedInstructions,
		&i.Capacity,
		&i.ID_2,
		&i.Name_2,
		&i.Description_2,
		&i.Volumetric,
		&i.IconPath_2,
		&i.Universal,
		&i.Metric,
		&i.Imperial,
		&i.Slug_2,
		&i.PluralName_2,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
		&i.WidthInMillimeters,
		&i.LengthInMillimeters,
		&i.HeightInMillimeters,
		&i.Shape,
		&i.CreatedAt_2,
		&i.LastUpdatedAt_2,
		&i.ArchivedAt_2,
	)
	return &i, err
}