CREATE TABLE student_guardian (
    student_id integer REFERENCES students(student_id),
    guardian_id integer REFERENCES guardians(guardian_id),
    PRIMARY KEY (student_id, guardian_id)
);