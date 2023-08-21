-- name: GetRandomValidVessel :one

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
    valid_vessels.width_in_millimeters::float,
    valid_vessels.length_in_millimeters::float,
    valid_vessels.height_in_millimeters::float,
    valid_vessels.shape,
    valid_vessels.created_at,
    valid_vessels.last_updated_at,
    valid_vessels.archived_at
FROM valid_vessels
	 JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE valid_vessels.archived_at IS NULL
	AND valid_measurement_units.archived_at IS NULL
	ORDER BY random() LIMIT 1;