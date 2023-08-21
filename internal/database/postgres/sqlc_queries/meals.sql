-- name: ArchiveMeal :exec

UPDATE meals SET archived_at = NOW() WHERE archived_at IS NULL AND created_by_user = $1 AND id = $2;


-- name: CreateMeal :exec

INSERT INTO meals (id,"name",description,min_estimated_portions,max_estimated_portions,eligible_for_meal_plans,created_by_user) VALUES ($1,$2,$3,$4,$5,$6,$7);

-- name: CheckMealExistence :one

SELECT EXISTS ( SELECT meals.id FROM meals WHERE meals.archived_at IS NULL AND meals.id = $1 );


-- name: GetMealsNeedingIndexing :many

SELECT meals.id
  FROM meals
 WHERE (meals.archived_at IS NULL)
       AND (
			(meals.last_indexed_at IS NULL)
			OR meals.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);

-- name: GetMeal :one

SELECT
	meals.id,
	meals.name,
	meals.description,
	meals.min_estimated_portions,
    meals.max_estimated_portions,
    meals.eligible_for_meal_plans,
	meals.created_at,
	meals.last_updated_at,
	meals.archived_at,
	meals.created_by_user,
	meal_components.recipe_id,
	meal_components.recipe_scale,
	meal_components.meal_component_type
FROM meals
	FULL OUTER JOIN meal_components ON meal_components.meal_id=meals.id
WHERE meals.archived_at IS NULL
	AND meal_components.archived_at IS NULL
	AND meals.id = $1;


-- name: UpdateMealLastIndexedAt :exec

UPDATE meals SET last_indexed_at = NOW() WHERE id = $1 AND archived_at IS NULL;
