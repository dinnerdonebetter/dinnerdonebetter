-- name: CheckUserIngredientPreferenceExistence :one

SELECT EXISTS ( SELECT user_ingredient_preferences.id FROM user_ingredient_preferences WHERE user_ingredient_preferences.archived_at IS NULL AND user_ingredient_preferences.id = $1 );
