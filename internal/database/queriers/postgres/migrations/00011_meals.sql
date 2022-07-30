CREATE TABLE IF NOT EXISTS meals (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL,
    "created_by_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS meals_created_by_user ON meals (created_by_user);

CREATE TABLE IF NOT EXISTS meal_recipes (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "meal_id" CHAR(27) NOT NULL REFERENCES meals("id") ON DELETE CASCADE,
    "recipe_id" CHAR(27) NOT NULL REFERENCES recipes("id") ON DELETE CASCADE,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "archived_on" BIGINT DEFAULT NULL
);

CREATE INDEX IF NOT EXISTS meal_recipes_meal_id on meal_recipes (meal_id);
CREATE INDEX IF NOT EXISTS meal_recipes_recipe_id on meal_recipes (recipe_id);