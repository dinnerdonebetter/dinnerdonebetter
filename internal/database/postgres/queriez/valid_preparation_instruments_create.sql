-- name: CreateValidPreparationInstrument :exec
INSERT INTO valid_preparation_instruments (id,notes,valid_preparation_id,valid_instrument_id) VALUES ($1,$2,$3,$4);
