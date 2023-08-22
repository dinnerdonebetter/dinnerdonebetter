-- name: ArchiveValidVessel :exec

UPDATE valid_vessels SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateValidVessel :exec

INSERT INTO valid_vessels (id,"name",plural_name,description,icon_path,usable_for_storage,slug,display_in_summary_lists,include_in_generated_instructions,capacity,capacity_unit,width_in_millimeters,length_in_millimeters,height_in_millimeters,shape)
    VALUES (sqlc.arg(id),sqlc.arg(name),sqlc.arg(plural_name),sqlc.arg(description),sqlc.arg(icon_path),sqlc.arg(usable_for_storage),sqlc.arg(slug),sqlc.arg(display_in_summary_lists),sqlc.arg(include_in_generated_instructions),sqlc.arg(capacity)::float,sqlc.arg(capacity_unit),sqlc.arg(width_in_millimeters)::float,sqlc.arg(length_in_millimeters)::float,sqlc.arg(height_in_millimeters)::float,sqlc.arg(shape));

-- name: CheckValidVesselExistence :one

SELECT EXISTS ( SELECT valid_vessels.id FROM valid_vessels WHERE valid_vessels.archived_at IS NULL AND valid_vessels.id = $1 );

-- name: GetValidVessels :many

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
  valid_vessels.archived_at,
  (
    SELECT
      COUNT(valid_vessels.id)
    FROM
      valid_vessels
    WHERE
      valid_vessels.archived_at IS NULL
      AND valid_vessels.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
      AND valid_vessels.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
      AND (
        valid_vessels.last_updated_at IS NULL
        OR valid_vessels.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
      )
      AND (
        valid_vessels.last_updated_at IS NULL
        OR valid_vessels.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
      )
  ) as filtered_count,
  (
    SELECT
      COUNT(valid_vessels.id)
    FROM
      valid_vessels
    WHERE
      valid_vessels.archived_at IS NULL
  ) as total_count
FROM
  valid_vessels
WHERE
  valid_vessels.archived_at IS NULL
  AND valid_vessels.created_at > (COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years')))
  AND valid_vessels.created_at < (COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years')))
  AND (
    valid_vessels.last_updated_at IS NULL
    OR valid_vessels.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
  )
  AND (
    valid_vessels.last_updated_at IS NULL
    OR valid_vessels.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
  )
GROUP BY
  valid_vessels.id
ORDER BY
  valid_vessels.id
OFFSET
    sqlc.narg(query_offset)
LIMIT
    sqlc.narg(query_limit);

-- name: GetValidVesselIDsNeedingIndexing :many

SELECT
	valid_vessels.id
FROM valid_vessels
	 JOIN valid_measurement_units ON valid_vessels.capacity_unit=valid_measurement_units.id
WHERE (valid_vessels.archived_at IS NULL AND valid_measurement_units.archived_at IS NULL)
   AND (
        (valid_vessels.last_indexed_at IS NULL)
        OR valid_vessels.last_indexed_at
            < now() - '24 hours'::INTERVAL
    );

-- name: GetValidVessel :one

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
	AND valid_vessels.id = $1;

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
  AND valid_vessels.id = ANY(sqlc.arg(ids)::text[]);

-- name: SearchForValidVessels :many

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
    AND valid_vessels.name ILIKE '%' || sqlc.arg(query)::text || '%'
	LIMIT 50;

-- name: UpdateValidVessel :exec

UPDATE valid_vessels
SET
    name = sqlc.arg(name),
    plural_name = sqlc.arg(plural_name),
    description = sqlc.arg(description),
    icon_path = sqlc.arg(icon_path),
    usable_for_storage = sqlc.arg(usable_for_storage),
    slug = sqlc.arg(slug),
    display_in_summary_lists = sqlc.arg(display_in_summary_lists),
    include_in_generated_instructions = sqlc.arg(include_in_generated_instructions),
    capacity = sqlc.arg(capacity)::float,
    capacity_unit = sqlc.arg(capacity_unit),
    width_in_millimeters = sqlc.arg(width_in_millimeters)::float,
    length_in_millimeters = sqlc.arg(length_in_millimeters)::float,
    height_in_millimeters = sqlc.arg(height_in_millimeters)::float,
    shape = sqlc.arg(shape),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);

-- name: UpdateValidVesselLastIndexedAt :exec

UPDATE valid_vessels SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
