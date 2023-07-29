UPDATE webhooks
   SET last_updated_at = now(), archived_at = now()
 WHERE archived_at IS NULL AND belongs_to_household = $1 AND id = $2;
