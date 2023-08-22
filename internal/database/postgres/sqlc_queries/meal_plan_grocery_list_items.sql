-- name: ArchiveMealPlanGroceryListItem :exec

UPDATE meal_plan_grocery_list_items SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1;

-- name: CreateMealPlanGroceryListItem :exec

INSERT INTO meal_plan_grocery_list_items
(id,belongs_to_meal_plan,valid_ingredient,valid_measurement_unit,minimum_quantity_needed,maximum_quantity_needed,quantity_purchased,purchased_measurement_unit,purchased_upc,purchase_price,status_explanation,status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);

-- name: CheckMealPlanGroceryListItemExistence :one

SELECT EXISTS ( SELECT meal_plan_grocery_list_items.id FROM meal_plan_grocery_list_items WHERE meal_plan_grocery_list_items.archived_at IS NULL AND meal_plan_grocery_list_items.id = sqlc.arg(meal_plan_grocery_list_item_id) AND meal_plan_grocery_list_items.belongs_to_meal_plan = sqlc.arg(meal_plan_id) );

-- name: GetMealPlanGroceryListItemsForMealPlan :many

SELECT
	meal_plan_grocery_list_items.id,
	meal_plan_grocery_list_items.belongs_to_meal_plan,
	meal_plan_grocery_list_items.valid_ingredient,
	meal_plan_grocery_list_items.valid_measurement_unit,
	meal_plan_grocery_list_items.minimum_quantity_needed,
	meal_plan_grocery_list_items.maximum_quantity_needed,
	meal_plan_grocery_list_items.quantity_purchased,
	meal_plan_grocery_list_items.purchased_measurement_unit,
	meal_plan_grocery_list_items.purchased_upc,
	meal_plan_grocery_list_items.purchase_price,
	meal_plan_grocery_list_items.status_explanation,
	meal_plan_grocery_list_items.status,
	meal_plan_grocery_list_items.created_at,
	meal_plan_grocery_list_items.last_updated_at,
	meal_plan_grocery_list_items.archived_at
FROM meal_plan_grocery_list_items
	FULL OUTER JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_grocery_list_items.archived_at IS NULL
  AND meal_plan_grocery_list_items.belongs_to_meal_plan = $1
  AND meal_plans.archived_at IS NULL
  AND meal_plans.id = $1
GROUP BY meal_plan_grocery_list_items.id
ORDER BY meal_plan_grocery_list_items.id;

-- name: GetMealPlanGroceryListItem :one

SELECT
	meal_plan_grocery_list_items.id,
	meal_plan_grocery_list_items.belongs_to_meal_plan,
	meal_plan_grocery_list_items.valid_ingredient,
	meal_plan_grocery_list_items.valid_measurement_unit,
	meal_plan_grocery_list_items.minimum_quantity_needed,
	meal_plan_grocery_list_items.maximum_quantity_needed,
	meal_plan_grocery_list_items.quantity_purchased,
	meal_plan_grocery_list_items.purchased_measurement_unit,
	meal_plan_grocery_list_items.purchased_upc,
	meal_plan_grocery_list_items.purchase_price,
	meal_plan_grocery_list_items.status_explanation,
	meal_plan_grocery_list_items.status,
	meal_plan_grocery_list_items.created_at,
	meal_plan_grocery_list_items.last_updated_at,
	meal_plan_grocery_list_items.archived_at
FROM meal_plan_grocery_list_items
	FULL OUTER JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_grocery_list_items.archived_at IS NULL
  AND meal_plan_grocery_list_items.id = $2
  AND meal_plan_grocery_list_items.belongs_to_meal_plan = $1;

-- name: UpdateMealPlanGroceryListItem :exec

UPDATE meal_plan_grocery_list_items
SET
	belongs_to_meal_plan = $1,
	valid_ingredient = $2,
	valid_measurement_unit = $3,
	minimum_quantity_needed = $4,
	maximum_quantity_needed = $5,
	quantity_purchased = $6,
	purchased_measurement_unit = $7,
	purchased_upc = $8,
	purchase_price = $9,
	status_explanation = $10,
	status = $11,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $12;
