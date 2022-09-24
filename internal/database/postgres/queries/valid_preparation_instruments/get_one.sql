SELECT
	valid_preparation_instruments.id,
	valid_preparation_instruments.notes,
	valid_preparations.id,
	valid_preparations.name,
	valid_preparations.description,
	valid_preparations.icon_path,
	valid_preparations.yields_nothing,
	valid_preparations.restrict_to_ingredients,
	valid_preparations.zero_ingredients_allowable,
	valid_preparations.past_tense,
	valid_preparations.created_at,
	valid_preparations.last_updated_at,
	valid_preparations.archived_at,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	valid_preparation_instruments.created_at,
	valid_preparation_instruments.last_updated_at,
	valid_preparation_instruments.archived_at
FROM
	valid_preparation_instruments
	    JOIN valid_instruments ON valid_preparation_instruments.valid_instrument_id = valid_instruments.id
	    JOIN valid_preparations ON valid_preparation_instruments.valid_preparation_id = valid_preparations.id
WHERE
	valid_preparation_instruments.archived_at IS NULL
  AND valid_preparation_instruments.id = $1;
