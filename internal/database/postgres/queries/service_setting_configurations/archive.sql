UPDATE service_setting_configurations SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
