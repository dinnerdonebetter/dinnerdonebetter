SELECT EXISTS (
	SELECT recipe_prep_tasks.id
	FROM recipe_prep_tasks
	WHERE recipe_prep_tasks.archived_at IS NULL
		AND recipe_prep_tasks.id = $1
);
