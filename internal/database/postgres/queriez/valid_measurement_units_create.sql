-- name: CreateValidMeasurementUnit :exec
INSERT INTO valid_measurement_units
(id,name,description,volumetric,icon_path,universal,metric,imperial,plural_name)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9);
