-- name: MealPlanOptionVoteExists :one
SELECT EXISTS ( SELECT meal_plan_option_votes.id FROM meal_plan_option_votes JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id WHERE meal_plan_option_votes.archived_on IS NULL AND meal_plan_option_votes.belongs_to_meal_plan_option = $1 AND meal_plan_option_votes.id = $2 AND meal_plan_options.archived_on IS NULL AND meal_plan_options.belongs_to_meal_plan = $3 AND meal_plan_options.id = $4 AND meal_plans.archived_on IS NULL AND meal_plans.id = $5 );

-- name: GetMealPlanOptionVote :one
SELECT meal_plan_option_votes.id, meal_plan_option_votes.rank, meal_plan_option_votes.abstain, meal_plan_option_votes.notes, meal_plan_option_votes.by_user, meal_plan_option_votes.created_on, meal_plan_option_votes.last_updated_on, meal_plan_option_votes.archived_on, meal_plan_option_votes.belongs_to_meal_plan_option FROM meal_plan_option_votes JOIN meal_plan_options ON meal_plan_option_votes.belongs_to_meal_plan_option=meal_plan_options.id JOIN meal_plans ON meal_plan_options.belongs_to_meal_plan=meal_plans.id WHERE meal_plan_option_votes.archived_on IS NULL AND meal_plan_option_votes.belongs_to_meal_plan_option = $1 AND meal_plan_option_votes.id = $2 AND meal_plan_options.archived_on IS NULL AND meal_plan_options.belongs_to_meal_plan = $3 AND meal_plan_options.id = $4 AND meal_plans.archived_on IS NULL AND meal_plans.id = $5;

-- name: GetTotalMealPlanOptionVotesCount :one
SELECT COUNT(meal_plan_option_votes.id) FROM meal_plan_option_votes WHERE meal_plan_option_votes.archived_on IS NULL;

-- name: MealPlanOptionVoteCreation :exec
INSERT INTO meal_plan_option_votes (id,rank,abstain,notes,by_user,belongs_to_meal_plan_option) VALUES ($1,$2,$3,$4,$5,$6);

-- name: UpdateMealPlanOptionVote :exec
UPDATE meal_plan_option_votes SET rank = $1, abstain = $2, notes = $3, by_user = $4, last_updated_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_meal_plan_option = $5 AND id = $6;

-- name: ArchiveMealPlanOptionVote :exec
UPDATE meal_plan_option_votes SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND belongs_to_meal_plan_option = $1 AND id = $2;
