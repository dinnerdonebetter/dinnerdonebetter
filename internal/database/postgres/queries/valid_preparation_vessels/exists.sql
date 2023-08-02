-- name: CheckValidPreparationVesselExistence :one

SELECT EXISTS ( SELECT valid_preparation_vessels.id FROM valid_preparation_vessels WHERE valid_preparation_vessels.archived_at IS NULL AND valid_preparation_vessels.id = $1 );
