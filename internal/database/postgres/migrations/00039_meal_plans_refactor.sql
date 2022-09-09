CREATE TABLE IF NOT EXISTS meal_plan_events (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "notes" TEXT NOT NULL DEFAULT '',
    "starts_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "ends_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "belongs_to_meal_plan" CHAR(27) NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

DELETE FROM meal_plans WHERE id IS NOT NULL;

ALTER TABLE meal_plans DROP COLUMN "starts_at";
ALTER TABLE meal_plans DROP COLUMN "ends_at";
ALTER TABLE meal_plan_options DROP COLUMN "belongs_to_meal_plan";
ALTER TABLE meal_plan_options ADD COLUMN "belongs_to_meal_plan_event" CHAR(27) REFERENCES meal_plan_events("id") ON DELETE CASCADE;
