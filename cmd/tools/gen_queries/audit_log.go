package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	auditLogsTableName = "audit_log"
	resourceTypeColumn = "resource_type"
	eventTypeColumn    = "event_type"
)

var (
	auditLogsColumns = []string{
		idColumn,
		resourceTypeColumn,
		"relevant_id",
		eventTypeColumn,
		"changes",
		belongsToUserColumn,
		belongsToHouseholdColumn,
		createdAtColumn,
	}
)

func buildAuditLogQueries() []*Query {
	insertColumns := filterForInsert(auditLogsColumns)
	fullSelectColumns := applyToEach(auditLogsColumns, func(_ int, s string) string {
		return fullColumnName(auditLogsTableName, s)
	})

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "CreateAuditLog",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				auditLogsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAuditLog",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT 
	%s
FROM %s
WHERE %s.%s = sqlc.arg(%s);`,
				strings.Join(applyToEach(fullSelectColumns, func(_ int, s string) string {
					parts := strings.Split(s, ".")
					return fmt.Sprintf("%s as %s_%s",
						s, strings.TrimSuffix(parts[0], "s"), parts[1],
					)
				}), ",\n\t"),
				auditLogsTableName,
				auditLogsTableName, idColumn, idColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAuditLogsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					auditLogsTableName,
					false,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToUserColumn, belongsToUserColumn),
				),
				buildTotalCountSelect(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToUserColumn, belongsToUserColumn),
				),
				auditLogsTableName,
				strings.TrimPrefix(buildFilterConditions(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToUserColumn, belongsToUserColumn),
				), "AND "),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAuditLogsForUserAndResourceType",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					auditLogsTableName,
					false,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToUserColumn, belongsToUserColumn),
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, resourceTypeColumn, resourceTypeColumn),
				),
				buildTotalCountSelect(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToUserColumn, belongsToUserColumn),
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, resourceTypeColumn, resourceTypeColumn),
				),
				auditLogsTableName,
				strings.TrimPrefix(buildFilterConditions(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToUserColumn, belongsToUserColumn),
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, resourceTypeColumn, resourceTypeColumn),
				), "AND "),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAuditLogsForHousehold",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					auditLogsTableName,
					false,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
				),
				buildTotalCountSelect(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
				),
				auditLogsTableName,
				strings.TrimPrefix(buildFilterConditions(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
				), "AND "),
				offsetLimitAddendum,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAuditLogsForHouseholdAndResourceType",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s,
	%s,
	%s
FROM %s
WHERE %s
%s;`,
				strings.Join(fullSelectColumns, ",\n\t"),
				buildFilterCountSelect(
					auditLogsTableName,
					false,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, resourceTypeColumn, resourceTypeColumn),
				),
				buildTotalCountSelect(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, resourceTypeColumn, resourceTypeColumn),
				),
				auditLogsTableName,
				strings.TrimPrefix(buildFilterConditions(
					auditLogsTableName,
					false,
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, belongsToHouseholdColumn, belongsToHouseholdColumn),
					fmt.Sprintf("%s.%s = sqlc.arg(%s)", auditLogsTableName, resourceTypeColumn, resourceTypeColumn),
				), "AND "),
				offsetLimitAddendum,
			)),
		},
	}
}
