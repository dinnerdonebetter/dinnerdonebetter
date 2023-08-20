-- name: GetValidIngredientStates :many

SELECT
	valid_ingredient_states.id,
	valid_ingredient_states.name,
	valid_ingredient_states.description,
	valid_ingredient_states.icon_path,
	valid_ingredient_states.slug,
	valid_ingredient_states.past_tense,
	valid_ingredient_states.attribute_type,
	valid_ingredient_states.created_at,
	valid_ingredient_states.last_updated_at,
	valid_ingredient_states.archived_at,
    (
        SELECT COUNT(valid_ingredient_states.id)
        FROM valid_ingredient_states
        WHERE valid_ingredient_states.archived_at IS NULL
        AND valid_ingredient_states.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
        AND valid_ingredient_states.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
        AND (valid_ingredient_states.last_updated_at IS NULL OR valid_ingredient_states.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
        AND (valid_ingredient_states.last_updated_at IS NULL OR valid_ingredient_states.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_states.id)
        FROM
            valid_ingredient_states
        WHERE
            valid_ingredient_states.archived_at IS NULL
    ) as total_count
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
    AND valid_ingredient_states.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
    AND valid_ingredient_states.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
    AND (valid_ingredient_states.last_updated_at IS NULL OR valid_ingredient_states.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years')))
    AND (valid_ingredient_states.last_updated_at IS NULL OR valid_ingredient_states.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years')))
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);
