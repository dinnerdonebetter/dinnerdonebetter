-- name: GetValidVesselsWithIDs :many

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
  AND valid_vessels.id = ANY(sqlc.arg(ids)::text[]);