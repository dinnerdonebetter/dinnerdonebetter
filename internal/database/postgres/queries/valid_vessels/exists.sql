-- name: CheckValidVesselExistence :one

SELECT EXISTS ( SELECT valid_vessels.id FROM valid_vessels WHERE valid_vessels.archived_at IS NULL AND valid_vessels.id = $1 );