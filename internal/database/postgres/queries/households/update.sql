UPDATE households
SET
	name = $1,
	contact_phone = $2,
	time_zone = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $4
	AND id = $5;
