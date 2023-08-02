-- name: UpdateUserIngredientPreference :exec

UPDATE user_ingredient_preferences
SET
	ingredient = $1,
	rating = $2,
	notes = $3,
	allergy = $4,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $5
	AND belongs_to_user = $6;
