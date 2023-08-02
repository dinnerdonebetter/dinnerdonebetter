-- name: GetValidIngredientStatesNeedingIndexing :many

SELECT valid_ingredient_states.id
  FROM valid_ingredient_states
 WHERE (valid_ingredient_states.archived_at IS NULL)
       AND (
			(
				valid_ingredient_states.last_indexed_at IS NULL
			)
			OR valid_ingredient_states.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);