package data

import "database/sql"

type ClassesModel struct {
	DB *sql.DB
}

type Classes struct {
	ClassID   int64  `json:"class_id"`
	FacultyID int64  `json:"faculty_id"`
	Term      string `json:"term"`
}
