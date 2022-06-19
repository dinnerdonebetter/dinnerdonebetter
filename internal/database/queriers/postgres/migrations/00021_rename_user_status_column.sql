ALTER TABLE users
    RENAME COLUMN reputation TO user_account_status;

ALTER TABLE users
    RENAME COLUMN reputation_explanation TO user_account_status_explanation;
