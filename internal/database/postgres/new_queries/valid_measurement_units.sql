-- name: ValidMeasurementUnitExists :one
SELECT EXISTS ( SELECT valid_measurement_units.id FROM valid_measurement_units WHERE valid_measurement_units.archived_on IS NULL AND valid_measurement_units.id = $1 );

-- name: GetValidMeasurementUnit :one
SELECT
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.plural_name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.universal,
    valid_measurement_units.metric,
    valid_measurement_units.imperial,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on
FROM valid_measurement_units
WHERE valid_measurement_units.archived_on IS NULL
  AND valid_measurement_units.id = $1;

-- name: GetRandomValidMeasurementUnit :one
SELECT
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.plural_name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.universal,
    valid_measurement_units.metric,
    valid_measurement_units.imperial,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on
FROM valid_measurement_units
WHERE valid_measurement_units.archived_on IS NULL
    ORDER BY random() LIMIT 1;

-- name: SearchForValidMeasurementUnits :many
SELECT
    valid_measurement_units.id,
    valid_measurement_units.name,
    valid_measurement_units.plural_name,
    valid_measurement_units.description,
    valid_measurement_units.volumetric,
    valid_measurement_units.universal,
    valid_measurement_units.metric,
    valid_measurement_units.imperial,
    valid_measurement_units.icon_path,
    valid_measurement_units.created_on,
    valid_measurement_units.last_updated_on,
    valid_measurement_units.archived_on
FROM valid_measurement_units
WHERE valid_measurement_units.name ILIKE $1
AND valid_measurement_units.archived_on IS NULL
LIMIT 50;

-- name: GetTotalValidMeasurementUnitCount :one
SELECT COUNT(valid_measurement_units.id) FROM valid_measurement_units WHERE valid_measurement_units.archived_on IS NULL;

-- name: CreateValidMeasurementUnit :exec
INSERT INTO valid_measurement_units (id,name,plural_name,description,volumetric,universal,metric,imperial,icon_path) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);

-- name: UpdateValidMeasurementUnit :exec
UPDATE valid_measurement_units SET
   name = $1,
   plural_name = $2,
   description = $3,
   volumetric = $4,
   universal = $5,
   metric = $6,
   imperial = $7,
   icon_path = $8,
   last_updated_on = extract(epoch FROM NOW())
WHERE archived_on IS NULL AND id = $9;

-- name: ArchiveValidMeasurementUnit :exec
UPDATE valid_measurement_units SET archived_on = extract(epoch FROM NOW()) WHERE archived_on IS NULL AND id = $1;
