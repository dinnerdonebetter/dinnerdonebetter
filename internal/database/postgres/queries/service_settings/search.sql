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
