ALTER TABLE recipe_steps ALTER COLUMN minimum_estimated_time_in_seconds DROP DEFAULT;
ALTER TABLE recipe_steps ALTER COLUMN minimum_estimated_time_in_seconds DROP NOT NULL;
ALTER TABLE recipe_steps ALTER COLUMN maximum_estimated_time_in_seconds DROP DEFAULT;
ALTER TABLE recipe_steps ALTER COLUMN maximum_estimated_time_in_seconds DROP NOT NULL;

UPDATE recipe_steps SET minimum_estimated_time_in_seconds = NULL where minimum_estimated_time_in_seconds = 0;
UPDATE recipe_steps SET maximum_estimated_time_in_seconds = NULL where maximum_estimated_time_in_seconds = 0;
