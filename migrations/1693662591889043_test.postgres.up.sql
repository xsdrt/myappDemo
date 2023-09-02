CREATE TABLE some_table (
    id serial PRIMARY KEY,
    some_field VARCHAR ( 255 ) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- add auto update of updated_at. If you already have this trigger
-- you can delete the next 7 lines
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER set_timestamp
    BEFORE UPDATE ON some_table
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();