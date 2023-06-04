UPDATE household_instrument_ownerships SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1 AND belongs_to_household = $2;
