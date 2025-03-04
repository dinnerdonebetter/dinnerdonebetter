// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: service_setting_configurations.sql

package generated

import (
	"context"
	"database/sql"
	"time"
)

const archiveServiceSettingConfiguration = `-- name: ArchiveServiceSettingConfiguration :execrows
UPDATE service_setting_configurations SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND id = $1
`

func (q *Queries) ArchiveServiceSettingConfiguration(ctx context.Context, db DBTX, id string) (int64, error) {
	result, err := db.ExecContext(ctx, archiveServiceSettingConfiguration, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

const checkServiceSettingConfigurationExistence = `-- name: CheckServiceSettingConfigurationExistence :one
SELECT EXISTS (
	SELECT service_setting_configurations.id
	FROM service_setting_configurations
	WHERE service_setting_configurations.archived_at IS NULL
		AND service_setting_configurations.id = $1
)
`

func (q *Queries) CheckServiceSettingConfigurationExistence(ctx context.Context, db DBTX, id string) (bool, error) {
	row := db.QueryRowContext(ctx, checkServiceSettingConfigurationExistence, id)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createServiceSettingConfiguration = `-- name: CreateServiceSettingConfiguration :exec
INSERT INTO service_setting_configurations (
	id,
	value,
	notes,
	service_setting_id,
	belongs_to_user,
	belongs_to_household
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6
)
`

type CreateServiceSettingConfigurationParams struct {
	ID                 string
	Value              string
	Notes              string
	ServiceSettingID   string
	BelongsToUser      string
	BelongsToHousehold string
}

func (q *Queries) CreateServiceSettingConfiguration(ctx context.Context, db DBTX, arg *CreateServiceSettingConfigurationParams) error {
	_, err := db.ExecContext(ctx, createServiceSettingConfiguration,
		arg.ID,
		arg.Value,
		arg.Notes,
		arg.ServiceSettingID,
		arg.BelongsToUser,
		arg.BelongsToHousehold,
	)
	return err
}

const getServiceSettingConfigurationByID = `-- name: GetServiceSettingConfigurationByID :one
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
	AND service_setting_configurations.id = $1
`

type GetServiceSettingConfigurationByIDRow struct {
	ServiceSettingCreatedAt     time.Time
	CreatedAt                   time.Time
	ArchivedAt                  sql.NullTime
	LastUpdatedAt               sql.NullTime
	ServiceSettingArchivedAt    sql.NullTime
	ServiceSettingLastUpdatedAt sql.NullTime
	ServiceSettingName          string
	ServiceSettingEnumeration   string
	ServiceSettingDescription   string
	ServiceSettingType          SettingType
	ID                          string
	BelongsToUser               string
	BelongsToHousehold          string
	ServiceSettingID            string
	Notes                       string
	Value                       string
	ServiceSettingDefaultValue  sql.NullString
	ServiceSettingAdminsOnly    bool
}

func (q *Queries) GetServiceSettingConfigurationByID(ctx context.Context, db DBTX, id string) (*GetServiceSettingConfigurationByIDRow, error) {
	row := db.QueryRowContext(ctx, getServiceSettingConfigurationByID, id)
	var i GetServiceSettingConfigurationByIDRow
	err := row.Scan(
		&i.ID,
		&i.Value,
		&i.Notes,
		&i.ServiceSettingID,
		&i.ServiceSettingName,
		&i.ServiceSettingType,
		&i.ServiceSettingDescription,
		&i.ServiceSettingDefaultValue,
		&i.ServiceSettingEnumeration,
		&i.ServiceSettingAdminsOnly,
		&i.ServiceSettingCreatedAt,
		&i.ServiceSettingLastUpdatedAt,
		&i.ServiceSettingArchivedAt,
		&i.BelongsToUser,
		&i.BelongsToHousehold,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getServiceSettingConfigurationForHouseholdBySettingName = `-- name: GetServiceSettingConfigurationForHouseholdBySettingName :one
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
	AND service_setting_configurations.belongs_to_household = $2
`

type GetServiceSettingConfigurationForHouseholdBySettingNameParams struct {
	Name               string
	BelongsToHousehold string
}

type GetServiceSettingConfigurationForHouseholdBySettingNameRow struct {
	ServiceSettingCreatedAt     time.Time
	CreatedAt                   time.Time
	ArchivedAt                  sql.NullTime
	LastUpdatedAt               sql.NullTime
	ServiceSettingArchivedAt    sql.NullTime
	ServiceSettingLastUpdatedAt sql.NullTime
	ServiceSettingName          string
	ServiceSettingEnumeration   string
	ServiceSettingDescription   string
	ServiceSettingType          SettingType
	ID                          string
	BelongsToUser               string
	BelongsToHousehold          string
	ServiceSettingID            string
	Notes                       string
	Value                       string
	ServiceSettingDefaultValue  sql.NullString
	ServiceSettingAdminsOnly    bool
}

func (q *Queries) GetServiceSettingConfigurationForHouseholdBySettingName(ctx context.Context, db DBTX, arg *GetServiceSettingConfigurationForHouseholdBySettingNameParams) (*GetServiceSettingConfigurationForHouseholdBySettingNameRow, error) {
	row := db.QueryRowContext(ctx, getServiceSettingConfigurationForHouseholdBySettingName, arg.Name, arg.BelongsToHousehold)
	var i GetServiceSettingConfigurationForHouseholdBySettingNameRow
	err := row.Scan(
		&i.ID,
		&i.Value,
		&i.Notes,
		&i.ServiceSettingID,
		&i.ServiceSettingName,
		&i.ServiceSettingType,
		&i.ServiceSettingDescription,
		&i.ServiceSettingDefaultValue,
		&i.ServiceSettingEnumeration,
		&i.ServiceSettingAdminsOnly,
		&i.ServiceSettingCreatedAt,
		&i.ServiceSettingLastUpdatedAt,
		&i.ServiceSettingArchivedAt,
		&i.BelongsToUser,
		&i.BelongsToHousehold,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getServiceSettingConfigurationForUserBySettingName = `-- name: GetServiceSettingConfigurationForUserBySettingName :one
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
	AND service_setting_configurations.belongs_to_user = $2
`

type GetServiceSettingConfigurationForUserBySettingNameParams struct {
	Name          string
	BelongsToUser string
}

type GetServiceSettingConfigurationForUserBySettingNameRow struct {
	ServiceSettingCreatedAt     time.Time
	CreatedAt                   time.Time
	ArchivedAt                  sql.NullTime
	LastUpdatedAt               sql.NullTime
	ServiceSettingArchivedAt    sql.NullTime
	ServiceSettingLastUpdatedAt sql.NullTime
	ServiceSettingName          string
	ServiceSettingEnumeration   string
	ServiceSettingDescription   string
	ServiceSettingType          SettingType
	ID                          string
	BelongsToUser               string
	BelongsToHousehold          string
	ServiceSettingID            string
	Notes                       string
	Value                       string
	ServiceSettingDefaultValue  sql.NullString
	ServiceSettingAdminsOnly    bool
}

func (q *Queries) GetServiceSettingConfigurationForUserBySettingName(ctx context.Context, db DBTX, arg *GetServiceSettingConfigurationForUserBySettingNameParams) (*GetServiceSettingConfigurationForUserBySettingNameRow, error) {
	row := db.QueryRowContext(ctx, getServiceSettingConfigurationForUserBySettingName, arg.Name, arg.BelongsToUser)
	var i GetServiceSettingConfigurationForUserBySettingNameRow
	err := row.Scan(
		&i.ID,
		&i.Value,
		&i.Notes,
		&i.ServiceSettingID,
		&i.ServiceSettingName,
		&i.ServiceSettingType,
		&i.ServiceSettingDescription,
		&i.ServiceSettingDefaultValue,
		&i.ServiceSettingEnumeration,
		&i.ServiceSettingAdminsOnly,
		&i.ServiceSettingCreatedAt,
		&i.ServiceSettingLastUpdatedAt,
		&i.ServiceSettingArchivedAt,
		&i.BelongsToUser,
		&i.BelongsToHousehold,
		&i.CreatedAt,
		&i.LastUpdatedAt,
		&i.ArchivedAt,
	)
	return &i, err
}

const getServiceSettingConfigurationsForHousehold = `-- name: GetServiceSettingConfigurationsForHousehold :many
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
	AND service_setting_configurations.belongs_to_household = $1
`

type GetServiceSettingConfigurationsForHouseholdRow struct {
	ServiceSettingCreatedAt     time.Time
	CreatedAt                   time.Time
	ArchivedAt                  sql.NullTime
	LastUpdatedAt               sql.NullTime
	ServiceSettingArchivedAt    sql.NullTime
	ServiceSettingLastUpdatedAt sql.NullTime
	ServiceSettingName          string
	ServiceSettingEnumeration   string
	ServiceSettingDescription   string
	ServiceSettingType          SettingType
	ID                          string
	BelongsToUser               string
	BelongsToHousehold          string
	ServiceSettingID            string
	Notes                       string
	Value                       string
	ServiceSettingDefaultValue  sql.NullString
	ServiceSettingAdminsOnly    bool
}

func (q *Queries) GetServiceSettingConfigurationsForHousehold(ctx context.Context, db DBTX, belongsToHousehold string) ([]*GetServiceSettingConfigurationsForHouseholdRow, error) {
	rows, err := db.QueryContext(ctx, getServiceSettingConfigurationsForHousehold, belongsToHousehold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetServiceSettingConfigurationsForHouseholdRow{}
	for rows.Next() {
		var i GetServiceSettingConfigurationsForHouseholdRow
		if err := rows.Scan(
			&i.ID,
			&i.Value,
			&i.Notes,
			&i.ServiceSettingID,
			&i.ServiceSettingName,
			&i.ServiceSettingType,
			&i.ServiceSettingDescription,
			&i.ServiceSettingDefaultValue,
			&i.ServiceSettingEnumeration,
			&i.ServiceSettingAdminsOnly,
			&i.ServiceSettingCreatedAt,
			&i.ServiceSettingLastUpdatedAt,
			&i.ServiceSettingArchivedAt,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getServiceSettingConfigurationsForUser = `-- name: GetServiceSettingConfigurationsForUser :many
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
	AND service_setting_configurations.belongs_to_user = $1
`

type GetServiceSettingConfigurationsForUserRow struct {
	ServiceSettingCreatedAt     time.Time
	CreatedAt                   time.Time
	ArchivedAt                  sql.NullTime
	LastUpdatedAt               sql.NullTime
	ServiceSettingArchivedAt    sql.NullTime
	ServiceSettingLastUpdatedAt sql.NullTime
	ServiceSettingName          string
	ServiceSettingEnumeration   string
	ServiceSettingDescription   string
	ServiceSettingType          SettingType
	ID                          string
	BelongsToUser               string
	BelongsToHousehold          string
	ServiceSettingID            string
	Notes                       string
	Value                       string
	ServiceSettingDefaultValue  sql.NullString
	ServiceSettingAdminsOnly    bool
}

func (q *Queries) GetServiceSettingConfigurationsForUser(ctx context.Context, db DBTX, belongsToUser string) ([]*GetServiceSettingConfigurationsForUserRow, error) {
	rows, err := db.QueryContext(ctx, getServiceSettingConfigurationsForUser, belongsToUser)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*GetServiceSettingConfigurationsForUserRow{}
	for rows.Next() {
		var i GetServiceSettingConfigurationsForUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Value,
			&i.Notes,
			&i.ServiceSettingID,
			&i.ServiceSettingName,
			&i.ServiceSettingType,
			&i.ServiceSettingDescription,
			&i.ServiceSettingDefaultValue,
			&i.ServiceSettingEnumeration,
			&i.ServiceSettingAdminsOnly,
			&i.ServiceSettingCreatedAt,
			&i.ServiceSettingLastUpdatedAt,
			&i.ServiceSettingArchivedAt,
			&i.BelongsToUser,
			&i.BelongsToHousehold,
			&i.CreatedAt,
			&i.LastUpdatedAt,
			&i.ArchivedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateServiceSettingConfiguration = `-- name: UpdateServiceSettingConfiguration :execrows
UPDATE service_setting_configurations SET
	value = $1,
	notes = $2,
	service_setting_id = $3,
	belongs_to_user = $4,
	belongs_to_household = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6
`

type UpdateServiceSettingConfigurationParams struct {
	Value              string
	Notes              string
	ServiceSettingID   string
	BelongsToUser      string
	BelongsToHousehold string
	ID                 string
}

func (q *Queries) UpdateServiceSettingConfiguration(ctx context.Context, db DBTX, arg *UpdateServiceSettingConfigurationParams) (int64, error) {
	result, err := db.ExecContext(ctx, updateServiceSettingConfiguration,
		arg.Value,
		arg.Notes,
		arg.ServiceSettingID,
		arg.BelongsToUser,
		arg.BelongsToHousehold,
		arg.ID,
	)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
