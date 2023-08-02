-- name: CheckHouseholdInstrumentOwnershipExistence :one

SELECT EXISTS ( SELECT household_instrument_ownerships.id FROM household_instrument_ownerships WHERE household_instrument_ownerships.archived_at IS NULL AND household_instrument_ownerships.id = $1 AND household_instrument_ownerships.belongs_to_household = $2 );