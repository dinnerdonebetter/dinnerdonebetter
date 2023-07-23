ALTER TABLE webhooks DROP COLUMN "topics";
ALTER TABLE webhooks DROP COLUMN "events";
ALTER TABLE webhooks DROP COLUMN "data_types";

CREATE TYPE webhook_event AS ENUM (
    'webhook_created',
    'webhook_updated',
    'webhook_archived'
);

CREATE TABLE IF NOT EXISTS webhook_trigger_events (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "trigger_event" webhook_event NOT NULL,
    "belongs_to_webhook" CHAR(27) NOT NULL REFERENCES webhooks("id") ON DELETE CASCADE,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    "archived_at" TIMESTAMP WITH TIME ZONE,
    UNIQUE("trigger_event", "belongs_to_webhook")
);