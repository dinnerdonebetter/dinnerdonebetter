UPDATE valid_measurement_units SET
    name = $1,
    description = $2,
    volumetric = $3,
    icon_path = $4,
    universal = $5,
    metric = $6,
    imperial = $7,
    plural_name = $8,
    last_updated_at = NOW()
WHERE archived_at IS NULL AND id = $9;
