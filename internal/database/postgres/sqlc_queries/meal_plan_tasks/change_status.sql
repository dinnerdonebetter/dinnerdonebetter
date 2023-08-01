UPDATE meal_plan_tasks SET completed_at = $4, status_explanation = $3, status = $2 WHERE id = $1;
