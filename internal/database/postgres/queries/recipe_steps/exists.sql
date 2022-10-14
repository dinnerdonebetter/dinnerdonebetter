SELECT EXISTS (
    SELECT recipe_steps.id
    FROM recipe_steps
    JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id
    WHERE recipe_steps.archived_at IS NULL
      AND recipe_steps.belongs_to_recipe = $1
      AND recipe_steps.id = $2
      AND recipes.archived_at IS NULL
      AND recipes.id = $1
);
