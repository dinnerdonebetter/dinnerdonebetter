SELECT recipes.id
  FROM recipes
 WHERE (recipes.archived_at IS NULL)
       AND (
			(recipes.last_indexed_at IS NULL)
			OR recipes.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);