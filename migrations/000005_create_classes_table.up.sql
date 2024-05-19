CREATE TABLE IF NOT EXISTS classes (
    class_id serial PRIMARY KEY,
    class_name text NOT NULL,
    schedule text NOT NULL,
    faculty_id integer REFERENCES faculty(faculty_id) ON DELETE CASCADE,
    term text NOT NULL
);