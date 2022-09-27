-- CREATE TYPE prep_step_status AS ENUM ('unfinished', 'postponed', 'canceled', 'finished');

ALTER TYPE prep_step_status RENAME VALUE 'postponed' to 'delayed';
ALTER TYPE prep_step_status ADD VALUE 'ignored';
