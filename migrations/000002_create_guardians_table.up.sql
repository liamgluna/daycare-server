CREATE TABLE IF NOT EXISTS guardians (
    guardian_id serial PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    gender text NOT NULL,
    relationship text NOT NULL,
    occupation text NOT NULL,
    contact bigint NOT NULL,
    student_id integer REFERENCES students(student_id) ON DELETE CASCADE
);