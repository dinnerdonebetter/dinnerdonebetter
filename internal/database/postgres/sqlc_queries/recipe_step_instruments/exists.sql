-- name: CheckRecipeStepInstrumentExistence :one

SELECT EXISTS ( SELECT recipe_step_instruments.id FROM recipe_step_instruments JOIN recipe_steps ON recipe_step_instruments.belongs_to_recipe_step=recipe_steps.id JOIN recipes ON recipe_steps.belongs_to_recipe=recipes.id WHERE recipe_step_instruments.archived_at IS NULL AND recipe_step_instruments.belongs_to_recipe_step = $1 AND recipe_step_instruments.id = $2 AND recipe_steps.archived_at IS NULL AND recipe_steps.belongs_to_recipe = $3 AND recipe_steps.id = $4 AND recipes.archived_at IS NULL AND recipes.id = $5 );
