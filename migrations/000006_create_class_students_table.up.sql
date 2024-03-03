CREATE TABLE IF NOT EXISTS class_students (
    class_id integer REFERENCES classes(class_id) ON DELETE CASCADE,
    student_id integer REFERENCES students(student_id) ON DELETE CASCADE,
    PRIMARY KEY (class_id, student_id)
);