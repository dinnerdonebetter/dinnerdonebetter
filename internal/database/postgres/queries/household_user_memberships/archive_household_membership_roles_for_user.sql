UPDATE household_membership_roles SET archived_at = NOW() WHERE belongs_to_household_user_membership = $1;
