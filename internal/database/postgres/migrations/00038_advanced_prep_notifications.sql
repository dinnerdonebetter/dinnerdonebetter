CREATE TABLE IF NOT EXISTS advanced_prep_steps (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_meal_plan_option" CHAR(27) NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    "satisfies_recipe_step" CHAR(27) NOT NULL REFERENCES recipe_steps("id") ON DELETE CASCADE,
    "cannot_complete_before" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "cannot_complete_after" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "created_at" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "completed_at" BIGINT DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS advanced_prep_notifications (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "advanced_prep_step_id" CHAR(27) NOT NULL REFERENCES advanced_prep_steps("id") ON DELETE CASCADE,
    "notification_sent_at" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "completed_at" BIGINT DEFAULT NULL
);
