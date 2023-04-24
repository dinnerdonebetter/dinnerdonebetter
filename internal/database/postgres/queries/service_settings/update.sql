UPDATE service_settings
SET
    "name" = $1,
    "type" = $2,
    description = $3,
    default_value = $4,
    admins_only = $5,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $6;
