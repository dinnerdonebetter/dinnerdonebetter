-- name: CheckValidPreparationExistence :one

SELECT EXISTS ( SELECT valid_preparations.id FROM valid_preparations WHERE valid_preparations.archived_at IS NULL AND valid_preparations.id = $1 );
