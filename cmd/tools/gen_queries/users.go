package main

import (
	"fmt"
	"strings"

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
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET %s = NOW() WHERE %s IS NULL AND id = sqlc.arg(id);`,
				usersTableName,
				archivedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "ArchiveUserMemberships",
				Type: ExecRowsType,
			},
			Content: buildUpdateHouseholdMembershipsQuery(belongsToUserColumn, []string{}),
		},
		{
			Annotation: QueryAnnotation{
				Name: "CreateUser",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`INSERT INTO users
(
    %s
) VALUES (
    %s
);`,
				strings.Join(insertColumns, ",\n\t"),
				strings.Join(applyToEach(insertColumns, func(s string) string {
					return fmt.Sprintf("sqlc.arg(%s)", s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetAdminUserByUsername",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.service_role = 'service_admin'
	AND users.username = sqlc.arg(username)
	AND users.two_factor_secret_verified_at IS NOT NULL;`,

				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByEmail",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.email_address = sqlc.arg(email_address);`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByEmailAddressVerificationToken",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.email_address_verification_token = sqlc.arg(email_address_verification_token);`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByID",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.id = sqlc.arg(id);`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserByUsername",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.username = sqlc.arg(username);`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetEmailVerificationTokenByUserID",
				Type: OneType,
			},
			Content: `SELECT
	users.email_address_verification_token
FROM users
WHERE users.archived_at IS NULL
    AND users.email_address_verified_at IS NULL
	AND users.id = sqlc.arg(id);
`,
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUsers",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
    %s,
    %s,
    %s
FROM users
WHERE users.archived_at IS NULL
	%s
LIMIT sqlc.narg(query_limit)
OFFSET sqlc.narg(query_offset);`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
				buildFilterCountSelect(
					usersTableName,
					true,
				),
				buildTotalCountSelect(
					usersTableName,
				),
				buildFilterConditions(
					usersTableName,
					true,
				),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserIDsNeedingIndexing",
				Type: ManyType,
			},
			Content: `SELECT users.id
FROM users
WHERE (users.archived_at IS NULL)
AND users.last_indexed_at IS NULL
OR (
    users.last_indexed_at < NOW() - '24 hours'::INTERVAL
);
`,
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserWithUnverifiedTwoFactor",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.id = sqlc.arg(id)
	AND users.two_factor_secret_verified_at IS NULL;`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "GetUserWithVerifiedTwoFactor",
				Type: OneType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.archived_at IS NULL
	AND users.id = sqlc.arg(id)
	AND users.two_factor_secret_verified_at IS NOT NULL;`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkEmailAddressAsVerified",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	email_address_verified_at = NOW(),
	%s = NOW()
WHERE %s IS NULL
    AND email_address_verified_at IS NULL
	AND id = sqlc.arg(id)
	AND email_address_verification_token = sqlc.arg(email_address_verification_token);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkTwoFactorSecretAsUnverified",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	two_factor_secret_verified_at = NULL,
	two_factor_secret = sqlc.arg(two_factor_secret),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "MarkTwoFactorSecretAsVerified",
				Type: ExecType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	two_factor_secret_verified_at = NOW(),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "SearchUsersByUsername",
				Type: ManyType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`SELECT
	%s
FROM users
WHERE users.username %s
AND users.archived_at IS NULL;`,
				strings.Join(applyToEach(usersColumns, func(s string) string {
					return fmt.Sprintf("%s.%s", usersTableName, s)
				}), ",\n\t"),
				`ILIKE '%' || sqlc.arg(username)::text || '%'`,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserAvatarSrc",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	avatar_src = sqlc.arg(avatar_src),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserDetails",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	first_name = sqlc.arg(first_name),
	last_name = sqlc.arg(last_name),
	birthday = sqlc.arg(birthday),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserEmailAddress",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	email_address = sqlc.arg(email_address),
	email_address_verified_at = NULL,
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserLastIndexedAt",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET last_indexed_at = NOW() WHERE id = sqlc.arg(id) AND %s IS NULL;`,
				usersTableName,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserPassword",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	hashed_password = sqlc.arg(hashed_password),
	password_last_changed_at = NOW(),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserTwoFactorSecret",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	two_factor_secret_verified_at = NULL,
	two_factor_secret = sqlc.arg(two_factor_secret),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
		{
			Annotation: QueryAnnotation{
				Name: "UpdateUserUsername",
				Type: ExecRowsType,
			},
			Content: buildRawQuery((&builq.Builder{}).Addf(`UPDATE %s SET
	username = sqlc.arg(username),
	%s = NOW()
WHERE %s IS NULL
	AND id = sqlc.arg(id);`,
				usersTableName,
				lastUpdatedAtColumn,
				archivedAtColumn,
			)),
		},
	}
}
