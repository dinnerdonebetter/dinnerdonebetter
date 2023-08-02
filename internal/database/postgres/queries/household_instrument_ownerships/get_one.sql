-- name: GetHouseholdInstrumentOwnership :one

SELECT
	household_instrument_ownerships.id,
	household_instrument_ownerships.notes,
	household_instrument_ownerships.quantity,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	household_instrument_ownerships.belongs_to_household,
	household_instrument_ownerships.created_at,
	household_instrument_ownerships.last_updated_at,
	household_instrument_ownerships.archived_at
FROM household_instrument_ownerships
INNER JOIN valid_instruments ON household_instrument_ownerships.valid_instrument_id = valid_instruments.id
WHERE household_instrument_ownerships.archived_at IS NULL
	AND household_instrument_ownerships.id = $1
	AND household_instrument_ownerships.belongs_to_household = $2;
