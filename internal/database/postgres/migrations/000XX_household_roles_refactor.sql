ALTER TABLE household_user_memberships DROP COLUMN "household_roles";

CREATE TYPE household_role AS ENUM (
    'household_admin',
    'household_member'
);

CREATE TABLE IF NOT EXISTS household_membership_roles (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "role" household_role NOT NULL,
    "belongs_to_household_user_membership" CHAR(27) NOT NULL REFERENCES household_user_memberships("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "archived_at" TIMESTAMP WITH TIME ZONE,
    UNIQUE("role", "belongs_to_household_user_membership", "archived_at")
);
