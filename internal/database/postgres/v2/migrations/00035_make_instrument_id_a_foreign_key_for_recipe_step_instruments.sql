ALTER TABLE recipe_step_instruments ALTER COLUMN instrument_id TYPE CHAR(27);
ALTER TABLE recipe_step_instruments ADD CONSTRAINT fk_valid_recipe_step_instrument_ids FOREIGN KEY (instrument_id) REFERENCES valid_instruments (id);
