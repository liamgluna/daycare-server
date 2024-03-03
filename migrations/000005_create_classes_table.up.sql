CREATE TABLE IF EXISTS classes (
    class_id serial PRIMARY KEY,
    faculty_id integer REFERENCES faculty(faculty_id),
    term text NOT NULL,
);