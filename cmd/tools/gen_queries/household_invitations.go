package main

const householdInvitationsTableName = "household_invitations"

var householdInvitationsColumns = []string{
	"id",
	"destination_household",
	"to_email",
	"to_user",
	"from_user",
	"status",
	"note",
	"status_note",
	"token",
	"expires_at",
	"to_name",
	createdAtColumn,
	lastUpdatedAtColumn,
	archivedAtColumn,
}
