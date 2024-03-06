CREATE TABLE IF NOT EXISTS faculty (
    faculty_id serial PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    email citext UNIQUE NOT NULL,
    contact bigint NOT NULL,
    position text NOT NULL
);