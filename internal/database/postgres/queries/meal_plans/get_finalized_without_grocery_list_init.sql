SELECT
	meal_plans.id,
	meal_plans.belongs_to_household
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.status = 'finalized'
	AND meal_plans.grocery_list_initialized IS FALSE;
