UPDATE recipes SET
    name = $1,
    slug = $2,
    source = $3,
    description = $4,
    inspired_by_recipe_id = $5,
	min_estimated_portions = $6,
	max_estimated_portions = $7,
    portion_name = $8,
    plural_portion_name = $9,
    seal_of_approval = $10,
    eligible_for_meals = $11,
    last_updated_at = NOW()
WHERE archived_at IS NULL
  AND created_by_user = $12
  AND id = $13;
