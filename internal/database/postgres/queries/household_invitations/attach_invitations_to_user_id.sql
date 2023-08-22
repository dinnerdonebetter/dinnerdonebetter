-- name: AttachHouseholdInvitationsToUserID :exec

UPDATE household_invitations SET
	to_user = sqlc.arg(user_id),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND to_email = LOWER(sqlc.arg(email_address));
