-- Prevent duplicate active rows in mealplanning bridge and enum tables.
-- The existing UNIQUE(..., archived_at) constraints do not prevent duplicates when
-- archived_at IS NULL (PostgreSQL treats NULL as distinct in standard UNIQUE).
-- Add partial unique indexes for active rows only.
--
-- For each table: de-duplicate (archive all but one per business key), then create index.

-- valid_preparation_instruments
WITH duplicates AS (
	SELECT id,
		ROW_NUMBER() OVER (PARTITION BY valid_preparation_id, valid_instrument_id ORDER BY id) AS rn
	FROM valid_preparation_instruments
	WHERE archived_at IS NULL
)
UPDATE valid_preparation_instruments
SET archived_at = NOW()
WHERE id IN (SELECT id FROM duplicates WHERE rn > 1);

CREATE UNIQUE INDEX idx_valid_preparation_instruments_prep_instrument_active
	ON valid_preparation_instruments (valid_preparation_id, valid_instrument_id)
	WHERE archived_at IS NULL;

-- valid_ingredient_preparations
WITH duplicates AS (
	SELECT id,
		ROW_NUMBER() OVER (PARTITION BY valid_preparation_id, valid_ingredient_id ORDER BY id) AS rn
	FROM valid_ingredient_preparations
	WHERE archived_at IS NULL
)
UPDATE valid_ingredient_preparations
SET archived_at = NOW()
WHERE id IN (SELECT id FROM duplicates WHERE rn > 1);

CREATE UNIQUE INDEX idx_valid_ingredient_preparations_prep_ing_active
	ON valid_ingredient_preparations (valid_preparation_id, valid_ingredient_id)
	WHERE archived_at IS NULL;

-- valid_ingredient_measurement_units
WITH duplicates AS (
	SELECT id,
		ROW_NUMBER() OVER (PARTITION BY valid_ingredient_id, valid_measurement_unit_id ORDER BY id) AS rn
	FROM valid_ingredient_measurement_units
	WHERE archived_at IS NULL
)
UPDATE valid_ingredient_measurement_units
SET archived_at = NOW()
WHERE id IN (SELECT id FROM duplicates WHERE rn > 1);

CREATE UNIQUE INDEX idx_valid_ingredient_measurement_units_ing_unit_active
	ON valid_ingredient_measurement_units (valid_ingredient_id, valid_measurement_unit_id)
	WHERE archived_at IS NULL;

-- valid_preparation_vessels
WITH duplicates AS (
	SELECT id,
		ROW_NUMBER() OVER (PARTITION BY valid_preparation_id, valid_vessel_id ORDER BY id) AS rn
	FROM valid_preparation_vessels
	WHERE archived_at IS NULL
)
UPDATE valid_preparation_vessels
SET archived_at = NOW()
WHERE id IN (SELECT id FROM duplicates WHERE rn > 1);

CREATE UNIQUE INDEX idx_valid_preparation_vessels_prep_vessel_active
	ON valid_preparation_vessels (valid_preparation_id, valid_vessel_id)
	WHERE archived_at IS NULL;

-- valid_ingredient_state_ingredients
WITH duplicates AS (
	SELECT id,
		ROW_NUMBER() OVER (PARTITION BY valid_ingredient, valid_ingredient_state ORDER BY id) AS rn
	FROM valid_ingredient_state_ingredients
	WHERE archived_at IS NULL
)
UPDATE valid_ingredient_state_ingredients
SET archived_at = NOW()
WHERE id IN (SELECT id FROM duplicates WHERE rn > 1);

CREATE UNIQUE INDEX idx_valid_ingredient_state_ingredients_ing_state_active
	ON valid_ingredient_state_ingredients (valid_ingredient, valid_ingredient_state)
	WHERE archived_at IS NULL;
