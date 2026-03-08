-- Add contaminates_equipment to valid_ingredients
ALTER TABLE valid_ingredients ADD COLUMN IF NOT EXISTS contaminates_equipment BOOLEAN DEFAULT FALSE NOT NULL;

-- Add source_isbn to recipes
ALTER TABLE recipes ADD COLUMN IF NOT EXISTS source_isbn TEXT DEFAULT ''::TEXT NOT NULL;
