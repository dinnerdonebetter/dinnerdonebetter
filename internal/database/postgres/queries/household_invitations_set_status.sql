UPDATE household_invitations SET
    status = $1,
    status_note = $2,
    last_updated_at = NOW(),
    archived_at = NOW()
WHERE archived_at IS NULL
  AND id = $3;
