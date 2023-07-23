ALTER TABLE valid_preparations ALTER COLUMN maximum_vessel_count DROP DEFAULT;
ALTER TABLE valid_preparations ALTER COLUMN maximum_vessel_count DROP NOT NULL;
UPDATE valid_preparations SET maximum_vessel_count = NULL WHERE maximum_vessel_count <= 0;

ALTER TABLE valid_preparations DROP COLUMN universal;
