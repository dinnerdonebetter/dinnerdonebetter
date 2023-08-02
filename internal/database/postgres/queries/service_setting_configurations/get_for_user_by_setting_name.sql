-- name: GetServiceSettingConfigurationsForUserBySettingName :many

SELECT
	service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
	service_settings.id,
    service_settings.name,
    service_settings.type,
    service_settings.description,
    service_settings.default_value,
    service_settings.enumeration,
    service_settings.admins_only,
    service_settings.created_at,
    service_settings.last_updated_at,
    service_settings.archived_at,
    service_setting_configurations.belongs_to_user,
    service_setting_configurations.belongs_to_household,
    service_setting_configurations.created_at,
    service_setting_configurations.last_updated_at,
    service_setting_configurations.archived_at
FROM service_setting_configurations
 JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
  AND service_setting_configurations.archived_at IS NULL
  AND service_settings.name = $1
  AND service_setting_configurations.belongs_to_user = $2;