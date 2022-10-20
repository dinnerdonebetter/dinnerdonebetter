CREATE TABLE IF NOT EXISTS recipe_media (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "belongs_to_recipe" CHAR(27) REFERENCES recipes("id") ON DELETE CASCADE,
    "belongs_to_recipe_step" CHAR(27) REFERENCES recipe_steps("id") ON DELETE CASCADE,
    "mime_type" TEXT NOT NULL,
    "internal_path" TEXT NOT NULL,
    "external_path" TEXT NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);