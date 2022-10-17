SELECT
    meal_plan_grocery_list_items.id,
    meal_plan_grocery_list_items.belongs_to_meal_plan_option,
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
         JOIN meal_plan_options ON meal_plan_grocery_list_items.belongs_to_meal_plan_option=meal_plan_options.id
         JOIN meal_plan_events ON meal_plan_options.belongs_to_meal_plan_event=meal_plan_events.id
         JOIN meal_plans ON meal_plan_events.belongs_to_meal_plan=meal_plans.id
WHERE meal_plan_grocery_list_items.completed_at IS NULL
  AND meal_plans.archived_at IS NULL
  AND meal_plans.id = $1
GROUP BY meal_plan_grocery_list_items.id
ORDER BY meal_plan_grocery_list_items.id
