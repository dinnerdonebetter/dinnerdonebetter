-- name: UpdateValidVessel :exec

UPDATE valid_vessels
SET
    name = sqlc.arg(name),
    plural_name = sqlc.arg(plural_name),
    description = sqlc.arg(description),
    icon_path = sqlc.arg(icon_path),
    usable_for_storage = sqlc.arg(usable_for_storage),
    slug = sqlc.arg(slug),
    display_in_summary_lists = sqlc.arg(display_in_summary_lists),
    include_in_generated_instructions = sqlc.arg(include_in_generated_instructions),
    capacity = sqlc.arg(capacity)::float,
    capacity_unit = sqlc.arg(capacity_unit),
    width_in_millimeters = sqlc.arg(width_in_millimeters)::float,
    length_in_millimeters = sqlc.arg(length_in_millimeters)::float,
    height_in_millimeters = sqlc.arg(height_in_millimeters)::float,
    shape = sqlc.arg(shape),
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = sqlc.arg(id);
