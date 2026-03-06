-- User Avatars Migration
-- Replace avatar_src (base64) with user_avatars join table linking to uploaded_media

CREATE TABLE IF NOT EXISTS user_avatars (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(belongs_to_user, archived_at)
);

CREATE TABLE IF NOT EXISTS recipe_images (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_recipe TEXT NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    uploaded_by_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS meal_images (
    id TEXT NOT NULL PRIMARY KEY,
    belongs_to_meal TEXT NOT NULL REFERENCES meals(id) ON DELETE CASCADE,
    uploaded_media_id TEXT NOT NULL REFERENCES uploaded_media(id) ON DELETE CASCADE,
    uploaded_by_user TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    archived_at TIMESTAMP WITH TIME ZONE
);

-- Drop legacy avatar_src column from users
ALTER TABLE users DROP COLUMN IF EXISTS avatar_src;
