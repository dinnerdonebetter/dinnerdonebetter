package main

import (
	"fmt"
	"github.com/cristalhq/builq"
	"strings"
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
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetServiceSetting",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchForServiceSettings",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
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
	AND service_settings.name %s
LIMIT 50;`,

				`ILIKE '%' || sqlc.arg(username)::text || '%'`,
			)),
		},
	}
}
