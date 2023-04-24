UPDATE service_setting_configurations
SET
    value = $1
    notes = $2
    service_setting_id = $3
    belongs_to_user = $4
    belongs_to_household = $5
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
