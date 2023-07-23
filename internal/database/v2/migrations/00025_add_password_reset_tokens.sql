CREATE TABLE IF NOT EXISTS password_reset_tokens (
    "id" CHAR(27) NOT NULL PRIMARY KEY,
    "token" TEXT NOT NULL,
    "expires_at" BIGINT NOT NULL,
    "created_on" BIGINT NOT NULL DEFAULT extract(epoch FROM NOW()),
    "last_updated_on" BIGINT DEFAULT NULL,
    "redeemed_on" BIGINT DEFAULT NULL,
    "belongs_to_user" CHAR(27) NOT NULL REFERENCES users("id") ON DELETE CASCADE
);