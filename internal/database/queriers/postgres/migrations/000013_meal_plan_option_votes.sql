CREATE TABLE IF NOT EXISTS meal_plan_option_votes (
	"id" CHAR(27) NOT NULL PRIMARY KEY,
	"points" INTEGER NOT NULL,
	"abstain" BOOLEAN NOT NULL,
	"notes" TEXT NOT NULL,
	"created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
	"last_updated_on" BIGINT DEFAULT NULL,
	"archived_on" BIGINT DEFAULT NULL,
	"belongs_to_meal_plan_option" CHAR(27) NOT NULL REFERENCES meal_plan_options("id") ON DELETE CASCADE
);