CREATE TABLE IF NOT EXISTS students (
    student_id serial PRIMARY KEY,
    first_name text NOT NULL,
    last_name text NOT NULL,
    gender text NOT NULL,
    date_of_birth date NOT NULL
);