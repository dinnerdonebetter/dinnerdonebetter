SELECT EXISTS ( SELECT valid_measurement_conversions.id FROM valid_measurement_conversions WHERE valid_measurement_conversions.archived_at IS NULL AND valid_measurement_conversions.id = $1 );
