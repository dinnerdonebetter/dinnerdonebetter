package main

import (
	"fmt"

	"github.com/cristalhq/builq"
)

const (
	usersTableName = "users"

	lastAcceptedTOSColumn           = "last_accepted_terms_of_service"
	lastAcceptedPrivacyPolicyColumn = "last_accepted_privacy_policy"
)

var usersColumns = []string{
	idColumn,
	"username",
	"avatar_src",
	"email_address",
	"hashed_password",
	"password_last_changed_at",
	"requires_password_change",
	"two_factor_secret",
	"two_factor_secret_verified_at",
	"service_role",
	"user_account_status",
	"user_account_status_explanation",
	"birthday",
	"email_address_verification_token",
	"email_address_verified_at",
	"first_name",
	"last_name",
	lastAcceptedTOSColumn,
	lastAcceptedPrivacyPolicyColumn,
	lastIndexedAtColumn,
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}

func buildUpdateHouseholdMembershipsQuery(ownershipColumn string, nowColumns []string) string {
	var updateQueryBuilder builq.Builder

	addendum := ""
	for _, col := range nowColumns {
		addendum = fmt.Sprintf(",\n\t%s = NOW()", col)
	}

	builder := updateQueryBuilder.Addf(
		`UPDATE %s 
SET %s = NOW()%s
WHERE %s IS NULL AND %s = sqlc.arg(id);`,
		householdUserMembershipsTableName,
		archivedAtColumn,
		addendum,
		archivedAtColumn,
		ownershipColumn,
	)

	return buildRawQuery(builder)
}

func buildUserUpdateQuery(columnName string, nowColumns []string) string {
	var updateQueryBuilder builq.Builder

	addendum := ""
	for _, col := range nowColumns {
		addendum = fmt.Sprintf(",\n\t%s = NOW()", col)
	}

	builder := updateQueryBuilder.Addf(
		`UPDATE users SET
	%s = NOW()%s
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);`,
		columnName,
		addendum,
	)

	return buildRawQuery(builder)
}

func buildUsersQueries() []*Query {
	insertColumns := filterForInsert(usersColumns,
		"password_last_changed_at",
		"email_address_verified_at",
		lastAcceptedTOSColumn,
		lastAcceptedPrivacyPolicyColumn,
		lastIndexedAtColumn,
	)

	return []*Query{
		{
			Annotation: QueryAnnotation{
				Name: "AcceptPrivacyPolicyForUser",
				Type: ExecType,
			},
			Content: buildUserUpdateQuery(lastAcceptedPrivacyPolicyColumn, nil),
		},
		{
			Annotation: QueryAnnotation{
				Name: "AcceptTermsOfServiceForUser",
				Type: ExecType,
			},
			Content: buildUserUpdateQuery(lastAcceptedTOSColumn, nil),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveUser",
				Type: ExecRowsType,
			},
			Content: formatQuery(
				buildArchiveQuery(usersTableName, ""),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveUserMemberships",
				Type: ExecRowsType,
			},
			Content: formatQuery(
				buildUpdateHouseholdMembershipsQuery(belongsToUserColumn, []string{}),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateUser",
				Type: ExecType,
			},
			Content: formatQuery(
				buildCreateQuery(usersTableName, insertColumns),
			),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAdminUserByUsername",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByEmail",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByEmailAddressVerificationToken",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByID",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByUsername",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetEmailVerificationTokenByUserID",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUsers",
				Type: ManyType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserIDsNeedingIndexing",
				Type: ManyType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserWithUnverifiedTwoFactor",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserWithVerifiedTwoFactor",
				Type: OneType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkEmailAddressAsVerified",
				Type: ExecType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkTwoFactorSecretAsUnverified",
				Type: ExecType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkTwoFactorSecretAsVerified",
				Type: ExecType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchUsersByUsername",
				Type: ManyType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUser",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserAvatarSrc",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserDetails",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserEmailAddress",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserPassword",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserTwoFactorSecret",
				Type: ExecRowsType,
			},
			Content: "",
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserUsername",
				Type: ExecRowsType,
			},
			Content: "",
		},
	}
}
