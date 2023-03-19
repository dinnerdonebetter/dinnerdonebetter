UPDATE recipes SET
    name = $1,
    slug = $2,
    source = $3,
    description = $4,
    inspired_by_recipe_id = $5,
    yields_portions = $6,
    portion_name = $7,
    plural_portion_name = $8,
    seal_of_approval = $9,
    last_updated_at = NOW()
WHERE archived_at IS NULL
  AND created_by_user = $10
  AND id = $11;
