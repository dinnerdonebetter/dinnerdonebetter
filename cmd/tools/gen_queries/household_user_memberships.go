package main

import (
	"github.com/cristalhq/builq"
)

const householdUserMembershipsTableName = "household_user_memberships"

var householdUserMembershipsColumns = []string{
	idColumn,
	belongsToHouseholdColumn,
	"belongs_to_user",
	"default_household",
	"household_role",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildHouseholdUserMembershipsQueries() []*Query {
	// insertColumns := filterForInsert(householdUserMembershipsColumns)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AddUserToHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateHouseholdUserMembershipForNewUser",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetDefaultHouseholdIDForUser",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetHouseholdUserMembershipsForUser",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkHouseholdUserMembershipAsUserDefault",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ModifyHouseholdUserPermissions",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "RemoveUserFromHousehold",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "TransferHouseholdMembership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "TransferHouseholdOwnership",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UserIsHouseholdMember",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(``)),
		},
	}
}
