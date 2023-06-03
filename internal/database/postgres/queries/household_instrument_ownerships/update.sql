UPDATE household_instrument_ownerships
SET
	notes = $1,
	quantity = $2,
	valid_instrument_id = $3,
	belongs_to_household = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $5;
