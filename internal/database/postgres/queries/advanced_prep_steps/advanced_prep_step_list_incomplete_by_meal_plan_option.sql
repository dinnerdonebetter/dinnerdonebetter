SELECT
    advanced_prep_steps.id,
    advanced_prep_steps.belongs_to_meal_plan_option,
    advanced_prep_steps.satisfies_recipe_step,
    advanced_prep_steps.cannot_complete_before,
    advanced_prep_steps.cannot_complete_after,
    advanced_prep_steps.created_at,
    advanced_prep_steps.completed_at
FROM advanced_prep_steps
WHERE advanced_prep_steps.belongs_to_meal_plan_option = $1
AND advanced_prep_steps.completed_at IS NULL;
