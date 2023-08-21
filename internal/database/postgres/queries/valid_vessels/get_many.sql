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
