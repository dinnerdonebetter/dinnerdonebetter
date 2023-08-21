-- name: GetValidPreparations :many

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
	valid_preparations.archived_at,
    (
        SELECT
            COUNT(valid_preparations.id)
        FROM
            valid_preparations
        WHERE
            valid_preparations.archived_at IS NULL
          AND valid_preparations.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_preparations.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (
                valid_preparations.last_updated_at IS NULL
                OR valid_preparations.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
          AND (
                valid_preparations.last_updated_at IS NULL
                OR valid_preparations.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_preparations.id)
        FROM
            valid_preparations
        WHERE
            valid_preparations.archived_at IS NULL
    ) as total_count
FROM valid_preparations
WHERE
    valid_preparations.archived_at IS NULL
  AND valid_preparations.created_at > (COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years')))
  AND valid_preparations.created_at < (COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years')))
  AND (
        valid_preparations.last_updated_at IS NULL
        OR valid_preparations.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
  AND (
        valid_preparations.last_updated_at IS NULL
        OR valid_preparations.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
GROUP BY
    valid_preparations.id
ORDER BY
    valid_preparations.id
OFFSET
    sqlc.narg(query_offset)
    LIMIT
    sqlc.narg(query_limit);
