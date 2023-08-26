-- name: ArchiveValidMeasurementUnit :execrows

UPDATE valid_measurement_units SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateValidMeasurementUnit :exec

INSERT INTO valid_measurement_units
(id,"name",description,volumetric,icon_path,universal,metric,imperial,plural_name,slug)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10);

-- name: CheckValidMeasurementUnitExistence :one

SELECT EXISTS ( SELECT valid_measurement_units.id FROM valid_measurement_units WHERE valid_measurement_units.archived_at IS NULL AND valid_measurement_units.id = $1 );

-- name: GetValidMeasurementUnitByID :one

SELECT valid_measurement_units.id,
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
       valid_measurement_units.archived_at
  FROM valid_measurement_units
 WHERE valid_measurement_units.archived_at IS NULL;

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

-- name: GetValidMeasurementUnitsNeedingIndexing :many

SELECT valid_measurement_units.id
  FROM valid_measurement_units
 WHERE (valid_measurement_units.archived_at IS NULL)
       AND (
			(
				valid_measurement_units.last_indexed_at IS NULL
			)
			OR valid_measurement_units.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);

-- name: GetValidMeasurementUnit :one

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
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
	AND valid_measurement_units.id = $1;

-- name: GetRandomValidMeasurementUnit :one

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
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
	ORDER BY random() LIMIT 1;

-- name: GetValidMeasurementUnitsWithIDs :many

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
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE valid_measurement_units.archived_at IS NULL
	AND valid_measurement_units.id = ANY(sqlc.arg(ids)::text[]);

-- name: SearchForValidMeasurementUnits :many

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
	valid_measurement_units.archived_at
FROM valid_measurement_units
WHERE (valid_measurement_units.name ILIKE '%' || sqlc.arg(query)::text || '%' OR valid_measurement_units.universal is TRUE)
AND valid_measurement_units.archived_at IS NULL
LIMIT 50;

-- name: SearchValidMeasurementUnitsByIngredientID :many

SELECT
	DISTINCT(valid_measurement_units.id),
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
	    AND (
	        valid_ingredient_measurement_units.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
	        OR valid_measurement_units.universal = true
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
FROM valid_measurement_units
	JOIN valid_ingredient_measurement_units ON valid_ingredient_measurement_units.valid_measurement_unit_id = valid_measurement_units.id
	JOIN valid_ingredients ON valid_ingredient_measurement_units.valid_ingredient_id = valid_ingredients.id
WHERE
	(
	    valid_ingredient_measurement_units.valid_ingredient_id = sqlc.arg(valid_ingredient_id)
	    OR valid_measurement_units.universal = true
	)
    AND valid_measurement_units.archived_at IS NULL
    AND valid_ingredients.archived_at IS NULL
    AND valid_ingredient_measurement_units.archived_at IS NULL
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
	LIMIT sqlc.narg(query_limit)
	OFFSET sqlc.narg(query_offset);

-- name: UpdateValidMeasurementUnit :execrows

UPDATE valid_measurement_units SET
	name = $1,
	description = $2,
	volumetric = $3,
	icon_path = $4,
	universal = $5,
	metric = $6,
	imperial = $7,
	slug = $8,
	plural_name = $9,
	last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $10;

-- name: UpdateValidMeasurementUnitLastIndexedAt :execrows

UPDATE valid_measurement_units SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
