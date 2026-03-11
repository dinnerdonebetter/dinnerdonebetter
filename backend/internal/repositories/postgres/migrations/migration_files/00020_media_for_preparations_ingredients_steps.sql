-- Media for Preparations, Ingredients, and Recipe Steps
-- preparation_media: optional for_ingredient_id allows ingredient-specific media (e.g., peel banana vs peel mango)
-- ingredient_media: media associated with valid ingredients
-- recipe_step_images: step-level images using uploaded_media (parallel to recipe_images)

CREATE TABLE IF NOT EXISTS preparation_media (
    id TEXT NOT NULL PRIMARY KEY,
    valid_preparation_id TEXT NOT NULL REFERENCES valid_preparations(id) ON DELETE CASCADE,
    for_ingredient_id TEXT REFERENCES valid_ingredients(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    index INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE NULLS NOT DISTINCT (valid_preparation_id, for_ingredient_id, uploaded_media_id, archived_at)
);

CREATE TABLE IF NOT EXISTS ingredient_media (
    id TEXT NOT NULL PRIMARY KEY,
    valid_ingredient_id TEXT NOT NULL REFERENCES valid_ingredients(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    index INTEGER DEFAULT 0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(valid_ingredient_id, uploaded_media_id, archived_at)
);

CREATE TABLE IF NOT EXISTS recipe_step_images (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe_step TEXT NOT NULL REFERENCES recipe_steps(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    uploaded_by_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW() NOT NULL,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_preparation_media_preparation ON preparation_media (valid_preparation_id) WHERE archived_at IS NULL;
CREATE INDEX idx_preparation_media_ingredient ON preparation_media (for_ingredient_id) WHERE archived_at IS NULL AND for_ingredient_id IS NOT NULL;

CREATE INDEX idx_ingredient_media_ingredient ON ingredient_media (valid_ingredient_id) WHERE archived_at IS NULL;

CREATE INDEX idx_recipe_step_images_step ON recipe_step_images (belongs_to_recipe_step) WHERE archived_at IS NULL;
