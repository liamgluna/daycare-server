package data

import "database/sql"

type StudentModel struct {
	DB *sql.DB
}

type Student struct {
	StudentID   int64  `json:"student_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth Date   `json:"date_of_birth"`
	Gender      string `json:"gender"`
}
