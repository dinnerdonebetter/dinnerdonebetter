-- name: CreateValidVessel :exec

INSERT INTO valid_vessels (id,"name",plural_name,description,icon_path,usable_for_storage,slug,display_in_summary_lists,include_in_generated_instructions,capacity,capacity_unit,width_in_millimeters,length_in_millimeters,height_in_millimeters,shape)
    VALUES (sqlc.arg(id),sqlc.arg(name),sqlc.arg(plural_name),sqlc.arg(description),sqlc.arg(icon_path),sqlc.arg(usable_for_storage),sqlc.arg(slug),sqlc.arg(display_in_summary_lists),sqlc.arg(include_in_generated_instructions),sqlc.arg(capacity)::float,sqlc.arg(capacity_unit),sqlc.arg(width_in_millimeters)::float,sqlc.arg(length_in_millimeters)::float,sqlc.arg(height_in_millimeters)::float,sqlc.arg(shape));
