UPDATE households SET belongs_to_user = $1 WHERE archived_at IS NULL AND belongs_to_user = $2 AND id = $3;
