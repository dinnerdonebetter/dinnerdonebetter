SELECT EXISTS(
	SELECT id
	FROM valid_ingredient_preparations
	WHERE valid_ingredient_id = $1
	AND valid_preparation_id = $2
	AND archived_at IS NULL
);