-- name: GetValidPreparationVessel :one

SELECT
    valid_preparation_vessels.id,
    valid_preparation_vessels.notes,
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
    valid_preparations.archived_at,
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
    valid_vessels.archived_at,
    valid_preparation_vessels.created_at,
    valid_preparation_vessels.last_updated_at,
    valid_preparation_vessels.archived_at
FROM
	valid_preparation_vessels
	 JOIN valid_vessels ON valid_preparation_vessels.valid_vessel_id = valid_vessels.id
	 LEFT JOIN valid_measurement_units ON valid_vessels.capacity_unit = valid_measurement_units.id
	 JOIN valid_preparations ON valid_preparation_vessels.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_vessels.archived_at IS NULL
	AND valid_vessels.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	AND valid_preparations.archived_at IS NULL
	AND valid_preparation_vessels.id = $1;
