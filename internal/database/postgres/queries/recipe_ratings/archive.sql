UPDATE recipe_ratings SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;
