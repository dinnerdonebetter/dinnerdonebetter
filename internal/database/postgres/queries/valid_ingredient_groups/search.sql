SELECT
	valid_ingredient_groups.id,
	valid_ingredient_groups.name,
	valid_ingredient_groups.description,
	valid_ingredient_groups.slug,
	valid_ingredient_groups.created_at,
	valid_ingredient_groups.last_updated_at,
	valid_ingredient_groups.archived_at,
	valid_ingredient_group_members.id,
    valid_ingredient_group_members.belongs_to_group,
    valid_ingredient_group_members.valid_ingredient,
    valid_ingredient_group_members.created_at,
    valid_ingredient_group_members.archived_at
FROM valid_ingredient_groups
 JOIN valid_ingredient_group_members ON valid_ingredient_group_members.belongs_to_group=valid_ingredient_groups.id
WHERE valid_ingredient_groups.name ILIKE $1
AND valid_ingredient_groups.archived_at IS NULL
AND valid_ingredient_group_members.archived_at IS NULL
LIMIT 50;
