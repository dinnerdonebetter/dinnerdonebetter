package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const serviceSettingsTableName = "service_settings"

var serviceSettingsColumns = []string{
	idColumn,
	"name",
	"type",
	"description",
	"default_value",
	"enumeration",
	"admins_only",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildServiceSettingQueries() []*Query {
	insertColumns := filterForInsert(serviceSettingsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveServiceSetting",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s
SET archived_at = NOW()
    WHERE id = sqlc.arg(id);`,
				serviceSettingsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateServiceSetting",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
    %s
) VALUES (
    %s
);`,
				serviceSettingsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckServiceSettingExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
    AND %s.id = sqlc.arg(id)
);`,
				serviceSettingsTableName,
				serviceSettingsTableName,
				serviceSettingsTableName,
				archivedAtColumn,
				serviceSettingsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSettings",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
    %s,
    %s
FROM service_settings
WHERE service_settings.archived_at IS NULL
    %s
OFFSET sqlc.narg(query_offset)
LIMIT sqlc.narg(query_limit);`,
				strings.Join(applyToEach(serviceSettingsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", serviceSettingsTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					serviceSettingsTableName,
					true,
				),
				buildTotalCountSelect(
					serviceSettingsTableName,
				),
				buildFilterConditions(
					serviceSettingsTableName,
					true,
				),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSetting",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE service_settings.archived_at IS NULL
	AND service_settings.id = sqlc.arg(id);`,
				strings.Join(applyToEach(serviceSettingsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", serviceSettingsTableName, s)
				}), ",\n\t"),
				serviceSettingsTableName,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForServiceSettings",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM %s
WHERE service_settings.archived_at IS NULL
	AND service_settings.name %s
LIMIT 50;`,
				strings.Join(applyToEach(serviceSettingsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", serviceSettingsTableName, s)
				}), ",\n\t"),
				serviceSettingsTableName,
				`ILIKE '%' || sqlc.arg(name_query)::text || '%'`,
			)),
		},
	}
}
