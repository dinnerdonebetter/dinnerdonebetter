package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	userNotificationsTableName      = "user_notifications"
	statusColumn                    = "status"
	userNotificationStatusUnread    = "unread"
	userNotificationStatusRead      = "read"
	userNotificationStatusDismissed = "dismissed"
)

var (
	userNotificationsColumns = []string{
		idColumn,
		"content",
		"status",
		belongsToUserColumn,
		createdAtColumn,
		lastUpdatedAtColumn,
	}
)

func buildUserNotificationQueries() []*Query {
	insertColumns := filterForInsert(userNotificationsColumns, "status")
	fullSelectColumns := applyToEach(userNotificationsColumns, func(_ int, s string) string {
		return fullColumnName(userNotificationsTableName, s)
	})

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "MarkUserNotificationAsDismissed",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = %s,
	%s = %s
WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				userNotificationsTableName,
				lastUpdatedAtColumn, currentTimeExpression,
				statusColumn, userNotificationStatusDismissed,
				idColumn, idColumn,
				belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateUserNotification",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				userNotificationsTableName,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(_ int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CheckUserNotificationExistence",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS(
	SELECT %s.%s
	FROM %s
	WHERE %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s)
);`,
				userNotificationsTableName, idColumn,
				userNotificationsTableName,
				userNotificationsTableName, idColumn, idColumn,
				userNotificationsTableName, belongsToUserColumn, belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserNotificationsForUser",
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
					userNotificationsTableName,
					true,
					"user_notifications.belongs_to_user = sqlc.arg(user_id)",
				),
				buildTotalCountSelect(
					userNotificationsTableName,
					"user_notifications.belongs_to_user = sqlc.arg(user_id)",
				),
				userNotificationsTableName,
				strings.TrimPrefix(buildFilterConditions(
					userNotificationsTableName,
					true,
					"user_notifications.belongs_to_user = sqlc.arg(user_id)",
				), "AND "),
				offsetLimitAddendum,
			)),
		},
	}
}
