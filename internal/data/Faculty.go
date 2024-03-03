package data

import "database/sql"

type FacultyModel struct {
	DB *sql.DB
}

type Faculty struct {
	FacultyID int64  `json:"faculty_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Contact   int64  `json:"contact"`
}
