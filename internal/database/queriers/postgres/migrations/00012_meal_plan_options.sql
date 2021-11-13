CREATE TABLE IF NOT EXISTS meal_plan_options (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"day" INTEGER NOT NULL CONSTRAINT valid_day CHECK (day > 0), CHECK (day < 7),
    "recipe_id" CHAR(27) NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
	"notes" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_meal_plan" CHAR(27) NOT NULL REFERENCES meal_plans("id") ON DELETE CASCADE,
    UNIQUE("recipe_id", "belongs_to_meal_plan")
);