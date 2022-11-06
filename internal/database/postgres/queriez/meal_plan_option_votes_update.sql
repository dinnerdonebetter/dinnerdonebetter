-- name: UpdateMealPlanOptionVote :exec
UPDATE meal_plan_option_votes SET rank = $1, abstain = $2, notes = $3, by_user = $4, last_updated_at = NOW() WHERE archived_at IS NULL AND belongs_to_meal_plan_option = $5 AND id = $6;
