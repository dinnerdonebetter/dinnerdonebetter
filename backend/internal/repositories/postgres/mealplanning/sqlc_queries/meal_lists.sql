-- name: ArchiveMealList :execrows
UPDATE meal_lists SET archived_at = NOW() WHERE archived_at IS NULL AND belongs_to_user = sqlc.arg(belongs_to_user) AND id = sqlc.arg(id);

-- name: CreateMealList :exec
INSERT INTO meal_lists (
	id,
	name,
	description,
	belongs_to_user
) VALUES (
	sqlc.arg(id),
	sqlc.arg(name),
	sqlc.arg(description),
	sqlc.arg(belongs_to_user)
);

-- name: GetMealLists :many
SELECT
	meal_lists.id,
	meal_lists.name,
	meal_lists.description,
	meal_lists.created_at,
	meal_lists.last_updated_at,
	meal_lists.archived_at,
	meal_lists.belongs_to_user,
	(
		SELECT COUNT(meal_lists.id)
		FROM meal_lists
		WHERE meal_lists.archived_at IS NULL
			AND
			meal_lists.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
			AND meal_lists.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
			AND (
				meal_lists.last_updated_at IS NULL
				OR meal_lists.last_updated_at > COALESCE(sqlc.narg(updated_before), (SELECT NOW() - '999 years'::INTERVAL))
			)
			AND (
				meal_lists.last_updated_at IS NULL
				OR meal_lists.last_updated_at < COALESCE(sqlc.narg(updated_after), (SELECT NOW() + '999 years'::INTERVAL))
			)
			AND (NOT COALESCE(sqlc.narg(include_archived), false)::boolean OR meal_lists.archived_at = NULL)
	) AS filtered_count,
	(
		SELECT COUNT(meal_lists.id)
		FROM meal_lists
		WHERE meal_lists.archived_at IS NULL
	) AS total_count
FROM meal_lists
	WHERE meal_lists.archived_at IS NULL
	AND meal_lists.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - '999 years'::INTERVAL))
	AND meal_lists.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + '999 years'::INTERVAL))
	AND (
		meal_lists.last_updated_at IS NULL
		OR meal_lists.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - '999 years'::INTERVAL))
	)
	AND (
		meal_lists.last_updated_at IS NULL
		OR meal_lists.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + '999 years'::INTERVAL))
	)
	AND meal_lists.id > COALESCE(sqlc.narg(cursor), '')
ORDER BY meal_lists.id ASC
LIMIT COALESCE(sqlc.narg(result_limit), 50);

-- name: UpdateMealList :execrows
UPDATE meal_lists SET
	name = sqlc.arg(name),
	description = sqlc.arg(description),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = sqlc.arg(belongs_to_user)
	AND id = sqlc.arg(id);
