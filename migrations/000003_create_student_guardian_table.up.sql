CREATE TABLE IF NOT EXISTS student_guardian (
    student_id integer REFERENCES students(student_id) ON DELETE CASCADE,
    guardian_id integer REFERENCES guardians(guardian_id) ON DELETE CASCADE,
    PRIMARY KEY (student_id, guardian_id)
);