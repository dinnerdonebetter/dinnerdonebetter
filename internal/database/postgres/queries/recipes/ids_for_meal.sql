SELECT
	recipes.id
FROM
	recipes
	    JOIN meal_recipes ON meal_recipes.recipe_id = recipes.id
	    JOIN meals ON meal_recipes.meal_id = meals.id
WHERE
	recipes.archived_at IS NULL
  AND meals.id = $1
GROUP BY
	recipes.id
ORDER BY
	recipes.id;
