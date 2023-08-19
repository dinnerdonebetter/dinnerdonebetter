// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: get_for_preparation.sql

package generated

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
)

const getValidPreparationVesselsForPreparation = `-- name: GetValidPreparationVesselsForPreparation :many
SELECT
    valid_preparation_vessels.id as valid_preparation_vessel_id,
    valid_preparation_vessels.notes as valid_preparation_vessel_notes,
    valid_preparations.id as valid_preparation_id,
    valid_preparations.name as valid_preparation_name,
    valid_preparations.description as valid_preparation_description,
    valid_preparations.icon_path as valid_preparation_icon_path,
    valid_preparations.yields_nothing as valid_preparation_yields_nothing,
    valid_preparations.restrict_to_ingredients as valid_preparation_restrict_to_ingredients,
    valid_preparations.minimum_ingredient_count as valid_preparation_minimum_ingredient_count,
    valid_preparations.maximum_ingredient_count as valid_preparation_maximum_ingredient_count,
    valid_preparations.minimum_instrument_count as valid_preparation_minimum_instrument_count,
    valid_preparations.maximum_instrument_count as valid_preparation_maximum_instrument_count,
    valid_preparations.temperature_required as valid_preparation_temperature_required,
    valid_preparations.time_estimate_required as valid_preparation_time_estimate_required,
    valid_preparations.condition_expression_required as valid_preparation_condition_expression_required,
    valid_preparations.consumes_vessel as valid_preparation_consumes_vessel,
    valid_preparations.only_for_vessels as valid_preparation_only_for_vessels,
    valid_preparations.minimum_vessel_count as valid_preparation_minimum_vessel_count,
    valid_preparations.maximum_vessel_count as valid_preparation_maximum_vessel_count,
    valid_preparations.slug as valid_preparation_slug,
    valid_preparations.past_tense as valid_preparation_past_tense,
    valid_preparations.created_at as valid_preparation_created_at,
    valid_preparations.last_updated_at as valid_preparation_last_updated_at,
    valid_preparations.archived_at as valid_preparation_archived_at,
    valid_vessels.id as valid_vessel_id,
    valid_vessels.name as valid_vessel_name,
    valid_vessels.plural_name as valid_vessel_plural_name,
    valid_vessels.description as valid_vessel_description,
    valid_vessels.icon_path as valid_vessel_icon_path,
    valid_vessels.usable_for_storage as valid_vessel_usable_for_storage,
    valid_vessels.slug as valid_vessel_slug,
    valid_vessels.display_in_summary_lists as valid_vessel_display_in_summary_lists,
    valid_vessels.include_in_generated_instructions as valid_vessel_include_in_generated_instructions,
    valid_vessels.capacity::float as valid_vessel_capacity,
    valid_measurement_units.id as valid_measurement_unit_id,
    valid_measurement_units.name as valid_measurement_unit_name,
    valid_measurement_units.description as valid_measurement_unit_description,
    valid_measurement_units.volumetric as valid_measurement_unit_volumetric,
    valid_measurement_units.icon_path as valid_measurement_unit_icon_path,
    valid_measurement_units.universal as valid_measurement_unit_universal,
    valid_measurement_units.metric as valid_measurement_unit_metric,
    valid_measurement_units.imperial as valid_measurement_unit_imperial,
    valid_measurement_units.slug as valid_measurement_unit_slug,
    valid_measurement_units.plural_name as valid_measurement_unit_plural_name,
    valid_measurement_units.created_at as valid_measurement_unit_created_at,
    valid_measurement_units.last_updated_at as valid_measurement_unit_last_updated_at,
    valid_measurement_units.archived_at as valid_measurement_unit_archived_at,
    valid_vessels.width_in_millimeters::float as valid_vessel_width_in_millimeters,
    valid_vessels.length_in_millimeters::float as valid_vessel_length_in_millimeters,
    valid_vessels.height_in_millimeters::float as valid_vessel_height_in_millimeters,
    valid_vessels.shape as valid_vessel_shape,
    valid_vessels.created_at as valid_vessel_created_at,
    valid_vessels.last_updated_at as valid_vessel_last_updated_at,
    valid_vessels.archived_at as valid_vessel_archived_at,
    valid_preparation_vessels.created_at as valid_preparation_vessel_created_at,
    valid_preparation_vessels.last_updated_at as valid_preparation_vessel_last_updated_at,
    valid_preparation_vessels.archived_at as valid_preparation_vessel_archived_at,
    (
        SELECT
            COUNT(valid_preparation_vessels.id)
        FROM
            valid_preparation_vessels
        WHERE
            valid_preparation_vessels.archived_at IS NULL
            AND valid_preparation_vessels.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
            AND valid_preparation_vessels.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
            AND (
                valid_preparation_vessels.last_updated_at IS NULL
                OR valid_preparation_vessels.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
            )
            AND (
                valid_preparation_vessels.last_updated_at IS NULL
                OR valid_preparation_vessels.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
            )
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_preparation_vessels.id)
        FROM
            valid_preparation_vessels
        WHERE
            valid_preparation_vessels.archived_at IS NULL
    ) as total_count
FROM
    valid_preparation_vessels
        JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id
        JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id
        LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id
WHERE
    valid_preparation_vessels.archived_at IS NULL
    AND valid_vessels.archived_at IS NULL
    AND valid_preparations.archived_at IS NULL
    AND valid_measurement_units.archived_at IS NULL
    AND valid_preparation_vessels.created_at > COALESCE($1, (SELECT NOW() - interval '999 years'))
    AND valid_preparation_vessels.created_at < COALESCE($2, (SELECT NOW() + interval '999 years'))
    AND (
        valid_preparation_vessels.last_updated_at IS NULL
        OR valid_preparation_vessels.last_updated_at > COALESCE($3, (SELECT NOW() - interval '999 years'))
    )
    AND (
        valid_preparation_vessels.last_updated_at IS NULL
        OR valid_preparation_vessels.last_updated_at < COALESCE($4, (SELECT NOW() + interval '999 years'))
    )
    AND valid_preparation_vessels.valid_preparation_id = ANY($5::text[])
OFFSET $6
LIMIT $7
`

type GetValidPreparationVesselsForPreparationParams struct {
	CreatedBefore sql.NullTime
	CreatedAfter  sql.NullTime
	UpdatedBefore sql.NullTime
	UpdatedAfter  sql.NullTime
	Ids           []string
	QueryOffset   sql.NullInt32
	QueryLimit    sql.NullInt32
}

type GetValidPreparationVesselsForPreparationRow struct {
	ValidPreparationCreatedAt                   time.Time
	ValidPreparationVesselCreatedAt             time.Time
	ValidVesselCreatedAt                        time.Time
	ValidMeasurementUnitArchivedAt              sql.NullTime
	ValidPreparationArchivedAt                  sql.NullTime
	ValidPreparationLastUpdatedAt               sql.NullTime
	ValidMeasurementUnitCreatedAt               sql.NullTime
	ValidMeasurementUnitLastUpdatedAt           sql.NullTime
	ValidPreparationVesselArchivedAt            sql.NullTime
	ValidPreparationVesselLastUpdatedAt         sql.NullTime
	ValidVesselArchivedAt                       sql.NullTime
	ValidVesselLastUpdatedAt                    sql.NullTime
	ValidVesselShape                            VesselShape
	ValidPreparationSlug                        string
	ValidPreparationIconPath                    string
	ValidPreparationDescription                 string
	ValidPreparationVesselID                    string
	ValidVesselSlug                             string
	ValidVesselIconPath                         string
	ValidVesselDescription                      string
	ValidPreparationPastTense                   string
	ValidPreparationName                        string
	ValidPreparationID                          string
	ValidPreparationVesselNotes                 string
	ValidVesselID                               string
	ValidVesselName                             string
	ValidVesselPluralName                       string
	ValidMeasurementUnitPluralName              sql.NullString
	ValidMeasurementUnitSlug                    sql.NullString
	ValidMeasurementUnitIconPath                sql.NullString
	ValidMeasurementUnitDescription             sql.NullString
	ValidMeasurementUnitName                    sql.NullString
	ValidMeasurementUnitID                      sql.NullString
	ValidVesselWidthInMillimeters               float64
	TotalCount                                  int64
	FilteredCount                               int64
	ValidVesselCapacity                         float64
	ValidVesselHeightInMillimeters              float64
	ValidVesselLengthInMillimeters              float64
	ValidPreparationMaximumIngredientCount      sql.NullInt32
	ValidPreparationMaximumInstrumentCount      sql.NullInt32
	ValidPreparationMaximumVesselCount          sql.NullInt32
	ValidPreparationMinimumVesselCount          int32
	ValidPreparationMinimumIngredientCount      int32
	ValidPreparationMinimumInstrumentCount      int32
	ValidMeasurementUnitUniversal               sql.NullBool
	ValidMeasurementUnitImperial                sql.NullBool
	ValidMeasurementUnitVolumetric              sql.NullBool
	ValidMeasurementUnitMetric                  sql.NullBool
	ValidPreparationOnlyForVessels              bool
	ValidPreparationTemperatureRequired         bool
	ValidPreparationTimeEstimateRequired        bool
	ValidPreparationConsumesVessel              bool
	ValidVesselUsableForStorage                 bool
	ValidPreparationConditionExpressionRequired bool
	ValidPreparationRestrictToIngredients       bool
	ValidPreparationYieldsNothing               bool
	ValidVesselDisplayInSummaryLists            bool
	ValidVesselIncludeInGeneratedInstructions   bool
}

func (q *Queries) GetValidPreparationVesselsForPreparation(ctx context.Context, db DBTX, arg *GetValidPreparationVesselsForPreparationParams) ([]*GetValidPreparationVesselsForPreparationRow, error) {
	rows, err := db.QueryContext(ctx, getValidPreparationVesselsForPreparation,
		arg.CreatedBefore,
		arg.CreatedAfter,
		arg.UpdatedBefore,
		arg.UpdatedAfter,
		pq.Array(arg.Ids),
		arg.QueryOffset,
		arg.QueryLimit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetValidPreparationVesselsForPreparationRow{}
	for rows.Next() {
		var i GetValidPreparationVesselsForPreparationRow
		if err := rows.Scan(
			&i.ValidPreparationVesselID,
			&i.ValidPreparationVesselNotes,
			&i.ValidPreparationID,
			&i.ValidPreparationName,
			&i.ValidPreparationDescription,
			&i.ValidPreparationIconPath,
			&i.ValidPreparationYieldsNothing,
			&i.ValidPreparationRestrictToIngredients,
			&i.ValidPreparationMinimumIngredientCount,
			&i.ValidPreparationMaximumIngredientCount,
			&i.ValidPreparationMinimumInstrumentCount,
			&i.ValidPreparationMaximumInstrumentCount,
			&i.ValidPreparationTemperatureRequired,
			&i.ValidPreparationTimeEstimateRequired,
			&i.ValidPreparationConditionExpressionRequired,
			&i.ValidPreparationConsumesVessel,
			&i.ValidPreparationOnlyForVessels,
			&i.ValidPreparationMinimumVesselCount,
			&i.ValidPreparationMaximumVesselCount,
			&i.ValidPreparationSlug,
			&i.ValidPreparationPastTense,
			&i.ValidPreparationCreatedAt,
			&i.ValidPreparationLastUpdatedAt,
			&i.ValidPreparationArchivedAt,
			&i.ValidVesselID,
			&i.ValidVesselName,
			&i.ValidVesselPluralName,
			&i.ValidVesselDescription,
			&i.ValidVesselIconPath,
			&i.ValidVesselUsableForStorage,
			&i.ValidVesselSlug,
			&i.ValidVesselDisplayInSummaryLists,
			&i.ValidVesselIncludeInGeneratedInstructions,
			&i.ValidVesselCapacity,
			&i.ValidMeasurementUnitID,
			&i.ValidMeasurementUnitName,
			&i.ValidMeasurementUnitDescription,
			&i.ValidMeasurementUnitVolumetric,
			&i.ValidMeasurementUnitIconPath,
			&i.ValidMeasurementUnitUniversal,
			&i.ValidMeasurementUnitMetric,
			&i.ValidMeasurementUnitImperial,
			&i.ValidMeasurementUnitSlug,
			&i.ValidMeasurementUnitPluralName,
			&i.ValidMeasurementUnitCreatedAt,
			&i.ValidMeasurementUnitLastUpdatedAt,
			&i.ValidMeasurementUnitArchivedAt,
			&i.ValidVesselWidthInMillimeters,
			&i.ValidVesselLengthInMillimeters,
			&i.ValidVesselHeightInMillimeters,
			&i.ValidVesselShape,
			&i.ValidVesselCreatedAt,
			&i.ValidVesselLastUpdatedAt,
			&i.ValidVesselArchivedAt,
			&i.ValidPreparationVesselCreatedAt,
			&i.ValidPreparationVesselLastUpdatedAt,
			&i.ValidPreparationVesselArchivedAt,
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
