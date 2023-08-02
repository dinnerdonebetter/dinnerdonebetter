-- name: CreateMealPlanGroceryListItem :exec

INSERT INTO meal_plan_grocery_list_items
(id,belongs_to_meal_plan,valid_ingredient,valid_measurement_unit,minimum_quantity_needed,maximum_quantity_needed,quantity_purchased,purchased_measurement_unit,purchased_upc,purchase_price,status_explanation,status)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12);
