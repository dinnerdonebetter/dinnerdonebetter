-- name: ValidIngredientMeasurementUnitPairIsValid :one

SELECT EXISTS(
	SELECT id
	FROM valid_ingredient_measurement_units
	WHERE valid_measurement_unit_id = $1
	AND valid_ingredient_id = $2
	AND archived_at IS NULL
);