CREATE TABLE IF NOT EXISTS student_attendance (
    student_id integer REFERENCES students(student_id) ON DELETE CASCADE,
    class_id integer REFERENCES classes(class_id) ON DELETE CASCADE,
    class_date date NOT NULL,
    present bool DEFAULT false,
    PRIMARY KEY (student_id, class_id, class_date)
);