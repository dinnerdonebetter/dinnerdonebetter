ALTER TABLE recipes RENAME COLUMN yields_portions TO min_estimated_portions;
ALTER TABLE recipes ALTER COLUMN min_estimated_portions TYPE NUMERIC(14, 2);
ALTER TABLE recipes ADD COLUMN "max_estimated_portions" NUMERIC(14, 2);
