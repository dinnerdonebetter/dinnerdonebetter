UPDATE user_ingredient_preferences SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1 AND belongs_to_user = $2;
