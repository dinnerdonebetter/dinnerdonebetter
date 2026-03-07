-- Prevent duplicate meals in meal lists: one meal per list (non-archived).
CREATE UNIQUE INDEX IF NOT EXISTS idx_meal_list_items_meal_list_meal_unique
ON meal_list_items (belongs_to_meal_list, meal_id)
WHERE archived_at IS NULL;

-- Prevent duplicate meal options per event: one meal per event (non-archived).
CREATE UNIQUE INDEX IF NOT EXISTS idx_meal_plan_options_event_meal_unique
ON meal_plan_options (belongs_to_meal_plan_event, meal_id)
WHERE archived_at IS NULL
  AND belongs_to_meal_plan_event IS NOT NULL;
