INSERT INTO meal_plan_tasks (id,status,status_explanation,creation_explanation,cannot_complete_before,cannot_complete_after,belongs_to_meal_plan_option,assigned_to_user)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8);
