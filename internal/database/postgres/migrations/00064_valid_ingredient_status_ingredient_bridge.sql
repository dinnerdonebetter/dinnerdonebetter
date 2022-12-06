-- CREATE TABLE IF NOT EXISTS valid_ingredient_states (
--     id CHAR(27) NOT NULL PRIMARY KEY,
--     "name" TEXT NOT NULL,
--     past_tense TEXT NOT NULL DEFAULT '',
--     slug TEXT NOT NULL DEFAULT '',
--     description TEXT NOT NULL DEFAULT '',
--     icon_path TEXT NOT NULL DEFAULT '',
--     created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
--     last_updated_at TIMESTAMP WITH TIME ZONE,
--     archived_at TIMESTAMP WITH TIME ZONE
-- );

CREATE TABLE IF NOT EXISTS valid_ingredient_state_ingredients (
    "id" TEXT NOT NULL PRIMARY KEY,
    "valid_ingredient" TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE,
    "valid_ingredient_state" TEXT NOT NULL REFERENCES valid_ingredient_states("id") ON DELETE CASCADE,
    "notes" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE,
    "archived_at" TIMESTAMP WITH TIME ZONE,
    UNIQUE("valid_ingredient", "valid_ingredient_state", "archived_at")
);