-- name: ArchiveServiceSetting :exec

UPDATE service_settings
SET archived_at = NOW()
    WHERE id = $1;

-- name: CreateServiceSetting :exec

INSERT INTO service_settings (id,name,type,description,default_value,admins_only,enumeration) VALUES
($1,$2,$3,$4,$5,$6,$7);

-- name: CheckServiceSettingExistence :one

SELECT EXISTS ( SELECT service_settings.id FROM service_settings WHERE service_settings.archived_at IS NULL AND service_settings.id = $1 );

-- name: GetServiceSettings :many

SELECT
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.admins_only,
    service_settings.enumeration,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at,
    (
        SELECT
            COUNT(service_settings.id)
        FROM
            service_settings
        WHERE
            service_settings.archived_at IS NULL
          AND service_settings.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
          AND service_settings.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
          AND (
                service_settings.last_updated_at IS NULL
                OR service_settings.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
            )
          AND (
                service_settings.last_updated_at IS NULL
                OR service_settings.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
            )
        OFFSET sqlc.narg(query_offset)
    ) AS filtered_count,
    (
        SELECT
            COUNT(service_settings.id)
        FROM
            service_settings
        WHERE
            service_settings.archived_at IS NULL
    ) AS total_count
FROM service_settings
WHERE service_settings.archived_at IS NULL
  AND service_settings.created_at > COALESCE(sqlc.narg(created_after), (SELECT NOW() - interval '999 years'))
  AND service_settings.created_at < COALESCE(sqlc.narg(created_before), (SELECT NOW() + interval '999 years'))
  AND (
    service_settings.last_updated_at IS NULL
   OR service_settings.last_updated_at > COALESCE(sqlc.narg(updated_after), (SELECT NOW() - interval '999 years'))
    )
  AND (
    service_settings.last_updated_at IS NULL
   OR service_settings.last_updated_at < COALESCE(sqlc.narg(updated_before), (SELECT NOW() + interval '999 years'))
    )
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);

-- name: GetServiceSetting :one

SELECT
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.admins_only,
    service_settings.enumeration,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at
FROM service_settings
WHERE service_settings.archived_at IS NULL
	AND service_settings.id = $1;

-- name: SearchForServiceSettings :many

SELECT
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.admins_only,
    service_settings.enumeration,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at
FROM service_settings
WHERE service_settings.archived_at IS NULL
	AND service_settings.name ILIKE '%' || sqlc.arg(name_query)::text || '%'
LIMIT 50;
