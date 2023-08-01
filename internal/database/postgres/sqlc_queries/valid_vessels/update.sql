UPDATE valid_vessels
SET
    name = $1,
    plural_name = $2,
    description = $3,
    icon_path = $4,
    usable_for_storage = $5,
    slug = $6,
    display_in_summary_lists = $7,
    include_in_generated_instructions = $8,
    capacity = $9,
    capacity_unit = $10,
    width_in_millimeters = $11,
    length_in_millimeters = $12,
    height_in_millimeters = $13,
    shape = $14,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $15;
