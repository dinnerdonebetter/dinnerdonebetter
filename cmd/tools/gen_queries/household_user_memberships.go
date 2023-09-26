package main

import (
	"fmt"
	"strings"

	"github.com/cristalhq/builq"
)

const (
	householdUserMembershipsTableName = "household_user_memberships"
)

var householdUserMembershipsColumns = []string{
	idColumn,
	belongsToHouseholdColumn,
	belongsToUserColumn,
	"default_household",
	"household_role",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdUserMembershipsQueries() []*Query {
	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AddUserToHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdUserMembershipsTableName,
				strings.Join(filterForInsert(householdUserMembershipsColumns, "default_household"), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(householdUserMembershipsColumns, "default_household"), func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHouseholdUserMembershipForNewUser",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO %s (
	%s
) VALUES (
	%s
);`,
				householdUserMembershipsTableName,
				strings.Join(filterForInsert(householdUserMembershipsColumns), ",\n\t"),
				strings.Join(applyToEach(filterForInsert(householdUserMembershipsColumns), func(i int, s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetDefaultHouseholdIDForUser",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT households.id
FROM households
	JOIN household_user_memberships ON household_user_memberships.belongs_to_household = households.id
WHERE household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user)
	AND household_user_memberships.default_household = TRUE;`)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdUserMembershipsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM household_user_memberships
	JOIN households ON households.id = household_user_memberships.belongs_to_household
WHERE household_user_memberships.archived_at IS NULL
	AND household_user_memberships.belongs_to_user = sqlc.arg(belongs_to_user);`,
				strings.Join(applyToEach(householdUserMembershipsColumns, func(i int, s string) string {
					return fmt.Sprintf("%s.%s", householdUserMembershipsTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkHouseholdUserMembershipAsUserDefault",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s
SET default_household = (belongs_to_user = sqlc.arg(user_id) AND belongs_to_household = sqlc.arg(household_id))
WHERE %s IS NULL
	AND %s = sqlc.arg(user_id);`,
				householdUserMembershipsTableName,
				archivedAtColumn,
				belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ModifyHouseholdUserPermissions",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	household_role = sqlc.arg(household_role)
WHERE %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				belongsToHouseholdColumn,
				belongsToHouseholdColumn,
				belongsToUserColumn,
				belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "RemoveUserFromHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s
SET %s = %s,
	default_household = 'false'
WHERE %s.%s IS NULL
	AND %s.%s = sqlc.arg(%s)
	AND %s.%s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				archivedAtColumn,
				currentTimeExpression,
				householdUserMembershipsTableName,
				archivedAtColumn,
				householdUserMembershipsTableName,
				belongsToHouseholdColumn,
				belongsToHouseholdColumn,
				householdUserMembershipsTableName,
				belongsToUserColumn,
				belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "TransferHouseholdMembership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(%s)
WHERE %s IS NULL
	AND %s = sqlc.arg(%s)
	AND %s = sqlc.arg(%s);`,
				householdUserMembershipsTableName,
				belongsToUserColumn,
				belongsToUserColumn,
				archivedAtColumn,
				belongsToHouseholdColumn,
				belongsToHouseholdColumn,
				belongsToUserColumn,
				belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "TransferHouseholdOwnership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	%s = sqlc.arg(new_owner)
WHERE %s IS NULL
    AND %s = sqlc.arg(old_owner)
    AND id = sqlc.arg(household_id);`,
				householdsTableName,
				belongsToUserColumn,
				archivedAtColumn,
				belongsToUserColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UserIsHouseholdMember",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT EXISTS (
    SELECT %s.id
    FROM %s
    WHERE %s.%s IS NULL
        AND %s.%s = sqlc.arg(%s)
        AND %s.%s = sqlc.arg(%s)
);`,
				householdUserMembershipsTableName,
				householdUserMembershipsTableName,
				householdUserMembershipsTableName,
				archivedAtColumn,
				householdUserMembershipsTableName,
				belongsToHouseholdColumn,
				belongsToHouseholdColumn,
				householdUserMembershipsTableName,
				belongsToUserColumn,
				belongsToUserColumn,
			)),
		},
	}
}
