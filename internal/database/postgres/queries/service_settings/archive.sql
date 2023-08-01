UPDATE service_settings
SET archived_at = NOW()
    WHERE id = $1;
