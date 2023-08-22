-- name: ArchiveValidIngredientState :exec

UPDATE valid_ingredient_states SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateValidIngredientState :exec

INSERT INTO valid_ingredient_states (id,"name",description,icon_path,past_tense,slug,attribute_type) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: CheckValidIngredientStateExistence :one

SELECT EXISTS ( SELECT valid_ingredient_states.id FROM valid_ingredient_states WHERE valid_ingredient_states.archived_at IS NULL AND valid_ingredient_states.id = $1 );

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

-- name: GetValidIngredientStatesNeedingIndexing :many

SELECT valid_ingredient_states.id
  FROM valid_ingredient_states
 WHERE (valid_ingredient_states.archived_at IS NULL)
       AND (
			(
				valid_ingredient_states.last_indexed_at IS NULL
			)
			OR valid_ingredient_states.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);

-- name: GetValidIngredientState :one

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
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.id = $1;

-- name: GetValidIngredientStatesWithIDs :many

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
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.id = ANY($1::text[]);

-- name: SearchForValidIngredientStates :many

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
	valid_ingredient_states.archived_at
FROM valid_ingredient_states
WHERE valid_ingredient_states.archived_at IS NULL
	AND valid_ingredient_states.name ILIKE '%' || sqlc.arg(query)::text || '%'
LIMIT 50;

-- name: UpdateValidIngredientState :exec

UPDATE valid_ingredient_states
SET
	name = $1,
	description = $2,
	icon_path = $3,
	slug = $4,
	past_tense = $5,
	attribute_type = $6,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $7;

-- name: UpdateValidIngredientStateLastIndexedAt :exec

UPDATE valid_ingredient_states SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
