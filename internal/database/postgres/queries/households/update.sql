UPDATE households
SET
	name = $1,
	contact_phone = $2,
	address_line_1 = $3,
	address_line_2 = $4,
	city = $5,
	state = $6,
	zip_code = $7,
	country = $8,
	latitude = $9,
    longitude = $10,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND belongs_to_user = $11
	AND id = $12;
