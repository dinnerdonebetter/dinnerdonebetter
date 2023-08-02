-- name: UpdateHouseholdInstrumentOwnership :exec

UPDATE household_instrument_ownerships
SET
	notes = $1,
	quantity = $2,
	valid_instrument_id = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $4
	AND household_instrument_ownerships.belongs_to_household = $5;
