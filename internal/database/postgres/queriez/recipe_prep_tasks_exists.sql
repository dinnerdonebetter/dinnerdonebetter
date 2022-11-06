-- name: RecipePrepTaskExists :exec
SELECT EXISTS (
    SELECT recipe_prep_tasks.id
    FROM recipe_prep_tasks
    JOIN recipes ON recipe_prep_tasks.belongs_to_recipe=recipes.id
    WHERE recipe_prep_tasks.archived_at IS NULL
      AND recipe_prep_tasks.belongs_to_recipe = $1
      AND recipe_prep_tasks.id = $2
      AND recipes.archived_at IS NULL
      AND recipes.id = $1
);
