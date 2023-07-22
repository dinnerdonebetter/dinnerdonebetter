SELECT valid_preparations.id
  FROM valid_preparations
 WHERE (valid_preparations.archived_at IS NULL)
       AND (
			(valid_preparations.last_indexed_at IS NULL)
			OR valid_preparations.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);