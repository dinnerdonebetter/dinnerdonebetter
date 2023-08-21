-- name: GetValidMeasurementUnits :many

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
    valid_measurement_units.archived_at,
    (
        SELECT
            COUNT(valid_measurement_units.id)
        FROM
            valid_measurement_units
        WHERE
            valid_measurement_units.archived_at IS NULL
            AND valid_measurement_units.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
            AND valid_measurement_units.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
            AND (
                valid_measurement_units.last_updated_at IS NULL
                OR valid_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
            )
            AND (
                valid_measurement_units.last_updated_at IS NULL
                OR valid_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
            )
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_measurement_units.id)
        FROM
            valid_measurement_units
        WHERE
            valid_measurement_units.archived_at IS NULL
    ) as total_count
FROM
    valid_measurement_units
WHERE
    valid_measurement_units.archived_at IS NULL
    AND valid_measurement_units.created_at > COALESCE(sqlc.narg(created_before), (SELECT NOW() - interval '999 years'))
    AND valid_measurement_units.created_at < COALESCE(sqlc.narg(created_after), (SELECT NOW() + interval '999 years'))
    AND (
        valid_measurement_units.last_updated_at IS NULL
        OR valid_measurement_units.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - interval '999 years'))
    )
    AND (
        valid_measurement_units.last_updated_at IS NULL
        OR valid_measurement_units.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + interval '999 years'))
    )
GROUP BY
    valid_measurement_units.id
ORDER BY
    valid_measurement_units.id
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);