CREATE TABLE IF NOT EXISTS valid_ingredient_lists (
    "id" TEXT NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL DEFAULT '',
    "slug" TEXT NOT NULL DEFAULT '',
    "description" TEXT NOT NULL DEFAULT '',
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE TABLE IF NOT EXISTS valid_ingredient_list_entries (
    "id" TEXT NOT NULL PRIMARY KEY,
    "belongs_to_list" TEXT NOT NULL REFERENCES valid_ingredient_lists("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "valid_ingredient" TEXT NOT NULL REFERENCES valid_ingredients("id") ON DELETE CASCADE ON UPDATE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-- CREATE TABLE IF NOT EXISTS valid_instrument_lists (
--     "id" TEXT NOT NULL PRIMARY KEY,
--     "name" TEXT NOT NULL DEFAULT '',
--     "description" TEXT NOT NULL DEFAULT '',
--     "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "last_updated_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL,
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
-- );
--
-- CREATE TABLE IF NOT EXISTS valid_instrument_list_entries (
--     "id" TEXT NOT NULL PRIMARY KEY,
--     "belongs_to_list" TEXT NOT NULL REFERENCES valid_instrument_lists("id") ON DELETE CASCADE ON UPDATE CASCADE,
--     "valid_instrument" TEXT NOT NULL REFERENCES valid_instruments("id") ON DELETE CASCADE ON UPDATE CASCADE,
--     "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
--     "archived_at" TIMESTAMP WITH TIME ZONE DEFAULT NULL
-- );
