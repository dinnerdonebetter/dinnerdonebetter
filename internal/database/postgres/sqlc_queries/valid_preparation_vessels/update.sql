UPDATE valid_preparation_vessels SET notes = $1, valid_preparation_id = $2, valid_vessel_id = $3, last_updated_at = NOW() WHERE archived_at IS NULL AND id = $4;
