-- name: CreateMeal :exec
INSERT INTO meals (id,name,description,created_by_user) VALUES ($1,$2,$3,$4);