DROP TABLE meal_plan_task_recipe_steps;
DROP TABLE meal_plan_tasks;

DROP TYPE prep_step_status;
CREATE TYPE prep_step_status AS ENUM ('unfinished', 'postponed', 'ignored', 'canceled', 'finished');

CREATE TABLE IF NOT EXISTS meal_plan_tasks (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_meal_plan_option" CHAR(27) NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    "belongs_to_recipe_prep_task" CHAR(27) NOT NULL REFERENCES recipe_prep_tasks("id") ON DELETE CASCADE,
    "creation_explanation" TEXT NOT NULL DEFAULT '',
    "status_explanation" TEXT NOT NULL DEFAULT '',
    "status" prep_step_status NOT NULL DEFAULT 'unfinished',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "assigned_to_user" CHAR(27) REFERENCES users("id") ON DELETE CASCADE,
    "completed_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);


