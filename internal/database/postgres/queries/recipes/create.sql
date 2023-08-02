-- name: CreateRecipe :exec

INSERT INTO recipes (id,"name",slug,"source",description,inspired_by_recipe_id,min_estimated_portions,max_estimated_portions,portion_name,plural_portion_name,seal_of_approval,eligible_for_meals,yields_component_type,created_by_user) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14);
