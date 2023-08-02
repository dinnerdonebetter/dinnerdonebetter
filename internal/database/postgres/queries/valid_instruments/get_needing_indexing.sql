-- name: GetValidInstrumentsNeedingIndexing :many

SELECT valid_instruments.id
  FROM valid_instruments
 WHERE (valid_instruments.archived_at IS NULL)
       AND (
			(valid_instruments.last_indexed_at IS NULL)
			OR valid_instruments.last_indexed_at
				< now() - '24 hours'::INTERVAL
		);