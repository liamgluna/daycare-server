package data

import "database/sql"

type ClassStudentsModel struct {
	DB *sql.DB
}

type ClassStudents struct {
	ClassID   int64 `json:"class_id"`
	StudentID int64 `json:"student_id"`
}
