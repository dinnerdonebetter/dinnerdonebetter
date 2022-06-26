ALTER TABLE valid_ingredients DROP COLUMN IF EXISTS "variant";
ALTER TABLE valid_ingredients DROP COLUMN IF EXISTS "animal_derived";

ALTER TABLE valid_ingredients ADD COLUMN "is_liquid" BOOLEAN default 'false';