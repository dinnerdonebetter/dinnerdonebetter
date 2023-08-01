SELECT meals.id
  FROM meals
 WHERE (meals.archived_at IS NULL)
       AND (
			(meals.last_indexed_at IS NULL)
			OR meals.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);