-- name: CreateHouseholdInstrumentOwnership :exec

INSERT INTO household_instrument_ownerships (id,notes,quantity,valid_instrument_id,belongs_to_household) VALUES ($1,$2,$3,$4,$5);
