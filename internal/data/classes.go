package data

import "database/sql"

type ClassesModel struct {
	DB *sql.DB
}

type Classes struct {
	ClassID   int64  `json:"class_id"`
	ClassName string `json:"class_name"`
	FacultyID int64  `json:"faculty_id"`
	Term      string `json:"term"`
}
