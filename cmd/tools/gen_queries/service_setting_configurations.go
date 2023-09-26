package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	serviceSettingConfigurationsTableName = "service_setting_configurations"

	serviceSettingIDColumn = "service_setting_id"
)

var serviceSettingConfigurationsColumns = []string{
	idColumn,
	"value",
	notesColumn,
	serviceSettingIDColumn,
	belongsToUserColumn,
	belongsToHouseholdColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildServiceSettingConfigurationQueries() []*Query {
	insertColumns := filterForInsert(serviceSettingConfigurationsColumns)

	selectColumnsWithServiceSettingColumns := mergeColumns(
		applyToEach(filterFromSlice(serviceSettingConfigurationsColumns, "service_setting_id"), func(i int, s string) string {
			return fmt.Sprintf("%s.%s", serviceSettingConfigurationsTableName, s)
		}),
		applyToEach(serviceSettingsColumns, func(i int, s string) string {
			return fmt.Sprintf("%s.%s as service_setting_%s", serviceSettingsTableName, s, s)
		}),
		3,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveServiceSettingConfiguration",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
    %s = %s
WHERE %s IS NULL
    AND %s = sqlc.arg(%s);`,
				serviceSettingConfigurationsTableName,
				archivedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateServiceSettingConfiguration",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				serviceSettingConfigurationsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckServiceSettingConfigurationExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.archived_at IS NULL
    AND %s.id = $1
);`,
				serviceSettingConfigurationsTableName,
				serviceSettingConfigurationsTableName,
				serviceSettingConfigurationsTableName,
				serviceSettingConfigurationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSettingConfigurationByID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
    JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
    AND service_setting_configurations.archived_at IS NULL
    AND service_setting_configurations.id = $1;`,
				strings.Join(selectColumnsWithServiceSettingColumns, ",\n\t"),
				serviceSettingConfigurationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSettingConfigurationForHouseholdBySettingName",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
    AND service_setting_configurations.archived_at IS NULL
    AND service_settings.name = $1
    AND service_setting_configurations.belongs_to_household = $2;`,
				strings.Join(selectColumnsWithServiceSettingColumns, ",\n\t"),
				serviceSettingConfigurationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSettingConfigurationForUserBySettingName",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
    AND service_setting_configurations.archived_at IS NULL
    AND service_settings.name = $1
    AND service_setting_configurations.belongs_to_user = $2;`,
				strings.Join(selectColumnsWithServiceSettingColumns, ",\n\t"),
				serviceSettingConfigurationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSettingConfigurationsForHousehold",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
    AND service_setting_configurations.archived_at IS NULL
    AND service_setting_configurations.belongs_to_household = $1;`,
				strings.Join(selectColumnsWithServiceSettingColumns, ",\n\t"),
				serviceSettingConfigurationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSettingConfigurationsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s
FROM %s
    JOIN service_settings ON service_setting_configurations.service_setting_id=service_settings.id
WHERE service_settings.archived_at IS NULL
    AND service_setting_configurations.archived_at IS NULL
    AND service_setting_configurations.belongs_to_user = $1;`,
				strings.Join(selectColumnsWithServiceSettingColumns, ",\n\t"),
				serviceSettingConfigurationsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateServiceSettingConfiguration",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE service_setting_configurations SET
    value = sqlc.arg(value),
    notes = sqlc.arg(notes),
    %s = sqlc.arg(%s),
    %s = sqlc.arg(%s),
    %s = sqlc.arg(%s),
	%s = %s
WHERE %s IS NULL
	AND %s = sqlc.arg(%s);`,
				serviceSettingIDColumn,
				serviceSettingIDColumn,
				belongsToUserColumn,
				belongsToUserColumn,
				belongsToHouseholdColumn,
				belongsToHouseholdColumn,
				lastUpdatedAtColumn,
				currentTimeExpression,
				archivedAtColumn,
				idColumn,
				idColumn,
			)),
		},
	}
}
