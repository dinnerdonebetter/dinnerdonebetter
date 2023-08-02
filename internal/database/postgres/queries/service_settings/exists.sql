-- name: CheckServiceSettingExistence :one

SELECT EXISTS ( SELECT service_settings.id FROM service_settings WHERE service_settings.archived_at IS NULL AND service_settings.id = $1 );
