UPDATE households
SET
	name = $1,
	contact_email = $2,
	contact_phone = $3,
	time_zone = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $5
	AND id = $6;
