SELECT EXISTS (
   SELECT household_user_memberships.id
	 FROM household_user_memberships
	WHERE household_user_memberships.archived_at IS NULL
	 AND household_user_memberships.belongs_to_household = $1
	 AND household_user_memberships.belongs_to_user = $2
);
