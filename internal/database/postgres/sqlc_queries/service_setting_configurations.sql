-- name: ArchiveServiceSettingConfiguration :execrows

UPDATE service_setting_configurations SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateServiceSettingConfiguration :exec

INSERT INTO service_setting_configurations (id,value,notes,service_setting_id,belongs_to_user,belongs_to_household) VALUES ($1,$2,$3,$4,$5,$6);

-- name: CheckServiceSettingConfigurationExistence :one

SELECT EXISTS ( SELECT service_setting_configurations.id FROM service_setting_configurations WHERE service_setting_configurations.archived_at IS NULL AND service_setting_configurations.id = $1 );

-- name: GetServiceSettingConfigurationByID :one

SELECT
	service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
	service_settings.id as service_setting_id,
    service_settings.name as service_setting_name,
    service_settings.type as service_setting_type,
    service_settings.description as service_setting_description,
    service_settings.default_value as service_setting_default_value,
    service_settings.enumeration as service_setting_enumeration,
    service_settings.admins_only as service_setting_admins_only,
    service_settings.created_at as service_setting_created_at,
    service_settings.last_updated_at as service_setting_last_updated_at,
    service_settings.archived_at as service_setting_archived_at,
    service_setting_configurations.belongs_to_user,
    service_setting_configurations.belongs_to_household,
    service_setting_configurations.created_at,
    service_setting_configurations.last_updated_at,
    service_setting_configurations.archived_at
FROM service_setting_configurations
 JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
  AND service_setting_configurations.archived_at IS NULL
  AND service_setting_configurations.id = $1;

-- name: GetServiceSettingConfigurationForHouseholdBySettingName :one

SELECT
    service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
    service_settings.id as service_setting_id,
    service_settings.name as service_setting_name,
    service_settings.type as service_setting_type,
    service_settings.description as service_setting_description,
    service_settings.default_value as service_setting_default_value,
    service_settings.enumeration as service_setting_enumeration,
    service_settings.admins_only as service_setting_admins_only,
    service_settings.created_at as service_setting_created_at,
    service_settings.last_updated_at as service_setting_last_updated_at,
    service_settings.archived_at as service_setting_archived_at,
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
  AND service_setting_configurations.belongs_to_household = $2;

-- name: GetServiceSettingConfigurationForUserBySettingName :one

SELECT
    service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
    service_settings.id as service_setting_id,
    service_settings.name as service_setting_name,
    service_settings.type as service_setting_type,
    service_settings.description as service_setting_description,
    service_settings.default_value as service_setting_default_value,
    service_settings.enumeration as service_setting_enumeration,
    service_settings.admins_only as service_setting_admins_only,
    service_settings.created_at as service_setting_created_at,
    service_settings.last_updated_at as service_setting_last_updated_at,
    service_settings.archived_at as service_setting_archived_at,
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

-- name: GetServiceSettingConfigurationsForHousehold :many

SELECT
    service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
    service_settings.id as service_setting_id,
    service_settings.name as service_setting_name,
    service_settings.type as service_setting_type,
    service_settings.description as service_setting_description,
    service_settings.default_value as service_setting_default_value,
    service_settings.enumeration as service_setting_enumeration,
    service_settings.admins_only as service_setting_admins_only,
    service_settings.created_at as service_setting_created_at,
    service_settings.last_updated_at as service_setting_last_updated_at,
    service_settings.archived_at as service_setting_archived_at,
    service_setting_configurations.belongs_to_user,
    service_setting_configurations.belongs_to_household,
    service_setting_configurations.created_at,
    service_setting_configurations.last_updated_at,
    service_setting_configurations.archived_at
FROM service_setting_configurations
 JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
  AND service_setting_configurations.archived_at IS NULL
  AND service_setting_configurations.belongs_to_household = $1;

-- name: GetServiceSettingConfigurationsForUser :many

SELECT
    service_setting_configurations.id,
    service_setting_configurations.value,
    service_setting_configurations.notes,
    service_settings.id as service_setting_id,
    service_settings.name as service_setting_name,
    service_settings.type as service_setting_type,
    service_settings.description as service_setting_description,
    service_settings.default_value as service_setting_default_value,
    service_settings.enumeration as service_setting_enumeration,
    service_settings.admins_only as service_setting_admins_only,
    service_settings.created_at as service_setting_created_at,
    service_settings.last_updated_at as service_setting_last_updated_at,
    service_settings.archived_at as service_setting_archived_at,
    service_setting_configurations.belongs_to_user,
    service_setting_configurations.belongs_to_household,
    service_setting_configurations.created_at,
    service_setting_configurations.last_updated_at,
    service_setting_configurations.archived_at
FROM service_setting_configurations
 JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
  AND service_setting_configurations.archived_at IS NULL
  AND service_setting_configurations.belongs_to_user = $1;

-- name: UpdateServiceSettingConfiguration :execrows

UPDATE service_setting_configurations
SET
    value = $1,
    notes = $2,
    service_setting_id = $3,
    belongs_to_user = $4,
    belongs_to_household = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
