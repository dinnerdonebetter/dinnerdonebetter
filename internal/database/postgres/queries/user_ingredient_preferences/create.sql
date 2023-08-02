-- name: CreateUserIngredientPreference :exec

INSERT INTO user_ingredient_preferences (id,ingredient,rating,notes,allergy,belongs_to_user) VALUES ($1,$2,$3,$4,$5,$6);
