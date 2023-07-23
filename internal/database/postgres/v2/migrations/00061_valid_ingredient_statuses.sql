CREATE TABLE IF NOT EXISTS valid_ingredient_states (
    id CHAR(27) NOT NULL PRIMARY KEY,
    "name" TEXT NOT NULL,
    past_tense TEXT NOT NULL DEFAULT '',
    slug TEXT NOT NULL DEFAULT '',
    description TEXT NOT NULL DEFAULT '',
    icon_path TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    last_updated_at TIMESTAMP WITH TIME ZONE,
    archived_at TIMESTAMP WITH TIME ZONE
);
