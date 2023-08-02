-- name: CheckServiceSettingConfigurationExistence :one

SELECT EXISTS ( SELECT service_setting_configurations.id FROM service_setting_configurations WHERE service_setting_configurations.archived_at IS NULL AND service_setting_configurations.id = $1 );
