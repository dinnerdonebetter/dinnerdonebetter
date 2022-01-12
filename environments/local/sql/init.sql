/*
Some Notes:

1. All passwords are `PrixFixe1!`, without backticks. 
 */

-- Users

INSERT INTO "users" (
    "id",
    "username",
    "email_address",
    "hashed_password",
    "two_factor_secret",
    "two_factor_secret_verified_on",
    "password_last_changed_on",
    "reputation",
    "service_roles",
    "created_on",
    "archived_on"
)
VALUES
(
    'admin',
    'admin',
    'admin@prixfixe.email',
    '$argon2id$v=19$m=65536,t=1,p=2$QdxGzEsSJc24mMaW4k3kzQ$uqwRs4CuwRJZKAIXjcR9G1V0Qpv38YtL9vK3wm7SZho',
    'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    '1641786300',
    NULL,
    'good',
    'service_admin',
    '1641786189',
    NULL
),
(
    'mom_jones',
    'momJones',
    'mom@jones.com',
    '$argon2id$v=19$m=65536,t=1,p=2$QdxGzEsSJc24mMaW4k3kzQ$uqwRs4CuwRJZKAIXjcR9G1V0Qpv38YtL9vK3wm7SZho',
    'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    '1641786300',
    NULL,
    'good',
    'service_user',
    '1641786189',
    NULL
),
(
    'dad_jones',
    'dadJones',
    'dad@jones.com',
    '$argon2id$v=19$m=65536,t=1,p=2$QdxGzEsSJc24mMaW4k3kzQ$uqwRs4CuwRJZKAIXjcR9G1V0Qpv38YtL9vK3wm7SZho',
    'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    '1641786300',
    NULL,
    'good',
    'service_user',
    '1641786189',
    NULL
),
(
    'sally_jones',
    'sallyJones',
    'sally@jones.com',
    '$argon2id$v=19$m=65536,t=1,p=2$QdxGzEsSJc24mMaW4k3kzQ$uqwRs4CuwRJZKAIXjcR9G1V0Qpv38YtL9vK3wm7SZho',
    'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    '1641786300',
    NULL,
    'good',
    'service_user',
    '1641786189',
    NULL
),
(
    'billy_jones',
    'billyJones',
    'billy@jones.com',
    '$argon2id$v=19$m=65536,t=1,p=2$QdxGzEsSJc24mMaW4k3kzQ$uqwRs4CuwRJZKAIXjcR9G1V0Qpv38YtL9vK3wm7SZho',
    'AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=',
    '1641786300',
    NULL,
    'good',
    'service_user',
    '1641786189',
    NULL
);

-- Households

INSERT INTO "households"
(
    "id",
    "name",
    "contact_email",
    "contact_phone",
    "billing_status",
    "payment_processor_customer_id",
    "subscription_plan_id",
    "belongs_to_user",
    "created_on",
    "last_updated_on",
    "archived_on"
)
VALUES
(
    'adminHousehold',
    'admin_default',
    '',
    '',
    'unpaid',
    '',
    NULL,
    'admin',
    '1641786189',
    NULL,
    NULL
),
(
    'jonesHousehold',
    'The Jones Household',
    '',
    '',
    'paid',
    '',
    NULL,
    'mom_jones',
    '1641786189',
    NULL,
    NULL
);

-- Household  Memberships

INSERT INTO "household_user_memberships" (
    "id",
    "belongs_to_user",
    "belongs_to_household",
    "default_household",
    "household_roles",
    "created_on",
    "last_updated_on",
    "archived_on"
) 
VALUES
(
    'adminMembership',
    'admin',
    'adminHousehold',
    true,
    'household_admin',
    '1641786189',
    NULL,
    NULL
),
(
    'momMembership',
    'mom_jones',
    'jonesHousehold',
    true,
    'household_admin',
    '1641786189',
    NULL,
    NULL
),
(
    'dadMembership',
    'dad_jones',
    'jonesHousehold',
    true,
    'household_member',
    '1641786189',
    NULL,
    NULL
),
(
    'sallyMembership',
    'sally_jones',
    'jonesHousehold',
    true,
    'household_member',
    '1641786189',
    NULL,
    NULL
),
(
    'billyMembership',
    'billy_jones',
    'jonesHousehold',
    true,
    'household_member',
    '1641786189',
    NULL,
    NULL
);

-- Valid Ingredients

INSERT INTO "valid_ingredients" (
    "id",
    "name",
    "variant",
    "description",
    "icon_path",
    "warning",
    "animal_derived",
    "animal_flesh",
    "contains_dairy",
    "contains_egg",
    "contains_fish",
    "contains_gluten",
    "contains_peanut",
    "contains_sesame",
    "contains_shellfish",
    "contains_soy",
    "contains_tree_nut",
    "contains_wheat",
    "volumetric",
    "created_on",
    "last_updated_on",
    "archived_on"
)
VALUES
(
    'chickenBreast',
    'chicken breast',
    '',
    '',
    '',
    '',
    true,
    true,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    '1641866784',
    NULL,
    NULL
),
(
    'coffee',
    'coffee',
    '',
    '',
    '',
    '',
    true,
    true,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    false,
    '1641866784',
    NULL,
    NULL
);

-- Valid Instruments

INSERT INTO "valid_instruments"
(
    "id",
    "name",
    "variant",
    "description",
    "icon_path",
    "created_on",
    "last_updated_on",
    "archived_on"
)
VALUES
(
    'spoon',
    'spoon',
    '',
    '',
    '',
    '1641870035',
    NULL,
    NULL
);

-- Valid Preparations

INSERT INTO "valid_preparations"
(
    "id",
    "name",
    "description",
    "icon_path",
    "created_on",
    "last_updated_on",
    "archived_on"
)
VALUES
(
    'grill',
    'grill',
    '',
    '',
    '1641866834',
    NULL,
    NULL
);

-- Recipes

INSERT INTO "recipes"
(
    "id",
    "name",
    "description",
    "inspired_by_recipe_id",
    "source",
    "created_by_user",
    "created_on",
    "last_updated_on",
    "archived_on"
)
VALUES
(
    'grilledChickenBreast',
    'grilled chicken breast',
    '',
    '',
    '',
    'mom_jones',
    '1641867047',
    NULL,
    NULL
),
(
    'friedChickenBreast',
    'fried chicken breast',
    '',
    '',
    '',
    'mom_jones',
    '1641867048',
    NULL,
    NULL
);
