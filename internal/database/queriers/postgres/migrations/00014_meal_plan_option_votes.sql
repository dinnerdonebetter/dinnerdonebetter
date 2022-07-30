CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"rank" INTEGER NOT NULL,
	"abstain" BOOLEAN NOT NULL,
	"notes" TEXT NOT NULL,
    "by_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_meal_plan_option" CHAR(27) NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE,
    UNIQUE("by_user", "belongs_to_meal_plan_option")
);

CREATE INDEX IF NOT EXISTS meal_plan_options_belongs_to_meal_plan_option on meal_plan_option_votes (belongs_to_meal_plan_option);
