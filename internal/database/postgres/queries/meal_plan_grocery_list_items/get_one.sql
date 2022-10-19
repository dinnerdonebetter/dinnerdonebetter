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
    meal_plan_grocery_list_items.completed_at
FROM meal_plan_grocery_list_items
    FULL OUTER JOIN meal_plans ON meal_plan_grocery_list_items.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_grocery_list_items.completed_at IS NULL
  AND meal_plan_grocery_list_items.id = $2
  AND meal_plan_grocery_list_items.belongs_to_meal_plan = $1;
