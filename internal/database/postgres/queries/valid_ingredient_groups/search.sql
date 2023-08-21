-- name: SearchForValidIngredientGroups :many

SELECT
    valid_ingredient_groups.id,
    valid_ingredient_groups.name,
    valid_ingredient_groups.description,
    valid_ingredient_groups.slug,
    valid_ingredient_groups.created_at,
    valid_ingredient_groups.last_updated_at,
    valid_ingredient_groups.archived_at,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
          AND valid_ingredient_groups.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_groups.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
        OFFSET sqlc.narg(query_offset)
    ) AS filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
    ) AS total_count
FROM valid_ingredient_groups
WHERE
    valid_ingredient_groups.archived_at IS NULL
  AND valid_ingredient_groups.name ILIKE '%' || sqlc.arg(name)::text || '%'
  AND valid_ingredient_groups.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
  AND valid_ingredient_groups.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
  AND (
    valid_ingredient_groups.last_updated_at IS NULL
   OR valid_ingredient_groups.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
  AND (
    valid_ingredient_groups.last_updated_at IS NULL
   OR valid_ingredient_groups.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
OFFSET sqlc.narg(query_offset)
    LIMIT sqlc.narg(query_limit);
