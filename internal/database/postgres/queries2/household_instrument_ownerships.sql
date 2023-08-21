-- name: ArchiveHouseholdInstrumentOwnership :exec

UPDATE household_instrument_ownerships SET archived_at = NOW() WHERE archived_at IS NULL AND id = $1 AND belongs_to_household = $2;


-- name: CreateHouseholdInstrumentOwnership :exec

INSERT INTO household_instrument_ownerships (id,notes,quantity,valid_instrument_id,belongs_to_household) VALUES ($1,$2,$3,$4,$5);


-- name: CheckHouseholdInstrumentOwnershipExistence :one

SELECT EXISTS ( SELECT household_instrument_ownerships.id FROM household_instrument_ownerships WHERE household_instrument_ownerships.archived_at IS NULL AND household_instrument_ownerships.id = $1 AND household_instrument_ownerships.belongs_to_household = $2 );

-- name: GetHouseholdInstrumentOwnerships :many

SELECT
  household_instrument_ownerships.id,
  household_instrument_ownerships.notes,
  household_instrument_ownerships.quantity,
  valid_instruments.id,
  valid_instruments.name,
  valid_instruments.plural_name,
  valid_instruments.description,
  valid_instruments.icon_path,
  valid_instruments.usable_for_storage,
  valid_instruments.display_in_summary_lists,
  valid_instruments.include_in_generated_instructions,
  valid_instruments.slug,
  valid_instruments.created_at,
  valid_instruments.last_updated_at,
  valid_instruments.archived_at,
  household_instrument_ownerships.belongs_to_household,
  household_instrument_ownerships.created_at,
  household_instrument_ownerships.last_updated_at,
  household_instrument_ownerships.archived_at,
  (
    SELECT
      COUNT(household_instrument_ownerships.id)
    FROM
      household_instrument_ownerships
    WHERE
      household_instrument_ownerships.belongs_to_household = $1
      AND household_instrument_ownerships.archived_at IS NULL
      AND household_instrument_ownerships.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
      AND household_instrument_ownerships.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
      AND (
        household_instrument_ownerships.last_updated_at IS NULL
        OR household_instrument_ownerships.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years'))
      )
      AND (
        household_instrument_ownerships.last_updated_at IS NULL
        OR household_instrument_ownerships.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years'))
      )
  ) as filtered_count,
  (
    SELECT
      COUNT(household_instrument_ownerships.id)
    FROM
      household_instrument_ownerships
    WHERE
      household_instrument_ownerships.belongs_to_household = $1
      AND household_instrument_ownerships.archived_at IS NULL
  ) as total_count
FROM
  household_instrument_ownerships
  JOIN valid_instruments ON valid_instruments.id = household_instrument_ownerships.valid_instrument_id
WHERE household_instrument_ownerships.belongs_to_household = $1
  AND household_instrument_ownerships.archived_at IS NULL
  AND household_instrument_ownerships.created_at > COALESCE($2, (SELECT NOW() - interval '999 years'))
  AND household_instrument_ownerships.created_at < COALESCE($3, (SELECT NOW() + interval '999 years'))
  AND (
    household_instrument_ownerships.last_updated_at IS NULL
    OR household_instrument_ownerships.last_updated_at > COALESCE($4, (SELECT NOW() - interval '999 years'))
  )
  AND (
    household_instrument_ownerships.last_updated_at IS NULL
    OR household_instrument_ownerships.last_updated_at < COALESCE($5, (SELECT NOW() + interval '999 years'))
  )
GROUP BY
  household_instrument_ownerships.id,
  valid_instruments.id
ORDER BY
  household_instrument_ownerships.id
OFFSET $6
LIMIT $7;

-- name: GetHouseholdInstrumentOwnership :one

SELECT
	household_instrument_ownerships.id,
	household_instrument_ownerships.notes,
	household_instrument_ownerships.quantity,
	valid_instruments.id,
	valid_instruments.name,
	valid_instruments.plural_name,
	valid_instruments.description,
	valid_instruments.icon_path,
	valid_instruments.usable_for_storage,
	valid_instruments.display_in_summary_lists,
	valid_instruments.include_in_generated_instructions,
	valid_instruments.slug,
	valid_instruments.created_at,
	valid_instruments.last_updated_at,
	valid_instruments.archived_at,
	household_instrument_ownerships.belongs_to_household,
	household_instrument_ownerships.created_at,
	household_instrument_ownerships.last_updated_at,
	household_instrument_ownerships.archived_at
FROM household_instrument_ownerships
INNER JOIN valid_instruments ON household_instrument_ownerships.valid_instrument_id = valid_instruments.id
WHERE household_instrument_ownerships.archived_at IS NULL
	AND household_instrument_ownerships.id = $1
	AND household_instrument_ownerships.belongs_to_household = $2;


-- name: UpdateHouseholdInstrumentOwnership :exec

UPDATE household_instrument_ownerships
SET
	notes = $1,
	quantity = $2,
	valid_instrument_id = $3,
	last_updated_at = NOW()
WHERE archived_at IS NULL
	AND id = $4
	AND household_instrument_ownerships.belongs_to_household = $5;
