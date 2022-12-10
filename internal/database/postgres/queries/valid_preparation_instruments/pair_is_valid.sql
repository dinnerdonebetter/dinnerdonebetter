SELECT EXISTS(
	SELECT id
	FROM valid_preparation_instruments
	WHERE valid_instrument_id = $1
	AND valid_preparation_id = $2
	AND archived_at IS NULL
);
