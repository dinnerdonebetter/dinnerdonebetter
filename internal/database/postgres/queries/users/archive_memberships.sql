UPDATE household_user_memberships SET
	archived_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $1;
