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
    valid_ingredient_group_members.archived_at,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
          AND valid_ingredient_groups.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
          AND valid_ingredient_groups.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years'))
            )
          AND (
                valid_ingredient_groups.last_updated_at IS NULL
                OR valid_ingredient_groups.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years'))
            )
        OFFSET COALESCE($1, 0)
    ) as filtered_count,
    (
        SELECT
            COUNT(valid_ingredient_groups.id)
        FROM
            valid_ingredient_groups
        WHERE
            valid_ingredient_groups.archived_at IS NULL
    ) as total_count
FROM valid_ingredient_groups
  JOIN valid_ingredient_group_members ON valid_ingredient_groups.id = valid_ingredient_group_members.belongs_to_group
WHERE
	valid_ingredient_groups.archived_at IS NULL
	AND valid_ingredient_group_members.archived_at IS NULL
	AND valid_ingredient_groups.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
	AND valid_ingredient_groups.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
	AND (
	    valid_ingredient_groups.last_updated_at IS NULL
	    OR valid_ingredient_groups.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years'))
	)
	AND (
	    valid_ingredient_groups.last_updated_at IS NULL
	    OR valid_ingredient_groups.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years'))
	)
	OFFSET COALESCE($1, 0);
