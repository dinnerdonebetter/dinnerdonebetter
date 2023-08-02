-- name: CreateMeal :exec

INSERT INTO meals (id,"name",description,min_estimated_portions,max_estimated_portions,eligible_for_meal_plans,created_by_user) VALUES ($1,$2,$3,$4,$5,$6,$7);