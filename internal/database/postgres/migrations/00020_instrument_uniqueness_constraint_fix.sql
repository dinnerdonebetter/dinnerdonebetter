ALTER TABLE valid_instruments DROP CONSTRAINT valid_instruments_name_key;
ALTER TABLE valid_instruments ADD CONSTRAINT valid_instruments_name_key unique(name, variant);
