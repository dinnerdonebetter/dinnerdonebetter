SELECT
	meal_plans.id,
	meal_plans.notes,
	meal_plans.status,
	meal_plans.voting_deadline,
    meal_plans.grocery_list_initialized,
    meal_plans.tasks_created,
	meal_plans.created_at,
	meal_plans.last_updated_at,
	meal_plans.archived_at,
	meal_plans.belongs_to_household
FROM meal_plans
WHERE meal_plans.archived_at IS NULL
	AND meal_plans.status = 'awaiting_votes'
	AND voting_deadline < now()
GROUP BY meal_plans.id
ORDER BY meal_plans.id;
