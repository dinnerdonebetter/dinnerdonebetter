-- name: CreateMealPlan :exec
INSERT INTO meal_plans (id,notes,status,voting_deadline,belongs_to_household) VALUES ($1,$2,$3,$4,$5);