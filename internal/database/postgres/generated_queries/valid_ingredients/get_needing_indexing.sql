SELECT valid_ingredients.id
  FROM valid_ingredients
 WHERE (valid_ingredients.archived_at IS NULL)
   AND (
		(valid_ingredients.last_indexed_at IS NULL)
		OR valid_ingredients.last_indexed_at < now() - '24 hours'::INTERVAL
       );
