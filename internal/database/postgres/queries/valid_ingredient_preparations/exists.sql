SELECT EXISTS ( SELECT valid_ingredient_preparations.id FROM valid_ingredient_preparations WHERE valid_ingredient_preparations.archived_at IS NULL AND valid_ingredient_preparations.id = $1 );
