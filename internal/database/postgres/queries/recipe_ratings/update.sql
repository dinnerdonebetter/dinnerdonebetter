UPDATE recipe_ratings
SET
	recipe_id = $1,
    taste = $2,
    difficulty = $3,
    cleanup = $4,
    instructions = $5,
    overall = $6,
    notes = $7,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $8;
