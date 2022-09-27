INSERT INTO meal_plan_tasks (id,belongs_to_meal_plan_option,satisfies_recipe_step,status,status_explanation,creation_explanation,cannot_complete_before,cannot_complete_after,assigned_to_user)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);
