SELECT users.id
  FROM users
 WHERE (users.archived_at IS NULL)
       AND (
			(
				users.last_indexed_at IS NULL
			)
			OR users.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);