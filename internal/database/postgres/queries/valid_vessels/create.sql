-- name: CreateValidVessel :exec

INSERT INTO valid_vessels (id,"name",plural_name,description,icon_path,usable_for_storage,slug,display_in_summary_lists,include_in_generated_instructions,capacity,capacity_unit,width_in_millimeters,length_in_millimeters,height_in_millimeters,shape) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);
