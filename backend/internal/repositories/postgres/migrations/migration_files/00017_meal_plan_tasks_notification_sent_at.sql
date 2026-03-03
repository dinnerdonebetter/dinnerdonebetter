-- Meal Plan Tasks: notification_sent_at column
-- Tracks when a push notification was sent for a meal plan task (for idempotency)

ALTER TABLE meal_plan_tasks
    ADD COLUMN IF NOT EXISTS notification_sent_at TIMESTAMP WITH TIME ZONE;
